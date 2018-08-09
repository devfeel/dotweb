## dotweb版本记录：

#### Version 1.5.5
* New feature: /dotweb/state add CurrentRequestCount data
* Update: improve 30% performance on app's metric
* 2018-08-09 15:00

#### Version 1.5.4
* New feature: Session.StoreConfig support CookieName, used to set custom cookie name which sessionid store, default is dotweb_sessionId
* Update: Config.SessionNode add CookieName, used to set custom cookie name which sessionid store
* Update: default log format update to "Time [LogLevel] [FileName:Line] Content"
* BugFixed: remove init session which exec on dotweb.initAppConfig
* 2018-08-02 15:00

#### Version 1.5.3
* New feature: HttpServer add Validator which be called by Context.Validate()
* New feature: Context add Validate(interface{}) used to validate data with HttpServer::Validator
* Update: use routerExpressSplit replace "_" when padding data to Router::RouterExpress
* 2018-07-12 12:00

#### Version 1.5.2
* New feature: dotweb.innerRenderer add cache mode, default is enabled
* New feature: dotweb.innerRenderer add NewInnerRendererNoCache() used to disabled cache
* Update for app run_mode: if it's develop run mode, the default renderer will use no cache mode
* 2018-06-22 14:00

#### Version 1.5.1
* Fixed Bug: double sprintf on logger.xlog
* 2018-06-15 14:00

#### Version 1.5
* New feature：UploadFile.RandomFileName used to get random file name from uuid
* New feature: encodes.base64x used to quickly use base64 EncodeString and DecodeString
* New feature: session.NewRedisConfig use to create redis session config include serverUrl and storeKeyPre
* if you set custom storeKeyPre, it will be store key start with custom set, default is "dotweb:session:"
* 2018-06-08 17:00

#### Version 1.4.9.7
* fixed bug for #114 dotweb: fix error found by vet
* 2018-06-04 09:00

#### Version 1.4.9.6
* 调整DevelopmentMode设置:默认开启Log,默认启用RequestLogMiddleware,默认Console Print日志内容
* fixed for #122 - dotweb没有打出 access log
* 合并const目录内容至consts文件,统一const定义文件
* 移除example/static
* 新增example/developmode,便于对developmode的理解
* 2018-06-04 09:00

#### Version 1.4.9.5
* framewrok.RedisClient增加Ping接口，用于检查服务是否可用
* GetSessionStore增加Redis ServerUrl的ping检查，若不可用直接panic异常信息
* 2018-05-23 13:30

#### Version 1.4.9.4
* 调整: dotweb.Classic移除自动UseRequestLog逻辑
* 调整：Session Redis模式时，新增sessionReExpire用于重置有效期，避免调用SessionUpdate导致不必要的redis数据交换
* 调整：Cache.Runtime调整map为sync.Map
* 调整：Session Redis模式时，gob.EncodeMap逻辑调整，移除自动注册interface{}
* 状态页：访问/dotweb/state时增加CurrentTime字段输出
* BUG： 修正Reponse自动释放时未重置body字段，导致内存溢出BUG
* 2018-05-10 13:30

#### Version 1.4.9.3
* 重要：go版本适配升级为1.9+
* 调整UploadFile.Size实现方法，直接返回Header.Size数据
* 调整dotweb.Classic签名为Classic(logPath string),若传入logPath为空，则默认以"bin-root/logs"为日志目录
* 新增dotweb.ClassicWithConf(config *config.Config),支持初始化传入config设置
* 调整默认Log目录由"bin-root"为"bin-root/logs"
* 2018-04-24 08:30

#### Version 1.4.9.2
* 修复BUG for #112
* 调整CharsetUTF8值为"charset=utf-8"
* 2018-03-28 09:00

#### Version 1.4.9.1
* 内置Vendor目录仅保留 golang.org/x/net包,移除redis与yaml包
* 完善RedisClient接口能力
* 2018-03-13 12:00

