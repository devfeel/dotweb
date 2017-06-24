package session

import (
	"testing"
	"github.com/devfeel/dotweb/test"
)

const (
	IP="0.0.0.0"
)

func TestGetSessionStore(t *testing.T) {
	defaultConfig:=NewDefaultRuntimeConfig()

	defaultSessionStore:=GetSessionStore(defaultConfig)

	test.Equal(t,SessionMode_Runtime,defaultConfig.StoreName)
	test.Equal(t,int64(DefaultSessionMaxLifeTime),defaultConfig.Maxlifetime)
	test.Equal(t,"",defaultConfig.ServerIP)

	test.NotNil(t,defaultSessionStore)

	defaultRedisConfig:=NewDefaultRedisConfig(IP)

	defaultRedisSessionStore:=GetSessionStore(defaultRedisConfig)

	test.Equal(t,SessionMode_Redis,defaultRedisConfig.StoreName)
	test.Equal(t,int64(DefaultSessionMaxLifeTime),defaultRedisConfig.Maxlifetime)
	test.Equal(t,IP,defaultRedisConfig.ServerIP)

	test.NotNil(t,defaultRedisSessionStore)
}

func TestNewDefaultSessionManager(t *testing.T) {
	defaultRedisConfig:=NewDefaultRedisConfig(IP)
	manager,err:=NewDefaultSessionManager(defaultRedisConfig)

	test.Nil(t,err)
	test.NotNil(t, manager)


	test.NotNil(t, manager.store)
	test.Equal(t,int64(DefaultSessionGCLifeTime),manager.GCLifetime)
	test.Equal(t,DefaultSessionCookieName,manager.CookieName)
	test.Equal(t,defaultRedisConfig,manager.storeConfig)


	sessionId:=manager.NewSessionID()

	test.Equal(t,32,len(sessionId))

	sessionState,err:=manager.GetSessionState(sessionId)
	test.Nil(t,err)
	test.NotNil(t, sessionState)
	test.Equal(t,sessionId,sessionState.sessionId)
}

