package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	//"github.com/devfeel/dotweb/session"
	"github.com/devfeel/dotweb/session"
	"strconv"
)

func main() {
	//初始化DotServer
	dotserver := dotweb.New()

	//设置dotserver日志目录
	dotserver.SetLogPath(file.GetCurrentDirectory())

	//设置Debug开关
	dotserver.SetEnabledDebug(true)
	//设置Session开关
	dotserver.SetEnabledSession(true)

	//设置Session配置
	//runtime mode
	dotserver.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	//dotserver.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6379", ""))

	//设置路由
	InitRoute(dotserver)

	//启动 监控服务
	//pprofport := 8081
	//go dotserver.StartPProfServer(pprofport)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := dotserver.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("welcome to dotweb")
}

func IndexReg(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("welcome to dotweb")
}

func DefaultError(ctx *dotweb.HttpContext) {
	panic("my panic error!")
}

func TestSession(ctx *dotweb.HttpContext) {
	ctx.Session().Set("username", "dotweb")
	userName := ctx.Session().GetString("username")

	ctx.WriteString("welcome to dotweb - sessionid=> " + ctx.SessionID +
		", session-len=>" + strconv.Itoa(ctx.Session().Count()) +
		",username=>" + userName)
}

func InitRoute(dotserver *dotweb.Dotweb) {
	dotserver.HttpServer.GET("/", Index)
	dotserver.HttpServer.GET("/error", DefaultError)
	dotserver.HttpServer.GET("/session", TestSession)
	dotserver.HttpServer.RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)

}
