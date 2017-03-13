package session

import (
	"github.com/devfeel/dotweb/framework/encodes/gob"
	"github.com/devfeel/dotweb/framework/redis"
	"sync"
)

const (
	defaultRedisKeyPre = "dotweb:session:"
)

// MemProvider Implement the provider interface
type RedisStore struct {
	lock        *sync.RWMutex // locker
	maxlifetime int64
	serverIp    string
	password    string
}

func getRedisKey(key string) string {
	return defaultRedisKeyPre + key
}

//create new redis store
func NewRedisStore(config *StoreConfig) *RedisStore {
	return &RedisStore{
		lock:        new(sync.RWMutex),
		serverIp:    config.StoreServerIP,
		password:    config.StorePassword,
		maxlifetime: config.Maxlifetime,
	}
}

// SessionRead get session state by sessionId
func (store *RedisStore) SessionRead(sessionId string) (*SessionState, error) {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := getRedisKey(sessionId)
	kvs, err := redisClient.Get(key)
	var kv map[interface{}]interface{}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = gob.DecodeGob([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}
	state := NewSessionState(store, sessionId, kv)
	go store.SessionUpdate(state)
	return state, nil
}

// SessionExist check session state exist by sessionId
func (store *RedisStore) SessionExist(sessionId string) bool {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := getRedisKey(sessionId)
	exists, err := redisClient.Exists(key)
	if err != nil {
		return false
	}
	return exists
}

//SessionUpdate update session state in store
func (store *RedisStore) SessionUpdate(state *SessionState) bool {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	bytes, err := gob.EncodeGob(state.values)
	if err != nil {
		return false
	}
	key := getRedisKey(state.SessionID())
	_, err = redisClient.SetWithExpire(key, string(bytes), store.maxlifetime)
	return err == nil
}

// SessionRemove delete session state in store
func (store *RedisStore) SessionRemove(sessionId string) error {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := getRedisKey(sessionId)
	_, err := redisClient.Del(key)
	return err
}

// SessionGC clean expired session states
// in redis store,not use
func (store *RedisStore) SessionGC() int {
	return 0
}

// SessionAll get count number
func (store *RedisStore) SessionCount() int {
	return 0
}
