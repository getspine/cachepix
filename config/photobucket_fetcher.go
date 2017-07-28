package config

type PhotobucketFetcherConfig struct {
	Prefix string `yaml:"prefix,omitempty"`
}

var DefaultPhotobucketFetcherConfig = &PhotobucketFetcherConfig{
	Prefix: "",
}
