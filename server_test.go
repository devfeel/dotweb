package dotweb

import (
	"github.com/devfeel/dotweb/session"
	"github.com/devfeel/dotweb/test"
	"testing"
)

//check httpServer
func TestNewHttpServer(t *testing.T) {
	server := NewHttpServer()

	test.NotNil(t, server.router)
	test.NotNil(t, server.stdServer)
	test.NotNil(t, server.ServerConfig)
	test.NotNil(t, server.SessionConfig)
	test.NotNil(t, server.lock_session)
	test.NotNil(t, server.binder)
	test.NotNil(t, server.Features)
	test.NotNil(t, server.pool)
	test.NotNil(t, server.pool.context)
	test.NotNil(t, server.pool.request)
	test.NotNil(t, server.pool.response)
	test.Equal(t, false, server.IsOffline())

	//t.Log("is offline:",server.IsOffline())
}

//session manager用来设置gc？
//总感觉和名字不是太匹配
func TestSesionConfig(t *testing.T) {
	server := NewHttpServer()
	//use default config
	server.SetSessionConfig(session.NewDefaultRuntimeConfig())

	//init
	server.InitSessionManager()

	//get session
	sessionManager := server.GetSessionManager()

	//EnabledSession flag is false
	test.Nil(t, sessionManager)

	//switch EnabledSession flag
	server.SessionConfig().EnabledSession = true
	sessionManager = server.GetSessionManager()

	test.NotNil(t, sessionManager)
	test.Equal(t, server.sessionManager.CookieName, session.DefaultSessionCookieName)
	test.Equal(t, server.sessionManager.GCLifetime, int64(session.DefaultSessionGCLifeTime))
}

func Index(ctx Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := ctx.WriteStringC(201, "index => ", ctx.RemoteIP(), "我是首页")
	return err
}
