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
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

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
	case "qualifications_pane":
		qualificationsPaneFragmentHandler(logger, qualificationStore, w, r)
	case "requirement_editor":
		requirementEditorFragmentHandler(logger, qualificationStore, w, r)
	default:
		membersPaneFragmentHandler(logger, memberStore, ranks, w, r)
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
		if r.URL.Query().Get("disabled") == "true" {
			member, err = memberStore.GetDisabledMember(selectedMemberQuery)
		} else {
			member, err = memberStore.GetMember(selectedMemberQuery)
		}
		if errors.Is(err, backend.ErrMemberNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if r.URL.Query().Get("disabled") != "true" {
			potentialSupervisors, err = memberStore.GetPotentialSupervisors(member, member.Grade)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
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
	adminPage := admin.AdminPage("Members", membersPane)
	HandleRenderError(r.Context(), logger, adminPage.Render(r.Context(), w))
}

func qualificationsPaneFragmentHandler(logger *slog.Logger, qualificationStore stores.QualificationStore, w http.ResponseWriter, r *http.Request) {
	logger = logger.With("handler", "qualificationsPaneFragmentHandler")

	// Potential Query Params
	selectedQualificationQuery := r.URL.Query().Get("selected")
	query := r.URL.Query().Get("qualification_query")
	newQual := r.URL.Query().Get("new") == "true"

	qualifications, err := qualificationStore.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var selectedQualification types.Qualification
	if selectedQualificationQuery != "" {
		selectedQualification, err = qualificationStore.GetQualification(selectedQualificationQuery)
		if errors.Is(err, backend.ErrQualificationNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// We need to get the name of any qualification requirement for display
		for i, requirement := range selectedQualification.InitialRequirements {
			if requirement.Type == types.QualificationType {
				q, err := qualificationStore.GetQualification(requirement.QualificationID)
				if errors.Is(err, backend.ErrQualificationNotFound) {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				selectedQualification.InitialRequirements[i].Name = q.Name
			}
		}
	}
	if query != "" {
		qualifications = types.FilterQualifications(qualifications, func(q types.Qualification) bool {
			return strings.Contains(strings.ToLower(q.Name), strings.ToLower(query))
		})
	}

	qualificationsPane := admin.QualificationPane(qualifications, selectedQualification, newQual, false, query)
	HandleRenderError(r.Context(), logger, qualificationsPane.Render(r.Context(), w))

	adminPage := admin.AdminPage("Qualifications", qualificationsPane)
	HandleRenderError(r.Context(), logger, adminPage.Render(r.Context(), w))
}

func requirementEditorFragmentHandler(logger *slog.Logger, qualificationStore stores.QualificationStore, w http.ResponseWriter, r *http.Request) {
	logger = logger.With("handler", "requirementEditorFragmentHandler")

	// Possible query parameters
	parentQualificationID := r.URL.Query().Get("qualification_id")
	requirementType := r.URL.Query().Get("type")
	requirementId := r.URL.Query().Get("requirement_id")
	initialOrRecurring := r.URL.Query().Get("requirement_type")

	var requirement types.Requirement
	var err error
	if requirementId != "" {
		requirement, err = qualificationStore.GetRequirement(requirementId)
		if err != nil && !errors.Is(err, backend.ErrRequirementNotFound) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	if requirement.Type == "" {
		requirement.Type = types.RequirementType(requirementType)
	}
	qualifications, err := qualificationStore.GetAllQualifications()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var requirementEditor templ.Component
	if initialOrRecurring == "initial" {
		requirementEditor = admin.InitialRequirementEditor(requirement, qualifications, parentQualificationID)
	} else if initialOrRecurring == "recurring" {
		requirementEditor = admin.RecurringRequirementEditor(requirement, qualifications, parentQualificationID)
	} else {
		logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid requirement type for initial or recurring", slog.String("type", initialOrRecurring))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	HandleRenderError(r.Context(), logger, requirementEditor.Render(r.Context(), w))
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
		potentialSupervisors, err := memberStore.GetPotentialSupervisors(m, grade)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		supervisorList := admin.SupervisorList(potentialSupervisors, m.SupervisorID)
		HandleRenderError(r.Context(), logger, supervisorList.Render(r.Context(), w))
	})
}

func updateMemberHandler(logger *slog.Logger, memberStore stores.MemberStore, ranks types.RankMap, organization string) http.Handler {
	logger = logger.With("route", "PUT /admin/updateMember")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "Error parsing form", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Updating member", slog.String("member", r.Form.Get("id")))
		fmt.Println(r.Form)
		member, err := memberStore.UpdateMember(types.Member{
			ID:           r.Form.Get("id"),
			FirstName:    r.Form.Get("first_name"),
			LastName:     r.Form.Get("last_name"),
			Username:     r.Form.Get("username"),
			Grade:        types.Grade(r.Form.Get("grade")),
			SupervisorID: r.Form.Get("supervisor"),
			Admin:        r.Form.Get("admin") == "on",
			Password:     r.Form.Get("password"),
			Disabled:     false,
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		members, err := memberStore.GetAllMembers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
			return
		}
		potentialSupervisors, err := memberStore.GetPotentialSupervisors(member, member.Grade)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
			return
		}
		loggedInMember, err := MemberFromContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
			return
		}
		// If the logged in member was the one that was updated, get the updated member
		if loggedInMember.ID == member.ID {
			loggedInMember = member
		}
		subordinates := memberStore.GetSubordinates(loggedInMember.ID)

		displayName := fmt.Sprintf("%s %s %s", ranks[loggedInMember.Grade], loggedInMember.FirstName, loggedInMember.LastName)

		membersPane := admin.MembersPane(
			members,
			member,
			ranks,
			potentialSupervisors,
			false,
			false,
			false,
			"",
		)

		adminPage := admin.AdminPage("Members", membersPane)

		rootTemplate := templates.Root(
			true,
			len(subordinates) > 0,
			loggedInMember.Admin,
			displayName,
			organization,
			adminPage)
		HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
	})
}

