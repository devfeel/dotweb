# dotweb
基于go语言开发的web framework

安装：

```
go get -u github.com/devfeel/dotweb
```

快速开始：

```golang
func StartServer() error {
	//初始化DotServer
	dotserver := dotweb.New()
	//设置dotserver日志目录
	dotserver.SetLogPath("/home/logs/wwwroot/")
	//设置路由
	dotserver.HttpServer.GET("/index", func(ctx *dotweb.HttpContext) {
		ctx.WriteString("welcome to my first web!")
	})
	//开始服务
	err := dotserver.StartServer(80)
	return err
}

```

