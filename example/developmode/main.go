package main

import (
	"github.com/devfeel/dotweb"
	"fmt"
	"strconv"
	"github.com/devfeel/dotweb/logger"
)

const loggerFileName = "develop-mode"

func main(){
	app := dotweb.New()
	//if use develop mode
	//1. Enabled Log
	//2. use RequestLogMiddleware
	//3. Enabled Console Print
	app.SetDevelopmentMode()

	//设置路由
	InitRoute(app.HttpServer)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

// Index index action
func Index(ctx dotweb.Context) error {
	ctx.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
	ctx.WriteString(ctx.Request().URL.Path)
	logger.Logger().Debug("Index:WriteString " + ctx.Request().URL.Path, loggerFileName)
	return nil
}

// InitRoute init routes
func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
}
