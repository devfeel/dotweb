package redis

import (
	"github.com/devfeel/dotweb/framework/redis"
	"strconv"
)

var (
	ZeroInt64 int64 = 0
)

// RedisCache is redis cache adapter.
// it contains serverIp for redis conn.
type RedisCache struct {
	serverIp string //connection string, like "redis://:password@10.0.1.11:6379/0"
}

// NewRedisCache returns a new *RedisCache.
func NewRedisCache(serverIp string) *RedisCache {
	cache := RedisCache{serverIp: serverIp}
	return &cache
}

// Exists check item exist in redis cache.
func (ca *RedisCache) Exists(key string) bool {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	exists, err := redisClient.Exists(key)
	if err != nil {
		return false
	}
	return exists
}

// Incr increase int64 counter in redis cache.
func (ca *RedisCache) Incr(key string) (int64, error) {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	val, err := redisClient.INCR(key)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// Decr decrease counter in redis cache.
func (ca *RedisCache) Decr(key string) (int64, error) {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	val, err := redisClient.DECR(key)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// Get cache from redis cache.
// if non-existed or expired, return nil.
func (ca *RedisCache) Get(key string) interface{} {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	reply, err := redisClient.GetObj(key)
	if err != nil {
		return nil
	} else {
		return reply
	}
}

//  returns value string format by given key
// if non-existed or expired, return "".
func (ca *RedisCache) GetString(key string) string {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	reply, err := redisClient.Get(key)
	if err != nil {
		return ""
	} else {
		return reply
	}
}

//  returns value int format by given key
// if non-existed or expired, return nil.
func (ca *RedisCache) GetInt(key string) int {
	v := ca.GetString(key)
	if v == "" {
		return 0
	} else {
		i, e := strconv.Atoi(v)
		if e != nil {
			return 0
		} else {
			return i
		}
	}
}

//  returns value int64 format by given key
// if non-existed or expired, return nil.
func (ca *RedisCache) GetInt64(key string) int64 {
	v := ca.GetString(key)
	if v == "" {
		return ZeroInt64
	} else {
		i, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			return ZeroInt64
		} else {
			return i
		}
	}
}

// Set cache to redis.
// ttl is second, if ttl is 0, it will be forever.
func (ca *RedisCache) Set(key string, value interface{}, ttl int64) error {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	_, err := redisClient.SetWithExpire(key, value, ttl)
	return err
}

// Delete item in redis cacha.
// if not exists, we think it's success
func (ca *RedisCache) Delete(key string) error {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	_, err := redisClient.Del(key)
	return err
}

// ClearAll will delete all item in redis cache.
// never error
func (ca *RedisCache) ClearAll() error {
	redisClient := redisutil.GetRedisClient(ca.serverIp)
	redisClient.FlushDB()
	return nil
}
