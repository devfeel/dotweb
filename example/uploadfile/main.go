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

	//设置路由
	InitRoute(app.HttpServer)

	//设置HttpModule
	//InitModule(app)

	//启动 监控服务
	//app.SetPProfConfig(true, 8081)

	// 开始服务
	port := 8080
	fmt.Println("dotweb.StartServer => " + strconv.Itoa(port))
	err := app.StartServer(port)
	fmt.Println("dotweb.StartServer error => ", err)
}

func InitRoute(server *dotweb.HttpServer) {
	server.Router().POST("/file", FileUpload)
}

func FileUpload(ctx dotweb.Context) error {
	upload, err := ctx.Request().FormFile("file")
	if err != nil {
		_, err := ctx.WriteString("FormFile error " + err.Error())
		return err
	} else {
		_, err = upload.SaveFile("d:\\" + upload.FileName())
		if err != nil {
			_, err := ctx.WriteString("SaveFile error => " + err.Error())
			return err
		} else {
			_, err := ctx.WriteString("SaveFile success || " + upload.FileName() + " || " + upload.GetFileExt() + " || " + fmt.Sprint(upload.Size()))

			return err
		}
	}

}
