// redisclient

// Package redisutil 命令的使用方式参考
// http://doc.redisfans.com/index.html
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
// redisURL: connection string, like "redis://:password@10.0.1.11:6379/0"
func newPool(redisURL string) *redis.Pool {

	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 20, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisURL)
			return c, err
		},
	}
}

// GetRedisClient 获取指定Address的RedisClient
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

// GetObj 获取指定key的内容, interface{}
func (rc *RedisClient) GetObj(key string) (interface{}, error) {
	// 从连接池里面获得一个连接
	conn := rc.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer conn.Close()
	reply, errDo := conn.Do("GET", key)
	return reply, errDo
}

// Get 获取指定key的内容, string
func (rc *RedisClient) Get(key string) (string, error) {
	val, err := redis.String(rc.GetObj(key))
	return val, err
}

// Exists 检查指定key是否存在
func (rc *RedisClient) Exists(key string) (bool, error) {
	// 从连接池里面获得一个连接
	conn := rc.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer conn.Close()

	reply, errDo := redis.Bool(conn.Do("EXISTS", key))
	return reply, errDo
}

// Del 删除指定key
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

// INCR 对存储在指定key的数值执行原子的加1操作
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

// DECR 对存储在指定key的数值执行原子的减1操作
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


// Append 如果 key 已经存在并且是一个字符串， APPEND 命令将 value 追加到 key 原来的值的末尾。
// 如果 key 不存在， APPEND 就简单地将给定 key 设为 value ，就像执行 SET key value 一样。
func (rc *RedisClient) Append(key string, val interface{}) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, errDo := conn.Do("APPEND", key, val)
	if errDo == nil && reply == nil {
		return 0, nil
	}
	val, err := redis.Uint64(reply, errDo)
	return val, err
}

// Set 设置指定Key/Value
func (rc *RedisClient) Set(key string, val interface{}) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("SET", key, val))
	return val, err
}


// Expire 设置指定key的过期时间
func (rc *RedisClient) Expire(key string, timeOutSeconds int64) (int64, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int64(conn.Do("EXPIRE", key, timeOutSeconds))
	return val, err
}

// SetWithExpire 设置指定key的内容
func (rc *RedisClient) SetWithExpire(key string, val interface{}, timeOutSeconds int64) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("SET", key, val, "EX", timeOutSeconds))
	return val, err
}

// SetNX  将 key 的值设为 value ，当且仅当 key 不存在。
// 若给定的 key 已经存在，则 SETNX 不做任何动作。 成功返回1, 失败返回0
func (rc *RedisClient) SetNX(key, value string) (interface{}, error){
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := conn.Do("SETNX", key, value)
	return val, err
}



//****************** hash 集合 ***********************

// HGet 获取指定hash的内容
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

// HGetAll 获取指定hash的所有内容
func (rc *RedisClient) HGetAll(hashID string) (map[string]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	reply, err := redis.StringMap(conn.Do("HGetAll", hashID))
	return reply, err
}

// HSet 设置指定hash的内容
func (rc *RedisClient) HSet(hashID string, field string, val string) error {
	conn := rc.pool.Get()
	defer conn.Close()
	_, err := conn.Do("HSET", hashID, field, val)
	return err
}

// HSetNX 设置指定hash的内容, 如果field不存在, 该操作无效
func (rc *RedisClient) HSetNX(hashID, field, value string) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := conn.Do("HSETNX", hashID, field, value)
	return val, err
}

// HExist 返回hash里面field是否存在
func (rc *RedisClient) HExist(hashID string, field string) (int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("HEXISTS", hashID, field))
	return val, err
}

// HIncrBy 增加 key 指定的哈希集中指定字段的数值
func (rc *RedisClient) HIncrBy(hashID string, field string, increment int)(int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("HINCRBY", hashID, field, increment))
	return val, err
}

