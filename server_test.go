package dotweb

import (
	"testing"

	"github.com/devfeel/dotweb/session"
	"github.com/devfeel/dotweb/test"
)

// check httpServer
func TestNewHttpServer(t *testing.T) {
	server := NewHttpServer()

	test.NotNil(t, server.router)
	test.NotNil(t, server.stdServer)
	test.NotNil(t, server.ServerConfig)
	test.NotNil(t, server.SessionConfig)
	test.NotNil(t, server.lock_session)
	test.NotNil(t, server.binder)
	test.NotNil(t, server.pool)
	test.NotNil(t, server.pool.context)
	test.NotNil(t, server.pool.request)
	test.NotNil(t, server.pool.response)
	test.Equal(t, false, server.IsOffline())

	// t.Log("is offline:",server.IsOffline())
}

func TestSesionConfig(t *testing.T) {
	server := NewHttpServer()
	server.DotApp = New()
	// use default config
	server.SetSessionConfig(session.NewDefaultRuntimeConfig())

	// init
	server.InitSessionManager()

	// get session
	sessionManager := server.GetSessionManager()

	// EnabledSession flag is false
	test.Nil(t, sessionManager)

	// switch EnabledSession flag
	server.SessionConfig().EnabledSession = true
	sessionManager = server.GetSessionManager()

	test.NotNil(t, sessionManager)
	test.Equal(t, server.sessionManager.StoreConfig().CookieName, session.DefaultSessionCookieName)
	test.Equal(t, server.sessionManager.GCLifetime, int64(session.DefaultSessionGCLifeTime))
}

func Index(ctx Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	err := ctx.WriteStringC(201, "index => ", ctx.RemoteIP(), "我是首页")
	return err
}
