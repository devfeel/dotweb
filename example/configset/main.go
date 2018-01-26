package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/framework/file"
	"strconv"
)

func main() {
	//初始化DotServer
	app := dotweb.New()

	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())

	app.SetDevelopmentMode()

	app.HttpServer.SetEnabledIgnoreFavicon(true)

	//引入自定义ConfigSet
	err := app.Config.IncludeConfigSet("d:/gotmp/userconf.xml", config.ConfigType_XML)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	//设置路由
	InitRoute(app.HttpServer)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err = app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

// ConfigSet
func ConfigSet(ctx dotweb.Context) error {
	vkey1 := ctx.ConfigSet().GetString("set1")
	vkey2 := ctx.ConfigSet().GetString("set2")
	ctx.WriteString(ctx.Request().Path(), "key1=", vkey1, "key2=", vkey2)
	return ctx.WriteString("\r\n")
}

// InitRoute
func InitRoute(server *dotweb.HttpServer) {
	server.GET("/c", ConfigSet)
}
