package server

import (
	"PORTal/backend"
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

func adminRootGetHandler(logger *slog.Logger, memberStore stores.MemberStore, qualificationStore stores.QualificationStore, ranks types.RankMap, organization string) http.Handler {
	logger = logger.With("route", "GET /admin")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") == "true" {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "HTMX Request, directing to fragment handler")
			adminPanelFragmentHandler(logger, memberStore, qualificationStore, ranks, w, r)
			return
		}

		// Determine which fragment is being requested and render based on that
		switch r.URL.Query().Get("fragment") {
		case "qualifications_pane":
			rootQualificationsPaneHandler(logger, memberStore, qualificationStore, organization, w, r)
			return
		case "requirement_editor":
			rootRequirementEditorHandler(logger, memberStore, qualificationStore, organization, w, r)
			return
		default:
			rootMembersPaneHandler(logger, memberStore, organization, ranks, w, r)
		}
	})
}

func rootMembersPaneHandler(
	logger *slog.Logger,
	memberStore stores.MemberStore,
	organization string,
	ranks types.RankMap,
	w http.ResponseWriter, r *http.Request,
) {
	logger.LogAttrs(r.Context(), slog.LevelInfo, "Serving root template with members pane content")

	// Potential Queries
	disabled := r.URL.Query().Get("disabled") == "true"
	query := r.URL.Query().Get("member_query")
	selectedMemberID := r.URL.Query().Get("selected")

	var members []types.Member
	var err error

	if disabled {
		members, err = memberStore.GetDisabledMembers()
		if err != nil {
			HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
		}
	} else {
		members, err = memberStore.GetAllMembers()
		if err != nil {
			HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
		}
	}
	if query != "" {
		members = types.FilterMembers(members, func(m types.Member) bool {
			return strings.Contains(
				strings.ToLower(fmt.Sprintf("%s %s", m.FirstName, m.LastName)),
				strings.ToLower(query),
			)
		})
	}
	loggedInMember, err := MemberFromContext(r.Context())
	if err != nil {
		HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
		return
	}
	var hasSubordinates bool
	if subordinates := memberStore.GetSubordinates(loggedInMember.ID); len(subordinates) == 0 {
		hasSubordinates = false
	} else {
		hasSubordinates = true
	}
	var selectedMember types.Member
	var potentialSupervisors []types.Member
	if selectedMemberID != "" {
		if disabled {
			selectedMember, err = memberStore.GetDisabledMember(selectedMemberID)
		} else {
			selectedMember, err = memberStore.GetMember(selectedMemberID)
		}
		// If the member is not found just continue. The editor will be blank.
		// If it's a different error, return an error page
		if err != nil && !errors.Is(err, backend.ErrMemberNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
			return
		}
		// Disabled members don't have supervisors
		if !disabled {
			potentialSupervisors, err = memberStore.GetPotentialSupervisors(selectedMember, selectedMember.Grade)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
				return

			}
		}
	}
	memberDisplayName := fmt.Sprintf("%s %s %s", ranks[loggedInMember.Grade], loggedInMember.FirstName, loggedInMember.LastName)
	HandleRenderError(r.Context(), logger, templates.Root(true, hasSubordinates, loggedInMember.Admin, memberDisplayName, organization,
		admin.AdminPage("Members", admin.MembersPane(
			members,
			selectedMember,
			ranks,
			potentialSupervisors,
			false,
			disabled,
			false,
			query))).Render(r.Context(), w))

}

func rootQualificationsPaneHandler(
	logger *slog.Logger,
	memberStore stores.MemberStore,
	qualificationStore stores.QualificationStore,
	organization string,
	w http.ResponseWriter, r *http.Request,
) {
	logger.LogAttrs(r.Context(), slog.LevelInfo, "Rendering root template with qualifications pane")

	// Potential queries
	selected := r.URL.Query().Get("selected")
	query := r.URL.Query().Get("qualification_query")
	newQualification := r.URL.Query().Get("new") == "true"

	qualifications, err := qualificationStore.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if query != "" {
		qualifications = types.FilterQualifications(qualifications, func(q types.Qualification) bool {
			return strings.Contains(strings.ToLower(q.Name), strings.ToLower(query))
		})
	}
	var selectedQualification types.Qualification
	if selected != "" {
		selectedQualification, err = qualificationStore.GetQualification(selected)
		if errors.Is(err, backend.ErrQualificationNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Printf("%+v\n", selectedQualification)
		for i, requirement := range selectedQualification.InitialRequirements {
			if requirement.Type == types.QualificationType {
				qual, err := qualificationStore.GetQualification(requirement.QualificationID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				selectedQualification.InitialRequirements[i].Name = qual.Name
			}
		}
	}

	loggedInMember, err := MemberFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var hasSubordinates bool
	if subordinates := memberStore.GetSubordinates(loggedInMember.ID); len(subordinates) > 0 {
		hasSubordinates = true
	}

	adminContent := admin.AdminPage("Qualifications",
		admin.QualificationPane(qualifications, selectedQualification, newQualification, false, query))
	rootTemplate := templates.Root(true, hasSubordinates, true, loggedInMember.Display(), organization, adminContent)
	HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
}

func rootRequirementEditorHandler(
	logger *slog.Logger,
	memberStore stores.MemberStore,
	qualificationStore stores.QualificationStore,
	organization string,
	w http.ResponseWriter,
	r *http.Request,
) {
	logger.LogAttrs(r.Context(), slog.LevelInfo, "Rendering root template with requirement editor")

	// Required Queries
	requirementID := r.URL.Query().Get("requirement_id")
	requirementType := r.URL.Query().Get("requirement_type")
	qualificationID := r.URL.Query().Get("qualification_id")

	// Optional Queries
	query := r.URL.Query().Get("qualification_query")

	// Validate required queries
	var redirect bool
	if requirementID == "" {
		logger.LogAttrs(r.Context(), slog.LevelWarn, "Missing requirement ID, redirecting to qualifications pane")
		redirect = true
	}
	if requirementType == "" {
		logger.LogAttrs(r.Context(), slog.LevelWarn, "Missing requirement type, redirecting to qualifications pane")
		redirect = true
	}
	if qualificationID == "" {
		logger.LogAttrs(r.Context(), slog.LevelWarn, "Missing parent qualification ID, redirecting to qualifications pane")
		redirect = true
	}
	if redirect {
		http.Redirect(w, r, "/admin?fragment=qualifications_pane", http.StatusFound)
		return
	}

	// Get and filter qualifications
	qualifications, err := qualificationStore.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if query != "" {
		qualifications = types.FilterQualifications(qualifications, func(q types.Qualification) bool {
			return strings.Contains(strings.ToLower(q.Name), strings.ToLower(query))
		})
	}

	// Get Selected parent qualification
	qualification, err := qualificationStore.GetQualification(qualificationID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get Selected requirement
	requirement, err := qualificationStore.GetRequirement(requirementID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get SelectedMember
	member, err := MemberFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var hasSubordinates bool
	if subordinates := memberStore.GetSubordinates(member.ID); len(subordinates) > 0 {
		hasSubordinates = true
	}

	adminContent := admin.AdminPage("Qualifications", admin.RequirementEditorPane(
		qualifications,
		qualification,
		requirement,
		requirementType == "initial",
		qualificationID,
		false,
		query,
	))

	rootTemplate := templates.Root(true, hasSubordinates, true, member.Display(), organization, adminContent)
	HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
}
