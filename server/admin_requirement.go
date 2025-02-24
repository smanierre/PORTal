package server

import (
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"PORTal/templates/pages/errorpages"
	"PORTal/types"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (s Server) AdminRequirementEditorGetHandler(w http.ResponseWriter, r *http.Request) {
	requirement, err := s.backend.GetRequirement(r.PathValue("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		err = components.Toast("Unable to get requirement.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}

	requirementEditor := admin.RequirementEditor(admin.RequirementEditorData{
		SelectedRequirement: requirement,
		NewRequirement:      false,
	})
	w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/requirements/%s", r.PathValue("id")))
	err = requirementEditor.Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin requirement editor fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminRequirementUpdateHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Updating requirement")
	id := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing form", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w) // TODO: Probably make this cleaner?
	}
	daysValidFor, err := strconv.Atoi(r.Form.Get("days_valid_for"))
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid value parsed for days valid for", slog.String("error", err.Error()))
		err = components.Toast("Invalid value provided for Expiration Days", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	req := types.Requirement{
		ID:           id,
		Name:         r.Form.Get("name"),
		Reference:    r.Form.Get("reference"),
		Notes:        r.Form.Get("notes"),
		DaysValidFor: daysValidFor,
		Type:         types.RequirementType(r.Form.Get("type")),
	}

	updatedReq, err := s.backend.UpdateRequirement(req)
	if err != nil {
		err = components.Toast("Unable to update requirement", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}

	err = admin.RequirementEditor(admin.RequirementEditorData{
		SelectedRequirement: updatedReq,
		NewRequirement:      false,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering updated requirement", slog.String("error", err.Error()))
		// TODO: Figure out toast for the modal?
	}
}

func (s Server) AdminNewRequirementHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX get request for new requirement, redirecting to qualification pane")
		http.Redirect(w, r, "/admin/qualifications", http.StatusFound)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty requirement editor")
	err := admin.RequirementEditor(admin.RequirementEditorData{
		SelectedRequirement: types.Requirement{},
		NewRequirement:      true,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering new requirement page/fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminRequirementAddHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new requirement")
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new requirement", slog.String("error", err.Error()))
		err = components.Toast("Unable to add new requirement.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	req := types.Requirement{}
	req.Name = r.Form.Get("name")
	req.Notes = r.Form.Get("notes")
	req.Type = types.RequirementType(r.Form.Get("type"))
	req.DaysValidFor, err = strconv.Atoi(r.Form.Get("days_valid_for"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = components.Toast("Unable to parse days valid for field, it must be a number.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	req.Reference = r.Form.Get("reference")
	req.Type = types.RequirementType(r.Form.Get("type"))

	req, err = s.backend.AddRequirement(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = components.Toast("Unable to create Requirement.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}

	err = components.RequirementItem(req.Name, req.ID).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirement item", slog.String("error", err.Error()))
	}
}

func (s Server) AdminRequirementDeleteHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Handling request to remove requirement")
	id := r.PathValue("id")
	err := s.backend.DeleteRequirement(id)
	if err != nil {
		err = components.Toast("Unable to remove requirement.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
	}
}
