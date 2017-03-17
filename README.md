# DotWeb
简约大方的go Web微型框架

## 安装：

```
go get -u github.com/devfeel/dotweb
```

## 快速开始：

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
## 特性
* 支持静态路由、参数路由
* 路由支持文件/目录服务
* 中间件支持
* 支持JSON/JSONP/HTML格式输出
* 统一的HTTP错误处理
* 统一的日志处理
* 支持Hijack与websocket

## 路由
特殊说明：集成github.com/julienschmidt/httprouter
#### 常规路由
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
10、HttpServer.RegisterRoute(routeMethod string, path string, handle HttpHandle)
```
接受两个参数，一个是URI路径，另一个是 HttpHandle 类型，设定匹配到该路径时执行的方法；
#### 静态路由
静态路由语法就是没有任何参数变量，pattern是一个固定的字符串。
```go
package main

import (
    "github.com/devfeel/dotweb"
)

func main() {
    dotserver := dotweb.New()
    dotserver.Get("/hello", func(ctx *dotweb.HttpContext) {
        ctx.WriteString("hello world!")
    })
    dotserver.StartServer(80)
}
```
测试：
curl http://127.0.0.1/hello
#### 参数路由
参数路由以冒号 : 后面跟一个字符串作为参数名称，可以通过 HttpContext的 GetRouterName 方法获取路由参数的值。
```go
package main

import (
    "github.com/devfeel/dotweb"
)

func main() {
    dotserver := dotweb.New()
    dotserver.Get("/hello/:name", func(ctx *dotweb.HttpContext) {
        ctx.WriteString("hello " + ctx.GetRouterName("name"))
    })
    dotserver.Get("/news/:category/:newsid", func(ctx *dotweb.HttpContext) {
    	category := ctx.GetRouterName("category")
	newsid := ctx.GetRouterName("newsid")
        ctx.WriteString("news info: category=" + category + " newsid=" + newsid)
    })
    dotserver.StartServer(80)
}
```
测试：
<br>curl http://127.0.0.1/hello/devfeel
<br>curl http://127.0.0.1/hello/category1/1


## 绑定
* HttpContext.Bind(interface{})
* 支持json、xml、Form数据
* 集成echo的bind实现模块
```go
type UserInfo struct {
		UserName string
		Sex      int
}

func(ctx *dotweb.HttpContext) TestBind{
        user := new(UserInfo)
        if err := ctx.Bind(user); err != nil {
        	 ctx.WriteString("err => " + err.Error())
        }else{
             ctx.WriteString("TestBind " + fmt.Sprint(user))
        }
}

```

## 中间件(拦截器)
#### RegisterModule
* 支持OnBeginRequest、OnEndRequest两类中间件
* 通过实现HttpModule.OnBeginRequest、HttpModule.OnEndRequest接口实现自定义中间件
* 通过设置HttpContext.End()提前终止请求

## 异常
#### 500错误
* 默认设置: 当发生未处理异常时，会根据DebugMode向页面输出默认错误信息或者具体异常信息，并返回 500 错误头
* 自定义: 通过DotServer.SetExceptionHandle(handler *ExceptionHandle)实现自定义异常处理逻辑
```go
type ExceptionHandle func(*HttpContext, interface{})
```

## Session
#### 支持runtime、redis两种
* 默认不开启Session支持
* runtime:基于内存存储实现session模块
* redis:基于Redis存储实现session模块,其中redis key以dotweb:session:xxxxxxxxxxxx组成
```go
//设置session支持
dotserver.SetEnabledSession(true)
//使用runtime模式
dotserver.SetSessionConfig(session.NewDefaultRuntimeConfig())
//使用redis模式
dotserver.SetSessionConfig(session.NewDefaultRedisConfig("127.0.0.1:6379", "xxxx"))
//HttpContext使用
ctx.Session().Set(key, value)
```

## Server Config
目前支持三个选项：Debug、Session、Gzip
* SetEnabledDebug 设置是否开启debug模式，会输出server端的debug日志，默认不开启
* SetEnabledSession 设置是否开启Session支持，目前支持runtime、redis两种模式，默认不开启
* SetEnabledGzip 设置是否开启Gzip支持，默认不开启

## 外部依赖
websocket - golang.org/x/net/websocket
<br>
redis - github.com/garyburd/redigo/redis


## 相关项目
#### <a href="https://github.com/devfeel/tokenserver" target="_blank">TokenServer</a>
项目简介：token服务，提供token一致性服务以及相关的全局ID生成服务等


## 如何联系
QQ群：193409346
