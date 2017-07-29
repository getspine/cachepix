package cachers

type Cacher interface {
	Init() error
	Get(url string) (bool, []byte, error)
	Name() string
	Set(url string, contents []byte) error
}
