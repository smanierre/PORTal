package rendered

import (
	"PORTal/server/core"
	"PORTal/templates"
	"PORTal/templates/pages/dashboard"
	"PORTal/templates/pages/errorpages"
	"log/slog"
	"net/http"
)

func dashboardGetHandler(logger *slog.Logger, c core.Core) http.Handler {
	logger = logger.With(slog.String("route", "GET /dashboard"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, err := c.Dashboard(r)
		if err != nil {
			RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
			return
		}
		content := dashboard.Dashboard(data)
		if r.Header.Get("Hx-Request") != "true" {
			tld, err := c.TopLevel(r, nil)
			if err != nil {
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
				return
			}
			err = templates.Root(tld, content).Render(r.Context(), w)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering root template with dashboard", slog.String("error", err.Error()))
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
				return
			}
		} else {
			err = content.Render(r.Context(), w)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering dashboard fragment", slog.String("error", err.Error()))
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
			}
		}
	})
}
