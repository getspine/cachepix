package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/ssalevan/photocache/app"
	"github.com/ssalevan/photocache/config"
)

func main() {
	log.Debug("Starting Photocache application...")

	// Loads the Photocache configuration.
	conf := config.NewPhotocacheConfig()

	// Launches the Photocache application.
	app.Launch(conf)
}
