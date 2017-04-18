// redisclient
package redisutil

import (
	"github.com/garyburd/redigo/redis"
	"sync"
)

type RedisClient struct {
	pool    *redis.Pool
	Address string
}

var (
	redisMap map[string]*RedisClient
	mapMutex *sync.RWMutex
)

const (
	defaultTimeout = 60 * 10 //默认10分钟
)

func init() {
	redisMap = make(map[string]*RedisClient)
	mapMutex = new(sync.RWMutex)
}

// 重写生成连接池方法
func newPool(redisIP string) *redis.Pool {

	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 20, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisIP)
			return c, err
		},
	}
}

//获取指定Address的RedisClient
func GetRedisClient(address string) *RedisClient {
	var redis *RedisClient
	var mok bool
	mapMutex.RLock()
	redis, mok = redisMap[address]
	mapMutex.RUnlock()
	if !mok {
		redis = &RedisClient{Address: address, pool: newPool(address)}
		mapMutex.Lock()
		redisMap[address] = redis
		mapMutex.Unlock()
	}
	return redis
}

//获取指定key的内容, interface{}
func (rc *RedisClient) GetObj(key string) (interface{}, error) {
	// 从连接池里面获得一个连接
	conn := rc.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer conn.Close()
	reply, errDo := conn.Do("GET", key)
	return reply, errDo
}

//获取指定key的内容, string
func (rc *RedisClient) Get(key string) (string, error) {
	val, err := redis.String(rc.GetObj(key))
	return val, err
}

//检查指定key是否存在
func (rc *RedisClient) Exists(key string) (bool, error) {
	// 从连接池里面获得一个连接
	conn := rc.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer conn.Close()
	reply, errDo := conn.Do("EXISTS", key)
	if errDo == nil && reply == nil {
		return false, nil
	}
	val, err := redis.Int(reply, errDo)
	return val > 0, err
}

//删除指定key
func (rc *RedisClient) Del(key string) (int64, error) {
	// 从连接池里面获得一个连接
	conn := rc.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer conn.Close()
	reply, errDo := conn.Do("DEL", key)
	if errDo == nil && reply == nil {
		return 0, nil
	}
	val, err := redis.Int64(reply, errDo)
	return val, err
}

//获取指定hashset的内容
func (rc *RedisClient) HGet(hashID string, field string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, errDo := conn.Do("HGET", hashID, field)
	if errDo == nil && reply == nil {
		return "", nil
	}
	val, err := redis.String(reply, errDo)
	return val, err
}

//对存储在指定key的数值执行原子的加1操作
func (rc *RedisClient) INCR(key string) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, errDo := conn.Do("INCR", key)
	if errDo == nil && reply == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}

//对存储在指定key的数值执行原子的减1操作
func (rc *RedisClient) DECR(key string) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, errDo := conn.Do("DECR", key)
	if errDo == nil && reply == nil {
		return 0, nil
	}
	val, err := redis.Int(reply, errDo)
	return val, err
}

//获取指定hashset的所有内容
func (rc *RedisClient) HGetAll(hashID string) (map[string]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, err := redis.StringMap(conn.Do("HGetAll", hashID))
	return reply, err
}

//设置指定hashset的内容
func (rc *RedisClient) HSet(hashID string, field string, val string) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", hashID, field, val)
	return err
}

//删除，并获得该列表中的最后一个元素，或阻塞，直到有一个可用
func (rc *RedisClient) BRPop(key string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.StringMap(conn.Do("BRPOP", key, defaultTimeout))
	if err != nil {
		return "", err
	} else {
		return val[key], nil
	}
}

//将所有指定的值插入到存于 key 的列表的头部
func (rc *RedisClient) LPush(key string, val string) (int64, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	ret, err := redis.Int64(conn.Do("LPUSH", key, val))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

func (rc *RedisClient) Set(key string, val interface{}) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := conn.Do("SET", key, val)
	return val, err
}

//设置指定key的内容
func (rc *RedisClient) SetWithExpire(key string, val interface{}, timeOutSeconds int64) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := conn.Do("SET", key, val, "EX", timeOutSeconds)
	return val, err
}

//设置指定key的过期时间
func (rc *RedisClient) Expire(key string, timeOutSeconds int64) (int64, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("EXPIRE", key, timeOutSeconds))
	return val, err
}

//删除当前数据库里面的所有数据
//这个命令永远不会出现失败
func (rc *RedisClient) FlushDB() {
	conn := rc.pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")
}
