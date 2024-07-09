package main

import (
	"PORTal/api"
	"PORTal/backends/sqlite"
	"context"
	"flag"
	"golang.org/x/crypto/bcrypt"
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
	if *dev {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Setting bcrypt cost to minimum for dev mode")
		sqlite.BcryptCost = bcrypt.MinCost
		l.LogAttrs(context.Background(), slog.LevelInfo, "Listening on localhost:8080 for dev mode")
		log.Fatal(http.ListenAndServe("localhost:8080", s))
	} else {
		l.LogAttrs(context.Background(), slog.LevelInfo, "Listening on :8080")
		log.Fatal(http.ListenAndServe(":8080", s))
	}
}
