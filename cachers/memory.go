package cachers

import (
	"time"

	"github.com/allegro/bigcache"

	"github.com/ssalevan/cachepix/config"
)

var Eternity = time.Duration(1<<63 - 1)

func NewMemoryCacher(conf *config.MemoryCacherConfig) *MemoryCacher {
	bigcacheConfig := bigcache.Config{
		// number of shards (must be a power of 2)
		Shards: conf.Shards,
		// time after which entry can be evicted
		LifeWindow: Eternity,
		// rps * lifeWindow, used only in initial memory allocation
		MaxEntriesInWindow: conf.MaxEntriesInWindow,
		// max entry size in bytes, used only in initial memory allocation
		MaxEntrySize: conf.MaxEntrySizeBytes,
		// prints information about additional memory allocation
		Verbose: conf.Verbose,
		// cache will not allocate more memory than this limit, value in MB
		// if value is reached then the oldest entries can be overridden for the new ones
		// 0 value means no size limit
		HardMaxCacheSize: conf.SizeMB,
		// callback fired when the oldest entry is removed because of its
		// expiration time or no space left for the new entry. Default value is nil which
		// means no callback and it prevents from unwrapping the oldest entry.
		OnRemove: nil,
	}

	if conf.EnableTTL {
		bigcacheConfig.LifeWindow = time.Duration(conf.TTLSecs) * time.Second
	}

	return &MemoryCacher{
		conf:         conf,
		bigcacheConf: bigcacheConfig,
	}
}

type MemoryCacher struct {
	conf         *config.MemoryCacherConfig
	bigcacheConf bigcache.Config

	bigcache *bigcache.BigCache
}

func (m *MemoryCacher) Init() error {
	var err error
	if m.bigcache == nil {
		m.bigcache, err = bigcache.NewBigCache(m.bigcacheConf)
	}
	return err
}

func (m *MemoryCacher) Get(url string) (bool, []byte, error) {
	contents, err := m.bigcache.Get(url)
	if _, notFound := err.(*bigcache.EntryNotFoundError); notFound {
		return false, contents, nil
	} else if err != nil {
		return false, contents, err
	}
	return true, contents, nil
}

func (m *MemoryCacher) Name() string {
	return "memory"
}

func (m *MemoryCacher) Set(url string, contents []byte) error {
	return m.bigcache.Set(url, contents)
}