#### Version 1.4.9
* 新增: ConcurrenceMap增加GetTimeDuration接口，用于直接返回time.Duration类型数据
* 新增: ReadonlyMap增加GetTimeDuration接口，用于直接返回time.Duration类型数据
* 新增: dotweb增加UseTimeoutHook接口，用于设置请求处理超过设置时间钩子函数
* 调整：dotweb内建集成三类中间件，无需单独引用Middleware库
* RequestLogMiddleware
* TimeoutHookMiddleware
* dotweb.SetExceptionHandle
* 新增 example/basemiddleware
* 集成vendor: gopkg.in/yaml.v2
* 2018-02-12 12:00

#### Version 1.4.8
* 调整：ItemContext更名为ItemMap，新增ConcurrenceMap、ReadonlyMap接口
* 调整：Dotweb.AppContext变更为Dotweb.Items
* 调整：HttpContext.AppContext变更为HttpContext.AppItems
* 调整：HttpContext.AppSetConfig变更为HttpContext.ConfigSet
* 调整：config.AppSet变更为config.ConfigSet
* 新增: config.ParseConfigSetXML\ParseConfigSetJSON\ParseConfigSetYaml，用于解析常规Key\Value格式的配置文件
* 新增：config.Config.IncludeConfigSet，用于向config.ConfigSet中导入Key\Value格式的配置文件，通过HttpContext.ConfigSet获取相关设置信息
* ParseConfigSetXML:支持xml格式文件解析，返回core.ConcurrenceMap
* ParseConfigSetJSON:支持json格式文件解析，返回core.ConcurrenceMap
* ParseConfigSetYaml:支持yaml格式文件解析，返回core.ConcurrenceMap
* 具体配置文件格式可参考example/configset
* 新增示例代码 example/configset
* 2018-01-24 22:00

#### Version 1.4.7
* BUG Fixed: 修复Middleware特定场景下无效问题
* 新增dotweb.IncludeDotwebGroup，用于自动集成dotweb相关路由
* /dotweb/state 增加 ServerVersion输出
* 2018-01-22 22:00

#### Version 1.4.6
* BUG Fixed: UploadFile废弃Bytes接口，新增ReadBytes接口，用于返回上传文件本身
* 需要特别注意，由于io.read具有一次性特性，UploadFile.SaveFile与UploadFile.ReadBytes只能使用其中一个，另外一个将无法正常获取数据
* 增加dotweb.Version，用于输出框架版本号
* 2018-01-21 09:00

#### Version 1.4.5
* 新增yaml格式配置文件支持，具体参考 example/config/dotweb.yaml
* config新增UnmarshalYaml\MarshalYaml\MarshalYamlString，提供针对Yaml的常规处理
* config新增UnmarshalXML\MarshalXML\MarshalXMLString，提供针对Xml的常规处理
* config新增UnmarshalJSON\MarshalJSON\MarshalJSONString，提供针对Json的常规处理
* UploadFile新增Bytes接口，用于返回上传文件本身
* 完善 example/config
* 移除 example/session，查看更多示例，请移步https://github.com/devfeel/dotweb-example
* 2018-01-20 23:40

#### Version 1.4.4
* 调整Bind模块内部获取req.Body为HttpContext.PostBody,避免因为Bind后无法再次获取Post内容
* 完善debug log 输出
* 更新example/static, 该示例目前为实现一个纯静态文件服务器功能
* * 2018-01-08 12:00

#### Version 1.4.3
* 调整dotweb内部路由注册逻辑，New模式默认不开启，Classic模式默认开启，可通过app.UseDotwebRouter手动开启
* 修复 issue #100, 解决特定场景下Exclude不生效问题
* Use filepath.Ext to replace file.GetFileExt, update for issue #99
* 移除 framework/file.GetFileExt 函数
* 同步更新example代码
* 2018-01-07 22:00

#### Version 1.4.2
* Context新增QueryInt\QueryInt64接口，用于简化获取Int类型的Get参数，如果参数未传入或不是合法整形，返回0
* Context接口调整：除Write外，其他WriteXXX接口，返回值从(int, error)调整为error
* 带来的变化：
```
	app.HttpServer.GET("/index", func(ctx dotweb.Context) error{
		_, err := ctx.WriteString("welcome to my first web!")
		return err
	})
```
简化为
```
	app.HttpServer.GET("/index", func(ctx dotweb.Context) error{
		return ctx.WriteString("welcome to my first web!")
	})
```
* 同步更新example代码
* 2018-01-03 22:00

