- [安装与配置](#install)
- [框架架构](#arch)
  - [生命周期](#life-circle)
  - [Context](#context)
- [路由](#router)
  - [基本路由](#basic-router)
  - [路由参数](#router-param)
  - [路由群组](#router-group)
- [控制器](#controller)
- [请求](#request)
  - [请求头](#request-header)
  - [Cookies](#request-cookie)
  - [上传文件](#upload)
- [响应](#response)
  - [响应头](#response-header)
  - [附加Cookie](#response-cookie)
  - [字符串响应](#response-string)
  - [JSON响应](#response-json)
  - [视图响应](#response-view)
  - [文件下载](#response-file)
  - [重定向](#response-redirect)
  - [同步异步](#sync-async)
- [视图](#view)
  - [传参](#view-param)
  - [视图组件](#view-unit)
- [中间件](#middleware)
  - [分类使用](#middleware-use)
  - [创建中间件](#middleware-create)
  - [中间件参数](#middleware-param)
- [数据库](#db)
  - [Mongodb](#db-mongodb)
  - [Mysql](#db-mysql)
  - [SqlServer](#db-sqlserver)
- [扩展包](#extensions)
- [常用方法](#functions)

<a name="install"></a>
### 安装与配置
安装：

```sh
$ go get -u github.com/devfeel/dotweb
```
`
注意：确保 GOPATH GOROOT 已经配置
`

导入：
```go
import "github.com/devfeel/dotweb"
```

<a name="arch"></a>
### 框架架构
<a name="http"></a>

- HTTP 服务器

```go
func main() {
	app := dotweb.New()
	err := app.StartServer(80)
	if err !=nil{
		fmt.Println("dotweb.StartServer error => ", err)
	}
}
```

<a name="life-circle"></a>

- 生命周期

<a name="context"></a>

- Context

<a name="router"></a>

### 路由
<a name="basic-router"></a>

- 基本路由

dotweb 框架中路由是基于httprouter演变而来。


```go
	//初始化DotServer
	app := dotweb.New()
	//设置dotserver日志目录
	app.SetLogPath(file.GetCurrentDirectory())
	//设置路由
	app.HttpServer.GET("/", Index)
	app.HttpServer.GET("/d/:x/y", Index)
	app.HttpServer.GET("/any", Any)

```

<a name="router-param"></a>

- 路由参数

API 参数以冒号 : 后面跟一个字符串作为参数名称，可以通过 GetRouterName 方法获取路由参数的值。

```go
	app.HttpServer.GET("/hello/:name", func(ctx dotweb.Context) error {
		return ctx.WriteString("hello " + ctx.GetRouterName("name"))
	})
```

URL 参数通过 QueryString、QueryInt 或 QueryInt64 方法获取

```go
	// url 为 http://localhost:8080/hello?name=billy时
	// 输出 hello 123
	app.HttpServer.GET("/hello", func(ctx dotweb.Context) error {
		return ctx.WriteString("hello " + ctx.QueryString("name"))
	})
```

表单参数通过 PostFormValue 方法获取

```go
	app.HttpServer.POST("/user", func(ctx dotweb.Context) error {
		name := ctx.PostFormValue("name")
		age := ctx.PostFormValue("age")
		return ctx.WriteString("name is " + name + ", age is " + age)
	})
```

<a name="router-group"></a>

- 路由群组

```go
	//设置路由组
	userCenterGroup := app.HttpServer.Group("/usercenter")
	userCenterGroup.GET("/userinfo", getUserInfo)
	userCenterGroup.GET("/account",getUserAccount)
```

<a name="controller"></a>

### 控制器

<a name="binding"></a>

- 数据解析绑定

模型绑定可以将请求体绑定给一个类型，目前支持绑定的类型有 JSON, XML 和标准表单数据 (foo=bar&boo=baz)。
要注意的是绑定时需要给字段设置绑定类型的标签。比如绑定 JSON 数据时，设置 `json:"fieldname"`。
使用绑定方法时，dotweb 会根据请求头中  Content-Type  来自动判断需要解析的类型。如果你明确绑定的类型。

```go
// Binding from JSON
type User struct {
	Name string `form:"name"`
	Age  int    `form:"age"`
}

func main() {
	app := dotweb.New()
	// 绑定普通表单或json格式
	app.HttpServer.POST("/user", func(ctx dotweb.Context) error {
		user := new(User)
		if err := ctx.Bind(user); err != nil {
			return ctx.WriteString("Bind err:" + err.Error())
		}
		return ctx.WriteString("Bind:" + fmt.Sprint(user))
	})
	// 只绑定JSON的例子 ({"user": "manu", "age": 12})
	app.HttpServer.POST("/userjson", func(ctx dotweb.Context) error {
		user := new(User)
		if err := ctx.BindJsonBody(user); err != nil {
			return ctx.WriteString("Bind err:" + err.Error())
		}
		return ctx.WriteString("Bind:" + fmt.Sprint(user))
	})
	app.StartServer(8888)

}
```
<a name="request"></a>

### 请求
<a name="request-header"></a>

- 请求头

<a name="request-params"></a>

- 请求参数

<a name="request-cookie"></a>

- Cookies

<a name="upload"></a>

- 上传文件

```go

```
<a name="response"></a>

### 响应
<a name="response-header"></a>

- 响应头

<a name="response-cookie"></a>

- 附加Cookie

<a name="response-string"></a>

- 字符串响应

```go
	ctx.WriteString("")
	ctx.WriteStringC(http.StatusOK, "")
```
<a name="response-json"></a>

- JSON/Byte/Html响应

```go
	ctx.WriteBlob([]byte)
	ctx.WriteBlobC(http.StatusOK,[]byte)
	ctx.WriteHtml(html)
	ctx.WriteHtmlC(http.StatusOK,html)
	ctx.WriteJson(user)
	ctx.WriteJsonC(http.StatusOK, user)
	ctx.WriteJsonBlob([]byte)
	ctx.WriteJsonBlobC(http.StatusOK, []byte)

```
<a name="response-view"></a>

- 视图响应

使用 View() 方法来加载模板文件，默认当前程序根目录模板路径
```go
func main() {
	app := dotweb.New()
	//set default template path, support multi path
	//模板查找顺序从最后一个插入的元素开始往前找
	app.HttpServer.GET("/", TestView)
	//设置模板路径
	app.HttpServer.Renderer().SetTemplatePath("views/")
	app.HttpServer.GET("/", func(ctx dotweb.Context) error {
		ctx.ViewData().Set("data", "测试信息")
		//加载模板
		err := ctx.View("testview.html")
		return err
	})
	app.StartServer(8888)
}
```

模板结构定义

```html
<html>
	<h1>
		{{ .data }}
	</h1>
</html>
```
不同文件夹下模板名字可以相同，此时需要 View() 指定模板路径

```go
	ctx.View("/test/testview.html")
	app.HttpServer.GET("/", func(ctx dotweb.Context) error {
		ctx.ViewData().Set("data", "图书信息")
		type BookInfo struct {
			Name string
			Size int64
		}
		m := make([]*BookInfo, 5)
		m[0] = &BookInfo{Name: "book0", Size: 1}
		m[1] = &BookInfo{Name: "book1", Size: 10}
		m[2] = &BookInfo{Name: "book2", Size: 100}
		m[3] = &BookInfo{Name: "book3", Size: 1000}
		m[4] = &BookInfo{Name: "book4", Size: 10000}
		ctx.ViewData().Set("Books", m)
		//加载test文件夹下testview模板
		err := ctx.View("/test/testview.html")
		//加载test1文件夹下testview模板
		err := ctx.View("test1/testview.html")
		return err
	})
```

views/test/testview.html
```html
<html>
	<h1>{{.data}}</h1>
	<br>
	<b>Books:</b>
	<br>
	{{range .Books}}
	BookName => {{.Name}}; Size => {{.Size}}
</html>

```

<a name="response-file"></a>

- 文件响应

```go
	//相对路径
	app.HttpServer.ServerFile("/src/*filepath", "./var/www")
	//等价
	app.HttpServer.ServerFile("/src/*filepath", "var/www")

```
<a name="response-redirect"></a>

- 重定向

```go
	app.HttpServer.GET("/redirect", func(ctx dotweb.Context) error {
		//内部重定向
		ctx.Redirect(http.StatusMovedPermanently, "src/1.html")
		//外部的重定向
		ctx.Redirect(http.StatusMovedPermanently, "https://www.baidu.com")
		return nil
	})
```
<a name="sync-async"></a>

- 同步异步

goroutine 机制可以方便地实现异步处理


<a name="view"></a>

### 视图
<a name="view-param"></a>

- 传参

<a name="view-unit"></a>

- 视图组件

<a name="middleware"></a>

### 中间件
<a name="middleware-use"></a>


- 分类使用方式

	支持粒度：App、Group、RouterNode

	启用顺序：App -> Group -> RouterNode，同级别下按Use的引入顺序执行
```go

	// 1.全局中间件
	app.Use(cors.DefaultMiddleware())

	// 2.单路由的中间件，可以Use多个
	app.HttpServer.POST("/",user).Use(...).Use(...)

	// 3.群组路由的中间件
	userCenterGroup := app.HttpServer.Group("/usercenter").Use(...)
	// 或者这样用：
	userCenterGroup.Use(...).Use(...)
```

<a name="middleware-create"></a>
- 自定义中间件

```go

//定义
func Handle(ctx dotweb.Context) error {
	//处理前
	m.Next(ctx)
	//处理后
	return nil
}

```

更多自定义中间件参考： https://github.com/devfeel/middleware


<a name="middleware-param"></a>
- 中间件参数

- 内置中间件

	app.UseRequestLog()//简单请求日志中间件

	app.UseTimeoutHook()//超时处理中间件



<a name="db"></a>

### 数据库

参考 https://github.com/devfeel/database

<a name="db-mongodb"></a>

- Mongodb

<a name="db-mysql"></a>

- Mysql

<a name="db-sqlserver"></a>

- SqlServer

<a name="extensions"></a>

### 扩展包

<a name="functions"></a>

### 常用方法
