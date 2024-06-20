package main

import (
	"PORTal/api"
	"PORTal/backends/sqlite"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	b, err := sqlite.New(l, "PORTal.db", 1.0)
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error creating new sqlite backend", slog.String("error", err.Error()))
		os.Exit(1)
	}
	s := api.New(l, b)
	l.LogAttrs(context.Background(), slog.LevelInfo, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}
