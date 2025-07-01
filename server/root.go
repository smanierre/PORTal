package server

import (
	"PORTal/server/stores"
	"PORTal/templates"
	"PORTal/templates/pages/errorpages"
	"github.com/a-h/templ"
	"log/slog"
	"net/http"
)

// Gets set during server creation so  memberstore and organization don't need to be passed on each call
var serveContentAsRoot func(component templ.Component, logger *slog.Logger, w http.ResponseWriter, r *http.Request)

func initializeServeContentAsRoot(memberStore stores.MemberStore, organization string) func(content templ.Component, logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
	return func(content templ.Component, logger *slog.Logger, w http.ResponseWriter, r *http.Request) {
		var showNav bool
		var subordinates bool
		loggedInMember, err := MemberFromContext(r.Context())
		if err != nil {
			logger.LogAttrs(r.Context(), slog.LevelInfo, "Error getting member from context", slog.String("error", err.Error()))
			// If it's not a request for the home page, respond with an internal server error as the member should be in context
			if r.URL.Path != "/" {
				HandleRenderError(r.Context(), logger, templates.Root(
					false, false, false, "", "", errorpages.GenericISE()).Render(r.Context(), w))
				return
			}
		}
		if r.URL.Path != "/" && loggedInMember.ID != "" {
			showNav = true
		}
		if subs := memberStore.GetSubordinates(loggedInMember.ID); len(subs) > 0 {
			subordinates = true
		}

		rootTemplate := templates.Root(showNav, subordinates, loggedInMember.Admin, loggedInMember.Display(), organization, content)
		HandleRenderError(r.Context(), logger, rootTemplate.Render(r.Context(), w))
	}
}