// HLen 返回哈希表 key 中域的数量, 当 key 不存在时，返回0
func (rc *RedisClient) HLen(hashID string) (int64, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("HLEN", hashID))
	return val, err
}

// HDel 设置指定hashset的内容, 如果field不存在, 该操作无效, 返回0
func (rc *RedisClient) HDel(args ...interface{}) (int64, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("HDEL", args...))
	return val, err
}

// HVals 返回哈希表 key 中所有域的值, 当 key 不存在时，返回空
func (rc *RedisClient) HVals(hashID string) (interface{}, error) {
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := redis.Strings(conn.Do("HVALS", hashID))
	return val, err
}

//****************** list ***********************

//将所有指定的值插入到存于 key 的列表的头部
func (rc *RedisClient) LPush(key string, value ...interface{}) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	ret, err := redis.Int(conn.Do("LPUSH", key, value))
	if err != nil {
		return -1, err
	} else {
		return ret, nil
	}
}

func (rc *RedisClient) LPushX(key string, value string) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.Int(conn.Do("LPUSHX", key, value))
	return resp, err
}

func (rc *RedisClient) LRange(key string, start int, stop int) ([]string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.Strings(conn.Do("LRANGE", key, start, stop))
	return resp, err
}

func (rc *RedisClient) LRem(key string, count int, value string) (int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.Int(conn.Do("LREM", key, count, value))
	return resp, err
}

func (rc *RedisClient) LSet(key string, index int, value string)(string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("LSET", key, index, value))
	return resp, err
}

func (rc *RedisClient) LTrim(key string, start int, stop int) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("LTRIM", key, start, stop))
	return resp, err
}

func (rc *RedisClient) RPop(key string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("RPOP", key))
	return resp, err
}

func (rc *RedisClient) RPush(key string, value ...interface{}) (int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{key}, value...)
	resp, err := redis.Int(conn.Do("RPUSH", args...))
	return resp, err
}

func (rc *RedisClient) RPushX(key string, value ...interface{}) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{key}, value...)
	resp, err := redis.Int(conn.Do("RPUSHX", args...))
	return resp, err
}

func (rc *RedisClient) RPopLPush(source string, destination string)(string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	resp, err := redis.String(conn.Do("RPOPLPUSH", source, destination))
	return resp, err
}


func (rc *RedisClient) BLPop(key ...interface{})(map[string]string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.StringMap(conn.Do("BLPOP", key, defaultTimeout))
	return val, err
}

//删除，并获得该列表中的最后一个元素，或阻塞，直到有一个可用
func (rc *RedisClient) BRPop(key ...interface{}) (map[string]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.StringMap(conn.Do("BRPOP", key, defaultTimeout))
	return val, err
}

func (rc *RedisClient) BRPopLPush(source string, destination string)(string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("BRPOPLPUSH", source, destination))
	return val, err
}

func (rc *RedisClient) LIndex(key string, index int)(string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("LINDEX", key, index))
	return val, err
}

func (rc *RedisClient) LInsertBefore(key string, pivot string, value string)(int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("LINSERT", key, "BEFORE", pivot, value))
	return val, err
}

func (rc *RedisClient) LInsertAfter(key string, pivot string, value string)(int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("LINSERT", key, "AFTER", pivot, value))
	return val, err
}

func (rc *RedisClient) LLen(key string)(int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("LLEN", key))
	return val, err
}

func (rc *RedisClient) LPop(key string)(string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("LPOP", key))
	return val, err
}

//****************** set 集合 ***********************

// SAdd 将一个或多个 member 元素加入到集合 key 当中，已经存在于集合的 member 元素将被忽略。
// 假如 key 不存在，则创建一个只包含 member 元素作成员的集合。
func (rc *RedisClient) SAdd(key string, member ...interface{}) (int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{key}, member...)
	val, err := redis.Int(conn.Do("SADD", args...))
	return val, err
}

