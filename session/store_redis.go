package session

import (
	"github.com/devfeel/dotweb/framework/encodes/gob"
	"github.com/devfeel/dotweb/framework/redis"
	"sync"
	"fmt"
)

const (
	defaultRedisKeyPre = "dotweb:session:"
)

// RedisStore Implement the SessionStore interface
type RedisStore struct {
	lock        *sync.RWMutex // locker
	maxlifetime int64
	serverIp    string //connection string, like "redis://:password@10.0.1.11:6379/0"
	storeKeyPre string //set custom redis key-pre; default is dotweb:session:
}

//create new redis store
func NewRedisStore(config *StoreConfig) (*RedisStore, error){
	store := &RedisStore{
		lock:        new(sync.RWMutex),
		serverIp:    config.ServerIP,
		maxlifetime: config.Maxlifetime,
	}
	//init redis key-pre
	if config.StoreKeyPre == ""{
		store.storeKeyPre = defaultRedisKeyPre
	}else{
		store.storeKeyPre = config.StoreKeyPre
	}
	redisClient := redisutil.GetRedisClient(store.serverIp)
	_, err:=redisClient.Ping()
	return store, err
}

func (store *RedisStore)getRedisKey(key string) string {
	return store.storeKeyPre + key
}

// SessionRead get session state by sessionId
func (store *RedisStore) SessionRead(sessionId string) (*SessionState, error) {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := store.getRedisKey(sessionId)
	kvs, err := redisClient.Get(key)
	if err != nil {
		return nil, err
	}
	var kv map[interface{}]interface{}
	if len(kvs) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = gob.DecodeMap([]byte(kvs))
		if err != nil {
			return nil, err
		}
	}
	state := NewSessionState(store, sessionId, kv)
	go store.sessionReExpire(state)
	return state, nil
}

// SessionExist check session state exist by sessionId
func (store *RedisStore) SessionExist(sessionId string) bool {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := store.getRedisKey(sessionId)
	exists, err := redisClient.Exists(key)
	if err != nil {
		return false
	}
	return exists
}

// sessionReExpire reset expire session key
func (store *RedisStore) sessionReExpire(state *SessionState) error {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := store.getRedisKey(state.SessionID())
	_, err := redisClient.Expire(key, store.maxlifetime)
	return err
}

//SessionUpdate update session state in store
func (store *RedisStore) SessionUpdate(state *SessionState) error {
	defer func(){
		//ignore error
		if err := recover(); err != nil {
			fmt.Println("SessionUpdate-Redis error", err)
			//TODO deal panic err
		}
	}()
	redisClient := redisutil.GetRedisClient(store.serverIp)
	bytes, err := gob.EncodeMap(state.values)
	if err != nil {
		return err
	}
	key := store.getRedisKey(state.SessionID())
	_, err = redisClient.SetWithExpire(key, string(bytes), store.maxlifetime)
	return err
}

// SessionRemove delete session state in store
func (store *RedisStore) SessionRemove(sessionId string) error {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := store.getRedisKey(sessionId)
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
