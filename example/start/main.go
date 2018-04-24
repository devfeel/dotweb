package main

import (
	"github.com/devfeel/dotweb"
	"fmt"
	"strconv"
)

func main(){
	app := dotweb.Classic(dotweb.DefaultLogPath)
	//app := dotweb.New()
	//开启development模式
	app.SetDevelopmentMode()

	//设置路由
	InitRoute(app.HttpServer)


	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

// Index index action
func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString(ctx.Request().URL.Path)
	//_, err := ctx.WriteStringC(201, "index => ", ctx.RemoteIP(), "我是首页")
	return nil
}

// InitRoute init routes
func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
}
