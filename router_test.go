package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
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
	r:=NewRouter(server)

	r.ServeHTTP(context.response.writer,context.request.Request)
}
