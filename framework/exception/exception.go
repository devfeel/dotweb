package exception

import (
	"fmt"
	"os"
	"runtime/debug"
)

// CatchError is the unified exception handler
func CatchError(title string, logtarget string, err interface{}) (errmsg string) {
	errmsg = fmt.Sprintln(err)
	stack := string(debug.Stack())
	os.Stdout.Write([]byte(title + " error! => " + errmsg + " => " + stack))
	return title + " error! => " + errmsg + " => " + stack
}
