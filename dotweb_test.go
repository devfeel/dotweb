package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
)

// 以下为功能测试

// 测试RunMode函数无配置文件时的返回值
func Test_RunMode_1(t *testing.T) {
	app := New()
	runMode := app.RunMode()
	t.Log("RunMode:", runMode)
}

// 测试RunMode函数有配置文件时的返回值
func Test_RunMode_2(t *testing.T) {
	runModes := []string{"dev", "development", "prod", "production"}

	app := New()
	for _, value := range runModes {
		app.Config.App.RunMode = value
		runMode := app.RunMode()
		t.Log("runModes value:", value, "RunMode:", runMode)
	}
}

//测试IsDevelopmentMode函数
func Test_IsDevelopmentMode_1(t *testing.T) {
	app := New()
	app.Config.App.RunMode = "development"
	b := app.IsDevelopmentMode()
	t.Log("Run IsDevelopmentMode :", b)
}

func Test_IsDevelopmentMode_2(t *testing.T) {
	app := New()
	app.Config.App.RunMode = "production"
	b := app.IsDevelopmentMode()
	t.Log("Run IsDevelopmentMode :", b)
}

//check httpServer
func TestNewHttpServer(t *testing.T) {
	server:=NewHttpServer()

	test.NotNil(t,server.router)
	test.NotNil(t,server.stdServer)
	test.NotNil(t,server.ServerConfig)
	test.NotNil(t,server.SessionConfig)
	test.NotNil(t,server.lock_session)
	test.NotNil(t,server.binder)
	test.NotNil(t,server.Features)
	test.NotNil(t,server.pool)
	test.NotNil(t,server.pool.context)
	test.NotNil(t,server.pool.request)
	test.NotNil(t,server.pool.response)
}
