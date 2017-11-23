package cache

import (
	"github.com/devfeel/dotweb/cache/redis"
	"github.com/devfeel/dotweb/cache/runtime"
)

type Cache interface {
	// Exist return true if value cached by given key
	Exists(key string) (bool, error)
	// Get returns value by given key
	Get(key string) (interface{}, error)
	// GetString returns value string format by given key
	GetString(key string) (string, error)
	// GetInt returns value int format by given key
	GetInt(key string) (int, error)
	// GetInt64 returns value int64 format by given key
	GetInt64(key string) (int64, error)
	// Set cache value by given key
	Set(key string, v interface{}, ttl int64) error
	// Incr increases int64-type value by given key as a counter
	// if key not exist, before increase set value with zero
	Incr(key string) (int64, error)
	// Decr decreases int64-type value by given key as a counter
	// if key not exist, before increase set value with zero
	Decr(key string) (int64, error)
	// Delete delete cache item by given key
	Delete(key string) error
	// ClearAll clear all cache items
	ClearAll() error
}

// NewRuntimeCache new runtime cache
func NewRuntimeCache() Cache {
	return runtime.NewRuntimeCache()
}

// NewRedisCache create new redis cache
// must set serverURL like "redis://:password@10.0.1.11:6379/0"
func NewRedisCache(serverURL string) Cache {
	return redis.NewRedisCache(serverURL)
}
