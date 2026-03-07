// redisclient

// Package redisutil provides Redis client utilities with go-redis/v9 backend.
// It maintains API compatibility with the previous redigo-based implementation.
package redisutil

import (
	"context"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps go-redis client with compatible API
type RedisClient struct {
	client   *redis.Client
	Address  string
	maxIdle  int
	maxActive int
}

var (
	redisMap map[string]*RedisClient
	mapMutex *sync.RWMutex
)

const (
	defaultTimeout   = 60 * 10 // defaults to 10 minutes
	defaultMaxIdle   = 10
	defaultMaxActive = 50
)

func init() {
	redisMap = make(map[string]*RedisClient)
	mapMutex = new(sync.RWMutex)
}

// parseRedisURL parses redis URL and returns options
func parseRedisURL(redisURL string) *redis.Options {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		// Return default options if parse fails
		return &redis.Options{
			Addr: redisURL,
		}
	}
	return opts
}

// newClient creates a new go-redis client
func newClient(redisURL string, maxIdle, maxActive int) *redis.Client {
	opts := parseRedisURL(redisURL)
	
	// Map maxIdle/maxActive to go-redis pool settings
	// go-redis uses MinIdleConns for min idle, PoolSize for max connections
	if maxIdle <= 0 {
		maxIdle = defaultMaxIdle
	}
	if maxActive <= 0 {
		maxActive = defaultMaxActive
	}
	
	opts.MinIdleConns = maxIdle
	opts.PoolSize = maxActive
	
	return redis.NewClient(opts)
}

// GetDefaultRedisClient returns the RedisClient of specified address
// use default maxIdle & maxActive
func GetDefaultRedisClient(address string) *RedisClient {
	return GetRedisClient(address, defaultMaxIdle, defaultMaxActive)
}

// GetRedisClient returns the RedisClient of specified address & maxIdle & maxActive
func GetRedisClient(address string, maxIdle, maxActive int) *RedisClient {
	if maxIdle <= 0 {
		maxIdle = defaultMaxIdle
	}
	if maxActive <= 0 {
		maxActive = defaultMaxActive
	}
	
	var rc *RedisClient
	var mok bool
	
	mapMutex.RLock()
	rc, mok = redisMap[address]
	mapMutex.RUnlock()
	
	if !mok {
		rc = &RedisClient{
			Address:   address,
			client:    newClient(address, maxIdle, maxActive),
			maxIdle:   maxIdle,
			maxActive: maxActive,
		}
		mapMutex.Lock()
		redisMap[address] = rc
		mapMutex.Unlock()
	}
	return rc
}

// GetObj returns the content specified by key
func (rc *RedisClient) GetObj(key string) (interface{}, error) {
	ctx := context.Background()
	return rc.client.Get(ctx, key).Result()
}

// Get returns the content as string specified by key
func (rc *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil // Key not exists, return empty string
	}
	return val, err
}

// Exists whether key exists
func (rc *RedisClient) Exists(key string) (bool, error) {
	ctx := context.Background()
	val, err := rc.client.Exists(ctx, key).Result()
	return val > 0, err
}

// Del deletes specified key
func (rc *RedisClient) Del(key string) (int64, error) {
	ctx := context.Background()
	return rc.client.Del(ctx, key).Result()
}

// INCR atomically increment the value by 1 specified by key
func (rc *RedisClient) INCR(key string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.Incr(ctx, key).Result()
	return int(val), err
}

// DECR atomically decrement the value by 1 specified by key
func (rc *RedisClient) DECR(key string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.Decr(ctx, key).Result()
	return int(val), err
}

// Append appends the string to original value specified by key.
func (rc *RedisClient) Append(key string, val interface{}) (interface{}, error) {
	ctx := context.Background()
	return rc.client.Append(ctx, key, toString(val)).Result()
}

// toString converts interface{} to string
func toString(val interface{}) string {
	switch v := val.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		return ""
	}
}

// Set put key/value into redis
func (rc *RedisClient) Set(key string, val interface{}) (interface{}, error) {
	ctx := context.Background()
	return rc.client.Set(ctx, key, val, 0).Result()
}

