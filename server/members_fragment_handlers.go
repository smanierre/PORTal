package server

import (
	"PORTal/server/stores"
	"PORTal/templates/pages/members"
	"PORTal/types"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

func membersFilterHandler(logger *slog.Logger, memberStore stores.MemberStore) http.Handler {
	logger = logger.With(slog.String("route", "GET /members/filter"))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Hx-Request") != "true" {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Non-HTMX request to members/filter, rendering members page")
			http.Redirect(w, r, "/members", http.StatusSeeOther)
			return
		}

		// Potential query params
		query := r.URL.Query().Get("query")
		selectedID := r.URL.Query().Get("selected")

		mems, err := memberStore.GetAllMembers()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if query != "" {
			mems = types.FilterMembers(mems, func(m types.Member) bool {
				return strings.Contains(strings.ToLower(fmt.Sprintf("%s %s", m.FirstName, m.LastName)), strings.ToLower(query))
			})
		}

		var m types.Member
		if selectedID != "" {
			m, err = memberStore.GetMember(selectedID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		membersList := members.MemberList(false, mems, m, query)
		HandleRenderError(r.Context(), logger, membersList.Render(r.Context(), w))
	})
}
