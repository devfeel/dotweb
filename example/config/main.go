package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置路由
	RegisterHandler(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//pprofport := 8081
	//go app.StartPProfServer(pprofport)

	appConfig := config.InitConfig("d:/dotweb.conf")

	fmt.Println("dotweb.StartServer => " + fmt.Sprint(appConfig))
	err := app.StartServerWithConfig(appConfig)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("index")
}

func DefaultError(ctx *dotweb.HttpContext) {
	panic("my panic error!")
}

func Redirect(ctx *dotweb.HttpContext) {
	ctx.Redirect("http://www.baidu.com")
}

func RegisterHandler(server *dotweb.HttpServer) {
	server.RegisterHandler("Index", Index)
	server.RegisterHandler("DefaultError", DefaultError)
	server.RegisterHandler("Redirect", Redirect)
}
