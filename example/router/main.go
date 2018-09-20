package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.Classic(file.GetCurrentDirectory())

	app.SetDevelopmentMode()

	app.HttpServer.SetEnabledAutoHEAD(true)
	//app.HttpServer.SetEnabledAutoOPTIONS(true)

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	flag := ctx.HttpServer().Router().MatchPath(ctx, "/d/:x/y")
	return ctx.WriteString("index - " + ctx.Request().Method + " - " + fmt.Sprint(flag))
}

func Any(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return ctx.WriteString("any - " + ctx.Request().Method)
}

func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
	server.GET("/d/:x/y", Index)
	server.Any("/any", Any)
}
