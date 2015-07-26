package stormpath_test

import (
	"encoding/json"
	"os"

	"github.com/garyburd/redigo/redis"
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cache", func() {
	Describe("RedisCache", func() {
		key := "key"
		redisServer := os.Getenv("REDIS_SERVER")
		redisConn, err := redis.Dial("tcp", redisServer+":6379")
		if err != nil {
			panic(err)
		}
		redisCache := RedisCache{redisConn}

		AfterEach(func() {
			redisConn.Do("FLUSHDB")
		})

		Describe("Exists", func() {
			It("should return false if the key doesn't exists", func() {
				r := redisCache.Exists(key)

				Expect(r).To(BeFalse())
			})
			It("should return true if the key does exists", func() {
				redisConn.Do("SET", key, 1)

				r := redisCache.Exists(key)

				Expect(r).To(BeTrue())
			})
		})

		Describe("Set", func() {
			It("should store a new object in the cache", func() {
				r := redisCache.Exists(key)

				Expect(r).To(BeFalse())

				redisCache.Set(key, 1)

				r = redisCache.Exists(key)

				Expect(r).To(BeTrue())
			})
			It("should update an existing object in the cache", func() {
				var r int

				redisCache.Set(key, 1)
				redisCache.Set(key, 2)

				cacheData, err := redisConn.Do("GET", key)

				json.Unmarshal(cacheData.([]byte), &r)
				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal(2))
			})
		})
		Describe("Get", func() {
			It("should load empty data if the key doesn't exists into the given interface", func() {
				var r int

				err := redisCache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal(0))
			})
			It("should load data from the cache into the given interface", func() {
				var r int

				redisCache.Set(key, 2)
				err := redisCache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal(2))
			})
		})
		Describe("Del", func() {
			It("should delete a given key from the cache", func() {
				redisCache.Set(key, 2)

				r := redisCache.Exists(key)
				Expect(r).To(BeTrue())

				redisCache.Del(key)

				r = redisCache.Exists(key)
				Expect(r).To(BeFalse())
			})
		})
	})
})

var _ = Describe("Cacheable", func() {
	Describe("Collection resource", func() {
		It("should not be cacheable", func() {
			var resources = []interface{}{
				&Applications{},
				&Accounts{},
				&Groups{},
				&Directories{},
				&AccountStoreMappings{},
			}
			for _, resource := range resources {
				c, ok := resource.(Cacheable)

				Expect(ok).To(BeTrue())
				Expect(c.IsCacheable()).To(BeFalse())
			}
		})
	})

	Describe("Single resource", func() {
		It("should be cacheable", func() {
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

				Expect(ok).To(BeTrue())
				Expect(c.IsCacheable()).To(BeTrue())
			}
		})
	})
})
