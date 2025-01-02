package server

import (
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"PORTal/types"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

var adminDropdownItems = []components.DropdownItem{
	{
		DisplayName: "Members",
		Value:       "members",
	},
	{
		DisplayName: "Qualifications",
		Value:       "qualifications",
	},
	{
		DisplayName: "Requirements",
		Value:       "requirements",
	},
	{
		DisplayName: "References",
		Value:       "references",
	},
}

func (s Server) AdminGetHandler(w http.ResponseWriter, r *http.Request) {
	member, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get member from context, redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusUnauthorized)
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to get members.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	data := admin.Data{
		MembersData: admin.MembersData{
			Members:          members,
			SwapTarget:       "#member-content",
			FragmentBasePath: "admin/members",
		},
		DropdownData: components.DropdownData{
			Items:            adminDropdownItems,
			SwapTarget:       "#admin-content",
			FragmentBasePath: "admin",
		}}
	if checkHTMXRequest(r) {

		w.Header().Set("HX-Push-URL", "/admin/members")
		err := s.templateRepo.RenderFragment(w, "admin", "content", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin fragme:hovernt", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", data)
		}
		return
	} else {
		err = s.templateRepo.Render(w, "admin", &templates.TplData{
			NavData: templates.NavData{
				Show:         true,
				OobSwap:      false,
				Member:       member,
				Subordinates: false,
			},
			ContentData: data})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin template", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "errors", "generic_ise", nil)
		}
	}
}

func (s Server) AdminMembersGetHandler(w http.ResponseWriter, r *http.Request) {
	mems, err := s.backend.GetAllMembers()
	if err != nil {
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to get members.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	data := admin.MembersData{
		Members:          mems,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/members",
		MemberEditorData: admin.MemberEditorData{
			Ranks: types.GetRanks(s.config.Service),
		},
	}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members")
		err = s.templateRepo.RenderFragment(w, "admin", "members", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		}
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (s Server) AdminMembersGetDisabledHandler(w http.ResponseWriter, r *http.Request) {
	m, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get member from context", slog.String("error", err.Error()))
		http.Redirect(w, r, "/dashboard", http.StatusUnauthorized)
		return
	}
	disabledMembers, err := s.backend.GetDisabledMembers()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get disabled members")
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to disable member.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	data := admin.Data{
		DropdownData: components.DropdownData{Items: adminDropdownItems},
		MembersData: admin.MembersData{
			Members:          disabledMembers,
			SwapTarget:       "#member-content",
			FragmentBasePath: "admin/members",
			MemberEditorData: admin.MemberEditorData{},
			OobSwap:          false,
			DisabledMembers:  true,
		},
	}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members/disabled")
		err = s.templateRepo.RenderFragment(w, "admin", "members", data.MembersData)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin disabled members fragment", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", data)
		}
		return
	}
	err = s.templateRepo.Render(w, "admin", &templates.TplData{
		NavData: templates.NavData{
			Show:         true,
			OobSwap:      false,
			Member:       m,
			Subordinates: false, // TODO: fix this
		},
		ContentData: data,
	})
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering disabled members page", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
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
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to get member.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to get members.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	data := admin.MemberEditorData{
		SelectedMember: member,
		Ranks:          types.GetRanks(s.config.Service),
		SupervisorListData: admin.SupervisorListData{
			PotentialSupervisors: member.GetPotentialSupervisors(members),
			SelectedMember:       member,
		},
	}
	if checkHTMXRequest(r) {
		if disabled {
			w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/members/disabled/%s", member.ID))
		} else {
			w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/members/%s", r.PathValue("id")))
		}
		err = s.templateRepo.RenderFragment(w, "admin", "member-editor", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		}
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (s Server) UpdateMember(w http.ResponseWriter, r *http.Request) {
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
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update member",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update member list, please refresh.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	membersData := admin.MembersData{
		MemberEditorData: admin.MemberEditorData{
			SelectedMember: updatedMember,
			Ranks:          types.GetRanks(s.config.Service),
			SupervisorListData: admin.SupervisorListData{
				PotentialSupervisors: updatedMember.GetPotentialSupervisors(members),
				SelectedMember:       updatedMember,
			},
			NewMember: false,
		},
		SelectedMember:   updatedMember,
		Members:          members,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/members",
		OobSwap:          true,
	}
	err = s.templateRepo.RenderFragment(w, "admin", "members", membersData)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}
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
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to disable member!",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to disable member.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	data := admin.MemberEditorData{}
	w.Header().Set("HX-Push-URL", "/admin/members")
	err = s.templateRepo.RenderFragment(w, "admin", "member-editor", data)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}
	membersData := admin.MembersData{
		Members:          members,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/members",
		OobSwap:          true,
	}
	err = s.templateRepo.RenderFragment(w, "admin", "members", membersData)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}

	err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{Message: "Member disabled!"})
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
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Error enabling member",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		}
		return
	}
	members, err := s.backend.GetDisabledMembers()
	if err != nil {
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update list of members, please refresh.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
		}
		return
	}
	data := admin.MemberEditorData{}
	w.Header().Set("HX-Push-URL", "/admin/members/disabled")
	err = s.templateRepo.RenderFragment(w, "admin", "member-editor", data)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}
	membersData := admin.MembersData{
		Members:          members,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/members",
		OobSwap:          true,
		DisabledMembers:  true,
	}
	err = s.templateRepo.RenderFragment(w, "admin", "members", membersData)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Error updating page after enabling member.",
			Danger:  true,
		})
		return
	}
	err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
		Message: "Enabled member!",
	})
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering success toast", slog.String("error", err.Error()))
		return
	}
}

