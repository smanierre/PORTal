package main

import (
	"PORTal/api"
	"PORTal/backends/sqlite"
	"context"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	dev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	b, err := sqlite.New(l, "PORTal.db", 1.0)
	if err != nil {
		l.LogAttrs(context.Background(), slog.LevelError, "Error creating new sqlite backend", slog.String("error", err.Error()))
		os.Exit(1)
	}
	s := api.New(l, b, *dev)
	l.LogAttrs(context.Background(), slog.LevelInfo, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", s))
}
