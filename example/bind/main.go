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

	//设置Debug开关
	app.SetEnabledDebug(true)

	//设置gzip开关
	//app.SetEnabledGzip(true)

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	//pprofport := 8081
	//go app.StartPProfServer(pprofport)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func TestBind(ctx *dotweb.HttpContext) {
	type UserInfo struct {
		UserName string
		Sex      int
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.Bind(user); err != nil {
		errstr = err.Error()
	} else {

	}

	ctx.WriteString("TestBind [" + errstr + "] " + fmt.Sprint(user))
}

func GetBind(ctx *dotweb.HttpContext) {
	type UserInfo struct {
		UserName string `form:"user"`
		Sex      int    `form:"sex"`
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.Bind(user); err != nil {
		errstr = err.Error()
	} else {

	}

	ctx.WriteString("GetBind [" + errstr + "] " + fmt.Sprint(user))
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().POST("/", TestBind)
	server.Router().GET("/getbind", GetBind)
}
