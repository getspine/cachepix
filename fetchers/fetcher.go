package fetchers

type Fetcher interface {
	Init() error
	MatchesURL(url string) bool
	Get(url string) ([]byte, error)
	Name() string
}
