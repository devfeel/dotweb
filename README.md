<div align=center>
    <img src="https://github.com/devfeel/dotweb/blob/master/logo.png" width="200"/>
</div>

# DotWeb
Simple and easy go web micro framework

Important: Now need go1.9+ version support.

Document: https://www.kancloud.cn/devfeel/dotweb/346608

Guide: https://github.com/devfeel/dotweb/blob/master/docs/GUIDE.md

[![Gitter](https://badges.gitter.im/devfeel/dotweb.svg)](https://gitter.im/devfeel-dotweb/wechat)
[![GoDoc](https://godoc.org/github.com/devfeel/dotweb?status.svg)](https://godoc.org/github.com/devfeel/dotweb)
[![Go Report Card](https://goreportcard.com/badge/github.com/devfeel/dotweb)](https://goreportcard.com/report/github.com/devfeel/dotweb)
[![Go Build Card](https://travis-ci.org/devfeel/dotweb.svg?branch=master)](https://travis-ci.org/devfeel/dotweb.svg?branch=master)
<a target="_blank" href="http://shang.qq.com/wpa/qunwpa?idkey=836e11667837ad674462a4a97fb21fba487cd3dff5b2e1ca0d7ea4c2324b4574"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Golang-Devfeel" title="Golang-Devfeel"></a>
## 1. Install

```
go get -u github.com/devfeel/dotweb
```

## 2. Getting Started
```go
func StartServer() error {
	//init DotApp
	app := dotweb.New()
	//set log path
	app.SetLogPath("/home/logs/wwwroot/")
	//set route
	app.HttpServer.GET("/index", func(ctx dotweb.Context) error{
		return ctx.WriteString("welcome to my first web!")
	})
	//begin server
	err := app.StartServer(80)
	return err
}

```
#### examples: https://github.com/devfeel/dotweb-example

## 3. Features
* 支持静态路由、参数路由、组路由
* 路由支持文件/目录服务，支持设置是否允许目录浏览
* HttpModule支持，支持路由之前全局级别的自定义代码能力
* 中间件支持，支持App、Group、Router级别的设置 - https://github.com/devfeel/middleware
* Feature支持，可绑定HttpServer全局启用
* 支持STRING/JSON/JSONP/HTML格式输出
* 统一的HTTP错误处理
* 统一的日志处理
* 支持Hijack与websocket
* 内建Cache支持
* 内建TLS支持
* 支持接入第三方模板引擎（需实现dotweb.Renderer接口）
* 模块可配置化，85%模块可通过配置维护
* 自集成基础统计数据，并支持按分钟为单位的间隔时间统计数据输出

#### Config Example
* [dotweb.conf](https://github.com/devfeel/dotweb/blob/master/example/config/dotweb.conf)
* [dotweb.json](https://github.com/devfeel/dotweb/blob/master/example/config/dotweb.json)
 
## 4. Performance

DotWeb | 1.9.2 | 16core16G |   |   |   |   |   |   |   |   |   |   |  
-- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | -- | --
cpu | 内存 | Samples | Average | Median | 90%Line | 95%Line | 99%Line | Min | Max | Error% | Throughput | Received KB/Sec | Send KB/Sec
40% | 39M | 15228356 | 19 | 4 | 43 | 72 | 204 | 0 | 2070 | 0.00% | 48703.47 | 7514.79 | 8656.28
40% | 42M | 15485189 | 18 | 4 | 41 | 63 | 230 | 0 | 3250 | 0.00% | 49512.99 | 7639.7 | 8800.16
40% | 44M | 15700385 | 18 | 3 | 41 | 64 | 233 | 0 | 2083 | 0.00% | 50203.32 | 7746.22 | 8922.86
  |   |   |   |   |   |   |   |   |   |   |   |   |  
ECHO | 1.9.2 | 16core16G |   |   |   |   |   |   |   |   |   |   |  
cpu | 内存 | Samples | Average | Median | 90%Line | 95%Line | 99%Line | Min | Max | Error% | Throughput | Received KB/Sec | Send KB/Sec
38% | 35M | 15307586 | 19 | 4 | 44 | 76 | 181 | 0 | 1808 | 0.00% | 48951.22 | 6166.71 | 9034.94
36% | 35M | 15239058 | 19 | 4 | 45 | 76 | 178 | 0 | 2003 | 0.00% | 48734.26 | 6139.37 | 8994.9
37% | 37M | 15800585 | 18 | 3 | 41 | 66 | 229 | 0 | 2355 | 0.00% | 50356.09 | 6343.68 | 9294.24
  |   |   |   |   |   |   |   |   |   |   |   |   |  
Gin  |  1.9.2 |  16core16G  |   |   |   |   |   |   |   |   |   |   |  
cpu | 内存 | Samples | Average | Median | 90%Line | 95%Line | 99%Line | Min | Max | Error% | Throughput | Received KB/Sec | Send KB/Sec
36% | 36M | 15109143 | 19 | 6 | 44 | 71 | 175 | 0 | 3250 | 0.00% | 48151.87 | 5877.91 | 8746.33
36% | 40M | 15255749 | 19 | 5 | 43 | 70 | 189 | 0 | 3079 | 0.00% | 48762.53 | 5952.45 | 8857.25
36% | 40M | 15385401 | 18 | 4 | 42 | 66 | 227 | 0 | 2312 | 0.00% | 49181.03 | 6003.54 | 8933.27


## 5. Router
#### 1) 常规路由
* 支持GET\POST\HEAD\OPTIONS\PUT\PATCH\DELETE 这几类请求方法
* 支持HiJack\WebSocket\ServerFile三类特殊应用
* 支持Any注册方式，默认兼容GET\POST\HEAD\OPTIONS\PUT\PATCH\DELETE方式
* 支持通过配置开启默认添加HEAD方式
* 支持注册Handler，以启用配置化
* 支持检查请求与指定路由是否匹配
```go
1、Router.GET(path string, handle HttpHandle)
2、Router.POST(path string, handle HttpHandle)
3、Router.HEAD(path string, handle HttpHandle)
4、Router.OPTIONS(path string, handle HttpHandle)
5、Router.PUT(path string, handle HttpHandle)
6、Router.PATCH(path string, handle HttpHandle)
7、Router.DELETE(path string, handle HttpHandle)
8、Router.HiJack(path string, handle HttpHandle)
9、Router.WebSocket(path string, handle HttpHandle)
10、Router.Any(path string, handle HttpHandle)
11、Router.RegisterRoute(routeMethod string, path string, handle HttpHandle)
12、Router.RegisterHandler(name string, handler HttpHandle)
13、Router.GetHandler(name string) (HttpHandle, bool)
14、Router.MatchPath(ctx Context, routePath string) bool
```
接受两个参数，一个是URI路径，另一个是 HttpHandle 类型，设定匹配到该路径时执行的方法；
#### 2) static router
静态路由语法就是没有任何参数变量，pattern是一个固定的字符串。
```go
package main

import (
    "github.com/devfeel/dotweb"
)

func main() {
    dotapp := dotweb.New()
    dotapp.HttpServer.GET("/hello", func(ctx dotweb.Context) error{
        return ctx.WriteString("hello world!")
    })
    dotapp.StartServer(80)
}
```
test：
curl http://127.0.0.1/hello
#### 3) parameter router
参数路由以冒号 : 后面跟一个字符串作为参数名称，可以通过 HttpContext的 GetRouterName 方法获取路由参数的值。
```go
package main

import (
    "github.com/devfeel/dotweb"
)

func main() {
    dotapp := dotweb.New()
    dotapp.HttpServer.GET("/hello/:name", func(ctx dotweb.Context) error{
        return ctx.WriteString("hello " + ctx.GetRouterName("name"))
    })
    dotapp.HttpServer.GET("/news/:category/:newsid", func(ctx dotweb.Context) error{
    	category := ctx.GetRouterName("category")
	    newsid := ctx.GetRouterName("newsid")
        return ctx.WriteString("news info: category=" + category + " newsid=" + newsid)
    })
    dotapp.StartServer(80)
}
```
test：
<br>curl http://127.0.0.1/hello/devfeel
<br>curl http://127.0.0.1/hello/category1/1
#### 4) group router
```go
    g := server.Group("/user")
	g.GET("/", Index)
	g.GET("/profile", Profile)
```
test：
<br>curl http://127.0.0.1/user
<br>curl http://127.0.0.1/user/profile


## 6. Binder
* HttpContext.Bind(interface{})
* Support data from json、xml、Form
```go
type UserInfo struct {
		UserName string `form:"user"`
		Sex      int    `form:"sex"`
}

func TestBind(ctx dotweb.HttpContext) error{
        user := new(UserInfo)
        if err := ctx.Bind(user); err != nil {
        	 return ctx.WriteString("err => " + err.Error())
        }else{
             return ctx.WriteString("TestBind " + fmt.Sprint(user))
        }
}

```

## 7. Middleware
#### Middleware
* 支持粒度：App、Group、RouterNode
* DotWeb.Use(m ...Middleware)
* Group.Use(m ...Middleware)
* RouterNode.Use(m ...Middleware)
* 启用顺序：App -> Group -> RouterNode，同级别下按Use的引入顺序执行
* 更多请参考：https://github.com/devfeel/middleware
* [JWT](https://github.com/devfeel/middleware/tree/master/jwt)   -  [example](https://github.com/devfeel/middleware/tree/master/example/jwt)
* [AccessLog](https://github.com/devfeel/middleware/tree/master/accesslog)   -  [example](https://github.com/devfeel/middleware/tree/master/example/accesslog)
* [CORS](https://github.com/devfeel/middleware/tree/master/cors)   -  [example](https://github.com/devfeel/middleware/tree/master/example/cors)
* [Gzip](https://github.com/devfeel/middleware/tree/master/gzip)   -  [example](https://github.com/devfeel/middleware/tree/master/example/gzip)
* [authorization based on Casbin](https://github.com/devfeel/middleware/tree/master/authz) - [example](https://github.com/devfeel/middleware/tree/master/example/authz) - [what's Casbin?](https://github.com/casbin/casbin)
* BasicAuth
* Recover
* HeaderOverride

```go
app.Use(NewAccessFmtLog("app"))

func InitRoute(server *dotweb.HttpServer) {
	server.GET("/", Index)
	server.GET("/use", Index).Use(NewAccessFmtLog("Router-use"))

	g := server.Group("/group").Use(NewAccessFmtLog("group"))
	g.GET("/", Index)
	g.GET("/use", Index).Use(NewAccessFmtLog("group-use"))
}

type AccessFmtLog struct {
	dotweb.BaseMiddlware
	Index string
}

func (m *AccessFmtLog) Handle(ctx dotweb.Context) error {
	fmt.Println(time.Now(), "[AccessFmtLog ", m.Index, "] begin request -> ", ctx.Request.RequestURI)
	err := m.Next(ctx)
	fmt.Println(time.Now(), "[AccessFmtLog ", m.Index, "] finish request ", err, " -> ", ctx.Request.RequestURI)
	return err
}

func NewAccessFmtLog(index string) *AccessFmtLog {
	return &AccessFmtLog{Index: index}
}
```

## 8. Server Config
#### HttpServer：
* HttpServer.EnabledSession

  设置是否开启Session支持，目前支持runtime、redis两种模式，默认不开启
* HttpServer.EnabledGzip

  设置是否开启Gzip支持，默认不开启
* HttpServer.EnabledListDir

  设置是否启用目录浏览，仅对Router.ServerFile有效，若设置该项，则可以浏览目录文件，默认不开启
* HttpServer.EnabledAutoHEAD

  设置是否自动启用Head路由，若设置该项，则会为除Websocket\HEAD外所有路由方式默认添加HEAD路由，默认不开启
* HttpServer.EnabledIgnoreFavicon

  设置是否忽略Favicon的请求，一般用于接口项目
* HttpServer.EnabledDetailRequestData

  设置是否启用详细请求数据统计,默认为false，若设置该项，将启用ServerStateInfo中DetailRequestUrlData的统计
* HttpServer.EnabledTLS

  设置是否启用TLS加密处理

#### Run Mode
* 新增development、production模式
* 默认development，通过DotWeb.SetDevelopmentMode\DotWeb.SetProductionMode开启相关模式
* 若设置development模式，未处理异常会输出异常详细信息，同时启用日志开关，同时启用日志console打印
* 未来会拓展更多运行模式的配置


## 9. Exception
#### 500 error
* Default: 当发生未处理异常时，会根据RunMode向页面输出默认错误信息或者具体异常信息，并返回 500 错误头
* User-defined: 通过DotServer.SetExceptionHandle(handler *ExceptionHandle)实现自定义异常处理逻辑
```go
type ExceptionHandle func(Context, error)
```
#### 404 error
* Default: 当发生404异常时，会默认使用http.NotFound处理
* User-defined: 通过DotWeb.SetNotFoundHandle(handler NotFoundHandle)实现自定义404处理逻辑
```go
type NotFoundHandle  func(http.ResponseWriter, *http.Request)
```

## Dependency
websocket - golang.org/x/net/websocket - 内置vendor
<br>
redis - github.com/garyburd/redigo - go get自动下载
<br>
yaml - gopkg.in/yaml.v2  - go get自动下载

## 相关项目
#### <a href="https://github.com/devfeel/longweb" target="_blank">LongWeb</a>
项目简介：http长连接网关服务，提供Websocket及长轮询服务

#### <a href="https://github.com/yulibaozi/yulibaozi.com" target="_blank">yulibaozi.com</a>
项目简介：基于dotweb与mapper的一款go的博客程序

#### <a href="https://github.com/devfeel/tokenserver" target="_blank">TokenServer</a>
项目简介：token服务，提供token一致性服务以及相关的全局ID生成服务等

#### <a href="https://github.com/gnuos/wechat-token" target="_blank">Wechat-token</a>
项目简介：微信Access Token中控服务器，用来统一管理各个公众号的access_token，提供统一的接口进行获取和自动刷新Access Token。

#### <a href="https://github.com/devfeel/dotweb-start" target="_blank">dotweb-start</a>
项目简介：基于dotweb、dotlog、mapper、dottask、cache、database的综合项目模板。

## 贡献名单
目前已经有几位朋友在为框架一起做努力，我们将在合适的时间向大家展现，谢谢他们的支持！

## Contact Us
#### QQ-Group：193409346 - <a target="_blank" href="http://shang.qq.com/wpa/qunwpa?idkey=836e11667837ad674462a4a97fb21fba487cd3dff5b2e1ca0d7ea4c2324b4574"><img border="0" src="http://pub.idqqimg.com/wpa/images/group.png" alt="Golang-Devfeel" title="Golang-Devfeel"></a>
#### Gitter：[![Gitter](https://badges.gitter.im/devfeel/dotweb.svg)](https://gitter.im/devfeel-dotweb/wechat)
