package app

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ssalevan/photocache/common"
	"github.com/ssalevan/photocache/config"
	cache "github.com/ssalevan/photocache/photocache"
)

const (
	photocacheKillWait time.Duration = 15 * time.Second
)

// Launch : starts a new Photocache process and loads a config
func Launch(cfg *config.PhotocacheConfig) {
	rand.Seed(time.Now().Unix())

	// Starts the safe process.
	photocache := cache.NewPhotocache(cfg)
	common.StartProcess(photocache)

	// Waits to receive an interruption signal.
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	select {
	case signal := <-signalChan:
		log.Infof("Stop requested (%v); shutting down...", signal)
	}

	// Interrupt received, kills Photocache.
	go photocache.Stop()
	select {
	case <-time.After(photocacheKillWait):
		log.Warningf("Forcibly dying after waiting for Photocache to stop...")
	case <-photocache.Stopped:
	}

	log.Debugf("Photocache shut down; exiting with a status code of 0")
	os.Exit(0)
}
