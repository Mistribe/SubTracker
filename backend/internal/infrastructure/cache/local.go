package cache

type Local interface {
	Set(key string, value interface{})
	Get(key string) (interface{}, bool)
}
