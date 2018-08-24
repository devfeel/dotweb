package main

import (
	"errors"
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/exception"
	"net/http"
	"strconv"
	"time"
	"github.com/devfeel/dotweb/session"
)

func main() {

	defer func() {
		var errmsg string
		if err := recover(); err != nil {
			errmsg = exception.CatchError("main", dotweb.LogTarget_HttpServer, err)
			fmt.Println("main error : ", errmsg)
		}
	}()

	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	//如果不设置，默认不启用，且默认为当前目录
	app.SetEnabledLog(true)

	//开启development模式
	app.SetDevelopmentMode()
	app.SetProductionMode()


	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//set virtual path
	app.HttpServer.SetVirtualPath("/1")

	//1.use default config
	//app.HttpServer.Features.SetEnabledCROS()
	//2.use user config
	//app.HttpServer.Features.SetEnabledCROS(true).SetOrigin("*").SetMethod("GET")

	//设置Session配置
	//runtime mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis no auth mode
	//app.HttpServer.SetSessionConfig(session.NewDefaultRedisConfig("redis://192.168.8.175:6379/0"))
	//redis auth mode
	//app.HttpServer.SetSessionConfig(session.NewDefaultRedisConfig("redis://:password@192.168.8.175:6379/0"))

	app.HttpServer.SetEnabledDetailRequestData(true)

	//设置路由
	InitRoute(app.HttpServer)

	//自定义404输出
	app.SetNotFoundHandle(func(ctx dotweb.Context) {
		ctx.Response().Write(http.StatusNotFound, []byte("is't app's not found!"))
	})

	app.SetExceptionHandle(func(ctx dotweb.Context, err error) {
		ctx.Response().SetContentType(dotweb.MIMEApplicationJSONCharsetUTF8)
		ctx.WriteJsonC(http.StatusInternalServerError, err.Error())
	})

	//设置超时钩子事件，当请求超过指定时间阀值，会自动调用传入的函数
	//不会终止请求，只作为旁路执行
	app.UseTimeoutHook(dotweb.DefaultTimeoutHookHandler, time.Second * 2)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	app.SetPProfConfig(true, 8081)

	//全局容器
	app.Items.Set("gstring", "gvalue")
	app.Items.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString(ctx.Request().URL.Path)
	//_, err := ctx.WriteStringC(201, "index => ", ctx.RemoteIP(), "我是首页")
	return nil
}

func Time(ctx dotweb.Context) error {
	minuteTimeLayout := "200601021504"
	if t, err := time.Parse(minuteTimeLayout, "201709251541"); err != nil {
		ctx.WriteString(err.Error())
	} else {
		now, _ := time.Parse(minuteTimeLayout, time.Now().Format(minuteTimeLayout))
		ctx.WriteString(t)
		ctx.WriteString(now)
		ctx.WriteString(t.Sub(now))
		//ctx.WriteString(t.Sub(time.Now()) > 5*time.Minute)
	}
	return nil
}

func IndexReg(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return ctx.WriteString("welcome to dotweb")
}

func IndexParam(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return ctx.WriteString("IndexParam", ctx.GetRouterName("id"))
}

func KeyPost(ctx dotweb.Context) error {
	username1 := ctx.PostFormValue("username")
	username2 := ctx.FormValue("username")
	username3 := ctx.PostFormValue("username")
	return ctx.WriteString("username:" + username1 + " - " + username2 + " - " + username3)
}

func JsonPost(ctx dotweb.Context) error {
	return ctx.WriteString("body:" + string(ctx.Request().PostBody()))
}

func DefaultError(ctx dotweb.Context) error {
	//panic("my panic error!")
	i := 0
	b := 2 / i
	return ctx.WriteString(b)
}

func Redirect(ctx dotweb.Context) error {
	err := ctx.Redirect(http.StatusMovedPermanently, "http://www.baidu.com")
	if err != nil {
		ctx.WriteString(err)
	}
	return err
}

func ReturnError(ctx dotweb.Context) error {
	return errors.New("return error")
}

func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
	server.GET("/time", Time)
	server.GET("/index", Index)
	server.GET("/id/:id", IndexParam)
	server.POST("/keypost", KeyPost)
	server.POST("/jsonpost", JsonPost)
	server.GET("/error", DefaultError)
	server.GET("/returnerr", ReturnError)
	server.GET("/redirect", Redirect)
	//server.Router().RegisterRoute(dotweb.RouteMethod_GET, "/index", IndexReg)
}
