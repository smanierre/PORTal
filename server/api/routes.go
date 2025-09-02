package api

import (
	"PORTal/server/stores"
	"PORTal/types"
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

//go:embed swagger
var specUI embed.FS

func RegisterRoutes(
	logger *slog.Logger,
	mux *http.ServeMux,
	memberStore stores.MemberStore,
	qualificationStore stores.QualificationStore,
	memberQualificationStore stores.MemberQualificationStore,
	sessionStore stores.SessionStore,
	domain, organization string,
	dev bool,
	ranks types.RankMap,
) {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering API Routes")
	// Register the swagger docs site
	subFS, err := fs.Sub(specUI, "swagger")
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelInfo, "Error getting API doc static files", slog.String("err", err.Error()))
	}
	mux.Handle("GET /api/docs/", http.StripPrefix("/api/docs", http.FileServer(http.FS(subFS))))
	mux.Handle("GET /api/spec", GetApiSpec(logger))

	a := api{
		logger:                   logger,
		memberStore:              memberStore,
		qualificationStore:       qualificationStore,
		memberQualificationStore: memberQualificationStore,
		sessionStore:             sessionStore,
		domain:                   domain,
		organization:             organization,
		dev:                      dev,
		ranks:                    ranks,
	}
	m := http.NewServeMux()
	h := HandlerFromMux(a, m)

	// All methods have to be registered separately, or they will collide with the GET "/" route
	mux.Handle("GET /api/", h)
	mux.Handle("PUT /api/", h)
	mux.Handle("DELETE /api/", h)
	mux.Handle("POST /api/", h)
}
