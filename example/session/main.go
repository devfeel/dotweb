package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/session"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())

	//设置Debug开关
	app.SetEnabledDebug(true)

	//设置gzip开关
	//app.SetEnabledGzip(true)

	//设置Session开关
	app.SetEnabledSession(true)

	//设置Session配置
	//runtime mode
	app.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	//app.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6379", ""))

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//pprofport := 8081
	//go app.StartPProfServer(pprofport)

	//全局容器
	app.AppContext.Set("gstring", "gvalue")
	app.AppContext.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func TestSession(ctx *dotweb.HttpContext) {
	type UserInfo struct {
		UserName string
		NickName string
	}
	user := UserInfo{UserName: "test", NickName: "testName"}
	ctx.Session().Set("username", user)
	userRead := ctx.Session().Get("username").(UserInfo)

	ctx.WriteString("welcome to dotweb - sessionid=> " + ctx.SessionID +
		", session-len=>" + strconv.Itoa(ctx.Session().Count()) +
		",username=>" + fmt.Sprintln(userRead))
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", TestSession)
}
