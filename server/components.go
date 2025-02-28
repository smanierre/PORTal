package server

import (
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"PORTal/types"
	"log/slog"
	"net/http"
)

func (s Server) RequirementItemComponent(w http.ResponseWriter, r *http.Request) {
	err := components.RequirementItem("", "").Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirementItem component", slog.String("error", err.Error()))
	}
}

func (s Server) RequirementEditorComponent(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty requirement editor")
	err := admin.RequirementEditor(admin.RequirementEditorData{
		SelectedRequirement: types.Requirement{},
		NewRequirement:      false,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirementEditor component", slog.String("error", err.Error()))
	}
}
