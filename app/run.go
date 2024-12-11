package app

import (
	"PORTal/api"
	"PORTal/backend"
	"PORTal/providers/sqlite"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
)

type Config struct {
	Backend backend.Config `yaml:"backend"`
	Api     api.Config     `yaml:"api"`
}

func (c Config) Merge(new Config) Config {
	if new.Backend.DbFile != "" {
		c.Backend.DbFile = new.Backend.DbFile
	}
	if new.Backend.BcryptCost != 0 {
		c.Backend.BcryptCost = new.Backend.BcryptCost
	}
	// Domain must be provided
	if new.Api.Domain == "" {
		log.Fatal("Domain must be provided in config file")
	}
	c.Api.Domain = new.Api.Domain
	if new.Api.Port != 0 {
		c.Api.Port = new.Api.Port
	}
	// Organization must be provided
	if new.Api.Organization == "" {
		log.Fatal("Organization must be provided in config file")
	}
	c.Api.Organization = new.Api.Organization
	return c
}

var DefaultConfig Config = Config{
	Backend: backend.Config{
		DbFile:     "PORTal.db",
		BcryptCost: 16,
	},
	Api: api.Config{
		Domain: "",
		Port:   8080,
	},
}

func New(config Config, dev bool, logDest io.Writer) App {
	config = DefaultConfig.Merge(config)
	l := slog.New(slog.NewTextHandler(logDest, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	provider, err := sqlite.New(l.With(slog.String("service", "sqlite_provider")), config.Backend.DbFile, 1)
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error creating provider", slog.String("error", err.Error()))
	}

	b := backend.New(
		l.With(slog.String("service", "backend")),
		provider,
		provider,
		provider,
		config.Backend,
		nil,
	)
	a := App{
		server: api.New(l.With(slog.String("service", "api_server")), b, dev, config.Api),
		config: config,
	}
	return a
}

type App struct {
	server api.Server
	config Config
}

func (a App) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.config.Api.Port), a.server))
}
