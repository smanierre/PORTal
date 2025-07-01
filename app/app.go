package app

import (
	"PORTal/backend/memberqualificationstore"
	"PORTal/backend/memberstore"
	"PORTal/backend/qualificationstore"
	"PORTal/backend/sessionstore"
	"PORTal/providers/sqlite/memberprovider"
	"PORTal/providers/sqlite/qualificationprovider"
	"PORTal/providers/sqlite/sessionprovider"
	"PORTal/server"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

func New(ctx context.Context, w io.Writer, args []string) (*App, error) {
	// Flag parsing and validation
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	dbFile := flags.String("db", "./PORTal.db", "Path to database file")
	hashCost := flags.Int("hashCost", 16, "Hash cost for passwords")
	sessionTimeout := flags.Duration("sessionTimeout", 8*time.Hour, "Session timeout")
	domain := flags.String("domain", "", "Domain name to be used for cookies")
	port := flags.String("port", "8080", "Port to listen on")
	organization := flags.String("organization", "", "Organization name to be used for login page")
	service := flags.String("service", "f", "Service identifier for ranks. f=Air Force, a=Army, m=Marines, n=Navy")
	dev := flags.Bool("dev", false, "Development mode")
	err := flags.Parse(args[1:])
	if err != nil {
		return nil, err
	}
	if *domain == "" {
		flags.PrintDefaults()
		return nil, errors.New("the domain flag is required")
	}
	if *organization == "" {
		flags.PrintDefaults()
		return nil, errors.New("the organization flag is required")
	}

	// Setup a text logger for dev, and JSON logger for production
	var l *slog.Logger
	if *dev {
		l = slog.New(slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.LevelDebug}))
	} else {
		l = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.LevelInfo}))
	}

	// Setup providers and stores
	memberProvider, err := memberprovider.New(*dbFile, l)
	if err != nil {
		return nil, err
	}
	memberStore := memberstore.New(memberProvider, *hashCost, l)

	qualificationProvider, err := qualificationprovider.New(*dbFile, l)
	if err != nil {
		return nil, err
	}
	qualificationStore := qualificationstore.New(qualificationProvider, l)

	mqStore := memberqualificationstore.New(memberProvider, qualificationProvider, memberProvider, l, nil)

	sessionProvider, err := sessionprovider.New(*dbFile, l)
	if err != nil {
		return nil, err
	}
	sessionStore := sessionstore.New(sessionProvider, memberStore, nil, *sessionTimeout, l)

	s := server.New(ctx, l, *port, *dev, *domain, *service, *organization, memberStore, qualificationStore, sessionStore, mqStore)
	a := &App{
		server:  s,
		context: ctx,
	}
	return a, nil
}

type App struct {
	context context.Context
	server  server.Server
}

func (a App) Run() error {
	var appErr error
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			appErr = err
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-a.context.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		if err := a.shutdown(shutdownCtx); err != nil {
			appErr = err
		}
	}()
	wg.Wait()
	return appErr
}

func (a App) shutdown(ctx context.Context) error {
	fmt.Fprintln(os.Stdout, "Shutting down server...")
	shutdownChan := make(chan error)
	go func(c chan<- error) {
		c <- nil
	}(shutdownChan)
	select {
	case err := <-shutdownChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
