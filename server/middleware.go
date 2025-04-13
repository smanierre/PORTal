package server

import (
	"PORTal/backend"
	"PORTal/server/serverutils"
	"PORTal/server/stores"
	"PORTal/templates"
	"PORTal/templates/pages"
	"PORTal/templates/pages/dashboard"
	"PORTal/templates/pages/errorpages"
	"PORTal/types"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func skipLoginMiddleware(next http.Handler, logger *slog.Logger, sessionStore stores.SessionStore) http.Handler {
	logger = logger.With(slog.String("source", "skipLoginMiddleware"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Running skipLoginMiddleware")
			sessionCookie, err := r.Cookie(serverutils.SessionCookieName)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, not skipping login")
				next.ServeHTTP(w, r)
				return
			}
			m, err := sessionStore.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelDebug, "Invalid session, not skipping login")
				next.ServeHTTP(w, r)
				return
			}
			logger.LogAttrs(r.Context(), slog.LevelDebug, "Valid session redirecting to the dashboard")
			r = r.WithContext(context.WithValue(r.Context(), serverutils.MemberContextKey, m))
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func sessionRequiredMiddleware(next http.Handler, logger *slog.Logger, organization string, sessionStore stores.SessionStore) http.Handler {
	logger = logger.With(slog.String("source", "sessionRequiredMiddleware"))
	// update this and add it to the chain in server.go
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Anything to the login page doesn't need to be authenticated
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		// Check for session
		sessionCookie, err := r.Cookie(serverutils.SessionCookieName)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, setting cookie to empty cookie")
			sessionCookie = &http.Cookie{}
		}
		m, err := sessionStore.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
		if errors.Is(err, backend.ErrSessionValidationFailed) {
			logger.LogAttrs(r.Context(), slog.LevelDebug, "Session validation failed, removing cookie and redirecting to login")
			serverutils.RemoveCookie(w, serverutils.SessionCookieName)
			w.Header().Set("HX-Push-URL", "/")
			if serverutils.CheckHTMXRequest(r) {
				err = pages.Login(organization).Render(r.Context(), w)
			} else {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Failed to render login", slog.String("error", err.Error()))
			}
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), serverutils.MemberContextKey, m))
		next.ServeHTTP(w, r)
	})
}

func assetMiddleware(next http.Handler, assetHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/assets/") {
			assetHandler.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func adminRequiredMiddleware(next http.Handler, ranks types.RankMap, memberStore stores.MemberStore, logger *slog.Logger) http.Handler {
	logger = logger.With(slog.String("source", "adminRequiredMiddleware"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/admin") {
			next.ServeHTTP(w, r)
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Request is for admin page, running middleware")
		// Runs after session validation middleware, so no need to verify session, member should be in context
		m, err := serverutils.GetMemberFromContext(r.Context())
		if err != nil || !m.Admin {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Member isn't in context or isn't admin, rendering the dashboard")
			// Ensure that member doesn't have the admin option in the Nav
			var err error
			if serverutils.CheckHTMXRequest(r) {
				err = dashboard.Dashboard().Render(r.Context(), w)
			} else {
				http.Redirect(w, r, "/dashboard", http.StatusFound)
			}
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard fragment", slog.String("error", err.Error()))
				_ = errorpages.GenericISE().Render(r.Context(), w)
			}
			var subordinateLength int
			subordinates, err := memberStore.GetSubordinates(m.ID)
			if err == nil {
				subordinateLength = len(subordinates)
			}

			displayName := fmt.Sprintf("%s %s %s", ranks[m.Grade], m.FirstName, m.LastName)
			if !serverutils.CheckHTMXRequest(r) {
				err = templates.Nav(
					true,
					true,
					subordinateLength > 0,
					m.Admin,
					displayName,
				).Render(r.Context(), w)
				if err != nil {
					logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering nav OOB", slog.String("error", err.Error()))
				}
			}
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Member is admin, continuing on")
		next.ServeHTTP(w, r)
	})
}
