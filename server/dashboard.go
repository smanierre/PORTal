package server

import (
	"PORTal/templates/pages/dashboard"
	"log/slog"
	"net/http"
)

func dashboardGetHandler(logger *slog.Logger) http.Handler {
	logger = logger.With(slog.String("route", "GET /dashboard"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dashboard := dashboard.Dashboard()
		if r.Header.Get("Hx-Request") != "true" {
			serveContentAsRoot(dashboard, logger, w, r)
		} else {
			HandleRenderError(r.Context(), logger, dashboard.Render(r.Context(), w))
		}
	})
}
