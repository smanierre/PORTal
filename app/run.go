package app

import (
	"PORTal/backend"
	"PORTal/providers/sqlite"
	"PORTal/server"
	"PORTal/templates"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type Config struct {
	Backend backend.Config `yaml:"backend"`
	Server  server.Config  `yaml:"server"`
}

func (c Config) Merge(new Config) Config {
	if new.Backend.DbFile != "" {
		c.Backend.DbFile = new.Backend.DbFile
	}
	if new.Backend.BcryptCost != 0 {
		c.Backend.BcryptCost = new.Backend.BcryptCost
	}
	// Domain must be provided
	if new.Server.Domain == "" {
		log.Fatal("Domain must be provided in config file")
	}
	c.Server.Domain = new.Server.Domain
	if new.Server.Port != 0 {
		c.Server.Port = new.Server.Port
	}
	// Organization must be provided
	if new.Server.Organization == "" {
		log.Fatal("Organization must be provided in config file")
	}
	c.Server.Organization = new.Server.Organization
	if new.Backend.SessionTimeout != 0 {
		c.Backend.SessionTimeout = new.Backend.SessionTimeout
	}
	if new.Server.Service != "" {
		c.Server.Service = new.Server.Service
	}
	return c
}

var DefaultConfig Config = Config{
	Backend: backend.Config{
		DbFile:         "PORTal.db",
		BcryptCost:     16,
		SessionTimeout: 4,
	},
	Server: server.Config{
		Domain:  "",
		Port:    8080,
		Service: "f",
	},
}

func New(config Config, dev bool, logDest io.Writer) App {
	config = DefaultConfig.Merge(config)
	var l *slog.Logger
	if dev {
		l = slog.New(slog.NewTextHandler(logDest, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	} else {
		l = slog.New(slog.NewJSONHandler(logDest, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	}
	provider, err := sqlite.New(l.With(slog.String("service", "sqlite_provider")), config.Backend.DbFile, 1)
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error creating provider", slog.String("error", err.Error()))
	}

	templates.Initialize(config.Server.Organization, config.Server.Service)

	b := backend.New(
		l.With(slog.String("service", "backend")),
		provider,
		provider,
		provider,
		provider,
		config.Backend,
		nil,
	)
	a := App{
		server: server.New(l.With(slog.String("service", "server")), b, dev, config.Server),
		config: config,
	}
	return a
}

type App struct {
	server server.Server
	config Config
}

func (a App) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.config.Server.Port), a.server))
}