func newMemberHandler(logger *slog.Logger, ranks types.RankMap) http.Handler {
	logger = logger.With("route", "GET /admin/newMember")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") != "true" {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request, redirecting to members pane")
			http.Redirect(w, r, "/admin?fragment=members_pane", http.StatusFound)
			return
		}
		membersEditor := admin.MemberEditor(types.Member{}, ranks, nil, true)
		HandleRenderError(r.Context(), logger, membersEditor.Render(r.Context(), w))
	})
}

func addMemberHandler(logger *slog.Logger, memberStore stores.MemberStore, ranks types.RankMap, organization string) http.Handler {
	logger = logger.With("route", "POST /admin/addMember")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Creating new member")
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelError, "Unable to parse form data for new member", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		m := types.Member{}
		m.Username = r.Form.Get("username")
		m.FirstName = r.Form.Get("first_name")
		m.LastName = r.Form.Get("last_name")
		m.Password = r.Form.Get("password")
		m.Grade = types.Grade(r.Form.Get("grade"))
		m.Admin = r.Form.Get("admin") == "on"
		m.SupervisorID = r.Form.Get("supervisor_id")

		m, err = memberStore.AddMember(m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		members, err := memberStore.GetAllMembers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		potentialSupervisors, err := memberStore.GetPotentialSupervisors(m, m.Grade)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		loggedInMember, err := MemberFromContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		subordinates := memberStore.GetSubordinates(loggedInMember.ID)
		displayName := fmt.Sprintf("%s %s %s", ranks[loggedInMember.Grade], loggedInMember.FirstName, loggedInMember.LastName)

		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin?selected=%s", m.ID))
		membersPane := admin.MembersPane(members, m, ranks, potentialSupervisors, false, false, false, "")
		adminPage := admin.AdminPage("Members", membersPane)
		rootTemplate := templates.Root(true, len(subordinates) > 0, loggedInMember.Admin, displayName, organization, adminPage)
		HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
	})
}

func disableMemberHandler(logger *slog.Logger, ranks types.RankMap, memberStore stores.MemberStore) http.Handler {
	logger = logger.With("route", "PUT /admin/disableMember")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		memberID := r.Form.Get("id")
		if memberID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = memberStore.DisableMember(memberID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		members, err := memberStore.GetAllMembers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Push-URL", "/admin?fragment=members_pane&disabled=false")
		membersPane := admin.MembersPane(members, types.Member{}, ranks, nil, false, false, false, "")
		HandleRenderError(r.Context(), logger, membersPane.Render(r.Context(), w))
	})
}

