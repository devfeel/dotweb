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

	app.SetDevelopmentMode()

	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置路由
	InitRoute(app.HttpServer)

	//set default template path
	app.HttpServer.Renderer().SetTemplatePath("d:/gotmp/")
	app.HttpServer.Renderer().RegisterTemplateFunc("echo", func(x string) interface{}{
		return "echo:" + x
	})

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

type UserInfo struct {
	UserName string
	Sex      bool
}

type BookInfo struct {
	Name string
	Size int64
}

func NotExistView(ctx dotweb.Context) error {
	err := ctx.View("1.html")
	return err
}

func TestView(ctx dotweb.Context) error {
	ctx.ViewData().Set("data", "图书信息")
	ctx.ViewData().Set("user", &UserInfo{UserName: "user1", Sex: true})
	m := make([]*BookInfo, 5)
	m[0] = &BookInfo{Name: "book0", Size: 1}
	m[1] = &BookInfo{Name: "book1", Size: 10}
	m[2] = &BookInfo{Name: "book2", Size: 100}
	m[3] = &BookInfo{Name: "book3", Size: 1000}
	m[4] = &BookInfo{Name: "book4", Size: 10000}
	ctx.ViewData().Set("Books", m)

	err := ctx.View("testview.html")
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", TestView)
	server.Router().GET("/noview", NotExistView)
}
