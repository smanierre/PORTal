package rendered

import (
	"PORTal/server/core"
	"log/slog"
	"net/http"
)

func LogoutHandler(logger *slog.Logger, c core.Core) http.Handler {
	logger = logger.With(slog.String("route", "GET /logout"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.Logout(r)
		c.RemoveSessionCookie(w)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})
}
