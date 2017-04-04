package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//注册HttpHandler
	RegisterHandler(app.HttpServer)

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
	ctx.Redirect(200, "http://www.baidu.com")
}

func RegisterHandler(server *dotweb.HttpServer) {
	server.Router().RegisterHandler("Index", Index)
	server.Router().RegisterHandler("DefaultError", DefaultError)
	server.Router().RegisterHandler("Redirect", Redirect)
}
