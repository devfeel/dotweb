package main

import (
	"fmt"
	"github.com/devfeel/dotweb"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/framework/json"
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

	//引入自定义ConfigSet
	err = app.Config.IncludeConfigSet("d:/gotmp/userconf.xml", config.ConfigType_XML)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	app.SetConfig(appConfig)

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

// ConfigSet
func ConfigSet(ctx dotweb.Context) error {
	vkey1 := ctx.ConfigSet().GetString("set1")
	vkey2 := ctx.ConfigSet().GetString("set2")
	return ctx.WriteString(ctx.Request().Path(), "key1=", vkey1, "key2=", vkey2)
}

func RegisterHandler(server *dotweb.HttpServer) {
	server.Router().RegisterHandler("Index", Index)
	server.Router().RegisterHandler("appset", GetAppSet)
	server.GET("/configser", ConfigSet)
}
