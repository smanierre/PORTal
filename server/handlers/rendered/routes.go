package rendered

import (
	"PORTal/server/core"
	"context"
	"log/slog"
	"net/http"
)

func RegisterRoutes(
	ctx context.Context,
	mux *http.ServeMux,
	logger *slog.Logger,
	c core.Core,
) {

	logger.LogAttrs(ctx, slog.LevelInfo, "Registering index, login, and logout routes...")
	mux.Handle("GET /", LoginGetHandler(logger, c))
	mux.Handle("POST /", LoginPostHandler(logger, c))
	mux.Handle("GET /logout", LogoutHandler(logger, c))

	logger.LogAttrs(ctx, slog.LevelInfo, "Registering admin routes...")
	mux.Handle("GET /admin", AdminGetHandler(logger, c))
	//mux.Handle("GET /admin/potentialSupervisors", getPotentialSupervisorsHandler(logger, memberStore))
	//mux.Handle("PUT /admin/updateMember", updateMemberHandler(logger, memberStore, ranks, organization))
	//mux.Handle("GET /admin/newMember", newMemberHandler(logger, ranks))
	//mux.Handle("POST /admin/addMember", addMemberHandler(logger, memberStore, ranks, organization))
	//mux.Handle("PUT /admin/disableMember", disableMemberHandler(logger, ranks, memberStore))
	//mux.Handle("PUT /admin/enableMember", enableMemberHandler(logger, ranks, memberStore))
	//mux.Handle("GET /admin/newQualification", newQualificationHandler(logger, qualificationStore))
	//mux.Handle("POST /admin/addQualification", addQualificationHandler(logger, qualificationStore))
	//mux.Handle("PUT /admin/updateQualification", updateQualificationHandler(logger, qualificationStore))
	//mux.Handle("GET /admin/newRequirement", newRequirementHandler(logger, qualificationStore))
	//mux.Handle("POST /admin/addRequirement", addRequirementHandler(logger, qualificationStore))
	//mux.Handle("PUT /admin/updateRequirement", updateRequirementHandler(logger, qualificationStore))
	//mux.Handle("DELETE /admin/removeRequirement", removeRequirementHandler(logger, qualificationStore))
	//
	//logger.LogAttrs(ctx, slog.LevelInfo, "Registering members routes...")
	//mux.Handle("GET /members", membersRootGetHandler(logger, memberStore))
	//mux.Handle("GET /members/filter", membersFilterHandler(logger, memberStore))
	//
	logger.LogAttrs(ctx, slog.LevelInfo, "Registering dashboard routes...")
	mux.Handle("GET /dashboard", dashboardGetHandler(logger, c))
}
