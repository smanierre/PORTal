package server

//
//import (
//	"PORTal/server/serverutils"
//	"PORTal/templates"
//	"PORTal/templates/components"
//	"PORTal/templates/components/webcomponents"
//	"PORTal/templates/pages/admin"
//	"PORTal/templates/pages/errorpages"
//	"PORTal/types"
//	"encoding/json"
//	"fmt"
//	"log/slog"
//	"net/http"
//	"strconv"
//	"strings"
//	"time"
//)
//
//func (s Server) AdminQualificationsPaneGetHandler(w http.ResponseWriter, r *http.Request) {
//	quals, err := s.backend.GetAllQualifications()
//	if err != nil {
//		if serverutils.CheckHTMXRequest(r) {
//			err = components.Toast("Unable to get qualifications.", true).Render(r.Context(), w)
//			if err != nil {
//				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//			}
//			return
//		} else {
//			err = errorpages.GenericISE().Render(r.Context(), w)
//			if err != nil {
//				s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering error page", slog.String("error", err.Error()))
//			}
//			return
//		}
//	}
//	qualificationPane := admin.QualificationPane(admin.QualificationPaneData{
//		Qualifications: quals,
//		OobSwap:        false,
//	})
//	if serverutils.CheckHTMXRequest(r) {
//		w.Header().Set("HX-Push-URL", "/admin/qualifications")
//		err = qualificationPane.Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualifications fragment", slog.String("error", err.Error()))
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//		}
//		return
//	}
//	m, err := serverutils.GetMemberFromContext(r.Context())
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error getting member from context, redirecting to dashboard")
//		http.Redirect(w, r, "/dashboard", http.StatusFound)
//		return
//	}
//	subordinates, err := s.backend.GetSubordinates(m.ID)
//	var subordinateLength int
//	if err == nil {
//		subordinateLength = len(subordinates)
//	}
//	err = templates.Root(templates.NavData{
//		Show:         true,
//		OobSwap:      false,
//		Member:       m,
//		Subordinates: subordinateLength > 0,
//	}, admin.AdminPage("Qualifications", qualificationPane)).Render(r.Context(), w)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualifications template", slog.String("error", err.Error()))
//		_ = errorpages.GenericISE().Render(r.Context(), w)
//		return
//	}
//}
//
//func (s Server) AdminQualificationEditorGetHandler(w http.ResponseWriter, r *http.Request) {
//	qual, err := s.backend.GetQualification(r.PathValue("id"))
//	if err != nil {
//		err = components.Toast("Unable to get qualification.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//		}
//		return
//	}
//	qualificationEditor := webcomponents.QualificationEditor(webcomponents.QualificationEditorData{
//		SelectedQualification: qual,
//		NewQualification:      false,
//	})
//	if serverutils.CheckHTMXRequest(r) {
//		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/qualifications/%s", r.PathValue("id")))
//		err = qualificationEditor.Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualification editor fragment", slog.String("error", err.Error()))
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//		}
//		return
//	}
//	m, err := serverutils.GetMemberFromContext(r.Context())
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "error getting member from context", slog.String("error", err.Error()))
//		_ = errorpages.GenericISE().Render(r.Context(), w)
//		return
//	}
//	quals, err := s.backend.GetAllQualifications()
//	if err != nil {
//		_ = errorpages.GenericISE().Render(r.Context(), w)
//		return
//	}
//	subordinates, err := s.backend.GetSubordinates(m.ID)
//	var subordinateLength int
//	if err == nil {
//		subordinateLength = len(subordinates)
//	}
//	err = templates.Root(templates.NavData{
//		Show:         true,
//		OobSwap:      false,
//		Member:       m,
//		Subordinates: subordinateLength > 0,
//	}, admin.AdminPage("Qualifications", admin.QualificationPane(admin.QualificationPaneData{
//		SelectedQualification: qual,
//		QualificationEditorData: webcomponents.QualificationEditorData{
//			SelectedQualification: qual,
//			NewQualification:      false,
//		},
//		OobSwap:        false,
//		Qualifications: quals,
//	}))).Render(r.Context(), w)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering member editor page", slog.String("error", err.Error()))
//		_ = errorpages.GenericISE().Render(r.Context(), w)
//	}
//}
//
//func (s Server) AdminQualificationUpdateHandler(w http.ResponseWriter, r *http.Request) {
//	if !serverutils.CheckHTMXRequest(r) {
//		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted qualification update without HTMX, redirecting to dashboard")
//		http.Redirect(w, r, "/dashboard", http.StatusFound)
//		return
//	}
//
//	qualificationID := r.PathValue("id")
//	err := r.ParseForm()
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing qualification update form", slog.String("error", err.Error()))
//		_ = errorpages.GenericISE().Render(r.Context(), w) // TOOD: Probably something more graceful
//		return
//	}
//
//	expirationDays, err := strconv.Atoi(r.Form.Get("expiration_interval"))
//	if err != nil {
//		w.WriteHeader(http.StatusBadRequest)
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error parsing expiration date", slog.String("error", err.Error()))
//		err = components.Toast("Unable to parse expiration date.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//		}
//		return
//	}
//
//	initialRequirementJSON := r.Form.Get("initial_requirements")
//	recurringRequirementJSON := r.Form.Get("recurring_requirements")
//	newInitialRequirements := []types.Requirement{}
//	newRecurringRequirements := []types.Requirement{}
//
//	err = json.NewDecoder(strings.NewReader(initialRequirementJSON)).Decode(&newInitialRequirements)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error decoding initial requirements from JSON", slog.String("error", err.Error()))
//		err = components.Toast("Unable to parse initial requirements.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast for invalid initial requirements", slog.String("error", err.Error()))
//		}
//		return
//	}
//	err = json.NewDecoder(strings.NewReader(recurringRequirementJSON)).Decode(&newRecurringRequirements)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error decoding recurring requirements from JSON", slog.String("error", err.Error()))
//		err = components.Toast("Unable to parse recurring requirements.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast for invalid recurring requirements", slog.String("error", err.Error()))
//		}
//		return
//	}
//
//	newQualification := types.Qualification{
//		ID:                    qualificationID,
//		Name:                  r.Form.Get("name"),
//		InitialRequirements:   newInitialRequirements,
//		RecurringRequirements: newRecurringRequirements,
//		Notes:                 r.Form.Get("notes"),
//		Expires:               r.Form.Get("expires") == "on",
//		ExpirationInterval:    time.Duration(expirationDays*24) * time.Hour,
//	}
//	updatedQualification, err := s.backend.UpdateQualification(newQualification)
//	if err != nil {
//		w.Header().Set("HX-Retarget", "#toast")
//		err = components.Toast("Unable to update qualification", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//		}
//		return
//	}
//	qualifications, err := s.backend.GetAllQualifications()
//	if err != nil {
//		err = components.Toast("Unable to update qualifications list, please refresh.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//		}
//		return
//	}
//	err = admin.QualificationPane(admin.QualificationPaneData{
//		Qualifications:        qualifications,
//		SelectedQualification: updatedQualification,
//		QualificationEditorData: webcomponents.QualificationEditorData{
//			SelectedQualification: updatedQualification,
//			NewQualification:      false,
//		},
//		OobSwap: false,
//	}).Render(r.Context(), w)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin qualifications fragment", slog.String("error", err.Error()))
//		_ = errorpages.GenericISE().Render(r.Context(), w)
//	}
//	err = components.Toast("Updated qualification!", false).Render(r.Context(), w)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering toast", slog.String("error", err.Error()))
//	}
//}
//
//func (s Server) AdminNewQualificationHandler(w http.ResponseWriter, r *http.Request) {
//	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty qualification editor")
//	if serverutils.CheckHTMXRequest(r) {
//		w.Header().Set("HX-Push-URL", "/admin/qualifications/add")
//		err := webcomponents.QualificationEditor(webcomponents.QualificationEditorData{
//			NewQualification: true,
//		}).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering Qualification editor", slog.String("error", err.Error()))
//		}
//		return
//	} else {
//		m, err := serverutils.GetMemberFromContext(r.Context())
//		if err != nil {
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//			return
//		}
//		qualifications, err := s.backend.GetAllQualifications()
//		if err != nil {
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//			return
//		}
//		subordinates, err := s.backend.GetSubordinates(m.ID)
//		var subordinateLength int
//		if err == nil {
//			subordinateLength = len(subordinates)
//		}
//		err = templates.Root(templates.NavData{
//			Show:         true,
//			OobSwap:      false,
//			Member:       m,
//			Subordinates: subordinateLength > 0,
//		}, admin.AdminPage("Qualifications", admin.QualificationPane(admin.QualificationPaneData{
//			Qualifications:        qualifications,
//			SelectedQualification: types.Qualification{},
//			QualificationEditorData: webcomponents.QualificationEditorData{
//				NewQualification: true,
//			},
//			OobSwap: false,
//		}))).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering qualification page with root template", slog.String("error", err.Error()))
//		}
//	}
//}
//
//func (s Server) AdminQualificationAddHandler(w http.ResponseWriter, r *http.Request) {
//	if !serverutils.CheckHTMXRequest(r) {
//		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to add qualification, redirecting to the dashboard")
//		http.Redirect(w, r, "/dashboard", http.StatusFound)
//		return
//	}
//	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new qualification")
//	err := r.ParseForm()
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new qualification", slog.String("error", err.Error()))
//		err = components.Toast("Unable to add new qualification.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//		}
//		return
//	}
//	qual := types.Qualification{}
//	qual.Name = r.Form.Get("name")
//	qual.Notes = r.Form.Get("notes")
//	qual.Expires = r.Form.Get("expires") == "on"
//	expirationDays, err := strconv.Atoi(r.Form.Get("expiration_interval"))
//	if err != nil {
//		err = components.Toast("Unable to parse expiration days field, it must be a number.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//		}
//		return
//	}
//	qual.ExpirationInterval = time.Duration(expirationDays) * 24 * time.Hour
//
//	initialRequirementJSON := r.Form.Get("initial_requirements")
//	recurringRequirementJSON := r.Form.Get("recurring_requirements")
//	newInitialRequirements := []types.Requirement{}
//	newRecurringRequirements := []types.Requirement{}
//
//	err = json.NewDecoder(strings.NewReader(initialRequirementJSON)).Decode(&newInitialRequirements)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error decoding initial requirements from JSON", slog.String("error", err.Error()))
//		err = components.Toast("Unable to parse initial requirements.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast for invalid initial requirements", slog.String("error", err.Error()))
//		}
//		return
//	}
//	err = json.NewDecoder(strings.NewReader(recurringRequirementJSON)).Decode(&newRecurringRequirements)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Error decoding recurring requirements from JSON", slog.String("error", err.Error()))
//		err = components.Toast("Unable to parse recurring requirements.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast for invalid recurring requirements", slog.String("error", err.Error()))
//		}
//		return
//	}
//	qual.InitialRequirements = newInitialRequirements
//	qual.RecurringRequirements = newRecurringRequirements
//	qual, err = s.backend.AddQualification(qual)
//	if err != nil {
//		err = components.Toast("Unable to create Qualification.", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//			_ = errorpages.GenericISE().Render(r.Context(), w)
//		}
//		return
//	}
//
//	qualifications, err := s.backend.GetAllQualifications()
//	if err != nil {
//		err = components.Toast("Unable to update page, please refresh", true).Render(r.Context(), w)
//		if err != nil {
//			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
//		}
//		return
//	}
//	w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/qualifications/%s", qual.ID))
//	err = admin.QualificationPane(admin.QualificationPaneData{
//		Qualifications:        qualifications,
//		SelectedQualification: qual,
//		OobSwap:               false,
//		QualificationEditorData: webcomponents.QualificationEditorData{
//			SelectedQualification: qual,
//			NewQualification:      false,
//		},
//	}).Render(r.Context(), w)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering QualificationPane for created Qualification", slog.String("error", err.Error()))
//		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
//		return
//	}
//	err = components.Toast("Qualification successfully created!", false).Render(r.Context(), w)
//	if err != nil {
//		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering reference creation toast", slog.String("error", err.Error()))
//	}
//}
