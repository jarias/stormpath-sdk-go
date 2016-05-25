package stormpath

import (
	"encoding/json"

	"github.com/siddontang/ledisdb/ledis"
)

//RedisCache is the default provided implementation of the Cache interface using Redis as backend
type LedisCache struct {
	DB *ledis.DB
}

//Exists returns true if the key exists in the cache false otherwise
func (cache LedisCache) Exists(key string) bool {
	exists, err := cache.DB.Exists([]byte(key))
	if err != nil {
		//Log the error but don't crash
		Logger.Printf("[ERROR] %s", err)
		return false
	}
	return exists == 1
}

//Set stores data in the the cache for the given key
func (cache LedisCache) Set(key string, data interface{}) {
	Logger.Printf("[DEBUG] Setting data from cache for key [%s]", key)
	jsonData, _ := json.Marshal(data)
	err := cache.DB.Set([]byte(key), jsonData)
	if err != nil {
		Logger.Printf("[ERROR] %s", err)
	}
}

//Get returns the data store under key it should return an error if any occur
func (cache LedisCache) Get(key string, result interface{}) error {
	Logger.Printf("[DEBUG] Geting data from cache for key [%s]", key)

	if !cache.Exists(key) {
		return nil
	}
	cacheData, err := cache.DB.Get([]byte(key))
	if err != nil {
		//Log the error and return an empty slice along with the error
		Logger.Printf("[ERROR] %s", err)
		return err
	}
	return json.Unmarshal(cacheData, result)
}

//Del deletes a key from the cache
func (cache LedisCache) Del(key string) {
	Logger.Printf("[DEBUG] Deleting data from cache for key [%s]", key)
	_, err := cache.DB.Del([]byte(key))
	if err != nil {
		Logger.Printf("[ERROR] %s", err)
	}
}
