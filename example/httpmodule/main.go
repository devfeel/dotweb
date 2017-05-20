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

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//设置Session配置
	//runtime mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	//app.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6379", ""))

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	InitModule(app)

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
	ctx.WriteString("index => " + ctx.Items().GetString("count"))
	_, err := ctx.WriteString("\r\n")
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
	server.Router().GET("/user", Index) //need login
	server.Router().GET("/login", Index)
	server.Router().GET("/reg", Index)
}

func InitModule(dotserver *dotweb.DotWeb) {
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			if ctx.HttpServer().Router().MatchPath(ctx, "/user") {
				//TODO:need login
			}
			ctx.Items().Set("count", 1)
			ctx.WriteString("OnBeginRequest => ", ctx.Items().GetString("count"))
			ctx.WriteString("\r\n")
			if ctx.QueryString("skip") == "1" {
				ctx.End()
			}
		},
		OnEndRequest: func(ctx dotweb.Context) {
			if ctx.Items().Exists("count") {
				ctx.WriteString("OnEndRequest => ", ctx.Items().GetString("count"))
			} else {
				ctx.WriteString("OnEndRequest => ", ctx.Items().Len())
			}
		},
	})
}
