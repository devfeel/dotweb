package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/exception"
	"strconv"
)

func main() {

	defer func() {
		var errmsg string
		if err := recover(); err != nil {
			errmsg = exception.CatchError("main", dotweb.LogTarget_HttpServer, err)
			fmt.Println("main error : ", errmsg)
		}
	}()

	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认不启用，且默认为当前目录
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	app.SetPProfConfig(true, 8081)

	//设置TLS
	app.HttpServer.SetEnabledTLS(true, "", "")

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

// Index index handler
func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString(ctx.Request().URL.Path)
	//_, err := ctx.WriteStringC(201, "index => ", ctx.RemoteIP(), "我是首页")
	return nil
}

// InitRoute init http server routers
func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
}
