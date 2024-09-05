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
	"os"
)

type Options struct {
	Port    int
	Dev     bool
	DbFile  string
	LogDest io.Writer
}

var DefaultOptions *Options = &Options{
	Port:    8080,
	Dev:     false,
	DbFile:  "PORTal.db",
	LogDest: os.Stdout,
}

func New(opts *Options) App {
	opts = DefaultOptions.Merge(opts)
	l := slog.New(slog.NewTextHandler(opts.LogDest, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	provider, err := sqlite.New(l.With(slog.String("service", "sqlite_provider")), opts.DbFile, 1)
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error creating provider", slog.String("error", err.Error()))
	}

	b := backend.New(
		l.With(slog.String("service", "backend")),
		provider,
		provider,
		provider,
		provider,
		nil,
	)
	l.LogAttrs(context.Background(), slog.LevelInfo, "Creating app...", slog.Any("options", opts))
	a := App{
		server:  api.New(l.With(slog.String("service", "api_server")), b, opts.Dev),
		options: opts,
	}
	return a
}

type App struct {
	server  api.Server
	options *Options
}

func (a App) Run() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", a.options.Port), a.server))
}

func (o Options) Merge(incoming *Options) *Options {

	if incoming.Dev != o.Dev {
		o.Dev = true
	}
	if incoming.LogDest != nil && incoming.LogDest != o.LogDest {
		o.LogDest = incoming.LogDest
	}
	if incoming.Port != 0 && incoming.Port != o.Port {
		o.Port = incoming.Port
	}
	if incoming.DbFile != "" && incoming.DbFile != o.DbFile {
		o.DbFile = incoming.DbFile
	}
	return &o
}

func (o Options) LogValue() slog.Value {
	return slog.StringValue(fmt.Sprintf("Dev: %t, Port: %d, DbFile: %s", o.Dev, o.Port, o.DbFile))
}
