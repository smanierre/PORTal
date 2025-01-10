package server

import (
	"PORTal/templates"
	"PORTal/templates/components"
	"PORTal/templates/pages/admin"
	"PORTal/templates/pages/errorpages"
	"log/slog"
	"net/http"
)

func (s Server) AdminGetHandler(w http.ResponseWriter, r *http.Request) {
	member, err := getMemberFromContext(r.Context())
	if err != nil {
		s.logger.LogAttrs(r.Context(), slog.LevelWarn, "Unable to get member from context, redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
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
	membersPane := admin.MembersPane(admin.MembersPaneData{
		Members:          members,
		SelectedMember:   member,
		MemberEditorData: admin.MemberEditorData{},
		OobSwap:          false,
		DisabledMembers:  false,
	})
	if checkHTMXRequest(r) {
		w.Header().Set("HX-Push-URL", "/admin/members")
		err := admin.AdminPage("Members", membersPane).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin fragment", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
		return
	} else {
		subordinates, err := s.backend.GetSubordinates(member.ID)
		var subordinateLength int
		if err == nil {
			subordinateLength = len(subordinates)
		}
		err = templates.Root(templates.NavData{
			Show:         true,
			OobSwap:      false,
			Member:       member,
			Subordinates: subordinateLength > 0,
		},
			admin.AdminPage("Members", membersPane)).Render(r.Context(), w)
		if err != nil {
			s.logger.LogAttrs(r.Context(), slog.LevelError, "Error rendering admin template", slog.String("error", err.Error()))
			_ = errorpages.GenericISE().Render(r.Context(), w)
		}
	}
}
