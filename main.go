package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/ssalevan/cachepix/app"
	"github.com/ssalevan/cachepix/config"
)

func main() {
	log.Debug("Starting Cachepix application...")

	// Loads the Cachepix configuration.
	conf := config.NewCachepixConfig()

	// Launches the Cachepix application.
	app.Launch(conf)
}
