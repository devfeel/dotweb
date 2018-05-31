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
	serverURL string //connection string, like "redis://:password@10.0.1.11:6379/0"
}

// NewRedisCache returns a new *RedisCache.
func NewRedisCache(serverURL string) *RedisCache {
	cache := RedisCache{serverURL: serverURL}
	return &cache
}

// Exists check item exist in redis cache.
func (ca *RedisCache) Exists(key string) (bool, error) {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	exists, err := redisClient.Exists(key)
	return exists, err
}

// Incr increase int64 counter in redis cache.
func (ca *RedisCache) Incr(key string) (int64, error) {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	val, err := redisClient.INCR(key)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// Decr decrease counter in redis cache.
func (ca *RedisCache) Decr(key string) (int64, error) {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	val, err := redisClient.DECR(key)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

// Get cache from redis cache.
// if non-existed or expired, return nil.
func (ca *RedisCache) Get(key string) (interface{}, error) {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	reply, err := redisClient.GetObj(key)
	return reply, err
}

//  returns value string format by given key
// if non-existed or expired, return "".
func (ca *RedisCache) GetString(key string) (string, error) {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	reply, err := redisClient.Get(key)
	return reply, err
}

//  returns value int format by given key
// if non-existed or expired, return nil.
func (ca *RedisCache) GetInt(key string) (int, error) {
	v, err := ca.GetString(key)
	if err != nil || v == "" {
		return 0, err
	} else {
		i, e := strconv.Atoi(v)
		if e != nil {
			return 0, err
		} else {
			return i, nil
		}
	}
}

//  returns value int64 format by given key
// if non-existed or expired, return nil.
func (ca *RedisCache) GetInt64(key string) (int64, error) {
	v, err := ca.GetString(key)
	if err != nil || v == "" {
		return ZeroInt64, err
	} else {
		i, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			return ZeroInt64, err
		} else {
			return i, nil
		}
	}
}

// Set cache to redis.
// ttl is second, if ttl is 0, it will be forever.
func (ca *RedisCache) Set(key string, value interface{}, ttl int64) error {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	var err error
	if ttl <= 0{
		_, err = redisClient.Set(key, value)
	}else{
		_, err = redisClient.SetWithExpire(key, value, ttl)
	}
	return err
}

// Delete item in redis cacha.
// if not exists, we think it's success
func (ca *RedisCache) Delete(key string) error {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	_, err := redisClient.Del(key)
	return err
}

// ClearAll will delete all item in redis cache.
// never error
func (ca *RedisCache) ClearAll() error {
	redisClient := redisutil.GetRedisClient(ca.serverURL)
	redisClient.FlushDB()
	return nil
}
