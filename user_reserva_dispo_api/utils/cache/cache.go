package cache

import (
	"encoding/json"
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
)

var cacheClient *memcache.Client

func Init_cache() {
	cacheClient = memcache.New("cache:11211") // Asegúrate de que la dirección IP y el puerto son correctos
	fmt.Println("Initialized cache", cacheClient)
	log.Info("Initialized cache")
}

// SetCache sets a value with a specific key in Memcached.
func SetCache(key string, value interface{}, ttl int32) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return cacheClient.Set(&memcache.Item{Key: key, Value: data, Expiration: ttl})
}

// GetCache retrieves a value based on a key from Memcached.
func GetCache(key string, v interface{}) error {
	item, err := cacheClient.Get(key)
	if err != nil {
		return err
	}
	return json.Unmarshal(item.Value, v)
}