// Expire specifies the expire duration for key
func (rc *RedisClient) Expire(key string, timeOutSeconds int64) (int64, error) {
	ctx := context.Background()
	val, err := rc.client.Expire(ctx, key, time.Duration(timeOutSeconds)*time.Second).Result()
	if err != nil {
		return 0, err
	}
	if val {
		return 1, nil
	}
	return 0, nil
}

// SetWithExpire set the key/value with specified duration
func (rc *RedisClient) SetWithExpire(key string, val interface{}, timeOutSeconds int64) (interface{}, error) {
	ctx := context.Background()
	return rc.client.Set(ctx, key, val, time.Duration(timeOutSeconds)*time.Second).Result()
}

// SetNX sets key/value only if key does not exists
func (rc *RedisClient) SetNX(key, value string) (interface{}, error) {
	ctx := context.Background()
	return rc.client.SetNX(ctx, key, value, 0).Result()
}

// ****************** hash set ***********************

// HGet returns content specified by hashID and field
func (rc *RedisClient) HGet(hashID string, field string) (string, error) {
	ctx := context.Background()
	val, err := rc.client.HGet(ctx, hashID, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// HGetAll returns all content specified by hashID
func (rc *RedisClient) HGetAll(hashID string) (map[string]string, error) {
	ctx := context.Background()
	return rc.client.HGetAll(ctx, hashID).Result()
}

// HSet set content with hashID and field
func (rc *RedisClient) HSet(hashID string, field string, val string) error {
	ctx := context.Background()
	return rc.client.HSet(ctx, hashID, field, val).Err()
}

// HSetNX set content with hashID and field, if the field does not exists
func (rc *RedisClient) HSetNX(hashID, field, value string) (interface{}, error) {
	ctx := context.Background()
	return rc.client.HSetNX(ctx, hashID, field, value).Result()
}

// HExist returns if the field exists in specified hashID
func (rc *RedisClient) HExist(hashID string, field string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.HExists(ctx, hashID, field).Result()
	if val {
		return 1, err
	}
	return 0, err
}

// HIncrBy increment the value specified by hashID and field
func (rc *RedisClient) HIncrBy(hashID string, field string, increment int) (int, error) {
	ctx := context.Background()
	val, err := rc.client.HIncrBy(ctx, hashID, field, int64(increment)).Result()
	return int(val), err
}

// HLen returns count of fields in hashID
func (rc *RedisClient) HLen(hashID string) (int64, error) {
	ctx := context.Background()
	return rc.client.HLen(ctx, hashID).Result()
}

// HDel delete content in hashset
func (rc *RedisClient) HDel(args ...interface{}) (int64, error) {
	ctx := context.Background()
	if len(args) == 0 {
		return 0, nil
	}
	
	// First arg is hashID, rest are fields
	hashID := toString(args[0])
	fields := make([]string, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		fields = append(fields, toString(args[i]))
	}
	
	return rc.client.HDel(ctx, hashID, fields...).Result()
}

// HVals return all the values in all fields specified by hashID
func (rc *RedisClient) HVals(hashID string) (interface{}, error) {
	ctx := context.Background()
	return rc.client.HVals(ctx, hashID).Result()
}

// ****************** list ***********************

// LPush insert the values into front of the list
func (rc *RedisClient) LPush(key string, value ...interface{}) (int, error) {
	ctx := context.Background()
	val, err := rc.client.LPush(ctx, key, value...).Result()
	return int(val), err
}

// LPushX inserts value at the head of the list only if key exists
func (rc *RedisClient) LPushX(key string, value string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.LPushX(ctx, key, value).Result()
	return int(val), err
}

// LRange returns elements from start to stop
func (rc *RedisClient) LRange(key string, start int, stop int) ([]string, error) {
	ctx := context.Background()
	return rc.client.LRange(ctx, key, int64(start), int64(stop)).Result()
}

// LRem removes count elements equal to value
func (rc *RedisClient) LRem(key string, count int, value string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.LRem(ctx, key, int64(count), value).Result()
	return int(val), err
}

// LSet sets the list element at index to value
func (rc *RedisClient) LSet(key string, index int, value string) (string, error) {
	ctx := context.Background()
	return rc.client.LSet(ctx, key, int64(index), value).Result()
}

// LTrim trims the list to the specified range
func (rc *RedisClient) LTrim(key string, start int, stop int) (string, error) {
	ctx := context.Background()
	return rc.client.LTrim(ctx, key, int64(start), int64(stop)).Result()
}

// RPop removes and returns the last element of the list
func (rc *RedisClient) RPop(key string) (string, error) {
	ctx := context.Background()
	return rc.client.RPop(ctx, key).Result()
}

// RPush inserts values at the tail of the list
func (rc *RedisClient) RPush(key string, value ...interface{}) (int, error) {
	ctx := context.Background()
	val, err := rc.client.RPush(ctx, key, value...).Result()
	return int(val), err
}

// RPushX inserts value at the tail of the list only if key exists
func (rc *RedisClient) RPushX(key string, value ...interface{}) (int, error) {
	ctx := context.Background()
	if len(value) == 0 {
		return 0, nil
	}
	val, err := rc.client.RPushX(ctx, key, value[0]).Result()
	return int(val), err
}

// RPopLPush removes the last element from one list and pushes it to another
func (rc *RedisClient) RPopLPush(source string, destination string) (string, error) {
	ctx := context.Background()
	return rc.client.RPopLPush(ctx, source, destination).Result()
}

// BLPop removes and returns the first element of the first non-empty list
func (rc *RedisClient) BLPop(key ...interface{}) (map[string]string, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	result, err := rc.client.BLPop(ctx, time.Duration(defaultTimeout)*time.Second, keys...).Result()
	if err != nil {
		return nil, err
	}
	// Convert []string to map[string]string
	if len(result) >= 2 {
		return map[string]string{result[0]: result[1]}, nil
	}
	return nil, nil
}

// BRPop removes and returns the last element of the first non-empty list
func (rc *RedisClient) BRPop(key ...interface{}) (map[string]string, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	result, err := rc.client.BRPop(ctx, time.Duration(defaultTimeout)*time.Second, keys...).Result()
	if err != nil {
		return nil, err
	}
	if len(result) >= 2 {
		return map[string]string{result[0]: result[1]}, nil
	}
	return nil, nil
}

// BRPopLPush pops from one list and pushes to another with blocking
func (rc *RedisClient) BRPopLPush(source string, destination string) (string, error) {
	ctx := context.Background()
	return rc.client.BRPopLPush(ctx, source, destination, time.Duration(defaultTimeout)*time.Second).Result()
}

// LIndex returns the element at index
func (rc *RedisClient) LIndex(key string, index int) (string, error) {
	ctx := context.Background()
	return rc.client.LIndex(ctx, key, int64(index)).Result()
}

// LInsertBefore inserts value before pivot
func (rc *RedisClient) LInsertBefore(key string, pivot string, value string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.LInsertBefore(ctx, key, pivot, value).Result()
	return int(val), err
}

// LInsertAfter inserts value after pivot
func (rc *RedisClient) LInsertAfter(key string, pivot string, value string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.LInsertAfter(ctx, key, pivot, value).Result()
	return int(val), err
}

// LLen returns the length of the list
func (rc *RedisClient) LLen(key string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.LLen(ctx, key).Result()
	return int(val), err
}

// LPop removes and returns the first element of the list
func (rc *RedisClient) LPop(key string) (string, error) {
	ctx := context.Background()
	return rc.client.LPop(ctx, key).Result()
}

// ****************** set ***********************

// SAdd add one or multiple members into the set
func (rc *RedisClient) SAdd(key string, member ...interface{}) (int, error) {
	ctx := context.Background()
	val, err := rc.client.SAdd(ctx, key, member...).Result()
	return int(val), err
}

// SCard returns cardinality of the set
func (rc *RedisClient) SCard(key string) (int, error) {
	ctx := context.Background()
	val, err := rc.client.SCard(ctx, key).Result()
	return int(val), err
}

// SPop removes and returns a random member from the set
func (rc *RedisClient) SPop(key string) (string, error) {
	ctx := context.Background()
	return rc.client.SPop(ctx, key).Result()
}

// SRandMember returns random count elements from set
func (rc *RedisClient) SRandMember(key string, count int) ([]string, error) {
	ctx := context.Background()
	return rc.client.SRandMemberN(ctx, key, int64(count)).Result()
}

// SRem removes multiple elements from set
func (rc *RedisClient) SRem(key string, member ...interface{}) (int, error) {
	ctx := context.Background()
	val, err := rc.client.SRem(ctx, key, member...).Result()
	return int(val), err
}

// SDiff returns the difference between sets
func (rc *RedisClient) SDiff(key ...interface{}) ([]string, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	return rc.client.SDiff(ctx, keys...).Result()
}

// SDiffStore stores the difference in a new set
func (rc *RedisClient) SDiffStore(destination string, key ...interface{}) (int, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	val, err := rc.client.SDiffStore(ctx, destination, keys...).Result()
	return int(val), err
}

// SInter returns the intersection of sets
func (rc *RedisClient) SInter(key ...interface{}) ([]string, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	return rc.client.SInter(ctx, keys...).Result()
}

// SInterStore stores the intersection in a new set
func (rc *RedisClient) SInterStore(destination string, key ...interface{}) (int, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	val, err := rc.client.SInterStore(ctx, destination, keys...).Result()
	return int(val), err
}

// SIsMember returns if member is a member of set
func (rc *RedisClient) SIsMember(key string, member string) (bool, error) {
	ctx := context.Background()
	return rc.client.SIsMember(ctx, key, member).Result()
}

// SMembers returns all members of the set
func (rc *RedisClient) SMembers(key string) ([]string, error) {
	ctx := context.Background()
	return rc.client.SMembers(ctx, key).Result()
}

// SMove moves member from one set to another
func (rc *RedisClient) SMove(source string, destination string, member string) (bool, error) {
	ctx := context.Background()
	return rc.client.SMove(ctx, source, destination, member).Result()
}

// SUnion returns the union of sets
func (rc *RedisClient) SUnion(key ...interface{}) ([]string, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	return rc.client.SUnion(ctx, keys...).Result()
}

// SUnionStore stores the union in a new set
func (rc *RedisClient) SUnionStore(destination string, key ...interface{}) (int, error) {
	ctx := context.Background()
	keys := make([]string, 0, len(key))
	for _, k := range key {
		if str, ok := k.(string); ok {
			keys = append(keys, str)
		}
	}
	val, err := rc.client.SUnionStore(ctx, destination, keys...).Result()
	return int(val), err
}

// ****************** Global functions ***********************

// Ping tests the client is ready for use
func (rc *RedisClient) Ping() (string, error) {
	ctx := context.Background()
	return rc.client.Ping(ctx).Result()
}

// DBSize returns count of keys in the database
func (rc *RedisClient) DBSize() (int64, error) {
	ctx := context.Background()
	return rc.client.DBSize(ctx).Result()
}

// FlushDB removes all data in the database
func (rc *RedisClient) FlushDB() {
	ctx := context.Background()
	rc.client.FlushDB(ctx)
}

// GetConn returns a connection from the pool
// Deprecated: This method exists for backwards compatibility but is not recommended.
// Use the RedisClient methods directly instead.
func (rc *RedisClient) GetConn() interface{} {
	// Return a wrapper that mimics redigo's Conn interface
	// This is for backwards compatibility only
	return &connWrapper{client: rc.client}
}

// connWrapper wraps go-redis client to provide a Conn-like interface
type connWrapper struct {
	client *redis.Client
}

// Do executes a command (simplified for backwards compatibility)
func (c *connWrapper) Do(commandName string, args ...interface{}) (interface{}, error) {
	ctx := context.Background()
	cmd := redis.NewCmd(ctx, append([]interface{}{commandName}, args...)...)
	c.client.Process(ctx, cmd)
	return cmd.Result()
}

// Close is a no-op for connection pooling
func (c *connWrapper) Close() error {
	return nil
}

// Err returns nil (go-redis handles errors differently)
func (c *connWrapper) Err() error {
	return nil
}
