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
	ranks := types.GetRanks(service)
	addRoutes(ctx, mux, logger, assetHandler, organization, memberStore, sessionStore, qualificationStore, ranks)

	// Handlers are applied last to first
	l.LogAttrs(ctx, slog.LevelInfo, "Applying global middlewares...")
	s.handler = adminRequiredMiddleware(mux, ranks, memberStore, logger)
	s.handler = skipLoginMiddleware(s.handler, logger, sessionStore)
	s.handler = sessionRequiredMiddleware(s.handler, logger, organization, sessionStore)
	s.handler = assetMiddleware(s.handler, assetHandler)

	l.LogAttrs(ctx, slog.LevelInfo, "Successfully registered routes")
	return s
}

func (s Server) ListenAndServe() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Starting server on port: %s", s.port))
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.handler)
}
