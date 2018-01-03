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

	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	app.SetDevelopmentMode()

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	app.HttpServer.SetEnabledIgnoreFavicon(true)

	//设置Session配置
	//runtime mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	//app.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6379", ""))

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	InitModule(app.HttpServer)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	//全局容器
	app.AppContext.Set("gstring", "gvalue")
	app.AppContext.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Items().Set("count", 2)
	ctx.WriteString(ctx.Request().Path() + ":Items.Count=> " + ctx.Items().GetString("count"))
	return ctx.WriteString("\r\n")
}

func WHtml(ctx dotweb.Context) error {
	ctx.WriteHtml("this is html response!")
	return nil
}

func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
	server.GET("/m", Index)
	server.GET("/h", WHtml)
}

func InitModule(dotserver *dotweb.HttpServer) {
	dotserver.RegisterModule(&dotweb.HttpModule{
		Name: "test change route",
		OnBeginRequest: func(ctx dotweb.Context) {
			if ctx.IsEnd() {
				return
			}
			if ctx.Request().Path() == "/" && ctx.QueryString("change") == "1" {
				//change route
				ctx.WriteString("变更访问路由测试")
				ctx.WriteString("\r\n")
				ctx.Request().URL.Path = "/m"
			}

			if ctx.Request().Path() == "/" {
				ctx.Items().Set("count", 1)
				ctx.WriteString("OnBeginRequest:Items.Count => ", ctx.Items().GetString("count"))
				ctx.WriteString("\r\n")
			}
			if ctx.QueryString("skip") == "1" {
				ctx.End()
			}
		},
		OnEndRequest: func(ctx dotweb.Context) {
			if ctx.IsEnd() {
				return
			}
			if ctx.Request().Path() == "/" {
				if ctx.Items().Exists("count") {
					ctx.WriteString("OnEndRequest:Items.Count => ", ctx.Items().GetString("count"))
				} else {
					ctx.WriteString("OnEndRequest:Items.Len => ", ctx.Items().Len())
				}
			}
		},
	})
}