#### Version 1.4.1
* dotweb新增ExcludeUse接口，用于设置指定排除路由的中间件，fixed for #94 建议增加exuse方法用来实现exclude uri功能
* 调整dotweb部分内部函数名称
* 重构Middleware实现,优化代码结构，使流程更清晰
* Node新增AppMiddlewares与GroupMiddlewares，用于存储App级、Group级的中间件实例
* Middleware新增Exclude相关接口
* Exclude(routers ...string)函数，用于指定不生效的路由
* HasExclude() bool函数，用于判定当前中间件实例是否存在排除路由
* ExistsExcludeRouter(router string) bool函数，用于判定当前中间件实例的排除路由中是否存在指定路由
* middleware执行优先级：
- 优先级1：app级别middleware
- 优先级2：group级别middleware
- 优先级：router级别middleware
* 更新example/middleware
* 2018-01-01 23:00:00

#### Version 1.4.0
* 强化Bind能力，添加BindJsonBody接口，不对Content-Type进行判定，直接获取req.Body进行json序列化
* 基础Bind机制：
 1)如果识别到Content-Type = application/json  会尝试使用json序列化到struct
 2)如果识别到Content-Type = application/xml  会尝试使用xml序列化到struct
 3)非以上类型，会默认识别form的key/value模式匹配
* Context.Request()新增PostValues接口，包含 POST, PATCH, or PUT请求集合，以map[string][]string返回
* 更正Readme部分文字描述
* 完善example/bind
* 2017-12-29 16:00:00

#### Version 1.3.9
* Request新增FormFiles接口，用于支持多文件上传
* fixed #92 支持多文件同时上传
* Server支持SetEnabledBindUseJsonTag接口，用于设置bind是否启用json标签,默认为false
* fixed #91 GET方式请求，根据传参执行bind方法，字段tag目前只支持form，可以考虑兼容json标签
* 感谢 @githubityu
* 完善 example/bind
* 2017-11-30 17:00:00

#### Version 1.3.8
* 完善redisutil,增加对hset、set的操作支持
* 感谢 @chacha923 的pr
* 2017-11-23 00:00

#### Version 1.3.7
* 新增Classic初始化App方式，调整New初始化内容
* 调整日志启用开关默认为false，可以通过SetEnabledLog(true)或者SetDevelopmentMode()启用
* dotweb.New()仅初始化必要的组件，移除打印Logo，移除启动日志
* dotweb.Classic()方式下，默认启用日志，默认启用RequestLog中间件，默认打印logo以及启动日志
* dotweb.SetDevelopmentMode()在原有启用屏幕打印的基础上，增加默认开启日志开关
* 新增 example/start
* 2017-11-10 21:30

#### Version 1.3.6
* 升级Redis库，当前版本为redigo V1.2.0
* 修改Redis工具类中，pool的创建方法由Dial调整为DialUrl
* 废弃：Redis连接字符串形式10.0.1.11:6379 不可再使用
* 变更：Redis连接字符串形式变更为：redis://:password@10.0.1.11:6379/0
* 更新 example/session
* 2017-11-08 09:00

#### Version 1.3.5
* 提升UUID性能，调整UUID算法，新增UUID包
* HttpServer新增设置是否启用唯一请求ID，默认不开启，开启后使用32位UUID，通过ctx.RequestID获取
* 提供HttpServer.SetEnabledRequestID()函数
* 提供uuid.V1\V2\V3\V4四类算法
* move example/exception to devfeel/dotweb-example
* move example/uploadfile to devfeel/dotweb-example
* add cookbook link
* 2017-11-02 16:00

#### Version 1.3.4
* 增强Render实现，新增多模板支持
* Context新增AddView接口，用于传入多模板
* 模板统一基于Renderer().SetTemplatePath设置的相对文件名
* 特别的：布局模板需放在第一个被解析，比如：
* ctx.AddView("_layout/layout.html")
  err := ctx.View("dashbord/index.html")
* 2017-10-31 10:00

