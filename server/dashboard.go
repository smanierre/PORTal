package server

import (
	"PORTal/templates"
	"PORTal/templates/pages/dashboard"
	"PORTal/types"
	"log/slog"
	"net/http"
)

func dashboardGetHandler(logger *slog.Logger, organization string) http.Handler {
	logger = logger.With(slog.String("route", "GET /dashboard"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := r.Context().Value(MemberContextKey)
		member, ok := m.(types.Member)
		if !ok {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "unable to get member from context, redirecting to login")
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		rootTemplate := templates.Root(
			true,
			false, member.Admin, member.Display(), organization, dashboard.Dashboard())

		HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
	})
}
