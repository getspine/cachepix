package config

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var DefaultPhotocacheConfig = &PhotocacheConfig{
	Loglevel: "info",

	EnableHTTPS:        false,
	HTTPListenPort:     12345,
	HTTPSListenPort:    12346,
	HealthcheckPort:    20026,
	HealthcheckTLSPort: 20027,
	SSLCert:            "/etc/photocache/photocache.crt",
	SSLKey:             "/etc/photocache/photocache.key",

	Cachers:  []string{"memory", "file"},
	Fetchers: []string{"photobucket"},

	FileCacher:   DefaultFileCacherConfig,
	MemoryCacher: DefaultMemoryCacherConfig,
	S3Cacher:     DefaultS3CacherConfig,

	PhotobucketFetcher: DefaultPhotobucketFetcherConfig,
}

func NewPhotocacheConfig() *PhotocacheConfig {
	newConfig := &PhotocacheConfig{
		FileCacher:   DefaultFileCacherConfig,
		MemoryCacher: DefaultMemoryCacherConfig,
		S3Cacher:     DefaultS3CacherConfig,

		PhotobucketFetcher: DefaultPhotobucketFetcherConfig,
	}
	newConfig.ReadConfig()

	configString, _ := json.MarshalIndent(newConfig, "", "  ")
	log.Debugf("Photocache config: %s", configString)
	return newConfig
}

type PhotocacheConfig struct {
	ConfigBase

	Loglevel string

	EnableHTTPS bool

	HTTPListenPort  int64
	HTTPSListenPort int64

	HealthcheckPort    int64
	HealthcheckTLSPort int64

	SSLCert string
	SSLKey  string

	Cachers  []string
	Fetchers []string

	FileCacher   *FileCacherConfig
	MemoryCacher *MemoryCacherConfig
	S3Cacher     *S3CacherConfig

	PhotobucketFetcher *PhotobucketFetcherConfig
}

func (p *PhotocacheConfig) ConfigureViper() {
	p.setConfig("loglevel", DefaultPhotocacheConfig.Loglevel)
	p.setConfig("enable_https", DefaultPhotocacheConfig.EnableHTTPS)
	p.setConfig("http_listen_port", DefaultPhotocacheConfig.HTTPListenPort)
	p.setConfig("https_listen_port", DefaultPhotocacheConfig.HTTPSListenPort)
	p.setConfig("healthcheck_port", DefaultPhotocacheConfig.HealthcheckPort)
	p.setConfig("healthcheck_tls_port", DefaultPhotocacheConfig.HealthcheckTLSPort)
	p.setConfig("ssl_cert", DefaultPhotocacheConfig.SSLCert)
	p.setConfig("ssl_key", DefaultPhotocacheConfig.SSLKey)
	p.setConfig("cachers", strings.Join(DefaultPhotocacheConfig.Cachers, ","))
	p.setConfig("fetchers", strings.Join(DefaultPhotocacheConfig.Fetchers, ","))

	DefaultFileCacherConfig.ConfigureViper()
	DefaultMemoryCacherConfig.ConfigureViper()
	DefaultS3CacherConfig.ConfigureViper()

	DefaultPhotobucketFetcherConfig.ConfigureViper()
}

func (p *PhotocacheConfig) ReadConfig() {
	p.Loglevel = viper.GetString("loglevel")
	p.EnableHTTPS = viper.GetBool("enable_https")
	p.HTTPListenPort = viper.GetInt64("http_listen_port")
	p.HTTPSListenPort = viper.GetInt64("https_listen_port")
	p.HealthcheckPort = viper.GetInt64("healthcheck_port")
	p.HealthcheckTLSPort = viper.GetInt64("healthcheck_tls_port")
	p.SSLCert = viper.GetString("ssl_cert")
	p.SSLKey = viper.GetString("ssl_key")
	p.Cachers = strings.Split(viper.GetString("cachers"), ",")
	p.Fetchers = strings.Split(viper.GetString("fetchers"), ",")

	p.FileCacher.ReadConfig()
	p.MemoryCacher.ReadConfig()
	p.S3Cacher.ReadConfig()

	p.PhotobucketFetcher.ReadConfig()
}
