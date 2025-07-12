//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../tools/config.yaml ../../../api.yaml
package api

import (
	"PORTal/server/stores"
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
)

//go:embed dist
var specUI embed.FS

func RegisterRoutes(logger *slog.Logger, mux *http.ServeMux, memberStore stores.MemberStore) {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering API Routes")
	// Register the swagger docs site
	subFS, err := fs.Sub(specUI, "dist")
	if err != nil {
		logger.LogAttrs(context.Background(), slog.LevelInfo, "Error getting API doc static files", slog.String("err", err.Error()))
	}
	mux.Handle("GET /api/docs/", http.StripPrefix("/api/docs", http.FileServer(http.FS(subFS))))
	mux.Handle("GET /api/spec", GetApiSpec(logger))

	a := api{
		logger:      logger,
		memberStore: memberStore,
	}
	m := http.NewServeMux()
	h := HandlerFromMux(a, m)

	// All methods have to be registered separately or they will collide with the GET "/" route
	mux.Handle("GET /api/", h)
	mux.Handle("PUT /api/", h)
	mux.Handle("DELETE /api/", h)
	mux.Handle("POST /api/", h)
}
