##dotweb版本记录：

#### Version 0.3.11
* 新增Group支持，可通过HttpServer.Group()或者NewGroup()创建，目前Group支持HEAD\GET\POST\PUT\OPTIONS\PATCH\DELETE路由方法
* Middleware支持三个级别Use：DotWeb.Use\Group.Use\RouterNode.Use
* 优化Router实现，移除router目录，整合xRouter和router.Router
* 废弃：移除RouterNode级别的Feature支持，原RouterNode.Features.SetEnabledCROS 可通过自实现Middleware实现
* 更新 example/middleware 目录
```go
func InitRoute(server *dotweb.HttpServer) {
	server.Router().GET("/", Index)
	server.Router().GET("/use", Index).Use(NewAccessFmtLog("Router-use"))

	g := server.Group("/group").Use(NewAccessFmtLog("group"))
	g.GET("/", Index)
	g.GET("/use", Index).Use(NewAccessFmtLog("group-use"))
}
```
* 2017-05-6 00:30

#### Version 0.3.10
* 新增Middleware支持，通过DotWeb.Use引入中间件，处理请求顺序在HttpModule之后，响应请求顺序在HttpModule之前
* 用户实现Middleware interface开发自定义middleware，为简化代码，可基于BaseMiddlware开发
* 内建支持RequestLogMiddleware，通过DotWeb.UseRequestLog()引入
* 优化路由注册日志，添加Handler名称，例如：main.Index
* 更新 example/middleware 目录
* 2017-04-23 09:00

#### Version 0.3.9.7
* 新增Json格式配置文件支持
* 更新 example/config 目录
* 2017-04-19 09:00

#### Version 0.3.9.6
* 优化精简HttpContext，移除部分非常用冗余方法，迁移至新增Request类型
* Request增加PostBody缓存，可通过Request.PostBody()重复获取
* 迁移：HttpContext.Proto()\Method()\FullRemoteIP()\Referer()\UserAgent()\Path()\Host()\IsAJAX 迁移至HttpContext.Request
* BUG修复：修复HttpContext释放时，Response释放不彻底问题
* 2017-04-12 09:00

#### Version 0.3.9.5
* 新增Feature模块，目前支持CROS特性，可通过HttpServer全局设置、RouterNode路由级别设置
* 新增CROS设置，设置是否允许跨域访问，并可设置AllowedOrigins、、AllowedMethods、AllowedHeaders、AllowCredentials、AllowedP3P属性
* PProfServer调整为属性设置，通过SetPProfConfig设置，移除StartPProf方法
* HttpContext新增RouterNode属性，可通过该属性获取路由信息
* 迁移：原DotWeb中SetSessionConfig函数迁移至HttpServer
* 废弃：DotWeb.SetSessionConfig、DotWeb.StartPProfServer
* 配置App节点新增EnabledPProf\PProfPort设置，用于设置PProfServer
* 完善配置文件与代码，完善代码，完善example
* 2017-04-09 20:00

#### Version 0.3.9.4
* 新增Logger接口，增加内置文件日志xlog实现
* 配置新增app节点，支持LogPath、EnabledLog、RunMode，DotWeb增加SetEnabledLog方法；
* 新增RunMode，目前支持development、production，相关特性后续完善，默认为development
* 迁移：原DotWeb中SetEnabledListDir、SetEnabledSession、SetEnabledGzip函数迁移至HttpServer
* 完善配置文件与代码，完善代码，完善example
* 2017-04-04 18:00

#### Version 0.3.9.3
* 调整Cache接口，Exists\Get\GetString\GetInt\GetInt64增加error返回值，影响Cache模块、Session模块
* Redis连接异常信息将从接口error返回，取消原来的panic设计
* 2017-04-04 16:00

#### Version 0.3.9.2
* Render增加SetTemplatePath接口，用于设置模板默认目录, 默认添加base、base/templates、base/views
* 模板查找顺序从最后一个插入的元素开始往前找
* SetTemplatePath(path ...string)
* 2017-04-01 10:30

#### Version 0.3.9.1
* Router增加MatchPath接口，用于检查指定请求与指定路由是否匹配，一般用于HttpModule中
* ctx.HttpServer.Router().MatchPath(ctx, "/user")
* 2017-03-31 22:30

