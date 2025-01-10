package server

import (
	"PORTal/backend"
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages"
	"PORTal/templates/pages/errorpages"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

func (s Server) LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	loginTpl := pages.Login()
	hx := checkHTMXRequest(r)
	if !hx {
		err := templates.Root(templates.NavData{}, loginTpl).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login template: %s", slog.String("error", err.Error()))
		}
	} else {
		err := loginTpl.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login fragment: %s", slog.String("error", err.Error()))
		}
	}
}

func (s Server) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Non-HTMX request on login post endpoint, unauthorized")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
	}
	m, err := s.backend.Login(r.FormValue("username"), r.FormValue("password"))
	if errors.Is(err, backend.ErrAuthenticationFailed) {
		w.Header().Set("HX-Retarget", "#loginError")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.WriteHeader(http.StatusUnauthorized)
		err = pages.LoginAuthError().Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login error fragment", slog.String("error", err.Error()))
		}
		return
	}
	if err != nil {
		w.Header().Set("HX-Retarget", "#loginError")
		w.Header().Set("HX-Reswap", "outerHTML")
		w.WriteHeader(http.StatusInternalServerError)
		err = pages.LoginInternalError().Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering ISE error fragment", slog.String("error", err.Error()))
		}
		return
	}
	sessionId, expiration := s.backend.CreateSession(m.ID, r.UserAgent(), strings.Split(r.RemoteAddr, ":")[0])
	w.Header().Set("HX-Push-URL", "/dashboard")
	http.SetCookie(w, s.makeCookie(SessionCookieName, sessionId, expiration))
	err = pages.Dashboard().Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard fragment", slog.String("error", err.Error()))
	}
	subordinates, err := s.backend.GetSubordinates(m.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = errorpages.GenericISE().Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering generic ise error fragment", slog.String("error", err.Error()))
		}
	}
	err = templates.Nav(templates.NavData{
		Show:         true,
		OobSwap:      true,
		Member:       m,
		Subordinates: len(subordinates) > 0,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering Nav OOB", slog.String("error", err.Error()))
		_ = components.Toast("Error updating nav, please refresh.", true).Render(r.Context(), w)
	}
}
