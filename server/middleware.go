package server

import (
	"PORTal/backend"
	"PORTal/templates"
	"PORTal/templates/pages"
	"PORTal/templates/pages/errorpages"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

func (s Server) skipLoginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == http.MethodGet {
			s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Running skipLoginMiddleware")
			sessionCookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, not skipping login")
				next.ServeHTTP(w, r)
				return
			}
			m, err := s.backend.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelDebug, "Invalid session, not skipping login")
				next.ServeHTTP(w, r)
				return
			}
			s.logger.LogAttrs(r.Context(), slog.LevelDebug, "Valid session redirecting to the dashboard")
			r = r.WithContext(context.WithValue(r.Context(), MemberContextKey, m))
			http.Redirect(w, r, "/dashboard", http.StatusTemporaryRedirect)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func (s Server) sessionRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Anything to the login page doesn't need to be authenticated
		if r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		// Check for session
		sessionCookie, err := r.Cookie(SessionCookieName)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelDebug, "No session cookie found, setting cookie to empty cookie")
			sessionCookie = &http.Cookie{}
		}
		m, err := s.backend.ValidateSession(sessionCookie.Value, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
		if errors.Is(err, backend.ErrSessionValidationFailed) {
			s.logger.LogAttrs(r.Context(), slog.LevelDebug, "Session validation failed, removing cookie and redirecting to login")
			removeCookie(w, SessionCookieName, s.config.Domain)
			w.Header().Set("HX-Push-URL", "/")
			if checkHTMXRequest(r) {
				err = pages.Login().Render(r.Context(), w)
			} else {
				http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
				return
			}
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Failed to render login", slog.String("error", err.Error()))
			}
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), MemberContextKey, m))
		next.ServeHTTP(w, r)
	})
}

func (s Server) assetMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/assets/") {
			s.mux.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s Server) adminRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/admin") {
			s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Admin not required, skipping")
			next.ServeHTTP(w, r)
			return
		}
		// Runs after session validation middleware, so no need to verify session, member should be in context
		m, err := getMemberFromContext(r.Context())
		if err != nil || !m.Admin {
			s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Member isn't in context or isn't admin, rendering the dashboard")
			// Ensure that member doesn't have the admin option in the Nav
			var err error
			if checkHTMXRequest(r) {
				err = pages.Dashboard().Render(r.Context(), w)
			} else {
				http.Redirect(w, r, "/dashboard", http.StatusFound)
			}
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard fragment", slog.String("error", err.Error()))
				_ = errorpages.GenericISE().Render(r.Context(), w)
			}
			var subordinateLength int
			subordinates, err := s.backend.GetSubordinates(m.ID)
			if err == nil {
				subordinateLength = len(subordinates)
			}

			if !checkHTMXRequest(r) {
				err = templates.Nav(templates.NavData{
					Show:         true,
					OobSwap:      true,
					Member:       m,
					Subordinates: subordinateLength > 0,
				}).Render(r.Context(), w)
				if err != nil {
					s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering nav OOB", slog.String("error", err.Error()))
				}
			}
			return
		}
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Member is admin, continuing on")
		next.ServeHTTP(w, r)
	})
}
