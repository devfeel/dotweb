package session

import (
	"fmt"
	"github.com/devfeel/dotweb/logger"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/devfeel/dotweb/framework/crypto"
)

const (
	DefaultSessionGCLifeTime  = 60      // second
	DefaultSessionMaxLifeTime = 20 * 60 // second
	DefaultSessionCookieName  = "dotweb_sessionId"
	DefaultSessionLength      = 20
	SessionMode_Runtime       = "runtime"
	SessionMode_Redis         = "redis"

	LogTarget_Session = "dotweb_session"
)

type (
	SessionStore interface {
		SessionRead(sessionId string) (*SessionState, error)
		SessionExist(sessionId string) bool
		SessionUpdate(state *SessionState) error
		SessionRemove(sessionId string) error
		SessionCount() int // get all active session length
		SessionGC() int    // gc session and return out of date state num
	}

	// session config info
	StoreConfig struct {
		StoreName       string
		Maxlifetime     int64  // session life time, with second
		CookieName      string // custom cookie name which sessionid store
		ServerIP        string // if use redis, connection string, like "redis://:password@10.0.1.11:6379/0"
		BackupServerUrl string // if use redis, if ServerIP is down, use this server, like "redis://:password@10.0.1.11:6379/0"
		StoreKeyPre     string // if use redis, set custom redis key-pre; default is dotweb:session:
		MaxIdle         int    // if use redis, set MaxIdle; default is 10
		MaxActive       int    // if use redis, set MaxActive; default is 50
	}

	SessionManager struct {
		GCLifetime int64 `json:"gclifetime"`

		appLog      logger.AppLog
		store       SessionStore
		storeConfig *StoreConfig
	}
)

// GetSessionStore create new session store with store config
func GetSessionStore(config *StoreConfig) SessionStore {
	switch config.StoreName {
	case SessionMode_Runtime:
		return NewRuntimeStore(config)
	case SessionMode_Redis:
		store, err := NewRedisStore(config)
		if err != nil {
			panic(fmt.Sprintf("redis session [%v] ping error -> %v", config.StoreName, err.Error()))
		} else {
			return store
		}
	default:
		panic("not support session store -> " + config.StoreName)
	}
	return nil
}

// NewDefaultRuntimeConfig create new store with default config and use runtime store
func NewDefaultRuntimeConfig() *StoreConfig {
	return NewStoreConfig(SessionMode_Runtime, DefaultSessionMaxLifeTime, "", "", 0, 0)
}

// NewDefaultRedisConfig create new store with default config and use redis store
func NewDefaultRedisConfig(serverIp string) *StoreConfig {
	return NewRedisConfig(serverIp, DefaultSessionMaxLifeTime, "", 0, 0)
}

// NewRedisConfig create new store with config and use redis store
// must set serverIp and storeKeyPre
func NewRedisConfig(serverIp string, maxlifetime int64, storeKeyPre string, maxIdle int, maxActive int) *StoreConfig {
	return NewStoreConfig(SessionMode_Redis, maxlifetime, serverIp, storeKeyPre, maxIdle, maxActive)
}

// NewStoreConfig create new store config
func NewStoreConfig(storeName string, maxlifetime int64, serverIp string, storeKeyPre string, maxIdle int, maxActive int) *StoreConfig {
	return &StoreConfig{
		StoreName:   storeName,
		Maxlifetime: maxlifetime,
		ServerIP:    serverIp,
		StoreKeyPre: storeKeyPre,
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
	}
}

// NewDefaultSessionManager create new session manager with default config info
func NewDefaultSessionManager(appLog logger.AppLog, config *StoreConfig) (*SessionManager, error) {
	return NewSessionManager(DefaultSessionGCLifeTime, appLog, config)
}

// NewSessionManager create new seesion manager
func NewSessionManager(gcLifetime int64, appLog logger.AppLog, config *StoreConfig) (*SessionManager, error) {
	if gcLifetime <= 0 {
		gcLifetime = DefaultSessionGCLifeTime
	}
	if config.CookieName == "" {
		config.CookieName = DefaultSessionCookieName
	}
	manager := &SessionManager{
		store:       GetSessionStore(config),
		appLog:      appLog,
		GCLifetime:  gcLifetime,
		storeConfig: config,
	}
	// enable GC
	go func() {
		time.AfterFunc(time.Duration(manager.GCLifetime)*time.Second, func() { manager.GC() })
	}()
	return manager, nil
}

// NewSessionID create new session id with DefaultSessionLength
func (manager *SessionManager) NewSessionID() string {
	val := cryptos.GetRandString(DefaultSessionLength)
	return val
}

// StoreConfig return store config
func (manager *SessionManager) StoreConfig() *StoreConfig {
	return manager.storeConfig
}

// GetClientSessionID get session id from client
// default mode is from cookie
func (manager *SessionManager) GetClientSessionID(req *http.Request) (string, error) {
	cookie, err := req.Cookie(manager.storeConfig.CookieName)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", nil
	}
	// TODO: check client validity
	// check ip & agent
	return url.QueryUnescape(cookie.Value)
}

func (manager *SessionManager) GetSessionState(sessionId string) (session *SessionState, err error) {
	session, err = manager.store.SessionRead(sessionId)
	if err != nil {
		session = NewSessionState(manager.store, sessionId, make(map[interface{}]interface{}))
	}
	return session, nil
}

// RemoveSessionState delete the session state associated with a specific session ID
func (manager *SessionManager) RemoveSessionState(sessionId string) error {
	return manager.store.SessionRemove(sessionId)
}

// GC loop gc session data
func (manager *SessionManager) GC() {
	num := manager.store.SessionGC()
	if num > 0 {
		manager.appLog.Debug("SessionManger.GC => "+strconv.Itoa(num), LogTarget_Session)
	}
	time.AfterFunc(time.Duration(manager.GCLifetime)*time.Second, func() { manager.GC() })
}
