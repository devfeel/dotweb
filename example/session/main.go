package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/session"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//设置Session配置
	//runtime mode
	//app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRedisConfig("redis://192.168.8.175:6379/1"))

	//设置路由
	InitRoute(app.HttpServer)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

type UserInfo struct {
	UserName string
	NickName string
}

func TestSession(ctx dotweb.Context) error {

	user := UserInfo{UserName: "test", NickName: "testName"}
	var userRead UserInfo

	ctx.WriteString("welcome to dotweb - CreateSession - sessionid=> "+ctx.SessionID(), "\r\n")
	err := ctx.Session().Set("username", user)
	if err != nil {
		ctx.WriteString("session set error => ", err, "\r\n")
	}
	c := ctx.Session().Get("username")
	if c != nil {
		userRead = c.(UserInfo)
	} else {
		ctx.WriteString("session read failed, get nil", "\r\n")
	}

	return ctx.WriteString("userinfo=>" + fmt.Sprintln(userRead))
	return err
}

func TestReadSession(ctx dotweb.Context) error {

	var userRead UserInfo

	ctx.WriteString("welcome to dotweb - ReadSession - sessionid=> "+ctx.SessionID(), "\r\n")

	c := ctx.Session().Get("username")
	if c != nil {
		userRead = c.(UserInfo)
	} else {
		ctx.WriteString("session read failed, get nil", "\r\n")
	}

	return ctx.WriteString("userinfo=>" + fmt.Sprintln(userRead))
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", TestSession)
	server.Router().GET("/read", TestReadSession)
}
