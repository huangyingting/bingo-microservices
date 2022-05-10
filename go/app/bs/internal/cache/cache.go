package cache

type ICache interface {
	Get(key string, val interface{}) error
	Set(key string, val interface{}) error
	Delete(key string) error
	BFAdd(key string, val string) (bool, error)
	BFExists(key string, val string) (bool, error)
}