#### Version 1.3.4
* 增强Render实现，新增多模板支持
* Context新增AddView接口，用于传入多模板
* 模板统一基于Renderer().SetTemplatePath设置的相对文件名
* 特别的：布局模板需放在第一个被解析，比如：
* ctx.AddView("_layout/layout.html")
  err := ctx.View("dashbord/index.html")
* 2017-10-31 10:00

#### Version 1.3.3
* fixed for #70 UPDATE default http port 80 to 8080
* 增加StandardHandle，统一自定义Handle场景
* 2017-10-14 20:00

#### Version 1.3.2
* 主要新增设置启用详细请求数据控制，新增MethodNotAllowedHandler自定义能力，完善数据统计逻辑
* 新增设置启用详细请求数据控制：
* 1) ServerStateInfo增加EnabledDetailRequestData，用于控制是否启用详细请求数据统计
* 2) config.ServerConfig增加EnabledDetailRequestData设置，支持配置文件控制
* 3) HttpServer增加SetEnabledDetailRequestData函数，用于代码设置是否启用详细请求数据统计
* 4) fixed #63 状态数据，当url较多时，导致内存占用过大
* 新增MethodNotAllowedHandler自定义能力
* 1) dotweb增加DefaultMethodNotAllowedHandler、SetMethodNotAllowedHandle函数
* 2) fixed #64 增加MethodNotAllowed自定义处理
* 完善数据统计逻辑
* 1) 404请求不计入详细请求数据统计
* 2017-10-13 12:00

#### Version 1.3.1
* 增加HttpServer.ServerFile，以支持server快捷设置静态文件目录
* 2017-10-11 22:00

#### Version 1.3
* 国庆迭代第三波
* 增加TLS支持
* 1)ServerConfig增加EnabledTLS\TLSCertFile\TLSKeyFile三个配置项
* 2)HttpServer增加SetEnabledTLS方法
* 3)增加example/tls目录
* 移除：移除DotWeb.StartWithConfig方法，使用配置方式时需显式调用dotweb.SetConfig方法
* 进一步完善默认日志输出
* 2017-10-07 22:00

#### Version 1.2.1
* 启动日志增加dotweb ASCII art logo
* 调整默认日志的时间毫秒段为4位
* 2017-10-05 22:00

#### Version 1.2
* 国庆特别版 - 2017年的国庆，恰逢中秋双节，祝愿国家繁荣昌盛，祝愿代码艺术家们双节快乐！
* 主要新增HttpModule模块、新增ServerConfig：EnabledIgnoreFavicon、完善HttpContext、完善debug日志输出
* 【新增】HttpModule模块，提供路由之前全局级别的自定义代码能力
* 1、新增module文件，将原dotweb文件中HttpMoule定义转移至module文件，新增Name属性
* 2、新增getIgnoreFaviconModule函数用于创建基础Module用于EnabledIgnoreFavicon特性
* 3、HttpServer新增RegisterModule函数，用于注册自定义HttpModule
* 4、增加代码示例-根据参数动态变更访问路由 - package: example/httpmodule
* 【新增】忽视favicon文件请求能力
* 1、通过ServerConfig：EnabledIgnoreFavicon设置，默认不启用
* 2、当启用该选项时，请求/facicon.ico文件时，将默认反馈空
* 3、配置文件server节点增加enabledignorefavicon设置项
* 【完善】HttpContext，使相关接口更清晰
* 1、新增WriteHtml、WriteHtmlC函数，默认设置Context-Type:text/html; charset=utf-8
* 3、注意，若在同一个请求期间，多次调用WriteHtml与WriteString，以最后一次调用的Context-Type为准
* 【完善】debug日志输出，使dotweb启动过程更清晰
* 1、新增Group创建日志
* 2、新增ServerConfig设置日志
* 3、新增设置Session参数日志
* 4、新增Session初始化日志
* 5、新增RegisterModule日志
* 2017-10-03 14:00

