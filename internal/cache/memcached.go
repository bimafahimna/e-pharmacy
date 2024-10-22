package cache

import (
	"encoding/json"

	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/config"
	"github.com/bimafahimna/E-Pharmacy-ServerSide/internal/pkg/logger"
	"github.com/bradfitz/gomemcache/memcache"
)

type memcached struct {
	client *memcache.Client
	lru    EvictAlgo
}

func NewCacheProvider(config config.CacheConfig) Provider {
	client := &memcached{client: memcache.New(config.ServerAddress), lru: constructor(config.MaxCapacity)}
	go func() {
		if err := client.Ping(); err != nil {
			logger.Log.Fatalf("failed connecting to memcached: %v", err)
		} else {
			logger.Log.Info("connected to memcached")
		}
	}()
	return client
}

func (c *memcached) Ping() error {
	return c.client.Ping()
}

func (c *memcached) Set(key string, value []byte, expiration int32) error {
	if err := c.client.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration}); err != nil {
		return err
	}
	c.lru.set(key)
	return nil
}
func (c *memcached) SetFromStruct(key string, data interface{}, expiration int32) error {
	value, err := json.Marshal(data)
	if err != nil {
		logger.Log.Errorf("failed to marshal data: %v", err)
		return err
	}
	return c.client.Set(&memcache.Item{Key: key, Value: value, Expiration: expiration})
}

func (c *memcached) Get(key string) ([]byte, error) {
	item, err := c.client.Get(key)
	if err != nil {
		return nil, err
	}
	c.lru.get(key)
	return item.Value, nil
}

func (c *memcached) FlushAll() error {
	return c.client.FlushAll()
}

func (c *memcached) DeleteAll() error {
	return c.client.DeleteAll()
}
