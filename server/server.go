package server

import (
	"PORTal/server/stores"
	"PORTal/types"
	"context"
	"embed"
	"fmt"
	"log/slog"
	"net/http"
	"time"
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

	mux := http.NewServeMux()
	l.LogAttrs(ctx, slog.LevelInfo, "Registering static asset routes...")
	// Static assets
	assetHandler := http.FileServerFS(assetsDir)
	ranks := types.GetRanks(service)
	addRoutes(ctx, mux, logger, assetHandler, organization, memberStore, sessionStore, qualificationStore, ranks)

	// Handlers are applied last to first
	l.LogAttrs(ctx, slog.LevelInfo, "Applying global middlewares...")
	s.handler = adminRequiredMiddleware(mux, logger)
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

// Util functions for the server package

func CheckHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") != ""
}

func GetSessionId(r *http.Request) (string, error) {
	c, err := r.Cookie(SessionCookieName)
	if err != nil {
		return "", err
	}
	return c.Value, nil
}

func HandleRenderError(ctx context.Context, logger *slog.Logger, err error) {
	if err != nil {
		logger.LogAttrs(ctx, slog.LevelError, "Error rendering template", slog.String("error", err.Error()))
	}
}

func MakeCookie(name, value string, expiration time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		Expires:  expiration,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func RemoveCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:    name,
		Domain:  domain,
		Expires: time.Now(),
		Path:    "/",
	})
}
