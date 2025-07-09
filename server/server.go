package server

import (
	"PORTal/server/core"
	"PORTal/server/handlers/api"
	"PORTal/server/handlers/rendered"
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

var domain string

// SetDomain should be called once on the initial setup of the server. Subsequent calls after it has been set will have no effect.
func SetDomain(d string) {
	if domain == "" {
		domain = d
	}
}

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
	mqStore stores.MemberQualificationStore,
) Server {
	l := logger.With(slog.String("source", "Server"))
	l.LogAttrs(ctx, slog.LevelInfo, fmt.Sprintf("Setting domain to: %s", domain))
	SetDomain(domain)
	logger.LogAttrs(context.Background(), slog.LevelInfo, "Creating new server")
	s := Server{
		logger: l,
		dev:    dev,
		port:   port,
	}
	c := core.New(
		logger.With(slog.String("source", "core")),
		memberStore,
		qualificationStore,
		mqStore,
		sessionStore,
		types.GetRanks(service),
		organization,
		service,
		domain,
	)

	mux := http.NewServeMux()
	l.LogAttrs(ctx, slog.LevelInfo, "Registering static asset routes...")
	// Static assets
	assetHandler := http.FileServerFS(assetsDir)
	rendered.RegisterRoutes(ctx, mux, logger, c)
	api.RegisterRoutes(logger, mux, c)

	// Handlers are applied last to first
	l.LogAttrs(ctx, slog.LevelInfo, "Applying global middlewares...")
	s.handler = adminRequiredMiddleware(mux, logger)
	s.handler = sessionRequiredMiddleware(s.handler, logger, c)
	s.handler = skipLoginMiddleware(s.handler, logger, c)
	s.handler = assetMiddleware(s.handler, assetHandler)

	l.LogAttrs(ctx, slog.LevelInfo, "Successfully registered routes")
	l.LogAttrs(ctx, slog.LevelInfo, "Setting up root serving func")
	return s
}

func (s Server) ListenAndServe() error {
	s.logger.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf("Starting server on port: %s", s.port))
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.handler)
}
