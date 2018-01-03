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

	//这里仅为示例，默认情况下，开启的模式就是development模式
	app.SetDevelopmentMode()

	//使用json标签
	app.HttpServer.SetEnabledBindUseJsonTag(true)
	//设置gzip开关
	//app.HttpServer.SetEnabledGzip(true)

	//设置路由
	InitRoute(app.HttpServer)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func TestBind(ctx dotweb.Context) error {
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

	return ctx.WriteString("TestBind [" + errstr + "] " + fmt.Sprint(user))
}

func GetBind(ctx dotweb.Context) error {
	//type UserInfo struct {
	//	UserName string `form:"user"`
	//	Sex      int    `form:"sex"`
	//}
	type UserInfo struct {
		UserName string `json:"user"`
		Sex      int    `json:"sex"`
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.Bind(user); err != nil {
		errstr = err.Error()
	} else {

	}

	return ctx.WriteString("GetBind [" + errstr + "] " + fmt.Sprint(user))
}

func PostJsonBind(ctx dotweb.Context) error{
	type UserInfo struct {
		UserName string `json:"user"`
		Sex      int    `json:"sex"`
	}
	user := new(UserInfo)
	errstr := "no error"
	if err := ctx.BindJsonBody(user); err != nil {
		errstr = err.Error()
	} else {

	}

	return ctx.WriteString("PostBind [" + errstr + "] " + fmt.Sprint(user))
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().POST("/", TestBind)
	server.Router().GET("/getbind", GetBind)
	server.Router().POST("/jsonbind", PostJsonBind)
}
