package server

//func adminRootGetHandler(logger *slog.Logger, core core.Core, organization string) http.Handler {
//	logger = logger.With("route", "GET /admin")
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		if r.Header.Get("Hx-Request") == "true" {
//			logger.LogAttrs(r.Context(), slog.LevelInfo, "HTMX Request, directing to fragment handler")
//			//adminPanelFragmentHandler(logger, memberStore, qualificationStore, ranks, w, r)
//			return
//		}
//
//		// Determine which fragment is being requested and render based on that
//		switch r.URL.Query().Get("fragment") {
//		case "qualifications_pane":
//			//rootQualificationsPaneHandler(logger, memberStore, qualificationStore, organization, w, r)
//			return
//		case "requirement_editor":
//			//rootRequirementEditorHandler(logger, memberStore, qualificationStore, organization, w, r)
//			return
//		default:
//			rootMembersPaneHandler(logger, core, w, r)
//		}
//	})
//}
//
//func rootMembersPaneHandler(
//	logger *slog.Logger,
//	core core.Core,
//	w http.ResponseWriter, r *http.Request,
//) {
//	data, err := core.AdminMemberPage(r)
//	if err != nil {
//		HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
//		return
//	}
//	tld, err := core.LoggedInTopLevel(r)
//	if err != nil {
//		HandleRenderError(r.Context(), logger, errorpages.GenericISE().Render(r.Context(), w))
//		return
//	}
//	HandleRenderError(r.Context(), logger, templates.Root(tld.ShowNav, tld.HasSubordinates, tld.Admin, tld.MemberDisplayName, tld.Organization,
//		admin.AdminPage("Members", admin.MembersPane(
//			data.Members,
//			data.SelectedMember,
//			core.Ranks,
//			data.PotentialSupervisors,
//			false,
//			data.DisabledMembers,
//			false,
//			data.MemberQuery))).Render(r.Context(), w))
//
//}
//
//func rootQualificationsPaneHandler(
//	logger *slog.Logger,
//	memberStore stores.MemberStore,
//	qualificationStore stores.QualificationStore,
//	organization string,
//	w http.ResponseWriter, r *http.Request,
//) {
//	logger.LogAttrs(r.Context(), slog.LevelInfo, "Rendering root template with qualifications pane")
//
//	// Potential queries
//	selected := r.URL.Query().Get("selected")
//	query := r.URL.Query().Get("qualification_query")
//	newQualification := r.URL.Query().Get("new") == "true"
//
//	qualifications, err := qualificationStore.GetAllQualifications()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	if query != "" {
//		qualifications = types.FilterQualifications(qualifications, func(q types.Qualification) bool {
//			return strings.Contains(strings.ToLower(q.Name), strings.ToLower(query))
//		})
//	}
//	var selectedQualification types.Qualification
//	if selected != "" {
//		selectedQualification, err = qualificationStore.GetQualification(selected)
//		if errors.Is(err, backend.ErrQualificationNotFound) {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		} else if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		fmt.Printf("%+v\n", selectedQualification)
//		for i, requirement := range selectedQualification.InitialRequirements {
//			if requirement.Type == types.QualificationType {
//				qual, err := qualificationStore.GetQualification(requirement.QualificationID)
//				if err != nil {
//					w.WriteHeader(http.StatusInternalServerError)
//					return
//				}
//				selectedQualification.InitialRequirements[i].Name = qual.Name
//			}
//		}
//	}
//
//	loggedInMember, err := MemberFromContext(r.Context())
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	var hasSubordinates bool
//	if subordinates := memberStore.GetSubordinates(loggedInMember.ID); len(subordinates) > 0 {
//		hasSubordinates = true
//	}
//
//	adminContent := admin.AdminPage("Qualifications",
//		admin.QualificationPane(qualifications, selectedQualification, newQualification, false, query))
//	rootTemplate := templates.Root(true, hasSubordinates, true, loggedInMember.Display("f"), organization, adminContent)
//	HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
//}
//
//func rootRequirementEditorHandler(
//	logger *slog.Logger,
//	memberStore stores.MemberStore,
//	qualificationStore stores.QualificationStore,
//	organization string,
//	w http.ResponseWriter,
//	r *http.Request,
//) {
//	logger.LogAttrs(r.Context(), slog.LevelInfo, "Rendering root template with requirement editor")
//
//	// Required Queries
//	requirementID := r.URL.Query().Get("requirement_id")
//	requirementType := r.URL.Query().Get("requirement_type")
//	qualificationID := r.URL.Query().Get("qualification_id")
//
//	// Optional Queries
//	query := r.URL.Query().Get("qualification_query")
//
//	// Validate required queries
//	var redirect bool
//	if requirementID == "" {
//		logger.LogAttrs(r.Context(), slog.LevelWarn, "Missing requirement ID, redirecting to qualifications pane")
//		redirect = true
//	}
//	if requirementType == "" {
//		logger.LogAttrs(r.Context(), slog.LevelWarn, "Missing requirement type, redirecting to qualifications pane")
//		redirect = true
//	}
//	if qualificationID == "" {
//		logger.LogAttrs(r.Context(), slog.LevelWarn, "Missing parent qualification ID, redirecting to qualifications pane")
//		redirect = true
//	}
//	if redirect {
//		http.Redirect(w, r, "/admin?fragment=qualifications_pane", http.StatusFound)
//		return
//	}
//
//	// Get and filter qualifications
//	qualifications, err := qualificationStore.GetAllQualifications()
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	if query != "" {
//		qualifications = types.FilterQualifications(qualifications, func(q types.Qualification) bool {
//			return strings.Contains(strings.ToLower(q.Name), strings.ToLower(query))
//		})
//	}
//
//	// Get Selected parent qualification
//	qualification, err := qualificationStore.GetQualification(qualificationID)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	// Get Selected requirement
//	requirement, err := qualificationStore.GetRequirement(requirementID)
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//
//	// Get SelectedMember
//	member, err := MemberFromContext(r.Context())
//	if err != nil {
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	var hasSubordinates bool
//	if subordinates := memberStore.GetSubordinates(member.ID); len(subordinates) > 0 {
//		hasSubordinates = true
//	}
//
//	adminContent := admin.AdminPage("Qualifications", admin.RequirementEditorPane(
//		qualifications,
//		qualification,
//		requirement,
//		requirementType == "initial",
//		qualificationID,
//		false,
//		query,
//	))
//
//	rootTemplate := templates.Root(true, hasSubordinates, true, member.Display("f"), organization, adminContent)
//	HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
//}