#### Version 1.1
* 主要新增集成基础统计数据
* 【新增】dotweb/state接口，提供基础统计数据，主要数据说明：
* 1、ServerStartTime 服务启动时间
* 2、TotalRequestCount 服务启动以来累计请求数（排除了"/dotweb/"下系统自有页面的访问数）
* 3、TotalErrorCount 服务启动以来累计错误数
* 4、IntervalRequestData 按1分钟为间隔，默认保存最近60分钟，每分钟的请求数（排除了"/dotweb/"下系统自有页面的访问数）
* 5、DetailRequestPageData 服务启动以来，每个页面的累计请求数（排除了"/dotweb/"下系统自有页面的访问数）
* 6、IntervalErrorData 按1分钟为间隔，默认保存最近60分钟，每分钟的错误数
* 7、DetailErrorPageData 服务启动以来，每个页面的累计错误数
* 8、DetailErrorData 服务启动以来，每个异常的累计数
* 9、DetailHttpCodeData 服务启动以来，每个Http状态码的累计数
* 10、可通过 {host}/dotweb/state 获取数据
* 【新增】dotweb/state/interval接口，提供按分钟级的基础数据查询
* 主要数据说明：
* 1、Time 表示查询时间的字符串，最小单位为分钟，例如：201709251200
* 2、RequestCount 单位时间内累计请求数（排除了"/dotweb/"下系统自有页面的访问数）
* 3、ErrorCount 单位时间内累计错误数
* 2017-09-25 13:00

#### Version 1.0.1
* 新增logger.SetEnabledConsole，当设置SetDevelopmentMode时自动开启
* SetEnabledConsole：自动向console输出dotweb基础日志，便于开发人员调试阶段实时查看
* 完善日志、注释
* 2017-08-31 08:30

#### Version 1.0.0
* dotweb正式进入1.0时代，感谢大家！
* HttpServer新增IndexPage(),返回当前Server默认首页名称，支持配置文件配置，如果未配置，自动读取默认值：DefaultIndexPage = "index.html"
* Context增加File、Attachment、Inline接口 感谢@esap issue #39
* Context.File(file string),将指定文件响应到客户端，如果指定为目录，则根据HttpServer.IndexPage，默认寻找该目录下同名文件，不存在则返回404
* Context.Attachment(file, name string),将指定文件响应到客户端下载，具体参考MIME协议扩展"Content-disposition"说明
* Context.Inline(file, name string),将指定文件响应到客户端输出，具体参考MIME协议扩展"Content-disposition"说明
* 新增example\file目录，提供对Context.File、Attachment、Inline三个接口的demo
* 2017-07-29 17:00


#### Version 0.3.21
* 新增DotWeb & HttpServer.ListenAndServer(addr string) 支持原生host监听
* 新增DotWeb Close() & Shutdown() 支持优雅关闭服务
* Context增加Context() & SetTimeoutContext() & WithContext()
* Request增加自动生成unique RequestID,通过RequestID()获取,默认置入dotweb.Context.context中，以RequestID为key
* 新增test目录，提供test公共方法，完善bind_test，感谢@zouyx
* 2017-06-14 23:10

#### Version 0.3.20
* 新增Gzip中间件，参考https://github.com/devfeel/middleware/tree/master/gzip
* 统一：所有中间件的New方式统一为package.Middleware()
* 2017-06-11 17:00

#### Version 0.3.19
* 新增AccessLog中间件，参考https://github.com/devfeel/middleware/tree/master/accesslog
* use dotweb.logger
* use debug level
* logtarget is "dotweb_accesslog", like "dotweb_accesslog_debug_2017_06_09.log"
* logformat: method userip proto status reqbytelen resbytelen timetaken
* log-example: [debug] 2017-06-09 08:38:10.416369 [middleware.go:49] / GET 127.0.0.1 HTTP/1.1 200 0 19 10
* 统一：所有中间件的New方式统一为package.Middleware()
* 2017-06-09 08:00

#### Version 0.3.18
* 新增cros中间件，参考 https://github.com/devfeel/middleware/tree/master/cros
* 支持较为丰富的option选项，可以根据实际业务场景自定义
* 支持设置Header：
* 1)Access-Control-Allow-Origin
* 2)Access-Control-Allow-Methods
* 3)Access-Control-Allow-Headers
* 4)Access-Control-Allow-Credentials
* 5)Access-Control-Expose-Headers
* 6)Access-Control-Max-Age
* 7)P3P
* HttpServer增加路由快捷方法:GET\POST\HEAD\OPTIONS\PUT\PATCH\DELETE\HiJack\WebSocket\ANY
* Rsponse增加QueryHeader方法
* 完善一些注释
* Demo: https://github.com/devfeel/middleware/tree/master/example/cros
* 2017-06-08 20:00

