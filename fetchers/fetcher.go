package fetchers

import (
	"net/http"
)

type Fetcher interface {
	MatchesURL(url string) bool
	Get(url string) ([]byte, error)
}
