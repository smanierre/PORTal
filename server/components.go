package server

import (
	"PORTal/templates/components"
	"PORTal/templates/components/webcomponents"
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
	quals, err := s.backend.GetAllQualifications()
	if err != nil {
		err = components.Toast("Unable to get qualifications for requirement editor.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast for requirement editor", slog.String("error", err.Error()))
		}
		return
	}
	err = webcomponents.RequirementEditor(webcomponents.RequirementEditorData{
		SelectedRequirement: types.Requirement{},
		NewRequirement:      false,
		Qualifications:      quals,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirementEditor component", slog.String("error", err.Error()))
	}
}
