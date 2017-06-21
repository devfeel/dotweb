package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置路由
	InitRoute(app.HttpServer)

	//设置自定义异常处理接口
	app.SetExceptionHandle(func(ctx dotweb.Context, err error) {
		ctx.WriteString("oh, 我居然出错了！ ", err.Error())
	})

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func DefaultError(ctx dotweb.Context) error {
	panic("my panic error!")
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/error", DefaultError)
}