func enableMemberHandler(logger *slog.Logger, ranks types.RankMap, memberStore stores.MemberStore) http.Handler {
	logger = logger.With("route", "PUT /admin/enableMember")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		memberID := r.Form.Get("id")
		if memberID == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = memberStore.EnableMember(memberID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		members, err := memberStore.GetDisabledMembers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Push-URL", "/admin?fragment=members_pane&disabled=true")
		membersPane := admin.MembersPane(members, types.Member{}, ranks, nil, false, true, false, "")
		HandleRenderError(r.Context(), logger, membersPane.Render(r.Context(), w))
	})
}

func newQualificationHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With("route", "GET /admin/newQualification")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") != "true" {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request, redirecting to qualification fragment")
			http.Redirect(w, r, "/admin?fragment=qualifications_pane", http.StatusFound)
		}
		qualificationEditor := admin.QualificationEditor(types.Qualification{}, true)
		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualificationList := admin.QualificationList(qualifications, types.Qualification{}, true, "")

		HandleRenderError(r.Context(), logger, qualificationEditor.Render(r.Context(), w))
		HandleRenderError(r.Context(), logger, qualificationList.Render(r.Context(), w))
	})
}

func addQualificationHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With("route", "POST /admin/addQualification")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Error parsing form", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		logger.LogAttrs(r.Context(), slog.LevelInfo, "Adding new qualification")
		qualification := types.Qualification{
			Name:                  r.Form.Get("name"),
			InitialRequirements:   nil,
			RecurringRequirements: nil,
			Notes:                 r.Form.Get("notes"),
		}

		qualification, err = qualificationStore.AddQualification(qualification)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("HX-Push-URL", fmt.Sprintf("/admin?fragment=qualifications_pane&selected=%s", qualification.ID))
		qualificationsPane := admin.QualificationPane(qualifications, qualification, false, false, "")
		HandleRenderError(r.Context(), logger, qualificationsPane.Render(r.Context(), w))
	})
}

func updateQualificationHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With("route", "PUT /admin/updateQualification")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "Error parsing form for qualification update", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		qualification := types.Qualification{
			ID:                    r.Form.Get("id"),
			Name:                  r.Form.Get("name"),
			InitialRequirements:   nil,
			RecurringRequirements: nil,
			Notes:                 r.Form.Get("notes"),
		}

		qualification, err = qualificationStore.UpdateQualification(qualification)
		if errors.Is(err, backend.ErrQualificationNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Check for any qualification requirements and update the name of them to display correctly on creation
		for i, q := range qualifications {
			for j, req := range q.InitialRequirements {
				if req.Type == types.QualificationType {
					qualReq, err := qualificationStore.GetQualification(req.QualificationID)
					if err != nil {
						logger.LogAttrs(r.Context(), slog.LevelWarn, "Couldn't get name for qualification, rendering without")
						continue
					}
					qualifications[i].InitialRequirements[j].Name = qualReq.Name
				}
			}
		}

		qualificationPane := admin.QualificationPane(qualifications, qualification, false, false, "")
		HandleRenderError(r.Context(), logger, qualificationPane.Render(r.Context(), w))
	})
}

func newRequirementHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With("route", "GET /admin/newRequirement")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Possible query parameters
		qualificationID := r.URL.Query().Get("id")
		initialOrRecurring := r.URL.Query().Get("requirement_type")

		if r.Header.Get("Hx-Request") != "true" {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request, redirecting to qualification editor")
			http.Redirect(w, r, fmt.Sprintf("/admin?fragment=qualifications_pane&selected=%s", qualificationID), http.StatusFound)
			return
		}

		logger.LogAttrs(r.Context(), slog.LevelInfo, fmt.Sprintf("Rendering new requirement editor for qualification: %s", qualificationID))

		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var requirementEditor templ.Component
		if initialOrRecurring == "initial" {
			requirementEditor = admin.InitialRequirementEditor(types.Requirement{}, qualifications, qualificationID)
		} else if initialOrRecurring == "recurring" {
			requirementEditor = admin.RecurringRequirementEditor(types.Requirement{}, qualifications, qualificationID)
		} else {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "Invalid requirement type for initial or recurring", slog.String("type", initialOrRecurring))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		HandleRenderError(r.Context(), logger, requirementEditor.Render(r.Context(), w))
	})
}

func addRequirementHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With("route", "POST /admin/addRequirement")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Error parsing form", slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		initial := r.URL.Query().Get("requirement_type") == "initial"
		var daysValidFor int
		if !initial {
			daysValidFor, err = strconv.Atoi(r.Form.Get("days_valid_for"))
			if types.RequirementType(r.Form.Get("type")) != types.QualificationType && err != nil {
				logger.LogAttrs(r.Context(), slog.LevelWarn, "Error parsing days_valid_for", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		requirement := types.Requirement{
			Name:            r.Form.Get("name"),
			Reference:       r.Form.Get("reference"),
			QualificationID: r.Form.Get("qualification_id"),
			Type:            types.RequirementType(r.Form.Get("type")),
			Notes:           r.Form.Get("notes"),
			DaysValidFor:    daysValidFor,
			Grade:           types.Grade(r.Form.Get("grade")),
			Initial:         initial,
		}
		requirement, err = qualificationStore.AddRequirement(requirement)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = qualificationStore.AssignRequirementToQualification(r.Form.Get("parent_qualification"), requirement.ID, initial)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualification, err := qualificationStore.GetQualification(r.Form.Get("parent_qualification"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualificationsPane := admin.QualificationPane(qualifications, qualification, false, false, "")
		HandleRenderError(r.Context(), logger, qualificationsPane.Render(r.Context(), w))
	})
}

func updateRequirementHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With("route", "PUT /admin/updateRequirement")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelInfo, fmt.Sprintf("Error parsing form: %s", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		initial := r.URL.Query().Get("requirement_type") == "initial"
		var daysValidFor int
		if !initial {
			daysValidFor, err = strconv.Atoi(r.Form.Get("days_valid_for"))
			if types.RequirementType(r.Form.Get("type")) != types.QualificationType && err != nil {
				logger.LogAttrs(r.Context(), slog.LevelWarn, "Error parsing days_valid_for", slog.String("error", err.Error()))
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		requirement := types.Requirement{
			ID:              r.Form.Get("requirement_id"),
			Name:            r.Form.Get("name"),
			Reference:       r.Form.Get("reference"),
			QualificationID: r.Form.Get("qualification_id"),
			Grade:           types.Grade(r.Form.Get("grade")),
			Notes:           r.Form.Get("notes"),
			DaysValidFor:    daysValidFor,
			Type:            types.RequirementType(r.Form.Get("type")),
			Initial:         initial,
		}
		requirement, err = qualificationStore.UpdateRequirement(requirement)
		if errors.Is(err, backend.ErrRequirementNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		selectedQualification, err := qualificationStore.GetQualification(r.Form.Get("parent_qualification"))
		if errors.Is(err, backend.ErrQualificationNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualificationsPane := admin.QualificationPane(qualifications, selectedQualification, false, false, "")
		HandleRenderError(r.Context(), logger, qualificationsPane.Render(r.Context(), w))
	})
}

func removeRequirementHandler(logger *slog.Logger, qualificationStore stores.QualificationStore) http.Handler {
	logger = logger.With(slog.String("route", "DELETE /admin/removeRequirement"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") != "true" {
			// TODO: What to do here?
			logger.LogAttrs(r.Context(), slog.LevelWarn, "Request to remove requirement should be an HTMX request but isn't.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err := r.ParseForm()
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelInfo, fmt.Sprintf("Error parsing form: %s", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		requirementID := r.Form.Get("requirement_id")
		if requirementID == "" {
			logger.LogAttrs(r.Context(), slog.LevelWarn, "Requirement id not provided")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		qualificationStore.DeleteRequirement(requirementID)

		qualifications, err := qualificationStore.GetAllQualifications()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		selectedQualification, err := qualificationStore.GetQualification(r.Form.Get("id"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		qualificationsPane := admin.QualificationPane(qualifications, selectedQualification, false, false, "")
		HandleRenderError(r.Context(), logger, qualificationsPane.Render(r.Context(), w))
	})
}
