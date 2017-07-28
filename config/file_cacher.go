package config

type FileCacherConfig struct {
	StorageDir string `yaml:"storage_dir,omitempty"`
	EnableTTL  bool   `yaml:"enable_ttl,omitempty"`
	TTLSecs    int64  `yaml:"ttl_secs,omitempty"`
}

var DefaultFileCacherConfig = &FileCacherConfig{
	StorageDir: "/var/lib/photocache",
	EnableTTL:  false,
	TTLSecs:    1209600, // Two weeks
}
