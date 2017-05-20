package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/session"
	"net/http"
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

	//设置gzip开关
	//app.SetEnabledGzip(true)

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//1.use default config
	//app.HttpServer.Features.SetEnabledCROS()
	//2.use user config
	//app.HttpServer.Features.SetEnabledCROS(true).SetOrigin("*").SetMethod("GET")

	//设置Session配置
	//runtime mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	//app.HttpServer.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6379", ""))

	//设置路由
	InitRoute(app.HttpServer)

	//自定义404输出
	app.NotFoundHandler = func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("is't app's not found!"))
	}

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	app.SetPProfConfig(true, 8081)

	//全局容器
	app.AppContext.Set("gstring", "gvalue")
	app.AppContext.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := ctx.WriteStringC(201, "index => ", ctx.RouterParams)
	return err
}

func IndexReg(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err := ctx.WriteString("welcome to dotweb")
	return err
}

func KeyPost(ctx dotweb.Context) error {
	username1 := ctx.PostFormValue("username")
	username2 := ctx.FormValue("username")
	username3 := ctx.PostFormValue("username")
	_, err := ctx.WriteString("username:" + username1 + " - " + username2 + " - " + username3)
	return err
}

func JsonPost(ctx dotweb.Context) error {
	_, err := ctx.WriteString("body:" + string(ctx.Request().PostBody()))
	return err
}

func DefaultError(ctx dotweb.Context) error {
	//panic("my panic error!")
	i := 0
	b := 2 / i
	_, err := ctx.WriteString(b)
	return err
}

func Redirect(ctx dotweb.Context) error {
	err := ctx.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	if err != nil {
		ctx.WriteString(err)
	}
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
	server.Router().POST("/keypost", KeyPost)
	server.Router().POST("/jsonpost", JsonPost)
	server.Router().GET("/error", DefaultError)
	server.Router().GET("/redirect", Redirect)
	server.Router().RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)
}

func InitModule(dotserver *dotweb.DotWeb) {
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			fmt.Println("BeginRequest1:", ctx)
		},
		OnEndRequest: func(ctx dotweb.Context) {
			fmt.Println("EndRequest1:", ctx)
		},
	})

	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			fmt.Println("BeginRequest2:", ctx)
		},
	})
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnEndRequest: func(ctx dotweb.Context) {
			fmt.Println("EndRequest3:", ctx)
		},
	})
}
