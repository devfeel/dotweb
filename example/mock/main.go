package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认不启用，且默认为当前目录
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()

	//设置Mock逻辑
	app.SetMock(AppMock())

	//设置路由
	InitRoute(app.HttpServer)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

// Index index handler
func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	err := ctx.WriteString("index  => ", ctx.Request().Url())
	return err
}

// InitRoute init app's route
func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
}

// AppMock create app Mock
func AppMock() dotweb.Mock{
	m := dotweb.NewStandardMock()
	m.RegisterString("/", "mock data")
	return m
}
