package server

import (
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"PORTal/types"
	"fmt"
	"log/slog"
	"net/http"
)

func (s Server) AdminGetHandler(w http.ResponseWriter, r *http.Request) {
	member, err := getMemberFromContext(r.Context())
	if err != nil {
		// TODO: Return error page?
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get member from context")
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		// TODO: RETURN ERROR PAGE
		return
	}
	data := admin.Data{
		MembersData: admin.MembersData{
			Members:          members,
			SwapTarget:       "#member-content",
			FragmentBasePath: "admin/member",
			SelectedMember:   member,
			MemberEditorData: admin.MemberEditorData{
				Members:              members,
				SelectedMember:       member,
				Ranks:                types.GetRanks(s.config.Service),
				PotentialSupervisors: member.GetPotentialSupervisors(members),
			},
		},
		DropdownData: components.DropdownData{
			Items: []components.DropdownItem{
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
			},
			SwapTarget:       "#admin-content",
			FragmentBasePath: "admin",
		}}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/member/%s", member.ID))
		err := s.templateRepo.RenderFragment(w, "admin", "content", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin fragment", slog.String("error", err.Error()))
			s.templateRepo.RenderFragment(w, "error", "generic_ise", data)
		}
		return
	} else {
		s.templateRepo.Render(w, "admin", &templates.TplData{
			NavData: templates.NavData{
				Show:         true,
				OobSwap:      false,
				Member:       member,
				Subordinates: false,
			},
			ContentData: data})
	}
}

func (s Server) AdminMembersGetHandler(w http.ResponseWriter, r *http.Request) {
	member, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get member from context")
		return
	}
	mems, err := s.backend.GetAllMembers()
	if err != nil {
		//TODO return error page/message
		return
	}
	data := admin.MembersData{
		Members:          mems,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/member",
		SelectedMember:   member,
		MemberEditorData: admin.MemberEditorData{
			Members:              mems,
			SelectedMember:       member,
			Ranks:                types.GetRanks(s.config.Service),
			PotentialSupervisors: member.GetPotentialSupervisors(mems),
		},
	}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/member/%s", member.ID))
		err = s.templateRepo.RenderFragment(w, "admin", "members", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
			s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		}
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (s Server) AdminMemberEditorGetHandler(w http.ResponseWriter, r *http.Request) {
	member, err := s.backend.GetMember(r.PathValue("id"))
	if err != nil {
		//TODO: Return error page/message
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		//TODO: Return error page/message
		return
	}
	data := admin.MemberEditorData{
		Members:              members,
		SelectedMember:       member,
		Ranks:                types.GetRanks(s.config.Service),
		PotentialSupervisors: member.GetPotentialSupervisors(members),
	}
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin/member/%s", r.PathValue("id")))
		err = s.templateRepo.RenderFragment(w, "admin", "member-editor", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
			s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
		}
		return
	}
	http.Redirect(w, r, "/admin", http.StatusFound)
}

func (s Server) UpdateMember(w http.ResponseWriter, r *http.Request) {
	if !checkHTMXRequest(r) {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Attempted member update without HTMX")
		// TODO: Unauthorized?
		return
	}
	memberID := r.PathValue("id")
	newMember := types.Member{
		ApiMember: types.ApiMember{
			ID:           memberID,
			FirstName:    r.FormValue("first_name"),
			LastName:     r.FormValue("last_name"),
			Grade:        types.Grade(r.FormValue("grade")),
			SupervisorID: r.FormValue("supervisor_id"),
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
		//TODO: Return an error page or message
		return
	}
	members, err := s.backend.GetAllMembers()
	if err != nil {
		//TODO: Return an error message
		return
	}
	data := admin.MemberEditorData{
		Members:              members,
		SelectedMember:       updatedMember,
		Ranks:                types.GetRanks(s.config.Service),
		PotentialSupervisors: updatedMember.GetPotentialSupervisors(members),
	}
	err = s.templateRepo.RenderFragment(w, "admin", "member-editor", data)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin member editor fragment", slog.String("error", err.Error()))
		s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}
	membersData := admin.MembersData{
		Members:          members,
		SwapTarget:       "#member-content",
		FragmentBasePath: "admin/member",
		OobSwap:          true,
	}

	err = s.templateRepo.RenderFragment(w, "admin", "members", membersData)
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin members fragment", slog.String("error", err.Error()))
		s.templateRepo.RenderFragment(w, "error", "generic_ise", nil)
	}
}