// SCard 返回集合 key 的基数(集合中元素的数量)。
// 返回值：
// 集合的基数。
// 当 key 不存在时，返回 0
func (rc *RedisClient) SCard(key string) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Int(conn.Do("SCARD", key))
	return val, err
}

// SPop 移除并返回集合中的一个随机元素。
// 如果只想获取一个随机元素，但不想该元素从集合中被移除的话，可以使用 SRANDMEMBER 命令。
// count 为 返回的随机元素的数量
func (rc *RedisClient) SPop(key string) (string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.String(conn.Do("SPOP", key))
	return val, err
}

// SRandMember 如果命令执行时，只提供了 key 参数，那么返回集合中的一个随机元素。
// 该操作和 SPOP 相似，但 SPOP 将随机元素从集合中移除并返回，而 SRANDMEMBER 则仅仅返回随机元素，而不对集合进行任何改动。
// count 为 返回的随机元素的数量
func (rc *RedisClient) SRandMember(key string, count int) ([]string, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Strings(conn.Do("SRANDMEMBER", key, count))
	return val, err
}

// SRem 移除集合 key 中的一个或多个 member 元素，不存在的 member 元素会被忽略。
// 当 key 不是集合类型，返回一个错误。
// 在 Redis 2.4 版本以前， SREM 只接受单个 member 值。
func (rc *RedisClient) SRem(key string, member ...interface{}) (int, error) {
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{key}, member...)
	val, err := redis.Int(conn.Do("SREM", args...))
	return val, err
}

func (rc *RedisClient) SDiff(key ...interface{}) ([]string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Strings(conn.Do("SDIFF", key...))
	return val, err
}

func (rc *RedisClient) SDiffStore(destination string, key ...interface{}) (int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{destination}, key...)
	val, err := redis.Int(conn.Do("SDIFFSTORE", args...))
	return val, err
}

func (rc *RedisClient) SInter(key ...interface{}) ([]string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Strings(conn.Do("SINTER", key...))
	return val, err
}

func (rc *RedisClient) SInterStore(destination string, key ...interface{})(int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{destination}, key...)
	val, err := redis.Int(conn.Do("SINTERSTORE", args...))
	return val, err
}

func (rc *RedisClient) SIsMember(key string, member string) (bool, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Bool(conn.Do("SISMEMBER", key, member))
	return val, err
}

func (rc *RedisClient) SMembers(key string) ([]string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Strings(conn.Do("SMEMBERS", key))
	return val, err
}

// smove is a atomic operate
func (rc *RedisClient) SMove(source string, destination string, member string) (bool, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Bool(conn.Do("SMOVE", source, destination, member))
	return val, err
}

func (rc *RedisClient) SUnion(key ...interface{}) ([]string, error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err := redis.Strings(conn.Do("SUNION", key...))
	return val, err
}

func (rc *RedisClient) SUnionStore(destination string, key ...interface{})(int, error){
	conn := rc.pool.Get()
	defer conn.Close()
	args := append([]interface{}{destination}, key...)
	val, err := redis.Int(conn.Do("SUNIONSTORE", args))
	return val, err
}

//****************** 全局操作 ***********************

// Ping 测试一个连接是否可用
func (rc *RedisClient) Ping()(string,error){
	conn := rc.pool.Get()
	defer conn.Close()
	val, err:=redis.String(conn.Do("PING"))
	return val, err
}

// DBSize 返回当前数据库的 key 的数量
func (rc *RedisClient) DBSize()(int64, error){
	conn := rc.pool.Get()
	defer conn.Close()

	val, err := redis.Int64(conn.Do("DBSIZE"))
	return val, err
}

// FlushDB 删除当前数据库里面的所有数据
// 这个命令永远不会出现失败
func (rc *RedisClient) FlushDB() {
	conn := rc.pool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")
}


// GetConn 返回一个从连接池获取的redis连接,
// 需要手动释放redis连接
func (rc *RedisClient) GetConn() redis.Conn{
	return rc.pool.Get()
}
