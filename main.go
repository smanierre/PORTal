package main

import (
	"PORTal/app"
	"flag"
)

func main() {
	dev := flag.Bool("dev", false, "development mode")
	flag.Parse()
	opts := &app.Options{
		Dev: *dev,
	}
	a := app.New(opts)
	a.Run()
}
