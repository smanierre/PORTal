package server

import (
	"PORTal/backend"
	"PORTal/server/serverutils"
	"PORTal/server/stores"
	"PORTal/templates"
	"PORTal/templates/pages/admin"
	"PORTal/templates/pages/errorpages"
	"PORTal/types"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func adminRootGetHandler(logger *slog.Logger, memberStore stores.MemberStore, qualificationStore stores.QualificationStore, ranks types.RankMap) http.Handler {
	logger = logger.With("route", "GET /admin")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") == "true" && r.URL.Query().Get("fragment") != "" {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "HTMX Request, directing to fragment handler")
			adminPanelFragmentHandler(logger, memberStore, qualificationStore, ranks, w, r)
			return
		}
		members, err := memberStore.GetAllMembers()
		if err != nil {
			serverutils.HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
		}
		member, err := serverutils.GetMemberFromContext(r.Context())
		if err != nil {
			serverutils.HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
		}
		var hasSubordinates bool
		if subordinates, err := memberStore.GetSubordinates(member.ID); err != nil || len(subordinates) == 0 {
			hasSubordinates = false
		} else {
			hasSubordinates = true
		}
		var selectedMember types.Member
		var potentialSupervisors []types.Member
		if r.URL.Query().Get("selected") != "" {
			selectedMember, err = memberStore.GetMember(r.URL.Query().Get("selected"))
			// If the member is not found just continue. The editor will be blank.
			// If it's a different error, return an error page
			if err != nil && !errors.Is(err, backend.ErrMemberNotFound) {
				w.WriteHeader(http.StatusInternalServerError)
				serverutils.HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
				return
			}
			potentialSupervisors, err = memberStore.GetPotentialSupervisors(selectedMember)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				serverutils.HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
				return
			}
		}
		serverutils.HandleRenderError(r.Context(), logger, templates.Root(templates.NavData{
			Show:         true,
			OobSwap:      false,
			Member:       member,
			Subordinates: hasSubordinates,
		}, admin.AdminPage("Members", admin.MembersPane(
			members,
			selectedMember,
			ranks,
			potentialSupervisors,
			false,
			false,
			false,
			""))).Render(r.Context(), w))
	})
}

func adminPanelFragmentHandler(
	logger *slog.Logger,
	memberStore stores.MemberStore,
	qualificationStore stores.QualificationStore,
	ranks types.RankMap,
	w http.ResponseWriter,
	r *http.Request,
) {
	fragment := r.URL.Query().Get("fragment")
	switch fragment {
	case "members_pane":
		membersPaneFragmentHandler(logger, memberStore, ranks, w, r)
	case "qualifications_pane":
		qualificationsPaneFragmentHandler(logger, qualificationStore, w, r)
	default:
		logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid fragment in query", slog.String("fragment", fragment))
	}
}

func membersPaneFragmentHandler(
	logger *slog.Logger,
	memberStore stores.MemberStore,
	ranks types.RankMap,
	w http.ResponseWriter,
	r *http.Request,
) {
	logger = logger.With("handler", "membersPaneFragmentHandler")
	memberFilterQuery := r.URL.Query().Get("member_query")
	disabledMembers := r.URL.Query().Get("disabled")
	var members []types.Member
	var err error
	if disabledMembers == "true" {
		members, err = memberStore.GetDisabledMembers()
	} else {
		members, err = memberStore.GetAllMembers()
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if memberFilterQuery != "" {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Filtering members", slog.String("query", memberFilterQuery))
		members = types.FilterMembers(members, func(m types.Member) bool {
			return strings.Contains(
				strings.ToLower(fmt.Sprintf("%s %s", m.FirstName, m.LastName)),
				strings.ToLower(memberFilterQuery),
			)
		})
	}
	var member types.Member
	var potentialSupervisors []types.Member
	selectedMemberQuery := r.URL.Query().Get("selected")
	if selectedMemberQuery != "" {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Getting selected member", slog.String("query", selectedMemberQuery))
		member, err = memberStore.GetMember(selectedMemberQuery)
		if errors.Is(err, backend.ErrMemberNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		potentialSupervisors, err = memberStore.GetPotentialSupervisors(member)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	membersPane := admin.MembersPane(
		members,
		member,
		ranks,
		potentialSupervisors,
		false,
		disabledMembers == "true",
		r.URL.Query().Get("new") == "true",
		memberFilterQuery,
	)
	serverutils.HandleRenderError(r.Context(), logger, membersPane.Render(r.Context(), w))
	dropdown := admin.AdminDropdown("Members", true)
	serverutils.HandleRenderError(r.Context(), logger, dropdown.Render(r.Context(), w))
}

func qualificationsPaneFragmentHandler(logger *slog.Logger, qualificationStore stores.QualificationStore, w http.ResponseWriter, r *http.Request) {
	logger = logger.With("handler", "qualificationsPaneFragmentHandler")
	qualifications, err := qualificationStore.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	selectedQualificationQuery := r.URL.Query().Get("selected")
	var selectedQualification types.Qualification
	if selectedQualificationQuery != "" {
		selectedQualification, err = qualificationStore.GetQualification(selectedQualificationQuery)
		if errors.Is(err, backend.ErrQualificationNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
	qualificationsPane := admin.QualificationPane(qualifications, selectedQualification, r.URL.Query().Get("new") == "true", false)
	serverutils.HandleRenderError(r.Context(), logger, qualificationsPane.Render(r.Context(), w))

	dropdown := admin.AdminDropdown("Qualifications", true)
	serverutils.HandleRenderError(r.Context(), logger, dropdown.Render(r.Context(), w))
}

func getPotentialSupervisorsHandler(logger *slog.Logger, memberStore stores.MemberStore) http.Handler {
	logger = logger.With("route", "GET /admin/potentialSupervisors")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Getting potential supervisors for grade", slog.String("grade", r.URL.Query().Get("grade")))
		grade := types.Grade(r.URL.Query().Get("grade"))
		if grade == "" {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "No grade provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id := r.URL.Query().Get("id")
		if id == "" {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "No id provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m, err := memberStore.GetMember(id)
		if errors.Is(err, backend.ErrMemberNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		potentialSupervisors, err := memberStore.GetPotentialSupervisors(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		supervisorList := admin.SupervisorList(potentialSupervisors, m)
		serverutils.HandleRenderError(r.Context(), logger, supervisorList.Render(r.Context(), w))
	})
}

func updateMemberHandler(logger *slog.Logger, memberStore stores.MemberStore) http.Handler {
	logger = logger.With("route", "PUT /admin/updateMember")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Updating member", slog.String("member", r.URL.Query().Get("member")))
		panic("implement me")
	})
}
