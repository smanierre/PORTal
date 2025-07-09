package rendered

import (
	"PORTal/server/core"
	"PORTal/templates"
	"PORTal/templates/pages/admin"
	"PORTal/templates/pages/errorpages"
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

func AdminGetHandler(logger *slog.Logger, c core.Core) http.Handler {
	logger = logger.With(slog.String("Source", "AdminGetHandler"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panel := r.URL.Query().Get("fragment")
		var content templ.Component
		var dropdownValue string
		var err error
		switch panel {
		case "qualifications_pane":
			var data core.AdminQualificationData
			dropdownValue = "Qualifications"
			data, err = c.AdminQualificationPage(r)
			if err != nil {
				break
			}
			content = admin.QualificationPane(data, false)
		default:
			var data core.AdminMembersData
			dropdownValue = "Members"
			data, err = c.AdminMemberPage(r)
			if err != nil {
				break

			}
			content = admin.MembersPane(data, false)
		}
		if r.Header.Get("Hx-Request") != "true" {
			tld, err := c.TopLevel(r, nil)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error getting top level data", slog.String("error", err.Error()))
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
				return
			}
			err = templates.Root(tld, admin.AdminPage(dropdownValue, content)).Render(r.Context(), w)
			if err != nil {
				logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering complete admin page", slog.String("error", err.Error()))
				RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
			}
			return
		}
		err = admin.AdminPage(dropdownValue, content).Render(r.Context(), w)
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin page fragment", slog.String("error", err.Error()))
			RenderErrorPage(w, r, logger, http.StatusInternalServerError, errorpages.GenericISE())
		}
	})
}
