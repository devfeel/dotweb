### dotweb
基于go语言开发的微Web框架

### 安装：

```
go get -u github.com/devfeel/dotweb
```

### 快速开始：

```go
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

### 路由规则
目前支持GET\POST\HEAD\OPTIONS\PUT\PATCH\DELETE 这几类请求方法

另外也支持HiJack\WebSocket\ServerFile三类特殊应用
```go
1、HttpServer.GET(path string, handle HttpHandle)
2、HttpServer.POST(path string, handle HttpHandle)
3、HttpServer.HEAD(path string, handle HttpHandle)
4、HttpServer.OPTIONS(path string, handle HttpHandle)
5、HttpServer.PUT(path string, handle HttpHandle)
6、HttpServer.PATCH(path string, handle HttpHandle)
7、HttpServer.DELETE(path string, handle HttpHandle)
8、HttpServer.HiJack(path string, handle HttpHandle)
9、HttpServer.WebSocket(path string, handle HttpHandle)
```
