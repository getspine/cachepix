package cachers

type Cacher interface {
	Get(url string) ([]byte, error)
	Hit(url string) bool
	Name() string
	Set(url string, contents []byte) error
}
