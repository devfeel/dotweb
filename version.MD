##dotweb版本记录：

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
