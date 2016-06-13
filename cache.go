package stormpath

//Cacheable determines if the implementor should be cached or not
type Cacheable interface {
	IsCacheable() bool
}

//Cache is a base interface for any cache provider
type Cache interface {
	Exists(key string) bool
	Set(key string, data []byte)
	Get(key string) []byte
	Del(key string)
}
