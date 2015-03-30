package stormpath

import (
	"encoding/json"

	"github.com/garyburd/redigo/redis"
)

//Cache is a base interface for any cache provider
type Cache interface {
	Exists(key string) bool
	Set(key string, data interface{})
	Get(key string, result interface{}) error
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
		Logger.Printf("[ERROR] %s", err)
		return false
	}
	return exists.(int64) == 1
}

//Set stores data in the the cache for the given key
func (r RedisCache) Set(key string, data interface{}) {
	Logger.Printf("[DEBUG] Setting data from cache for key [%s]", key)
	jsonData, _ := json.Marshal(data)
	_, err := r.Conn.Do("SETEX", key, 30, string(jsonData))
	if err != nil {
		Logger.Printf("[ERROR] %s", err)
	}
}

//Get returns the data store under key it should return an error if any occur
func (r RedisCache) Get(key string, result interface{}) error {
	Logger.Printf("[DEBUG] Geting data from cache for key [%s]", key)
	cacheData, err := r.Conn.Do("GET", key)
	if err != nil {
		//Log the error and return an empty slice along with the error
		Logger.Printf("[ERROR] %s", err)
		return err
	}
	return json.Unmarshal(cacheData.([]byte), result)
}

//Del deletes a key from the cache
func (r RedisCache) Del(key string) {
	Logger.Printf("[DEBUG] Deleting data from cache for key [%s]", key)
	_, err := r.Conn.Do("DEL", key)
	if err != nil {
		Logger.Printf("[ERROR] %s", err)
	}
}
