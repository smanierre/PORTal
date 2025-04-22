package server

import (
	"PORTal/backend"
	"PORTal/server/stores"
	"PORTal/templates/pages"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

func skipLoginMiddleware(next http.Handler, logger *slog.Logger, sessionStore stores.SessionStore) http.Handler {
	logger = logger.With(slog.String("source", "skipLoginMiddleware"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Running skipLoginMiddleware")
			sessionCookie, err := r.Cookie(SessionCookieName)
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
			r = r.WithContext(context.WithValue(r.Context(), MemberContextKey, m))
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
		sessionCookie, err := r.Cookie(SessionCookieName)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, setting cookie to empty cookie")
			sessionCookie = &http.Cookie{}
		}
		m, err := sessionStore.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
		if errors.Is(err, backend.ErrSessionValidationFailed) {
			logger.LogAttrs(r.Context(), slog.LevelDebug, "Session validation failed, removing cookie and redirecting to login")
			RemoveCookie(w, SessionCookieName)
			w.Header().Set("HX-Push-URL", "/")
			if CheckHTMXRequest(r) {
				err = pages.Login(organization).Render(r.Context(), w)
			} else {
				http.Redirect(w, r, "/", http.StatusFound)
				return
			}
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Failed to render login", slog.String("error", err.Error()))
			}
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), MemberContextKey, m))
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

func adminRequiredMiddleware(next http.Handler, logger *slog.Logger) http.Handler {
	logger = logger.With(slog.String("source", "adminRequiredMiddleware"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/admin") {
			next.ServeHTTP(w, r)
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Request is for admin page, running middleware")
		// Runs after session validation middleware, so no need to verify session, member should be in context
		m, err := MemberFromContext(r.Context())
		if err != nil || !m.Admin {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Member isn't in context or isn't admin, redirecting to the dashboard")
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Member is admin, continuing on")
		next.ServeHTTP(w, r)
	})
}
