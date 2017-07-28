package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/ssalevan/photocache/app"
	"github.com/ssalevan/photocache/config"
)

// Contain command-line arguments.
var configFilePath string
var logPath string
var debug bool

func init() {
	flag.StringVar(&configFilePath, "config", "/etc/photocache/config.yml",
		"Full path of the configuration JSON file")
}

func main() {
	flag.Parse()

	// Attempts to load the Photocache configuration.
	conf, err := config.FromFile(configFilePath)
	if err != nil {
		log.Errorf("Could not load configuration: %v", err)
		os.Exit(-1)
	}

	// Launches the Photocache application.
	app.Launch(conf)
}
