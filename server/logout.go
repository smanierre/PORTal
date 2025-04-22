package server

import (
	"PORTal/server/stores"
	"log/slog"
	"net/http"
)

func LogoutHandler(logger *slog.Logger, sessionStore stores.SessionStore) http.Handler {
	logger = logger.With(slog.String("route", "GET /logout"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := GetSessionId(r)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "User not logged in, rendering login page")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		sessionStore.DeleteSession(sessionID)
		RemoveCookie(w, SessionCookieName)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
