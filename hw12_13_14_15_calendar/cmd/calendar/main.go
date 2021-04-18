package main

import (
	"flag"
	"log"
	"os"

	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/app"
	"github.com/yodaskilledme/homew_otus/hw12_13_14_15_calendar/internal/config"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()
	if configFile == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	conf, err := config.ParseConfig(configFile)
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	calendar := app.New(*conf)
	err = calendar.Run()
	if err != nil {
		log.Fatalf("can't run app: %v", err)
	}
}
