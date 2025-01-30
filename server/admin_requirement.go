package server

import (
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"PORTal/templates/pages/errorpages"
	"PORTal/types"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

func (s Server) AdminRequirementsPaneGetHandler(w http.ResponseWriter, r *http.Request) {
	reqs, err := s.backend.GetAllRequirements()
	if err != nil {
		if checkHTMXRequest(r) {
			err = components.Toast("Unable to get requirements.", true).Render(r.Context(), w)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			}
			return
		} else {
			err = errorpages.GenericISE().Render(r.Context(), w)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering error page", slog.String("error", err.Error()))
			}
			return
		}
	}
	requirementsPane := admin.RequirementsPane(admin.RequirementPaneData{
		SelectedRequirement: types.Requirement{},
		Requirements:        reqs,
		OobSwap:             false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/requirements")
		err = requirementsPane.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin requirements fragment", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	m, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error getting member from context", slog.String("error", err.Error()))
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	subordinates, err := s.backend.GetSubordinates(m.ID)
	var subordinateLength int
	if err == nil {
		subordinateLength = len(subordinates)
	}
	err = templates.Root(templates.NavData{
		Show:         true,
		OobSwap:      false,
		Member:       m,
		Subordinates: subordinateLength > 0,
	}, admin.AdminPage("Requirements", requirementsPane)).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin reference page", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminRequirementEditorGetHandler(w http.ResponseWriter, r *http.Request) {
	requirement, err := s.backend.GetRequirement(r.PathValue("id"))
	if err != nil {
		err = components.Toast("Unable to get requirement.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	references, err := s.backend.GetReferences()
	if err != nil {
		if checkHTMXRequest(r) {
			err = components.Toast("Unable to get reference list.", true).Render(r.Context(), w)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
				_ = errorpages.GenericISE().Render(r.Context(), w)
			}
			return
		}
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
	requirementEditor := admin.RequirementEditor(admin.RequirementEditorData{
		SelectedRequirement: requirement,
		References:          references,
		NewRequirement:      false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/requirements/%s", r.PathValue("id")))
		err = requirementEditor.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin requirement editor fragment", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	m, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "error getting member from context", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
	requirements, err := s.backend.GetAllRequirements()
	if err != nil {
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
	subordinates, err := s.backend.GetSubordinates(m.ID)
	var subordinateLength int
	if err == nil {
		subordinateLength = len(subordinates)
	}
	err = templates.Root(templates.NavData{
		Show:         true,
		OobSwap:      false,
		Member:       m,
		Subordinates: subordinateLength > 0,
	}, admin.AdminPage("Requirements", admin.RequirementsPane(admin.RequirementPaneData{
		SelectedRequirement: requirement,
		Requirements:        requirements,
		OobSwap:             false,
		RequirementEditorData: admin.RequirementEditorData{
			SelectedRequirement: requirement,
			References:          references,
			NewRequirement:      false,
		},
	}))).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirement editor page", slog.String("error", err.Error()))
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
		Reference:    types.Reference{ID: r.Form.Get("reference")},
		Notes:        r.Form.Get("notes"),
		DaysValidFor: daysValidFor,
	}

	updatedReq, err := s.backend.UpdateRequirement(req)
	if err != nil {
		err = components.Toast("Unable to update requirement", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}

	requirements, err := s.backend.GetAllRequirements()
	if err != nil {
		err = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}

	references, err := s.backend.GetReferences()
	if err != nil {
		err = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}

	err = admin.RequirementsPane(admin.RequirementPaneData{
		SelectedRequirement: updatedReq,
		Requirements:        requirements,
		OobSwap:             false,
		RequirementEditorData: admin.RequirementEditorData{
			SelectedRequirement: updatedReq,
			References:          references,
			NewRequirement:      false,
		},
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirements pane", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
	err = components.Toast("Successfully updated requirement", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirement updated toast", slog.String("error", err.Error()))
	}
}

func (s Server) AdminNewRequirementHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty requirement editor")
	references, err := s.backend.GetReferences()
	if checkHTMXRequest(r) {
		if err != nil {
			err = components.Toast("Unable to get references for new requirement.", true).Render(r.Context(), w)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			}
			return
		}
		w.Header().Set("HX-Push-URL", "/admin/requirements/add")
		err = admin.RequirementEditor(admin.RequirementEditorData{
			SelectedRequirement: types.Requirement{},
			References:          references,
			NewRequirement:      true,
		}).Render(r.Context(), w)
	} else {
		if err != nil {
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		m, err := getMemberFromContext(r.Context())
		if err != nil {
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
		requirements, err := s.backend.GetAllRequirements()
		if err != nil {
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
		subordinates, err := s.backend.GetSubordinates(m.ID)
		var subordinateLength int
		if err == nil {
			subordinateLength = len(subordinates)
		}
		err = templates.Root(templates.NavData{
			Show:         true,
			OobSwap:      false,
			Member:       m,
			Subordinates: subordinateLength > 0,
		}, admin.AdminPage("Requirements", admin.RequirementsPane(admin.RequirementPaneData{
			SelectedRequirement: types.Requirement{},
			Requirements:        requirements,
			OobSwap:             false,
			RequirementEditorData: admin.RequirementEditorData{
				SelectedRequirement: types.Requirement{},
				References:          references,
				NewRequirement:      true,
			},
		}))).Render(r.Context(), w)
	}
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering new requirement page/fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminRequirementAddHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to add requirement, redirecting to the dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new requirement")
	err := r.ParseForm()
	if err != nil {
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
	req.DaysValidFor, err = strconv.Atoi(r.Form.Get("days_valid_for"))
	if err != nil {
		err = components.Toast("Unable to parse days valid for field, it must be a number.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	req.Reference = types.Reference{ID: r.Form.Get("reference")}

	req, err = s.backend.AddRequirement(req)
	if err != nil {
		err = components.Toast("Unable to create Requirement.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	reqs, err := s.backend.GetAllRequirements()
	if err != nil {
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	refs, err := s.backend.GetReferences()
	if err != nil {
		err = components.Toast("Unable to update list of references, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
	}
	w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/requirements/%s", req.ID))
	err = admin.RequirementsPane(admin.RequirementPaneData{
		SelectedRequirement: req,
		Requirements:        reqs,
		OobSwap:             false,
		RequirementEditorData: admin.RequirementEditorData{
			SelectedRequirement: req,
			References:          refs,
			NewRequirement:      false,
		},
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering RequirementPane for created Requirement", slog.String("error", err.Error()))
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	err = components.Toast("Requirement successfully created!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering requirement creation toast", slog.String("error", err.Error()))
	}
}
