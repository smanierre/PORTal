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
	"time"
)

func (s Server) AdminQualificationsPaneGetHandler(w http.ResponseWriter, r *http.Request) {
	quals, err := s.backend.GetAllQualifications()
	if err != nil {
		if checkHTMXRequest(r) {
			err = components.Toast("Unable to get qualifications.", true).Render(r.Context(), w)
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
	qualificationPane := admin.QualificationPane(admin.QualificationPaneData{
		Qualifications: quals,
		OobSwap:        false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/qualifications")
		err = qualificationPane.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualifications fragment", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	m, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error getting member from context, redirecting to dashboard")
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
	}, admin.AdminPage("Qualifications", qualificationPane)).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualifications template", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
}

func (s Server) AdminQualificationEditorGetHandler(w http.ResponseWriter, r *http.Request) {
	qual, err := s.backend.GetQualification(r.PathValue("id"))
	if err != nil {
		err = components.Toast("Unable to get qualification.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	requirements, err := s.backend.GetAllRequirements()
	if err != nil {
		if checkHTMXRequest(r) {
			err = components.Toast("Unable to get Qualification Editor.", true).Render(r.Context(), w)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
				_ = errorpages.GenericISE().Render(r.Context(), w)
			}
		} else {
			err = errorpages.GenericISE().Render(r.Context(), w)
			if err != nil {
				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering generic error page", slog.String("error", err.Error()))
			}
		}
		return
	}
	initialRequirements, recurringRequirements := getPotentialInitialAndRecurringRequirements(requirements)
	qualificationEditor := admin.QualificationEditor(admin.QualificationEditorData{
		SelectedQualification:          qual,
		PotentialInitialRequirements:   initialRequirements,
		PotentialRecurringRequirements: recurringRequirements,
		NewQualification:               false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/qualifications/%s", r.PathValue("id")))
		err = qualificationEditor.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualification editor fragment", slog.String("error", err.Error()))
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
	quals, err := s.backend.GetAllQualifications()
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
	}, admin.AdminPage("Qualifications", admin.QualificationPane(admin.QualificationPaneData{
		SelectedQualification: qual,
		QualificationEditorData: admin.QualificationEditorData{
			SelectedQualification:          qual,
			PotentialInitialRequirements:   initialRequirements,
			PotentialRecurringRequirements: recurringRequirements,
			NewQualification:               false,
		},
		OobSwap:        false,
		Qualifications: quals,
	}))).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering member editor page", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminQualificationUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted qualification update without HTMX, redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}

	qualificationID := r.PathValue("id")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing qualification update form", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w) // TOOD: Probably something more graceful
		return
	}

	expirationDays, err := strconv.Atoi(r.Form.Get("expiration_interval"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing expiration date", slog.String("error", err.Error()))
		err = components.Toast("Unable to parse expiration date.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}

	initialRequirementIds := r.Form["initial_requirements"]
	recurringRequirementIds := r.Form["recurring_requirements"]
	newInitialRequirements := []types.Requirement{}
	newRecurringRequirements := []types.Requirement{}
	for _, id := range initialRequirementIds {
		newInitialRequirements = append(newInitialRequirements, types.Requirement{ID: id})
	}
	for _, id := range recurringRequirementIds {
		newRecurringRequirements = append(newRecurringRequirements, types.Requirement{ID: id})
	}

	newQualification := types.Qualification{
		ID:                    qualificationID,
		Name:                  r.Form.Get("name"),
		InitialRequirements:   newInitialRequirements,
		RecurringRequirements: newRecurringRequirements,
		Notes:                 r.Form.Get("notes"),
		Expires:               r.Form.Get("expires") == "on",
		ExpirationInterval:    time.Duration(expirationDays*24) * time.Hour,
	}
	updatedQualification, err := s.backend.UpdateQualification(newQualification)
	if err != nil {
		w.Header().Set("HX-Retarget", "#toast")
		err = components.Toast("Unable to update qualification", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	qualifications, err := s.backend.GetAllQualifications()
	if err != nil {
		err = components.Toast("Unable to update qualifications list, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	requirements, err := s.backend.GetAllRequirements()
	if err != nil {
		err = components.Toast("Unable to update page, please refresh", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
	}
	initialRequirements, recurringRequirements := getPotentialInitialAndRecurringRequirements(requirements)
	err = admin.QualificationPane(admin.QualificationPaneData{
		Qualifications:        qualifications,
		SelectedQualification: updatedQualification,
		QualificationEditorData: admin.QualificationEditorData{
			SelectedQualification:          updatedQualification,
			PotentialInitialRequirements:   initialRequirements,
			PotentialRecurringRequirements: recurringRequirements,
			NewQualification:               false,
		},
		OobSwap: false,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualifications fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
	err = components.Toast("Updated qualification!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering toast", slog.String("error", err.Error()))
	}
}

func (s Server) AdminNewQualificationHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty qualification editor")
	requirements, err := s.backend.GetAllRequirements()
	if err != nil {
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
	initialRequirements, recurringRequirements := getPotentialInitialAndRecurringRequirements(requirements)
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/qualifications/add")
		err = admin.QualificationEditor(admin.QualificationEditorData{
			NewQualification:               true,
			PotentialInitialRequirements:   initialRequirements,
			PotentialRecurringRequirements: recurringRequirements,
		}).Render(r.Context(), w)
	} else {
		m, err := getMemberFromContext(r.Context())
		if err != nil {
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
		qualifications, err := s.backend.GetAllQualifications()
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
		}, admin.AdminPage("Qualifications", admin.QualificationPane(admin.QualificationPaneData{
			Qualifications:        qualifications,
			SelectedQualification: types.Qualification{},
			QualificationEditorData: admin.QualificationEditorData{
				PotentialInitialRequirements:   initialRequirements,
				PotentialRecurringRequirements: recurringRequirements,
				NewQualification:               true,
			},
			OobSwap: false,
		}))).Render(r.Context(), w)
	}
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering new qualification fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminQualificationAddHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to add qualification, redirecting to the dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new qualification")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new qualification", slog.String("error", err.Error()))
		err = components.Toast("Unable to add new qualification.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	qual := types.Qualification{}
	qual.Name = r.Form.Get("name")
	qual.Notes = r.Form.Get("notes")
	qual.Expires = r.Form.Get("expires") == "on"
	expirationDays, err := strconv.Atoi(r.Form.Get("expiration_interval"))
	if err != nil {
		err = components.Toast("Unable to parse expiration days field, it must be a number.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	qual.ExpirationInterval = time.Duration(expirationDays) * 24 * time.Hour

	for _, id := range r.Form["initial_requirements"] {
		qual.InitialRequirements = append(qual.InitialRequirements, types.Requirement{ID: id})
	}
	for _, id := range r.Form["recurring_requirements"] {
		qual.RecurringRequirements = append(qual.RecurringRequirements, types.Requirement{ID: id})
	}
	qual, err = s.backend.AddQualification(qual)
	if err != nil {
		err = components.Toast("Unable to create Qualification.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
	}
	requirements, err := s.backend.GetAllRequirements()
	if err != nil {
		err = components.Toast("Unable to update page, please refresh", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			return
		}
	}
	initialRequirements, recurringRequirements := getPotentialInitialAndRecurringRequirements(requirements)
	qualifications, err := s.backend.GetAllQualifications()
	if err != nil {
		err = components.Toast("Unable to update page, please refresh", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			return
		}
	}
	w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/qualifications/%s", qual.ID))
	err = admin.QualificationPane(admin.QualificationPaneData{
		Qualifications:        qualifications,
		SelectedQualification: qual,
		OobSwap:               false,
		QualificationEditorData: admin.QualificationEditorData{
			SelectedQualification:          qual,
			PotentialInitialRequirements:   initialRequirements,
			PotentialRecurringRequirements: recurringRequirements,
			NewQualification:               false,
		},
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering QualificationPane for created Qualification", slog.String("error", err.Error()))
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	err = components.Toast("Qualification successfully created!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering reference creation toast", slog.String("error", err.Error()))
	}
}

func getPotentialInitialAndRecurringRequirements(allReqs []types.Requirement) ([]types.Requirement, []types.Requirement) {
	var initialRequirements, recurringRequirements []types.Requirement
	for _, v := range allReqs {
		if v.DaysValidFor > 0 {
			recurringRequirements = append(recurringRequirements, v)
		} else {
			initialRequirements = append(initialRequirements, v)
		}
	}
	return initialRequirements, recurringRequirements
}
