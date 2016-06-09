package stormpath_test

import (
	"testing"
	"time"

	. "github.com/jarias/stormpath-sdk-go"
	"github.com/stretchr/testify/assert"
)

const key = "key"

func createTestLocalCache() *LocalCache {
	return NewLocalCache(5*time.Second, 2*time.Second)
}

func TestLocalCacheKeyTTL(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	cache.Set(key, []byte("hello"))

	time.Sleep(6 * time.Second)

	assert.False(t, cache.Exists(key))
}

func TestLocalCacheKeyTTI(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	cache.Set(key, []byte("hello"))

	cache.Get(key)
	time.Sleep(3 * time.Second)

	assert.False(t, cache.Exists(key))
}

func TestLocalCacheKeyNoExists(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	r := cache.Exists(key)

	assert.False(t, r)
}

func TestLocalCacheKeyExists(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	cache.Set(key, []byte("hello"))

	r := cache.Exists(key)

	assert.True(t, r)
}

func TestLocalCacheSetObject(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	assert.False(t, cache.Exists(key))

	cache.Set(key, []byte("hello"))

	assert.True(t, cache.Exists(key))
}

func TestLocalCacheUpdateObject(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	assert.False(t, cache.Exists(key))

	cache.Set(key, []byte("hello"))
	cache.Set(key, []byte("bye"))

	r := cache.Get(key)

	assert.Equal(t, []byte("bye"), r)
}

func TestLocalCacheGetObjectNoExists(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	r := cache.Get(key)

	assert.Empty(t, r)
}

func TestLocalCacheGetObject(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	cache.Set(key, []byte("hello"))

	r := cache.Get(key)

	assert.Equal(t, []byte("hello"), r)
}

func TestLocalCacheDeleteObject(t *testing.T) {
	t.Parallel()
	cache := createTestLocalCache()

	cache.Set(key, []byte("hello"))

	assert.True(t, cache.Exists(key))

	cache.Del(key)

	assert.False(t, cache.Exists(key))
}

func TestNonCacheableResources(t *testing.T) {
	var resources = []interface{}{
		&Applications{},
		&Accounts{},
		&Groups{},
		&Directories{},
		&AccountStoreMappings{},
	}
	for _, resource := range resources {
		c, ok := resource.(Cacheable)

		assert.True(t, ok)
		assert.False(t, c.IsCacheable())
	}
}

func TestCacheableResources(t *testing.T) {
	var resources = []interface{}{
		&Application{},
		&Account{},
		&Group{},
		&Directory{},
		&AccountStoreMapping{},
		&Tenant{},
	}

	for _, resource := range resources {
		c, ok := resource.(Cacheable)

		assert.True(t, ok)
		assert.True(t, c.IsCacheable())
	}
}
