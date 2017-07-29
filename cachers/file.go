package cachers

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ssalevan/photocache/config"
)

func NewFileCacher(conf *config.FileCacherConfig) *FileCacher {
	return &FileCacher{
		conf: conf,
	}
}

type FileCacher struct {
	conf *config.FileCacherConfig
}

func (f *FileCacher) Init() error {
	return nil
}

func (f *FileCacher) Get(url string) (bool, []byte, error) {
	fileLocation := path.Join(f.conf.StorageDir, url)

	file, err := os.Stat(fileLocation)
	if err == nil {
		// If this file has exceeded its TTL, deletes it and forces a refetch.
		cutoffTime := time.Now().UTC().Add(time.Duration(-1*f.conf.TTLSecs) * time.Second)
		if f.conf.EnableTTL && file.ModTime().Before(cutoffTime) {
			err = os.Remove(fileLocation)
			if err != nil {
				log.Errorf("Unable to remove file %s: %v", fileLocation, err)
			}
			return false, []byte{}, nil
		}
	} else if os.IsNotExist(err) {
		return false, []byte{}, nil
	} else {
		return false, []byte{}, err
	}

	data, err := ioutil.ReadFile(fileLocation)
	return true, data, err
}

func (f *FileCacher) Name() string {
	return "file"
}

func (f *FileCacher) Set(url string, contents []byte) error {
	fileLocation := path.Join(f.conf.StorageDir, url)

	// If the file exists, updates atime/mtime for TTL purposes.
	_, err := os.Stat(fileLocation)
	if err == nil {
		currentTime := time.Now().UTC()
		return os.Chtimes(fileLocation, currentTime, currentTime)
	}

	fileDir := filepath.Dir(fileLocation)
	err = os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileLocation, contents, os.FileMode(f.conf.Mode))
	return err
}
