package exception

import (
	"fmt"
	"os"
	"runtime"
)

//统一异常处理
func CatchError(title string, logtarget string, err interface{}) (errmsg string) {
	errmsg = fmt.Sprintln(err)
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, true)
	stack := string(buf[:n])
	os.Stdout.Write([]byte(title + " error! => " + errmsg + " => " + stack))
	return title + " error! => " + errmsg + " => " + stack
}
