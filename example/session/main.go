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

	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置Session开关
	app.HttpServer.SetEnabledSession(true)

	//设置Session配置
	//runtime mode
	//app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
	//redis mode
	app.HttpServer.SetSessionConfig(session.NewDefaultRedisConfig("192.168.8.175:6381"))

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	//全局容器
	app.AppContext.Set("gstring", "gvalue")
	app.AppContext.Set("gint", 1)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func TestSession(ctx dotweb.Context) error {
	type UserInfo struct {
		UserName string
		NickName string
	}
	user := UserInfo{UserName: "test", NickName: "testName"}
	var userRead UserInfo

	ctx.WriteString("welcome to dotweb - sessionid=> "+ctx.SessionID(), "\r\n")
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

	_, err = ctx.WriteString("userinfo=>" + fmt.Sprintln(userRead))
	return err
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", TestSession)
}
