package config

import (
	"github.com/spf13/viper"
)

var DefaultMemoryCacherConfig = &MemoryCacherConfig{
	MaxEntriesInWindow: 1000 * 10 * 60,
	MaxEntrySizeBytes:  500,
	Shards:             1024,
	SizeMB:             64,
	EnableTTL:          false,
	TTLSecs:            1209600, // two weeks
	Verbose:            false,
}

type MemoryCacherConfig struct {
	ConfigBase

	MaxEntriesInWindow int
	MaxEntrySizeBytes  int
	Shards             int
	SizeMB             int
	EnableTTL          bool
	TTLSecs            int64
	Verbose            bool
}

func (m *MemoryCacherConfig) ConfigureViper() {
	m.setConfig("memory_cacher.max_entries_in_window",
		DefaultMemoryCacherConfig.MaxEntriesInWindow)
	m.setConfig("memory_cacher.max_entry_size_bytes",
		DefaultMemoryCacherConfig.MaxEntrySizeBytes)
	m.setConfig("memory_cacher.shards",
		DefaultMemoryCacherConfig.Shards)
	m.setConfig("memory_cacher.size_mb",
		DefaultMemoryCacherConfig.SizeMB)
	m.setConfig("memory_cacher.enable_ttl",
		DefaultMemoryCacherConfig.EnableTTL)
	m.setConfig("memory_cacher.ttl_secs",
		DefaultMemoryCacherConfig.TTLSecs)
	m.setConfig("memory_cacher.verbose",
		DefaultMemoryCacherConfig.Verbose)
}

func (m *MemoryCacherConfig) ReadConfig() {
	m.MaxEntriesInWindow = viper.GetInt("memory_cacher.max_entries_in_window")
	m.MaxEntrySizeBytes = viper.GetInt("memory_cacher.max_entry_size_bytes")
	m.Shards = viper.GetInt("memory_cacher.shards")
	m.SizeMB = viper.GetInt("memory_cacher.size_mb")
	m.EnableTTL = viper.GetBool("memory_cacher.enable_ttl")
	m.TTLSecs = viper.GetInt64("memory_cacher.ttl_secs")
	m.Verbose = viper.GetBool("memory_cacher.verbose")
}
