package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())

	app.HttpServer.SetEnabledListDir(false)

	//设置路由
	InitRoute(app.HttpServer)

	// 开始服务
	port := 80
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}
func InitRoute(server *dotweb.HttpServer) {
	server.Router().ServerFile("/*filepath", "/devfeel/dotweb/public")
}
