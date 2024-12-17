package server

import (
	"PORTal/backend"
	"PORTal/templates"
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
				err = s.templateRepo.RenderFragment(w, "login", "content", templates.LoginData{Organization: s.config.Organization})
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
