package api

import (
	"PORTal/server/core"
	"context"
	"log/slog"
	"net/http"
)

func RegisterRoutes(logger *slog.Logger, mux *http.ServeMux, c core.Core) {
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Registering API Routes")

	mux.Handle("GET /api/potentialSupervisors", PotentialSupervisors(logger, c))
}
