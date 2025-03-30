package server

import (
	"PORTal/server/serverutils"
	"PORTal/server/stores"
	"PORTal/types"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
)

//go:embed assets
var assetsDir embed.FS

type Server struct {
	logger  *slog.Logger
	handler http.Handler
	port    string
	dev     bool
}

func New(
	ctx context.Context,
	logger *slog.Logger,
	port string,
	dev bool,
	domain string,
	service string,
	organization string,
	memberStore stores.MemberStore,
	qualificationStore stores.QualificationStore,
	sessionStore stores.SessionStore,
) Server {
	l := logger.With(slog.String("source", "Server"))
	l.LogAttrs(ctx, slog.LevelInfo, fmt.Sprintf("Setting domain to: %s", domain))
	serverutils.SetDomain(domain)
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating new server")
	s := Server{
		logger: l,
		dev:    dev,
		port:   port,
	}

	mux := http.NewServeMux()
	l.LogAttrs(ctx, slog.LevelInfo, "Registering static asset routes...")
	// Static assets
	assetHandler := http.FileServerFS(assetsDir)

	addRoutes(ctx, mux, logger, assetHandler, organization, memberStore, sessionStore, qualificationStore, types.GetRanks(service))

	// Handlers are applied last to first
	l.LogAttrs(ctx, slog.LevelInfo, "Applying global middlewares...")
	s.handler = adminRequiredMiddleware(mux, memberStore, logger)
	s.handler = skipLoginMiddleware(s.handler, logger, sessionStore)
	s.handler = sessionRequiredMiddleware(s.handler, logger, organization, sessionStore)
	s.handler = assetMiddleware(s.handler, assetHandler)
	//
	//	// Index is the login page
	//	s.mux.Handle("GET /", http.HandlerFunc(s.LoginGetHandler))
	//s.mux.Handle("POST /", http.HandlerFunc(s.LoginPostHandler))
	//s.mux.Handle("GET /logout", http.HandlerFunc(s.LogoutHandler))
	//
	//// Dashboard
	//s.mux.Handle("GET /dashboard", http.HandlerFunc(s.DashboardGetHandler))
	//
	//// Admin page
	//s.mux.Handle("GET /admin", http.HandlerFunc(s.AdminGetHandler))
	//
	//// Admin Member Routes
	//s.mux.Handle("GET /admin/members", http.HandlerFunc(s.AdminMembersPaneGetHandler))
	//s.mux.Handle("GET /admin/members/disabled", http.HandlerFunc(s.AdminMembersPaneGetDisabledHandler))
	//s.mux.Handle("GET /admin/members/{id}", http.HandlerFunc(s.AdminMemberEditorGetHandler))
	//s.mux.Handle("POST /admin/members/{id}", http.HandlerFunc(s.AdminMemberUpdateHandler))
	//s.mux.Handle("POST /admin/members/{id}/disable", http.HandlerFunc(s.AdminMemberDisableHandler))
	//s.mux.Handle("POST /admin/members/{id}/enable", http.HandlerFunc(s.AdminMemberEnableHandler))
	//s.mux.Handle("GET /admin/members/add", http.HandlerFunc(s.AdminNewMemberHandler))
	//s.mux.Handle("POST /admin/members/add", http.HandlerFunc(s.AdminMemberAddHandler))
	//s.mux.Handle("GET /admin/members/{id}/potentialSupervisors/{grade}", http.HandlerFunc(s.AdminGetPotentialSupervisorsHandler))
	//
	//// Admin Qualification Routes
	//s.mux.Handle("GET /admin/qualifications", http.HandlerFunc(s.AdminQualificationsPaneGetHandler))
	//s.mux.Handle("GET /admin/qualifications/{id}", http.HandlerFunc(s.AdminQualificationEditorGetHandler))
	//s.mux.Handle("POST /admin/qualifications/{id}", http.HandlerFunc(s.AdminQualificationUpdateHandler))
	//s.mux.Handle("GET /admin/qualifications/add", http.HandlerFunc(s.AdminNewQualificationHandler))
	//s.mux.Handle("POST /admin/qualifications/add", http.HandlerFunc(s.AdminQualificationAddHandler))
	//
	//// Component Routes
	//s.mux.Handle("GET /components/requirementItem", http.HandlerFunc(s.RequirementItemComponent))
	//s.mux.Handle("GET /components/requirementEditor", http.HandlerFunc(s.RequirementEditorComponent))
	l.LogAttrs(ctx, slog.LevelInfo, "Successfully registered routes")
	return s
}

func (s Server) ListenAndServe() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Starting server on port: %s", s.port))
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.handler)
}
