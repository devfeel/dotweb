package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/framework/file"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())

	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	//app.SetCache(cache.NewRuntimeCache())
	app.SetCache(cache.NewRedisCache("127.0.0.1:6379"))

	err := app.Cache().Set("g", "gv", 20)
	if err != nil {
		fmt.Println("Cache Set ", err)
	}

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err = app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

type UserInfo struct {
	UserName string
	Sex      int
}

func One(ctx dotweb.Context) error {
	g, err := ctx.Cache().GetString("g")
	if err != nil {
		g = err.Error()
	}
	_, err = ctx.Cache().Incr("count")
	_, err = ctx.WriteString("One [" + g + "] " + fmt.Sprint(err))
	return err
}

func Two(ctx dotweb.Context) error {
	g, err := ctx.Cache().GetString("g")
	if err != nil {
		g = err.Error()
	}
	_, err = ctx.Cache().Incr("count")
	c, _ := ctx.Cache().GetString("count")
	_, err = ctx.WriteString("Two [" + g + "] [" + c + "] " + fmt.Sprint(err))
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/1", One)
	server.Router().GET("/2", Two)
}
