package main

import (
	"PORTal/app"
	"context"
	"fmt"
	"os"
	"os/signal"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()
	a, err := app.New(ctx, os.Stdout, os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating application: %s\n", err)
		os.Exit(1)
	}
	if err := a.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
