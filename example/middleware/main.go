package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"net/http"
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

	//设置gzip开关
	//app.SetEnabledGzip(true)

	//设置路由
	InitRoute(app.HttpServer)

	//InitModule(app)

	app.UseRequestLog()
	app.Use(
		NewAccessFmtLog("app"),
	//NewSimpleAuth("admin"),
	)

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
	//fmt.Println(time.Now(), "Index Handler")
	_, err := ctx.WriteString("index  => ", fmt.Sprint(ctx.RouterNode().Middlewares()))
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
	server.Router().GET("/use", Index).Use(NewAccessFmtLog("Router-use"))

	g := server.Group("/group").Use(NewAccessFmtLog("group")).Use(NewSimpleAuth("admin"))
	g.GET("/", Index)
	g.GET("/use", Index).Use(NewAccessFmtLog("group-use"))
}

func InitModule(dotserver *dotweb.DotWeb) {
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			fmt.Println(time.Now(), "HttpModule BeginRequest1:", ctx.Request().RequestURI)
		},
		OnEndRequest: func(ctx dotweb.Context) {
			fmt.Println(time.Now(), "HttpModule EndRequest1:", ctx.Request().RequestURI)
		},
	})

	dotserver.RegisterModule(&dotweb.HttpModule{
		OnBeginRequest: func(ctx dotweb.Context) {
			fmt.Println(time.Now(), "HttpModule BeginRequest2:", ctx.Request().RequestURI)
		},
	})
	dotserver.RegisterModule(&dotweb.HttpModule{
		OnEndRequest: func(ctx dotweb.Context) {
			fmt.Println(time.Now(), "HttpModule EndRequest3:", ctx.Request().RequestURI)
		},
	})
}

type AccessFmtLog struct {
	dotweb.BaseMiddlware
	Index string
}

func (m *AccessFmtLog) Handle(ctx dotweb.Context) error {
	fmt.Println(time.Now(), "[AccessFmtLog ", m.Index, "] begin request -> ", ctx.Request().RequestURI)
	err := m.Next(ctx)
	fmt.Println(time.Now(), "[AccessFmtLog ", m.Index, "] finish request ", err, " -> ", ctx.Request().RequestURI)
	return err
}

func NewAccessFmtLog(index string) *AccessFmtLog {
	return &AccessFmtLog{Index: index}
}

type SimpleAuth struct {
	dotweb.BaseMiddlware
	exactToken string
}

func (m *SimpleAuth) Handle(ctx dotweb.Context) error {
	fmt.Println(time.Now(), "[SimpleAuth] begin request -> ", ctx.Request().RequestURI)
	var err error
	if ctx.QueryString("token") != m.exactToken {
		ctx.Write(http.StatusUnauthorized, []byte("sorry, Unauthorized"))
	} else {
		err = m.Next(ctx)
	}
	fmt.Println(time.Now(), "[SimpleAuth] finish request ", err, " -> ", ctx.Request().RequestURI)
	return err
}

func NewSimpleAuth(exactToken string) *SimpleAuth {
	return &SimpleAuth{exactToken: exactToken}
}
