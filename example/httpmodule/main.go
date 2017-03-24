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
	InitModule(app)

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

func Index(ctx *dotweb.HttpContext) {
	ctx.WriteString("index => " + ctx.Items().GetString("count"))
	ctx.WriteString("\r\n")
	ctx.Items().Set("count", 2)
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
}

func InitModule(dotserver *dotweb.DotWeb) {
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx *dotweb.HttpContext) {
			ctx.Items().Set("count", 1)
			ctx.WriteString("OnBeginRequest => ", ctx.Items().GetString("count"))
			ctx.WriteString("\r\n")
			if ctx.QueryString("skip") == "1" {
				ctx.End()
			}
		},
		OnEndRequest: func(ctx *dotweb.HttpContext) {
			if ctx.Items().Exists("count") {
				ctx.WriteString("OnEndRequest => ", ctx.Items().GetString("count"))
			} else {
				ctx.WriteString("OnEndRequest => ", ctx.Items().Len())
			}
		},
	})
}
