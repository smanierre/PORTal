package server

import (
	"PORTal/templates"
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
	err = s.templateRepo.RenderFragment(w, "nav", "nav", templates.NavData{
		OobSwap: true,
	})
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering nav OOB", slog.String("error", err.Error()))
	}
	err = s.templateRepo.RenderFragment(w, "login", "content", templates.LoginData{
		Organization: s.config.Organization,
	})
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering login fragment", slog.String("error", err.Error()))
	}
}
