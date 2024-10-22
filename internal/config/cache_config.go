package config

import (
	"log"
	"os"
	"strconv"
)

type CacheConfig struct {
	ServerAddress string
	MaxCapacity   int
}

func initCacheConfig() CacheConfig {
	maxCapacity, err := strconv.Atoi(os.Getenv("CACHE_MAX_CAPACITY"))
	if err != nil {
		log.Fatal("failedto parse CACHE_MAX_CAPACITY")
	}
	return CacheConfig{
		ServerAddress: os.Getenv("MEMCACHE_SERVER_ADDRESS"),
		MaxCapacity:   maxCapacity,
	}
}
