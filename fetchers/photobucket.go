package fetchers

import (
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
		conf:       conf,
		httpClient: &http.Client{},
	}
}

type PhotobucketFetcher struct {
	conf       *config.PhotobucketFetcherConfig
	httpClient *http.Client
}

func (f *PhotobucketFetcher) Init() error {
	return nil
}

func (f *PhotobucketFetcher) MatchesURL(url string) bool {
	return URLEmbeddedRegex.MatchString(url) && strings.HasPrefix(url, f.conf.Prefix)
}

func (f *PhotobucketFetcher) Get(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s", url), nil)
	if err != nil {
		return []byte{}, err
	}

	albumURL := strings.Replace(URLEmbeddedRegex.FindString(url), "i", "s", 1)
	req.Header.Set("Referer", fmt.Sprintf("http://%s", albumURL))

	response, err := f.httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	return ioutil.ReadAll(response.Body)
}

func (f *PhotobucketFetcher) Name() string {
	return "photobucket"
}
