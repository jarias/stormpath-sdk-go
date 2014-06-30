package stormpath

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jarias/stormpath-sdk-go/logger"
)

//Cache is a base interface for any cache provider
type Cache interface {
	Exists(key string) bool
	Set(key string, data []byte)
	Get(key string) ([]byte, error)
	Del(key string)
}

//RedisCache is the default provided implementation of the Cache interface using Redis as backend
type RedisCache struct {
	Conn redis.Conn
}

//Exists returns true if the key exists in the cache false otherwise
func (r RedisCache) Exists(key string) bool {
	exists, err := r.Conn.Do("EXISTS", key)
	if err != nil {
		//Log the error but don't crash or pass it along if the cache is not working the reset should
		logger.ERROR.Println(err)
		return false
	}
	return exists.(int64) == 1
}

//Set stores data in the the cache for the given key
func (r RedisCache) Set(key string, data []byte) {
	logger.CACHE.Printf("Setting data from cache for key [%s]", key)
	_, err := r.Conn.Do("SETEX", key, 30, string(data))
	if err != nil {
		logger.ERROR.Println(err)
	}
}

//Get returns the data store under key it should return an error if any occur
func (r RedisCache) Get(key string) ([]byte, error) {
	logger.CACHE.Printf("Geting data from cache for key [%s]", key)
	cacheData, err := r.Conn.Do("GET", key)
	if err != nil {
		//Log the error and return an empty slice along with the error
		logger.ERROR.Println(err)
		return []byte{}, err
	}
	return cacheData.([]byte), err
}

//Del deletes a key from the cache
func (r RedisCache) Del(key string) {
	logger.CACHE.Printf("Deleting data from cache for key [%s]", key)
	_, err := r.Conn.Do("DEL", key)
	if err != nil {
		logger.ERROR.Println(err)
	}
}
