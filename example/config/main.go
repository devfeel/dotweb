package main

import (
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

	//appConfig := config.InitConfig("d:/gotmp/dotweb.conf")
	//json config
	appConfig := config.InitConfig("d:/gotmp/dotweb.json.conf", "json")
	fmt.Println(jsonutil.GetJsonString(appConfig))

	RegisterMiddlewares(app)

	err := app.SetConfig(appConfig)
	if err != nil {
		fmt.Println("dotweb.SetConfig error => " + fmt.Sprint(err))
	}

	fmt.Println("dotweb.StartServer => " + fmt.Sprint(appConfig))
	err = app.StartServer(appConfig.Server.Port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func Index(ctx *dotweb.HttpContext) {
	ctx.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString("index => ", fmt.Sprint(ctx.RouterNode.Middlewares()))
}

func DefaultError(ctx *dotweb.HttpContext) {
	panic("my panic error!")
}

func Redirect(ctx *dotweb.HttpContext) {
	ctx.Redirect(200, "http://www.baidu.com")
}

func Login(ctx *dotweb.HttpContext) {
	ctx.WriteString("login => ", fmt.Sprint(ctx.RouterNode.Middlewares()))
}

func Logout(ctx *dotweb.HttpContext) {
	ctx.WriteString("logout => ", fmt.Sprint(ctx.RouterNode.Middlewares()))
}

func RegisterHandler(server *dotweb.HttpServer) {
	server.Router().RegisterHandler("Index", Index)
	server.Router().RegisterHandler("DefaultError", DefaultError)
	server.Router().RegisterHandler("Redirect", Redirect)
	server.Router().RegisterHandler("Login", Login)
	server.Router().RegisterHandler("Logout", Logout)
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

func (m *AccessFmtLog) Handle(ctx *dotweb.HttpContext) error {
	fmt.Println(time.Now(), "[AccessFmtLog ", m.Index, "] begin request -> ", ctx.Request.RequestURI)
	err := m.Next(ctx)
	fmt.Println(time.Now(), "[AccessFmtLog ", m.Index, "] finish request ", err, " -> ", ctx.Request.RequestURI)
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

func (m *SimpleAuth) Handle(ctx *dotweb.HttpContext) error {
	fmt.Println(time.Now(), "[SimpleAuth] begin request -> ", ctx.Request.RequestURI)
	var err error
	if ctx.QueryString("token") != m.exactToken {
		ctx.Write(http.StatusUnauthorized, []byte("sorry, Unauthorized"))
	} else {
		err = m.Next(ctx)
	}
	fmt.Println(time.Now(), "[SimpleAuth] finish request ", err, " -> ", ctx.Request.RequestURI)
	return err
}

func NewSimpleAuth() dotweb.Middleware {
	return &SimpleAuth{exactToken: "admin"}
}
