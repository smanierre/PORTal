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

func (s Server) AdminReferenceGetHandler(w http.ResponseWriter, r *http.Request) {
	refs, err := s.backend.GetReferences()
	if err != nil {
		if checkHTMXRequest(r) {
			err = components.Toast("Unable to get references.", true).Render(r.Context(), w)
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
	referencesPane := admin.ReferencesPane(admin.ReferencePaneData{
		SelectedReference:   types.Reference{},
		References:          refs,
		OobSwap:             false,
		ReferenceEditorData: admin.ReferenceEditorData{},
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/references")
		err = referencesPane.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin references fragment", slog.String("error", err.Error()))
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
	}, admin.AdminPage("References", referencesPane)).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin reference page", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminReferenceEditorGetHandler(w http.ResponseWriter, r *http.Request) {
	reference, err := s.backend.GetReference(r.PathValue("id"))
	if err != nil {
		err = components.Toast("Unable to get reference.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	referenceEditor := admin.ReferenceEditor(admin.ReferenceEditorData{
		SelectedReference: reference,
		NewReference:      false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/references/%s", r.PathValue("id")))
		err = referenceEditor.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin reference editor fragment", slog.String("error", err.Error()))
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
	references, err := s.backend.GetReferences()
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
	}, admin.AdminPage("References", admin.ReferencesPane(admin.ReferencePaneData{
		SelectedReference: reference,
		References:        references,
		OobSwap:           false,
		ReferenceEditorData: admin.ReferenceEditorData{
			SelectedReference: reference,
			NewReference:      false,
		},
	}))).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering reference editor page", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminReferenceUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted reference update without HTMX, redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	referenceID := r.PathValue("id")
	volume, err := strconv.Atoi(r.FormValue("volume"))
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Invalid value passed as form value", slog.String("error", err.Error()))
		err = components.Toast("Invalid value for form value", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	// TODO: Handle ui and checking for no volume
	newReference := types.Reference{
		ID:        referenceID,
		Name:      r.FormValue("name"),
		Volume:    volume,
		Paragraph: r.FormValue("paragraph"),
	}
	updatedReference, err := s.backend.UpdateReference(newReference)
	if err != nil {
		err = components.Toast("Unable to update reference", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	references, err := s.backend.GetReferences()
	if err != nil {
		err = components.Toast("Unable to update reference list, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	err = admin.ReferencesPane(admin.ReferencePaneData{
		SelectedReference: updatedReference,
		References:        references,
		OobSwap:           false,
		ReferenceEditorData: admin.ReferenceEditorData{
			SelectedReference: updatedReference,
			NewReference:      false,
		},
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin references fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
	err = components.Toast("Reference updated successfully!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering reference update success toast", slog.String("error", err.Error()))
	}
}

func (s Server) AdminNewReferenceHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty reference editor")
	var err error
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/references/add")
		err = admin.ReferenceEditor(admin.ReferenceEditorData{
			SelectedReference: types.Reference{},
			NewReference:      true,
		}).Render(r.Context(), w)
	} else {
		m, err := getMemberFromContext(r.Context())
		if err != nil {
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
		references, err := s.backend.GetReferences()
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
		}, admin.AdminPage("References", admin.ReferencesPane(admin.ReferencePaneData{
			SelectedReference: types.Reference{},
			References:        references,
			OobSwap:           false,
			ReferenceEditorData: admin.ReferenceEditorData{
				SelectedReference: types.Reference{},
				NewReference:      true,
			},
		}))).Render(r.Context(), w)
	}
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering new reference fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminReferenceAddHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to add reference, redirecting to the dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new reference")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new reference", slog.String("error", err.Error()))
		err = components.Toast("Unable to add new reference.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	ref := types.Reference{}
	ref.Name = r.Form.Get("name")
	ref.Volume, err = strconv.Atoi(r.Form.Get("volume"))
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Invalid volume supplied by user")
		w.WriteHeader(http.StatusBadRequest)
		_ = components.Toast("Invalid value for volume (Must be a number).", true).Render(r.Context(), w)
		return
	}
	ref.Paragraph = r.Form.Get("paragraph")

	ref, err = s.backend.AddReference(ref)
	if err != nil {
		err = components.Toast("Unable to create Reference.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
	}
	refs, err := s.backend.GetReferences()
	if err != nil {
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/references/%s", ref.ID))
	err = admin.ReferencesPane(admin.ReferencePaneData{
		SelectedReference: ref,
		References:        refs,
		OobSwap:           false,
		ReferenceEditorData: admin.ReferenceEditorData{
			SelectedReference: ref,
			NewReference:      false,
		},
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering ReferencePane for created Reference", slog.String("error", err.Error()))
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	err = components.Toast("Reference successfully created!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering reference creation toast", slog.String("error", err.Error()))
	}
}
