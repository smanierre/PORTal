package server

import (
	"PORTal/templates"
	"PORTal/templates/pages"
	"log/slog"
	"net/http"
)

func (s Server) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := getSessionId(r)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "User not logged in, rendering login page")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	s.backend.DeleteSession(sessionID)
	removeCookie(w, SessionCookieName, s.config.Domain)
	if !checkHTMXRequest(r) {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	w.Header().Set("HX-Push-URL", "/")
	err = templates.Nav(templates.NavData{
		OobSwap: true,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering nav OOB", slog.String("error", err.Error()))
	}
	err = pages.Login().Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login fragment", slog.String("error", err.Error()))
	}
}
