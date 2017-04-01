##dotweb版本记录：

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
