package stormpath

//Cacheable determines if the implementor should be cached or not
type Cacheable interface {
	IsCacheable() bool
}

//cacheResource stores a resource in the cache if the resource allows caching
func cacheResource(key string, resource interface{}, cache Cache) {
	c, ok := resource.(Cacheable)

	if ok && c.IsCacheable() {
		cache.Set(key, resource)
	}
}

//Cache is a base interface for any cache provider
type Cache interface {
	Exists(key string) bool
	Set(key string, data interface{})
	Get(key string, result interface{}) error
	Del(key string)
}
