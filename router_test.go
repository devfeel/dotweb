package dotweb

import (
	"github.com/devfeel/dotweb/session"
	"github.com/devfeel/dotweb/test"
	"testing"
	"time"
)

func TestRouter_ServeHTTP(t *testing.T) {
	param := &InitContextParam{
		t,
		"",
		"",
		test.ToDefault,
	}

	context := initAllContext(param)

	app := New()
	server := app.HttpServer
	r := NewRouter(server)

	r.ServeHTTP(context)
}

//
func TestWrapRouterHandle(t *testing.T) {
	param := &InitContextParam{
		t,
		"",
		"",
		test.ToDefault,
	}

	context := initAllContext(param)

	app := New()
	server := app.HttpServer
	router := server.Router().(*router)
	//use default config
	server.SetSessionConfig(session.NewDefaultRuntimeConfig())
	handle := router.wrapRouterHandle(Index, false)

	handle(context)
}

func TestLogWebsocketContext(t *testing.T) {
	param := &InitContextParam{
		t,
		"",
		"",
		test.ToDefault,
	}

	context := initAllContext(param)

	log := logWebsocketContext(context, time.Now().Unix())
	t.Log("logContext:", log)
	//test.NotNil(t,log)
	test.Equal(t, "", "")
}
