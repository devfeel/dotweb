package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/session"
	"net/http"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()

	//设置gzip开关
	//app.SetEnabledGzip(true)

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//1.use default config
	//app.HttpServer.Features.SetEnabledCROS()
	//2.use user config
	//app.HttpServer.Features.SetEnabledCROS(true).SetOrigin("*").SetMethod("GET")

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
	app.SetPProfConfig(true, 8081)

	//全局容器
	app.AppContext.Set("gstring", "gvalue")
	app.AppContext.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("index => ", ctx.RouterParams)
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

func Redirect(ctx *dotweb.HttpContext) {
	ctx.Redirect(http.StatusNotFound, "http://www.baidu.com")
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index).SetEnabledCROS()
	server.Router().POST("/keypost", KeyPost)
	server.Router().POST("/jsonpost", JsonPost)
	server.Router().GET("/error", DefaultError)
	server.Router().GET("/redirect", Redirect)
	server.Router().RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)
}

func InitModule(dotserver *dotweb.DotWeb) {
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
