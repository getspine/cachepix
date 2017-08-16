package config

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type Config interface {
	ConfigureViper()
	ReadConfig()
}

type ConfigBase struct{}

func (c *ConfigBase) setConfig(key string, value interface{}) {
	viper.SetDefault(key, value)
	viper.BindEnv(key, "PCACHE_"+strings.ToUpper(strings.Replace(key, ".", "_", -1)))
}

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/cachepix/")
	viper.AddConfigPath("$HOME/.cachepix")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix("pcache")

	DefaultCachepixConfig.ConfigureViper()

	err := viper.ReadInConfig()
	if err != nil {
		log.Infof("No config file found, using defaults and environment overrides:%v",
			err)
	}

	log.SetLevel(log.InfoLevel)

	configuredLevel, err := log.ParseLevel(viper.GetString("loglevel"))
	if err != nil {
		log.Errorf("Could not parse loglevel: %s", viper.GetString("loglevel"))
	} else {
		// Sets master logrus loglevel from configured value.
		log.SetLevel(configuredLevel)
	}
}
