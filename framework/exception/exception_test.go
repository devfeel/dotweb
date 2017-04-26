package exception

import (
	"errors"
	"testing"

	"github.com/devfeel/dotweb"
)

func Test_CatchError_1(t *testing.T) {
	err := errors.New("runtime error: slice bounds out of range.")
	errMsg := CatchError("httpserver::RouterHandle", dotweb.LogTarget_HttpServer, err)
	t.Log(errMsg)
}
