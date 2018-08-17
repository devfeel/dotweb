package session

import (
	"github.com/devfeel/dotweb/framework/encodes/gob"
	"github.com/devfeel/dotweb/framework/redis"
	"sync"
	"fmt"
	"github.com/devfeel/dotweb/framework/hystrix"
	"strings"
)

const (
	defaultRedisKeyPre = "dotweb:session:"
	HystrixErrorCount = 20
)

// RedisStore Implement the SessionStore interface
type RedisStore struct {
	hystrix hystrix.Hystrix
	lock        *sync.RWMutex // locker
	maxlifetime int64
	serverIp    string //connection string, like "redis://:password@10.0.1.11:6379/0"
	backupServerUrl string //backup connection string, like "redis://:password@10.0.1.11:6379/0"
	storeKeyPre string //set custom redis key-pre; default is dotweb:session:
}

//create new redis store
func NewRedisStore(config *StoreConfig) (*RedisStore, error){
	store := &RedisStore{
		lock:        	new(sync.RWMutex),
		serverIp:    	config.ServerIP,
		backupServerUrl:config.BackupServerUrl,
		maxlifetime: 	config.Maxlifetime,
	}
	store.hystrix = hystrix.NewHystrix(store.checkRedisAlive, nil)
	store.hystrix.SetMaxFailedNumber(HystrixErrorCount)
	store.hystrix.Do()
	//init redis key-pre
	if config.StoreKeyPre == ""{
		store.storeKeyPre = defaultRedisKeyPre
	}else{
		store.storeKeyPre = config.StoreKeyPre
	}
	redisClient := store.getRedisClient()
	_, err:=redisClient.Ping()
	if store.checkConnErrorAndNeedRetry(err){
		store.hystrix.TriggerHystrix()
		redisClient = store.getBackupRedis()
		_, err = redisClient.Ping()
	}
	return store, err
}

func (store *RedisStore)getRedisKey(key string) string {
	return store.storeKeyPre + key
}

// SessionRead get session state by sessionId
func (store *RedisStore) SessionRead(sessionId string) (*SessionState, error) {
	redisClient := store.getRedisClient()
	key := store.getRedisKey(sessionId)
	kvs, err := redisClient.Get(key)
	if store.checkConnErrorAndNeedRetry(err){
		redisClient = store.getBackupRedis()
		kvs, err = redisClient.Get(key)
	}
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
	redisClient := store.getRedisClient()
	key := store.getRedisKey(sessionId)
	exists, err := redisClient.Exists(key)
	if store.checkConnErrorAndNeedRetry(err){
		redisClient = store.getBackupRedis()
		exists, err = redisClient.Exists(key)
	}
	if err != nil {
		return false
	}
	return exists
}

// sessionReExpire reset expire session key
func (store *RedisStore) sessionReExpire(state *SessionState) error {
	redisClient := store.getRedisClient()
	key := store.getRedisKey(state.SessionID())
	_, err := redisClient.Expire(key, store.maxlifetime)
	if store.checkConnErrorAndNeedRetry(err){
		redisClient = store.getBackupRedis()
		_, err = redisClient.Expire(key, store.maxlifetime)
	}
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
	redisClient := store.getRedisClient()
	bytes, err := gob.EncodeMap(state.values)
	if err != nil {
		return err
	}
	key := store.getRedisKey(state.SessionID())
	_, err = redisClient.SetWithExpire(key, string(bytes), store.maxlifetime)
	if store.checkConnErrorAndNeedRetry(err){
		redisClient = store.getBackupRedis()
		_, err = redisClient.SetWithExpire(key, string(bytes), store.maxlifetime)
	}
	return err
}

// SessionRemove delete session state in store
func (store *RedisStore) SessionRemove(sessionId string) error {
	redisClient := redisutil.GetRedisClient(store.serverIp)
	key := store.getRedisKey(sessionId)
	_, err := redisClient.Del(key)
	if store.checkConnErrorAndNeedRetry(err){
		redisClient = store.getBackupRedis()
		_, err = redisClient.Del(key)
	}
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


// getRedisClient get alive redis client
func (store *RedisStore) getRedisClient() *redisutil.RedisClient{
	if store.hystrix.IsHystrix(){
		if store.backupServerUrl != "" {
			return store.getBackupRedis()
		}
	}
	return store.getDefaultRedis()
}

func (store *RedisStore) getDefaultRedis() *redisutil.RedisClient{
	return redisutil.GetRedisClient(store.serverIp)
}

func (store *RedisStore) getBackupRedis() *redisutil.RedisClient{
	return redisutil.GetRedisClient(store.backupServerUrl)
}

// checkConnErrorAndNeedRetry check err is Conn error and is need to retry
func (store *RedisStore) checkConnErrorAndNeedRetry(err error) bool{
	if err == nil{
		return false
	}
	if strings.Index(err.Error(), "no such host") >= 0 ||
		strings.Index(err.Error(), "No connection could be made because the target machine actively refused it") >= 0 ||
		strings.Index(err.Error(), "A connection attempt failed because the connected party did not properly respond after a period of time") >= 0 {
		store.hystrix.GetCounter().Inc(1)
		//if is hystrix, not to retry, because in getReadRedisClient already use backUp redis
		if store.hystrix.IsHystrix(){
			return false
		}
		if store.backupServerUrl == ""{
			return false
		}
		return true
	}
	return false
}

// checkRedisAlive check redis is alive use ping
// if set readonly redis, check readonly redis
// if not set readonly redis, check default redis
func (store *RedisStore) checkRedisAlive() bool{
	isAlive := false
	var redisClient *redisutil.RedisClient
	redisClient = store.getDefaultRedis()
	for i := 0;i<=5;i++ {
		reply, err := redisClient.Ping()
		if err != nil {
			isAlive = false
			break
		}
		if reply != "PONG" {
			isAlive = false
			break
		}
		isAlive = true
		continue
	}
	return isAlive
}