func (s Server) AdminGetNewMember(w http.ResponseWriter, r *http.Request) {
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Sending back empty member editor")
	var err error
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members/add")
		err = s.templateRepo.RenderFragment(w, "admin", "member-editor", admin.MemberEditorData{
			NewMember: true,
			Ranks:     types.GetRanks(s.config.Service),
		})
	} else {
		m, err := getMemberFromContext(r.Context())
		if err != nil {
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
			return
		}
		members, err := s.backend.GetAllMembers()
		if err != nil {
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
			return
		}
		err = s.templateRepo.Render(w, "admin", &templates.TplData{
			NavData: templates.NavData{
				Show:         true,
				OobSwap:      false,
				Member:       m,
				Subordinates: false, // TODO: Populate this
			},
			ContentData: admin.Data{
				DropdownData: components.DropdownData{
					Items:            adminDropdownItems,
					SwapTarget:       "#admin-content",
					FragmentBasePath: "admin",
				},
				MembersData: admin.MembersData{
					Members:          members,
					SelectedMember:   types.Member{},
					SwapTarget:       "#member-content",
					FragmentBasePath: "admin/members",
					MemberEditorData: admin.MemberEditorData{NewMember: true, Ranks: types.GetRanks(s.config.Service)},
					OobSwap:          false,
					DisabledMembers:  false,
				},
			},
		})
	}
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering new member fragment", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}
}

func (s Server) AddMember(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to add member, redirecting to the dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}
	s.logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new member")
	err := r.ParseForm()
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new member", slog.String("error", err.Error()))
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to add new member.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
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
		err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to create member.",
			Danger:  true,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering danger toast", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
			return
		}
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		_ = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update page, please refresh.",
			Danger:  true,
		})
		return
	}
	err = s.templateRepo.RenderFragment(w, "admin", "members", admin.MembersData{
		Members:          members,
		SelectedMember:   m,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/members",
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
	})
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering member_editor for created member", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update page, please refresh.",
			Danger:  true,
		})
		return
	}
	err = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{Message: "User Successfully create!"})
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering user creation toast", slog.String("error", err.Error()))
	}
}

func (s Server) AdminGetPotentialSupervisors(w http.ResponseWriter, r *http.Request) {
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
		_ = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update list of potential supervisors based on selected rank.",
			Danger:  true,
		})
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		_ = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update list of potential supervisors based on selected rank.",
			Danger:  true,
		})
		return
	}
	data := admin.SupervisorListData{
		PotentialSupervisors: member.GetPotentialSupervisors(members),
		SelectedMember:       member,
	}
	err = s.templateRepo.RenderFragment(w, "admin", "supervisor_list", data)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering list of potential supervisors", slog.String("error", err.Error()))
		_ = s.templateRepo.RenderFragment(w, "admin", "toast", components.ToastData{
			Message: "Unable to update list of potential supervisors based on selected rank.",
			Danger:  true,
		})
		return
	}
}
