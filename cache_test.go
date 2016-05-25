package stormpath_test

import (
	"testing"

	"os"

	. "github.com/jarias/stormpath-sdk-go"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
	"github.com/stretchr/testify/assert"
)

const key = "key"

func createTestLedisCache() LedisCache {
	cfg := lediscfg.NewConfigDefault()
	cfg.DataDir = os.TempDir() + "/stormpath-go-sdk-ledisCache/var/" + randomName()
	l, err := ledis.Open(cfg)
	if err != nil {
		panic(err)
	}
	db, err := l.Select(0)
	if err != nil {
		panic(err)
	}

	return LedisCache{db}
}

func TestLedisCacheKeyNoExists(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	r := cache.Exists(key)

	assert.False(t, r)
}

func TestLedisCacheKeyExists(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	cache.DB.Set([]byte(key), []byte("hello"))

	r := cache.Exists(key)

	assert.True(t, r)
}

func TestLedisCacheSetObject(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	assert.False(t, cache.Exists(key))

	cache.Set(key, "hello")

	assert.True(t, cache.Exists(key))
}

func TestLedisCacheUpdateObject(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	assert.False(t, cache.Exists(key))

	cache.Set(key, "hello")
	cache.Set(key, "bye")

	var r string
	err := cache.Get(key, &r)

	assert.NoError(t, err)
	assert.Equal(t, "bye", r)
}

func TestLedisCacheGetObjectNoExists(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	var r string

	err := cache.Get(key, &r)

	assert.NoError(t, err)
	assert.Empty(t, r)
}

func TestLedisCacheGetObject(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	cache.Set(key, "hello")

	var r string

	err := cache.Get(key, &r)

	assert.NoError(t, err)
	assert.Equal(t, "hello", r)
}

func TestLedisCacheDeleteObject(t *testing.T) {
	t.Parallel()
	cache := createTestLedisCache()

	cache.Set(key, "hello")

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
