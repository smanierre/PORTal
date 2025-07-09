package server

//func membersRootGetHandler(logger *slog.Logger, memberStore stores.MemberStore) http.Handler {
//	logger = logger.With(slog.String("route", "GET /members"))
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//
//		// Potential Query Params
//		selected := r.URL.Query().Get("selected")
//		var member types.Member
//		var err error
//		if selected != "" {
//			member, err = memberStore.GetMember(selected)
//			if err != nil {
//				w.WriteHeader(http.StatusInternalServerError)
//				return
//			}
//		}
//		ms, err := memberStore.GetAllMembers()
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//		content := members.MemberPage(ms, member, "")
//		if r.Header.Get("Hx-Request") == "true" {
//			HandleRenderError(r.Context(), logger, content.Render(r.Context(), w))
//		} else {
//			serveContentAsRoot(content, logger, w, r)
//		}
//	})
//}
