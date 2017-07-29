package config

import (
	"github.com/spf13/viper"
)

var DefaultPhotobucketFetcherConfig = &PhotobucketFetcherConfig{
	Prefix: "",
}

type PhotobucketFetcherConfig struct {
	ConfigBase

	Prefix string
}

func (p *PhotobucketFetcherConfig) ConfigureViper() {
	p.setConfig("photobucket_fetcher.prefix",
		DefaultPhotobucketFetcherConfig.Prefix)
}

func (p *PhotobucketFetcherConfig) ReadConfig() {
	p.Prefix = viper.GetString("photobucket_fetcher.prefix")
}
