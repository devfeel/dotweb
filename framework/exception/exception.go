package exception

import (
	"fmt"
	"os"
	"runtime"
)

const logLevel_Error = "error"

//统一异常处理
func CatchError(title string, logtarget string, err interface{}) (errmsg string) {
	errmsg = fmt.Sprintln(err)
	os.Stdout.Write([]byte(title + " error! => " + errmsg + " => "))
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, true)
	return title + " error! => " + errmsg + " => " + string(buf[:n])
}
