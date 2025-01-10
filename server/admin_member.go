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
	"strings"
)

func (s Server) AdminMembersPaneGetHandler(w http.ResponseWriter, r *http.Request) {
	mems, err := s.backend.GetAllMembers()
	if err != nil {
		if checkHTMXRequest(r) {
			err = components.Toast("Unable to get members.", true).Render(r.Context(), w)
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
	memberPane := admin.MembersPane(admin.MembersPaneData{
		Members:          mems,
		SelectedMember:   types.Member{},
		MemberEditorData: admin.MemberEditorData{},
		OobSwap:          false,
		DisabledMembers:  false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members")
		err = memberPane.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
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
	}, admin.AdminPage("Members", memberPane)).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members template", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
		return
	}
}

func (s Server) AdminMembersPaneGetDisabledHandler(w http.ResponseWriter, r *http.Request) {
	m, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get member from context", slog.String("error", err.Error()))
		http.Redirect(w, r, "/dashboard", http.StatusUnauthorized)
		return
	}
	disabledMembers, err := s.backend.GetDisabledMembers()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get disabled members")
		err = components.Toast("Unable to disable member.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	memberPane := admin.MembersPane(admin.MembersPaneData{
		Members:          disabledMembers,
		SelectedMember:   types.Member{},
		MemberEditorData: admin.MemberEditorData{},
		OobSwap:          false,
		DisabledMembers:  true,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members/disabled")
		err = memberPane.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin disabled members fragment", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
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
	}, admin.AdminPage("Members", memberPane)).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering disabled members page", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminMemberEditorGetHandler(w http.ResponseWriter, r *http.Request) {
	var member types.Member
	var err error
	disabled := strings.Contains(r.Header.Get("HX-Current-Url"), "disabled")
	if disabled {
		member, err = s.backend.GetDisabledMember(r.PathValue("id"))
	} else {
		member, err = s.backend.GetMember(r.PathValue("id"))
	}
	if err != nil {
		err = components.Toast("Unable to get member.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		err = components.Toast("Unable to get members.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	memberEditor := admin.MemberEditor(admin.MemberEditorData{
		SelectedMember: member,
		Ranks:          types.GetRanks(s.config.Service),
		SupervisorListData: admin.SupervisorListData{
			PotentialSupervisors: member.GetPotentialSupervisors(members),
			SelectedMember:       member,
		},
		NewMember: false,
	})
	if checkHTMXRequest(r) {
		if disabled {
			w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/members/disabled/%s", member.ID))
		} else {
			w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/members/%s", r.PathValue("id")))
		}
		err = memberEditor.Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
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
	}, admin.AdminPage("Members", admin.MembersPane(admin.MembersPaneData{
		Members:        members,
		SelectedMember: member,
		MemberEditorData: admin.MemberEditorData{
			SelectedMember: member,
			Ranks:          types.GetRanks(s.config.Service),
			SupervisorListData: admin.SupervisorListData{
				PotentialSupervisors: member.GetPotentialSupervisors(members),
				SelectedMember:       member,
			},
			NewMember: false,
		},
		OobSwap:         false,
		DisabledMembers: false,
	}))).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering member editor page", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminMemberUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted member update without HTMX, redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	memberID := r.PathValue("id")
	newMember := types.Member{
		ApiMember: types.ApiMember{
			ID:           memberID,
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Grade:        types.Grade(r.FormValue("grade")),
			SupervisorID: r.FormValue("supervisor"),
		},
		Password: r.FormValue("password"),
		//TODO: Handle thisDisabled: ,
	}
	var forceNoAdmin bool
	if r.FormValue("admin") == "on" {
		newMember.Admin = true
	} else {
		forceNoAdmin = true
	}
	var forceNoSupervisor bool
	if newMember.SupervisorID == "" {
		forceNoSupervisor = true
	}
	updatedMember, err := s.backend.UpdateMember(newMember, forceNoSupervisor, forceNoAdmin)
	if err != nil {
		err = components.Toast("Unable to update member", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		err = components.Toast("Unable to update member list, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	err = admin.MembersPane(admin.MembersPaneData{
		Members:        members,
		SelectedMember: updatedMember,
		MemberEditorData: admin.MemberEditorData{
			SelectedMember: updatedMember,
			Ranks:          types.GetRanks(s.config.Service),
			SupervisorListData: admin.SupervisorListData{
				PotentialSupervisors: updatedMember.GetPotentialSupervisors(members),
				SelectedMember:       updatedMember,
			},
			NewMember: false,
		},
		OobSwap:         false,
		DisabledMembers: false,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}

	subordinates, err := s.backend.GetSubordinates(updatedMember.ID)
	var subordinateLength int
	if err == nil {
		subordinateLength = len(subordinates)
	}
	// Update nav incase members were added/removed or rank/name changed
	err = templates.Nav(templates.NavData{
		Show:         true,
		OobSwap:      true,
		Member:       updatedMember,
		Subordinates: subordinateLength > 0,
	}).Render(r.Context(), w)
}

func (s Server) AdminMemberDisableHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted member update without HTMX")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	memberID := r.PathValue("id")
	err := s.backend.DisableMember(memberID)
	if err != nil {
		err = components.Toast("Unable to disable member!", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = components.Toast("Unable to disable member.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	// TODO: See if I can just render the content data
	w.Header().Set("HX-Push-URL", "/admin/members")
	err = admin.MemberEditor(admin.MemberEditorData{}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
	err = admin.MembersPane(admin.MembersPaneData{
		Members:          members,
		SelectedMember:   types.Member{},
		MemberEditorData: admin.MemberEditorData{},
		OobSwap:          true,
		DisabledMembers:  false,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}

	err = components.Toast("Member disabled!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering member disable toast", slog.String("error", err.Error()))
	}
}

func (s Server) AdminMemberEnableHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted member update without HTMX")
		http.Redirect(w, r, "/dashboard", http.StatusUnauthorized)
		return
	}
	memberID := r.PathValue("id")
	err := s.backend.EnableMember(memberID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = components.Toast("Error enabling member", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	members, err := s.backend.GetDisabledMembers()
	if err != nil {
		err = components.Toast("Unable to update list of members, please refresh.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	w.Header().Set("HX-Push-URL", "/admin/members/disabled")
	// TODO: See if I can just render the member pane instead of both
	err = admin.MemberEditor(admin.MemberEditorData{}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
	err = admin.MembersPane(admin.MembersPaneData{
		Members:         members,
		OobSwap:         true,
		DisabledMembers: true,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		_ = components.Toast("Error updating page after enabling member.", true).Render(r.Context(), w)
		return
	}
	err = components.Toast("Enabled member!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering success toast", slog.String("error", err.Error()))
		return
	}
}

func (s Server) AdminNewMemberHandler(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty member editor")
	var err error
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members/add")
		err = admin.MemberEditor(admin.MemberEditorData{
			NewMember: true,
			Ranks:     types.GetRanks(s.config.Service),
		}).Render(r.Context(), w)
	} else {
		m, err := getMemberFromContext(r.Context())
		if err != nil {
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
		members, err := s.backend.GetAllMembers()
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
		}, admin.AdminPage("Members", admin.MembersPane(admin.MembersPaneData{
			Members:        members,
			SelectedMember: types.Member{},
			MemberEditorData: admin.MemberEditorData{
				Ranks:     types.GetRanks(s.config.Service),
				NewMember: true,
			},
			OobSwap:         false,
			DisabledMembers: false,
		}))).Render(r.Context(), w)
	}
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering new member fragment", slog.String("error", err.Error()))
		_ = errorpages.GenericISE().Render(r.Context(), w)
	}
}

func (s Server) AdminMemberAddHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to add member, redirecting to the dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new member")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new member", slog.String("error", err.Error()))
		err = components.Toast("Unable to add new member.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	}
	m := types.Member{}
	m.Username = r.Form.Get("username")
	m.FirstName = r.Form.Get("first_name")
	m.LastName = r.Form.Get("last_name")
	m.Password = r.Form.Get("password")
	m.Grade = types.Grade(r.Form.Get("grade"))
	m.Admin = r.Form.Get("admin") == "true"
	m.SupervisorID = r.Form.Get("supervisor_id")

	m, err = s.backend.AddMember(m)
	if err != nil {
		err = components.Toast("Unable to create member.", true).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
			return
		}
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/members/%s", m.ID))
	err = admin.MembersPane(admin.MembersPaneData{
		Members:        members,
		SelectedMember: m,
		MemberEditorData: admin.MemberEditorData{
			SelectedMember: m,
			Ranks:          types.GetRanks(s.config.Service),
			SupervisorListData: admin.SupervisorListData{
				PotentialSupervisors: m.GetPotentialSupervisors(members),
				SelectedMember:       m,
			},
			NewMember: false,
		},
		OobSwap:         false,
		DisabledMembers: false,
	}).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering member_editor for created member", slog.String("error", err.Error()))
		_ = components.Toast("Unable to update page, please refresh.", true).Render(r.Context(), w)
		return
	}
	err = components.Toast("Member successfully created!", false).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering user creation toast", slog.String("error", err.Error()))
	}
}

func (s Server) AdminGetPotentialSupervisorsHandler(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX requrest to get potential supervisors, redirecting to the dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	grade := r.PathValue("grade")
	memberID := r.PathValue("id")
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, fmt.Sprintf("Getting potential supervisors for grade %s", grade))
	member, err := s.backend.GetMember(memberID)
	member.Grade = types.Grade(grade)
	if err != nil {
		_ = components.Toast("Unable to update list of potential supervisors based on selected rank.", true).Render(r.Context(), w)
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		_ = components.Toast("Unable to update list of potential supervisors based on selected rank.", true).Render(r.Context(), w)
		return
	}
	data := admin.SupervisorListData{
		PotentialSupervisors: member.GetPotentialSupervisors(members),
		SelectedMember:       member,
	}
	err = admin.SupervisorList(data).Render(r.Context(), w)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering list of potential supervisors", slog.String("error", err.Error()))
		_ = components.Toast("Unable to update list of potential supervisors based on selected rank.", true).Render(r.Context(), w)
		return
	}
}
