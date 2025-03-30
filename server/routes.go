package server

import (
	"PORTal/server/stores"
	"PORTal/types"
	"context"
	"log/slog"
	"net/http"
)

func addRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	logger *slog.Logger,
	assetHandler http.Handler,
	organization string,
	memberStore stores.MemberStore,
	sessionStore stores.SessionStore,
	qualificationStore stores.QualificationStore,
	ranks types.RankMap,
) {
	mux.Handle("GET /assets/", assetHandler)

	logger.LogAttrs(ctx, slog.LevelInfo, "Registering index, login, and logout routes...")
	mux.Handle("/", LoginDirectorHandler(logger, organization, memberStore, sessionStore))
	mux.Handle("GET /logout", LogoutHandler(logger, sessionStore))

	logger.LogAttrs(ctx, slog.LevelInfo, "Registering admin routes...")
	mux.Handle("GET /admin", adminRootGetHandler(logger, memberStore, qualificationStore, ranks))
	mux.Handle("GET /admin/potentialSupervisors", getPotentialSupervisorsHandler(logger, memberStore))
	mux.Handle("PUT /admin/updateMember", updateMemberHandler(logger, memberStore))

	logger.LogAttrs(ctx, slog.LevelInfo, "Registering dashboard routes...")
	mux.Handle("GET /dashboard", dashboardGetHandler(logger))
}