#### Version 0.3.17
* 新增jwt中间件，参考demo：example/jwt
* 支持较为丰富的option选项，可以根据实际业务场景自定义
* 可配置通过Cookie、Header、QueryString三类token存储方式
* 支持自定义非法token处理函数
* 支持token通过jwt基础验证后的扩展处理接口
* 新增依赖 引入 github.com/dgrijalva/jwt-go
* 新增 example/jwt 目录
* 2017-06-06 13:00

#### Version 0.3.16
* 配置文件加载优化，加载解析配置文件失败的情况下返回error对象
× 配置文件匹配完善，依次按三个目录匹配：1按绝对路径检查，2进程根目录下寻找，3进程根目录/config/下寻找
* 配置文件增加AppSet，用于应用自定义简单配置节点，配置文件增加appset节点，通过Context.AppSetConfig()访问
* 更新 example/config 目录
```xml
<appset>
    <set key="email-host" value="mail.xx.com" />
    <set key="email-account" value="username" />
    <set key="email-password" value="password" />
    <set key="limit-per-ip" value="10" />
</appset>
```
```go
//查询方式
ctx.AppSetConfig().GetString("email-host"))
ctx.AppSetConfig().GetInt("limit-per-ip"))
```
* 2017-05-22 17:30

#### Version 0.3.15
* 重要更新：HttpHandle增加error返回值，有error再也不用panic啦！
* 重要更新：新增Context接口
* ExceptionHandle func(*HttpContext, error) 调整为 ExceptionHandle func(Context, error)
* HttpHandle func(*HttpContext) 调整为 HttpHandle func(Context) error
* Middleware.Handle(ctx *HttpContext) error调整为Handle(ctx Context) error
* Middleware.Next(ctx *HttpContext)调整为Next(ctx Context)
* 为使Context设计与理解更清晰，做以下重要变更：
* 迁移Context.QueryStrings()至Context.Request().QueryStrings()
* 迁移Context.RawQuery()至Context.Request().RawQuery()
* 迁移Context.FormFile()至Context.Request().FormFile()
* 迁移Context.FormValues()至Context.Request().FormValues()
* 迁移Context.PostBody()至Context.Request().PostBody()
* 迁移Context.QueryHeader()至Context.Request().QueryHeader()
* 迁移Context.Url()至Context.Request().Url()
* 迁移Context.ContentType()至Context.Request().ContentType()
* 迁移Context.DelHeader()至Context.Response().DelHeader()
* 迁移Context.SetHeader()至Context.Response().SetHeader()
* 迁移Context.SetContentType()至Context.Response().SetContentType()
* 迁移Context.SetStatusCode()至Context.Response().SetStatusCode()
* 废弃Context.PostString(key string) string
* 通过在handler内返回相应的error，框架自动命中并通过设置的ExceptionHandle处理
* 更新 example/下所有代码，以支持error返回值，支持Context接口
* 2017-05-20 15:00

#### Version 0.3.14
* 配置文件支持Middleware，支持App\Group\Router三个级别
* DotWeb增加RegisterMiddlewareFunc\GetMiddlewareFunc方法，用于创建中间件的方法注册
* 默认日志增加文件名与行号输出 - logger/xlog.go - [贡献者: @淡定]
* 更新 example/config 目录
* 2017-05-17 15:00

#### Version 0.3.12
* ExceptionHandle定义由ExceptionHandle func(*HttpContext, interface{})调整为ExceptionHandle func(*HttpContext, error)
* HttpContext.SetCookie定义由SetCookie(cookie http.Cookie)调整为SetCookie(cookie *http.Cookie)
* HttpModule去除*Server属性
* Group新增RegisterRoute(method, path string, h HttpHandle) 方法
* 配置文件增加Group支持，可参考 example/config目录下xml与json示例配置文件
* 更新 example/main 文件、example/config 目录
* 2017-05-15 22:10

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
