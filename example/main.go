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
	InitRoute(app)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//pprofport := 8081
	//go app.StartPProfServer(pprofport)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	//ctx.WriteString("welcome to dotwebwelcome to dotwebwelcome to dotwebwelcome to dotweb")
	ctx.WriteString("")
}

func IndexReg(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("welcome to dotweb")
}

func KeyPost(ctx *dotweb.HttpContext) {
	username1 := ctx.PostString("username")
	username2 := ctx.FormValue("username")
	username3 := ctx.PostFormValue("username")
	ctx.WriteString("username:" + username1 + " - " + username2 + " - " + username3)
}

func JsonPost(ctx *dotweb.HttpContext) {
	ctx.WriteString("body:" + string(ctx.PostBody()))
}

func DefaultError(ctx *dotweb.HttpContext) {
	panic("my panic error!")
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

func InitRoute(dotserver *dotweb.Dotweb) {
	dotserver.HttpServer.GET("/", Index)
	dotserver.HttpServer.POST("/keypost", KeyPost)
	dotserver.HttpServer.POST("/jsonpost", JsonPost)
	dotserver.HttpServer.GET("/error", DefaultError)
	dotserver.HttpServer.GET("/session", TestSession)
	dotserver.HttpServer.RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)
}

func InitModule(dotserver *dotweb.Dotweb) {
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx *dotweb.HttpContext) {
			fmt.Println("BeginRequest1:", ctx)
		},
		OnEndRequest: func(ctx *dotweb.HttpContext) {
			fmt.Println("EndRequest1:", ctx)
		},
	})

	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx *dotweb.HttpContext) {
			fmt.Println("BeginRequest2:", ctx)
		},
	})
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnEndRequest: func(ctx *dotweb.HttpContext) {
			fmt.Println("EndRequest3:", ctx)
		},
	})
}
