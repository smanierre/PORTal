package server

import (
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"log/slog"
	"net/http"
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
		http.Redirect(w, r, "/dashboard", http.StatusFound)
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
		err := s.templateRepo.RenderFragment(w, "admin_members", "content", data)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin fragment", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "error", "generic_ise", data)
		}
		return
	} else {
		err = s.templateRepo.Render(w, "admin_members", &templates.TplData{
			NavData: templates.NavData{
				Show:         true,
				OobSwap:      false,
				Member:       member,
				Subordinates: false,
			},
			ContentData: data,
		})
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin template", slog.String("error", err.Error()))
			_ = s.templateRepo.RenderFragment(w, "errors", "generic_ise", nil)
		}
	}
}
