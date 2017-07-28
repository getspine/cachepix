package cache

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/ssalevan/photocache/cachers"
	"github.com/ssalevan/photocache/common"
	"github.com/ssalevan/photocache/config"
	"github.com/ssalevan/photocache/fetchers"
)

// Photocache : Interface that ties process and config
type Photocache struct {
	common.BackgroundProcess

	config *config.Config

	cachers  []cachers.Cacher
	fetchers []fetchers.Fetcher
}

// NewPhotocache : crates a new Photocache process with a config
func NewPhotocache(config *config.Config) *Photocache {
	photocache := &Photocache{
		config: config,
	}
	photocache.InitProcess("Photocache")
	photocache.Init()
	return photocache
}

func (p *Photocache) Init() {
	for _, fetcherName := range p.config.Fetchers {
		if fetcherName == "photobucket" {
			p.fetchers = append(p.fetchers,
				fetchers.NewPhotobucketFetcher(p.config.PhotobucketFetcher))
		}
	}

	for _, cacherName := range p.config.Cachers {
		if cacherName == "file" {
			p.cachers = append(p.cachers, cachers.NewFileCacher(p.config.FileCacher))
		}
	}
}

func (r *Photocache) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	mediaURL := req.URL.Path[1:]

	var data []byte
	var err error

	// Detects whether there is a provided cacher for the given thing.
	cacheHit := false
	for _, cacher := range r.cachers {
		if cacher.Hit(mediaURL) {
			data, err = cacher.Get(mediaURL)
			if err != nil {
				log.Errorf("Cacher %s returned error: %v", cacher.Name(), err)
				err = nil
			} else {
				cacheHit = true
				break
			}
		}
	}

	// If nothing found in the local cache, reaches out to the photo sharing service to
	// retrieve the provided image.
	if !cacheHit {
		foundFetcher := false
		for _, fetcher := range r.fetchers {
			if fetcher.MatchesURL(mediaURL) {
				foundFetcher = true
				data, err = fetcher.Get(mediaURL)
				if err != nil {
					log.Errorf("Caught error during fetch of %s: %v", mediaURL, err)
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Temporarily unable to retrieve image, try again later"))
					return
				}
			}
		}
		if !foundFetcher {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Cannot retrieve provided URL"))
			return
		}
	}

	// Caches response in all configured caching mechanisms.
	for _, cacher := range r.cachers {
		err = cacher.Set(mediaURL, data)
		if err != nil {
			log.Errorf("Could not set cache for %s on %s cacher: %v",
				mediaURL, cacher.Name(), err)
		}
	}

	// Sends the image's data back to the requester.
	fileExtension := filepath.Ext(req.URL.Path)
	w.Header().Set("Content-Type", mime.TypeByExtension(fileExtension))
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	if _, err := w.Write(data); err != nil {
		log.Errorf("Unable to write data for HTTP request %s: %v", mediaURL, err)
	}
}

// Run : starts the actual routing
func (r *Photocache) Run() {
	r.Wg.Add(1)
	defer r.Wg.Done()

	healthcheckServer := common.NewHealthcheckServer()

	// Healthcheck servers are started on both HTTP and HTTPS to satisfy ELB constraints.
	go func() {
		log.Infof("Starting Photocache healthcheck server on port: %d",
			r.config.HealthcheckPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", r.config.HealthcheckPort),
			healthcheckServer)
		if err != nil {
			log.Fatalf("Could not start healthcheck server: %v", err)
		}
	}()

	go func() {
		log.Infof("Starting Photocache HTTP server on port: %d", r.config.HTTPListenPort)
		err := http.ListenAndServe(fmt.Sprintf(":%d", r.config.HTTPListenPort), r)
		if err != nil {
			log.Fatalf("Could not start HTTP server: %v", err)
		}
	}()

	if r.config.EnableHTTPS {
		go func() {
			log.Infof("Starting Photocache healthcheck TLS server on port: %d",
				r.config.HealthcheckTLSPort)
			err := http.ListenAndServeTLS(fmt.Sprintf(":%d", r.config.HealthcheckTLSPort),
				r.config.SSLCert, r.config.SSLKey, healthcheckServer)
			if err != nil {
				log.Fatalf("Could not start healthcheck TLS server: %v", err)
			}
		}()

		go func() {
			log.Infof("Starting Photocache HTTPS server on port: %d", r.config.HTTPSListenPort)
			err := http.ListenAndServeTLS(fmt.Sprintf(":%d", r.config.HTTPSListenPort),
				r.config.SSLCert, r.config.SSLKey, r)
			if err != nil {
				log.Fatalf("Could not start HTTPS server: %v", err)
			}
		}()
	}

	// Waits for the done signal.
	select {
	case <-r.Done:
	}
}
