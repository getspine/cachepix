package fetchers

type Fetcher interface {
	MatchesURL(url string) bool
	Get(url string) ([]byte, error)
	Name() string
}
