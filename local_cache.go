package stormpath

import (
	"sync"
	"time"
)

type cacheItem struct {
	sync.RWMutex
	data    []byte
	expires *time.Time
}

func (item *cacheItem) touch(duration time.Duration) {
	item.Lock()
	expiration := time.Now().Add(duration)
	item.expires = &expiration
	item.Unlock()
}

func (item *cacheItem) expired() bool {
	var value bool
	item.RLock()
	if item.expires == nil {
		value = true
	} else {
		value = item.expires.Before(time.Now())
	}
	item.RUnlock()
	return value
}

type LocalCache struct {
	mutex sync.RWMutex
	ttl   time.Duration
	tti   time.Duration
	items map[string]*cacheItem
}

func (cache *LocalCache) Set(key string, data []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	item := &cacheItem{data: data}
	item.touch(cache.ttl)
	cache.items[key] = item
}

func (cache *LocalCache) Get(key string) []byte {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item, exists := cache.items[key]
	if exists && !item.expired() {
		item.touch(cache.tti)

		return item.data
	}
	return []byte{}
}

func (cache *LocalCache) Del(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	delete(cache.items, key)
}

func (cache *LocalCache) Exists(key string) bool {
	_, exists := cache.items[key]
	return exists
}

// Count returns the number of items in the cache
// (helpful for tracking memory leaks)
func (cache *LocalCache) Count() int {
	cache.mutex.RLock()
	count := len(cache.items)
	cache.mutex.RUnlock()
	return count
}

func (cache *LocalCache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *LocalCache) startCleanupTimer() {
	ttlTicker := time.Tick(cache.ttl)
	go (func() {
		for {
			select {
			case <-ttlTicker:
				cache.cleanup()
			}
		}
	})()

	ttiTicker := time.Tick(cache.tti)
	go (func() {
		for {
			select {
			case <-ttiTicker:
				cache.cleanup()
			}
		}
	})()
}

func NewLocalCache(ttl time.Duration, tti time.Duration) *LocalCache {
	cache := &LocalCache{
		ttl:   ttl,
		tti:   tti,
		items: map[string]*cacheItem{},
	}
	cache.startCleanupTimer()
	return cache
}
