package main

import (
	"PORTal/app"
	"flag"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	dev := flag.Bool("dev", false, "development mode")
	configPath := flag.String("config", "config.yml", "Path to the yaml config file")
	flag.Parse()
	f, err := os.Open(*configPath)
	if err != nil {
		panic(err)
	}
	var config app.Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	a := app.New(config, *dev, os.Stdout)
	a.Run()
}
