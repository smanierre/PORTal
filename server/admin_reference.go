package server

import (
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"log/slog"
	"net/http"
)

func (s Server) AdminReferenceGetHandler(w http.ResponseWriter, r *http.Request) {
	refs, err := s.backend.GetReferences()
	if err != nil {
		if checkHTMXRequest(r) {
			err = s.templateRepo.RenderFragment(w, "admin_references", "toast", components.ToastData{
				Message: "Unable to get references.",
				Danger:  true,
			})
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			}
			return
		} else {
			err = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering error page", slog.String("error", err.Error()))
			}
			return
		}
	}
	data := admin.ReferenceContentData{
		References:          refs,
		SwapTarget:          "#member-content",
		FragmentBasePath:    "admin/references",
		ReferenceEditorData: admin.ReferenceEditorData{},
	}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/references")
		err = s.templateRepo.RenderFragment(w, "admin_references", "references", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin references fragment", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		}
		return
	}
	m, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error getting member from context", slog.String("error", err.Error()))
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	tplData := &templates.TplData{
		NavData: templates.NavData{
			Show:         true,
			OobSwap:      false,
			Member:       m,
			Subordinates: false,
		},
		ContentData: admin.ReferenceRootData{
			DropdownData: components.DropdownData{
				Items:            adminDropdownItems,
				SwapTarget:       "#admin-content",
				FragmentBasePath: "admin",
			},
			ReferencesData: data,
		},
	}
	err = s.templateRepo.Render(w, "admin_references", tplData)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin reference page", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "errors", "generic_ise", nil)
	}
}
