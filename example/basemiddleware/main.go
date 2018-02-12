package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"strconv"
	"time"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认不启用，且默认为当前目录
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()

	//启用超时处理，这里设置为3秒
	app.UseTimeoutHook(
		func(ctx dotweb.Context) {
			fmt.Println(ctx.Items().GetTimeDuration(dotweb.ItemKeyHandleDuration)/time.Millisecond)
		}, time.Second * 3)
	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

// Index
func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	//fmt.Println(time.Now(), "Index Handler")
	err := ctx.WriteString("index  => ", ctx.Request().Url())
	fmt.Println(ctx.RouterNode().GroupMiddlewares())
	return err
}

// Wait10Second
func Wait10Second(ctx dotweb.Context) error{
	time.Sleep(time.Second * 10)
	ctx.WriteString("HandleDuration:", fmt.Sprint(ctx.Items().Get(dotweb.ItemKeyHandleStartTime)))
	return nil
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
	server.Router().GET("/index", Index)
	server.Router().GET("/wait", Wait10Second)
}


