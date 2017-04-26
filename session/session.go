package session

import (
	"github.com/devfeel/dotweb/framework/crypto"
	"github.com/devfeel/dotweb/logger"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	DefaultSessionGCLifeTime  = 60      //second
	DefaultSessionMaxLifeTime = 20 * 60 //second
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
		SessionCount() int //get all active session length
		SessionGC() int    //gc session and return out of date state num
	}

	//session config info
	StoreConfig struct {
		StoreName   string
		Maxlifetime int64
		ServerIP    string
	}

	SessionManager struct {
		store       SessionStore
		CookieName  string `json:"cookieName"`
		GCLifetime  int64  `json:"gclifetime"`
		storeConfig *StoreConfig
	}
)

//create new session store with store config
func GetSessionStore(config *StoreConfig) SessionStore {
	switch config.StoreName {
	case SessionMode_Runtime:
		return NewRuntimeStore(config)
	case SessionMode_Redis:
		return NewRedisStore(config)
	default:
		panic("not support session store -> " + config.StoreName)
	}
	return nil

}

//create new store with default config and use runtime store
func NewDefaultRuntimeConfig() *StoreConfig {
	return NewStoreConfig(SessionMode_Runtime, DefaultSessionMaxLifeTime, "")
}

//create new store with default config and use redis store
func NewDefaultRedisConfig(serverIp string) *StoreConfig {
	return NewStoreConfig(SessionMode_Redis, DefaultSessionMaxLifeTime, serverIp)
}

//create new store config
func NewStoreConfig(storeName string, maxlifetime int64, serverIp string) *StoreConfig {
	return &StoreConfig{
		StoreName:   storeName,
		Maxlifetime: maxlifetime,
		ServerIP:    serverIp,
	}
}

//create new session manager with default config info
func NewDefaultSessionManager(config *StoreConfig) (*SessionManager, error) {
	return NewSessionManager(DefaultSessionCookieName, DefaultSessionGCLifeTime, config)
}

//create new seesion manager
func NewSessionManager(cookieName string, gcLifetime int64, config *StoreConfig) (*SessionManager, error) {
	if gcLifetime <= 0 {
		gcLifetime = DefaultSessionGCLifeTime
	}
	if cookieName == "" {
		cookieName = DefaultSessionCookieName
	}
	manager := &SessionManager{
		store:       GetSessionStore(config),
		GCLifetime:  gcLifetime,
		storeConfig: config,
		CookieName:  cookieName,
	}
	//开启GC
	go func() {
		time.AfterFunc(time.Duration(manager.GCLifetime)*time.Second, func() { manager.GC() })
	}()
	return manager, nil
}

//create new session id with DefaultSessionLength
func (manager *SessionManager) NewSessionID() string {
	val := cryptos.GetRandString(DefaultSessionLength)
	return val
}

//get session id from client
//default mode is from cookie
func (manager *SessionManager) GetClientSessionID(req *http.Request) (string, error) {
	cookie, err := req.Cookie(manager.CookieName)
	if err != nil {
		return "", err
	}
	if cookie.Value == "" {
		return "", nil
	}
	//TODO: check client validity
	//check ip & agent
	return url.QueryUnescape(cookie.Value)
}

func (manager *SessionManager) GetSessionState(sessionId string) (session *SessionState, err error) {
	session, err = manager.store.SessionRead(sessionId)
	if err != nil {
		session = NewSessionState(manager.store, sessionId, make(map[interface{}]interface{}))
	}
	return session, nil
}

//GC loop gc session data
func (manager *SessionManager) GC() {
	num := manager.store.SessionGC()
	if num > 0 {
		logger.Logger().Debug("SessionManger.GC => "+strconv.Itoa(num), LogTarget_Session)
	}
	time.AfterFunc(time.Duration(manager.GCLifetime)*time.Second, func() { manager.GC() })
}