#### Version 0.3.9
* 增强静态文件处理，增加设置目录浏览开关，通过Router.ServerFile使用
* Router增加Any路由方式，如果设置为Any，则自动匹配：GET\HEAD\POST\PUT\PATCH\OPTIONS\DELETE
* ServerConfig增加EnabledListDir选项，仅对Router.ServerFile有效，若设置该项，则可以浏览目录文件，默认不开启
* ServerConfig增加EnabledAutoHEAD选项，如果设置该选项，则会为除Websocket\HEAD外所有路由方式默认添加HEAD路由，默认不开启
* 优化部分代码，增加example\static、example\router
* 2017-03-26 19:30

#### Version 0.3.8
* 新增模板渲染支持，新增Renderer接口，需实现Render(io.Writer, string, interface{}, *HttpContext) error
* 支持接入第三方模板引擎，默认使用http/template内建引擎
* HttpContext增加ViewData，用于向模板传递值
* HttpContext增加View()方法，传入模板文件名，用于简化模板使用：View(name string) error
* 增加example\render
* 2017-03-26 11:10

#### Version 0.3.7
* 新增Cache模块，支持runtime、redis，常规支持Get\Set\GetString\GetInt\GetInt64\Incr\Decr\Exists\Clear操作
* Cache通过app.SetCache启用，通过HttpContext.Cache()使用
* ItemContext增加Remove、Once接口，影响HttpContext.AppContext、HttpContext.Items
* 增加example\cache，优化部分代码
* 2017-03-25 22:30

#### Version 0.3.6.5
* App、Server、Session配置代码梳理优化
* 新增Server接口，新增OfflineServer实现
* 新增Router接口，将原HttpServer路由操作转移至HttpServer.Router()
* 2017-03-24 14:00

#### Version 0.3.6
* 新增Config支持, 支持Server、Session、Router三类配置
* 新增HttpServer.RegisterHandler、HttpServer.GetHandler用于注册HttpHandle
* 新增Offline配置，可设置Offline欢迎语或者默认跳转维护页
* 新增example\config
* 2017-03-23 22:00

#### Version 0.3.5.1
* 新增ItemContext支持 - 单次请求完整通道内共享容器,通过HttpContext.Items().Get\Set操作
* 可在OnBeginRequest、UserHandler、OnEndRequest三个环节使用
* 新增example\httpmodule
* 2017-03-21 15:00

#### Version 0.3.5
* 新增AppContext支持 - 进程内共享容器,通过HttpContext.AppContext.Get\Set操作
* 新增文件上传支持, 新增HttpContext.FormFile("key"),返回*UploadFile
* 新增example\appcontext & example\uploadfile
* 2017-03-21 15:00

#### Version 0.3.4
* 新增Bind支持 - HttpContext.Bind(interface{}) 支持json、xml、Form数据
* 2017-03-16 11:00

#### Version 0.3.3
* 新增Gzip支持，通过app.SetEnabledGzip(true)开启，默认关闭
* 2017-03-14 13:00

#### Version 0.3.2
* 新增Session模块实现 - runtime\redis 两种模式
* 2017-03-13 14:00

#### Version 0.3.1
* 新增HttpServer.RegisterRoute函数，新增example目录
* 2017-03-09 10:00

#### Version 0.3
* 新增Hijack支持
* 2017-01-10 16:00

#### Version 0.2
* 新增WebSocket协议支持
* 2016-11-29 16:00

#### Version 0.1.4
* 新增HttpContext对Cookie的操作函数：WriteCookie\WriteCookieObj\RemoveCookie\ReadCookie\ReadCookieObj
* 启用pool优化对象管理，节约GC压力
* 2016-11-23 15:00

#### Version 0.1.3
* 将fasthttp引擎变更为原生nethttp引擎
* 升级go版本从1.6->1.7.3
* 2016-11-15 16:00

#### Version 0.1.2
* 增加HttpModule、ExceptionHandle、NotFound、Error特性支持
* 2016-11-14 15:00

#### Version 0.1.1
* 扩展HttpContext支持方法
* 2016-11-10 11:00

#### Version 0.1
* 初始版本，基于fasthttp引擎
* 2016-11-09 11:00
