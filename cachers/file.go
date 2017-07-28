package cachers

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"time"

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

func (f *FileCacher) Get(url string) ([]byte, error) {
	fileLocation := path.Join(f.conf.StorageDir, url)
	_, err := os.Stat(fileLocation)
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadFile(fileLocation)
}

func (f *FileCacher) Hit(url string) bool {
	fileLocation := path.Join(f.conf.StorageDir, url)

	file, err := os.Stat(fileLocation)
	if err == nil {
		cutoffTime := time.Now().UTC().Add(time.Duration(-1*f.conf.TTLSecs) * time.Second)
		if f.conf.EnableTTL && file.ModTime().Before(cutoffTime) {
			return false
		}
		return true
	}
	return false
}

func (f *FileCacher) Name() string {
	return "file"
}

func (f *FileCacher) Set(url string, contents []byte) error {
	fileLocation := path.Join(f.conf.StorageDir, url)
	fileDir := filepath.Dir(fileLocation)
	err := os.MkdirAll(fileDir, os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileLocation, contents, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}
