package app

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	cache "github.com/ssalevan/cachepix/cachepix"
	"github.com/ssalevan/cachepix/common"
	"github.com/ssalevan/cachepix/config"
)

const (
	cachepixKillWait time.Duration = 15 * time.Second
)

// Launch : starts a new Cachepix process and loads a config
func Launch(cfg *config.CachepixConfig) {
	rand.Seed(time.Now().Unix())

	// Starts the safe process.
	cachepix := cache.NewCachepix(cfg)
	common.StartProcess(cachepix)

	// Waits to receive an interruption signal.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case signal := <-signalChan:
		log.Infof("Stop requested (%v); shutting down...", signal)
	}

	// Interrupt received, kills Cachepix.
	go cachepix.Stop()
	select {
	case <-time.After(cachepixKillWait):
		log.Warningf("Forcibly dying after waiting for Cachepix to stop...")
	case <-cachepix.Stopped:
	}

	log.Debugf("Cachepix shut down; exiting with a status code of 0")
	os.Exit(0)
}
