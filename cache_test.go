package stormpath_test

import (
	. "github.com/jarias/stormpath-sdk-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

var _ = Describe("Cache", func() {
	Describe("LedisCache", func() {
		key := "key"
		cfg := lediscfg.NewConfigDefault()
		l, err := ledis.Open(cfg)
		if err != nil {
			panic(err)
		}
		db, err := l.Select(0)
		if err != nil {
			panic(err)
		}

		ledisCache := LedisCache{db}

		AfterEach(func() {
			ledisCache.DB.FlushAll()
		})

		Describe("Exists", func() {
			It("should return false if the key doesn't exists", func() {
				r := ledisCache.Exists(key)

				Expect(r).To(BeFalse())
			})
			It("should return true if the key does exists", func() {
				db.Set([]byte(key), []byte("hello"))

				r := ledisCache.Exists(key)

				Expect(r).To(BeTrue())
			})
		})

		Describe("Set", func() {
			It("should store a new object in the cache", func() {
				r := ledisCache.Exists(key)

				Expect(r).To(BeFalse())

				ledisCache.Set(key, "hello")

				r = ledisCache.Exists(key)

				Expect(r).To(BeTrue())
			})
			It("should update an existing object in the cache", func() {
				ledisCache.Set(key, "hello")
				ledisCache.Set(key, "bye")

				var r string
				ledisCache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal("bye"))
			})
		})
		Describe("Get", func() {
			It("should load empty data if the key doesn't exists into the given interface", func() {
				var r string

				err := ledisCache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(BeEmpty())
			})
			It("should load data from the cache into the given interface", func() {
				var r string

				ledisCache.Set(key, "hello")
				err := ledisCache.Get(key, &r)

				Expect(err).NotTo(HaveOccurred())
				Expect(r).To(Equal("hello"))
			})
		})
		Describe("Del", func() {
			It("should delete a given key from the cache", func() {
				ledisCache.Set(key, []byte("hello"))

				r := ledisCache.Exists(key)
				Expect(r).To(BeTrue())

				ledisCache.Del(key)

				r = ledisCache.Exists(key)
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
