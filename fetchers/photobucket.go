package fetchers

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/ssalevan/photocache/config"
)

var URLEmbeddedRegex = regexp.MustCompile(`(.*.photobucket.com/albums/.*/.*/)`)

func NewPhotobucketFetcher(conf *config.PhotobucketFetcherConfig) *PhotobucketFetcher {
	return &PhotobucketFetcher{
		conf: conf,
	}
}

type PhotobucketFetcher struct {
	conf *config.PhotobucketFetcherConfig
}

func (f *PhotobucketFetcher) MatchesURL(url string) bool {
	return URLEmbeddedRegex.MatchString(url)
}

func (f *PhotobucketFetcher) Get(url string) ([]byte, error) {
	if !URLEmbeddedRegex.matches(url) || !strings.HasPrefix(url, f.conf.Prefix) {
		return []byte{}, errors.New("Invalid URL")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s", url), nil)
	if err != nil {
		return []byte{}, err
	}

	albumURL := URLEmbeddedRegex.FindString(url)
	req.Header.Set("Referer", albumURL)

	response, err := http.Client.Do(req)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(response.Body)
}
