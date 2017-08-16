package config

import (
	"github.com/spf13/viper"
)

type FileCacherConfig struct {
	ConfigBase

	StorageDir string `yaml:"storage_dir,omitempty"`
	EnableTTL  bool   `yaml:"enable_ttl,omitempty"`
	TTLSecs    int64  `yaml:"ttl_secs,omitempty"`
	Mode       int    `yaml:"mode,omitempty"`
}

var DefaultFileCacherConfig = &FileCacherConfig{
	StorageDir: "/var/lib/cachepix",
	EnableTTL:  false,
	TTLSecs:    1209600, // Two weeks
	Mode:       0644,
}

func (f *FileCacherConfig) ConfigureViper() {
	f.setConfig("file_cacher.storage_dir", DefaultFileCacherConfig.StorageDir)
	f.setConfig("file_cacher.enable_ttl", DefaultFileCacherConfig.EnableTTL)
	f.setConfig("file_cacher.ttl_secs", DefaultFileCacherConfig.TTLSecs)
	f.setConfig("file_cacher.mode", DefaultFileCacherConfig.Mode)
}

func (f *FileCacherConfig) ReadConfig() {
	f.StorageDir = viper.GetString("file_cacher.storage_dir")
	f.EnableTTL = viper.GetBool("file_cacher.enable_ttl")
	f.TTLSecs = viper.GetInt64("file_cacher.ttl_secs")
	f.Mode = viper.GetInt("file_cacher.mode")
}
