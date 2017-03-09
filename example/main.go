package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"strconv"
)

func main() {
	//初始化DotServer
	dotserver := dotweb.New()

	//设置dotserver日志目录
	dotserver.SetLogPath(file.GetCurrentDirectory())

	//设置路由
	InitRoute(dotserver)

	//启动监控服务
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

func InitRoute(dotserver *dotweb.Dotweb) {
	dotserver.HttpServer.GET("/", Index)
	dotserver.HttpServer.RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)
}
