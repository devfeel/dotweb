package main

import (
	"errors"
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/framework/json"
	"net/http"
	"time"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//注册HttpHandler
	RegisterHandler(app.HttpServer)

	//xml config
	//appConfig, err := config.InitConfig("d:/gotmp/dotweb.conf")
	//json config
	//appConfig, err := config.InitConfig("d:/gotmp/dotweb.json", "json")
	//yaml config
	appConfig, err := config.InitConfig("d:/gotmp/dotweb.yaml", "yaml")
	if err != nil {
		fmt.Println("dotweb.InitConfig error => " + fmt.Sprint(err))
		return
	}
	fmt.Println(jsonutil.GetJsonString(appConfig))

	RegisterMiddlewares(app)

	err = app.SetConfig(appConfig)
	if err != nil {
		fmt.Println("dotweb.SetConfig error => " + fmt.Sprint(err))
		return
	}

	fmt.Println("dotweb.StartServer => " + fmt.Sprint(appConfig))
	err = app.Start()
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	return ctx.WriteString("index => ", fmt.Sprint(ctx.RouterNode().Middlewares()))
}

func GetAppSet(ctx dotweb.Context) error {
	key := ctx.QueryString("key")
	return ctx.WriteString(ctx.Request().Url(), " => key = ", ctx.ConfigSet().GetString(key))
}

func DefaultPanic(ctx dotweb.Context) error {
	panic("my panic error!")
	return nil
}

func DefaultError(ctx dotweb.Context) error {
	err := errors.New("my return error")
	return err
}

func Redirect(ctx dotweb.Context) error {
	return ctx.Redirect(200, "http://www.baidu.com")
}

func Login(ctx dotweb.Context) error {
	return ctx.WriteString("login => ", fmt.Sprint(ctx.RouterNode().Middlewares()))
}

func Logout(ctx dotweb.Context) error {
	return ctx.WriteString("logout => ", fmt.Sprint(ctx.RouterNode().Middlewares()))
}

func RegisterHandler(server *dotweb.HttpServer) {
	server.Router().RegisterHandler("Index", Index)
	server.Router().RegisterHandler("Error", DefaultError)
	server.Router().RegisterHandler("Panic", DefaultPanic)
	server.Router().RegisterHandler("Redirect", Redirect)
	server.Router().RegisterHandler("Login", Login)
	server.Router().RegisterHandler("Logout", Logout)
	server.Router().RegisterHandler("appset", GetAppSet)
}

func RegisterMiddlewares(app *dotweb.DotWeb) {
	//集中注册middleware
	app.RegisterMiddlewareFunc("applog", NewAppAccessFmtLog)
	app.RegisterMiddlewareFunc("grouplog", NewGroupAccessFmtLog)
	app.RegisterMiddlewareFunc("urllog", NewUrlAccessFmtLog)
	app.RegisterMiddlewareFunc("simpleauth", NewSimpleAuth)
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

func NewAppAccessFmtLog() dotweb.Middleware {
	return &AccessFmtLog{Index: "app"}
}

func NewGroupAccessFmtLog() dotweb.Middleware {
	return &AccessFmtLog{Index: "group"}
}

func NewUrlAccessFmtLog() dotweb.Middleware {
	return &AccessFmtLog{Index: "url"}
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

func NewSimpleAuth() dotweb.Middleware {
	return &SimpleAuth{exactToken: "admin"}
}
