package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/jinzhu/copier"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	EnableHTTPS bool `yaml:"enable_https,omitempty"`

	HTTPListenPort  int64 `yaml:"http_listen_port,omitempty"`
	HTTPSListenPort int64 `yaml:"https_listen_port,omitempty"`

	FileServerPort int64 `yaml:"file_server_port,omitempty"`

	HealthcheckPort    int64 `yaml:"healthcheck_port,omitempty"`
	HealthcheckTLSPort int64 `yaml:"healthcheck_tls_port,omitempty"`

	SSLCert string `yaml:"ssl_cert,omitempty"`
	SSLKey  string `yaml:"ssl_key,omitempty"`

	Cachers  []string `yaml:"cachers,omitempty"`
	Fetchers []string `yaml:"fetchers,omitempty"`

	FileCacher         *FileCacherConfig         `yaml:"file_cacher,omitempty"`
	PhotobucketFetcher *PhotobucketFetcherConfig `yaml:"photobucket_fetcher,omitempty"`
}

var DefaultConfig = &Config{
	EnableHTTPS:        false,
	HTTPListenPort:     80,
	HTTPSListenPort:    443,
	FileServerPort:     20025,
	HealthcheckPort:    20026,
	HealthcheckTLSPort: 20027,
	SSLCert:            "/etc/photocache/photocache.crt",
	SSLKey:             "/etc/photocache/photocache.key",
	StaticDir:          "/var/lib/photocache",
	Cachers:            []string{"file"},
	Fetchers:           []string{"photobucket"},
	FileCacher:         DefaultFileCacherConfig,
	PhotobucketFetcher: DefaultPhotobucketFetcherConfig,
}

func (config *Config) FromFile(filePath string) error {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	newConf := &Config{}
	err = yaml.Unmarshal(content, newConf)
	if err != nil {
		return err
	}
	err = copier.Copy(config, newConf)
	if err != nil {
		return err
	}
	return nil
}

func FromFile(filePath string) (*Config, error) {
	conf := &Config{}
	err := copier.Copy(conf, DefaultConfig)
	if err != nil {
		return nil, err
	}

	_, err := os.Stat(filePath)
	if err == nil {
		err := conf.FromFile(filePath)
		if err != nil {
			return nil, err
		}
	}

	return conf, nil
}
