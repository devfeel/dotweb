package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/framework/log"
	"github.com/devfeel/dotweb/session"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
)

type Dotweb struct {
	HttpServer       *HttpServer
	SessionConfig    *session.StoreConfig
	Modules          []*HttpModule
	logpath          string
	ExceptionHandler ExceptionHandle
}

type ExceptionHandle func(*HttpContext, interface{})

/*
* 创建DotServer实例，返回指针
 */
func New() *Dotweb {
	application := &Dotweb{
		HttpServer: NewHttpServer(),
		Modules:    make([]*HttpModule, 0, 10),
	}
	application.HttpServer.setDotweb(application)
	return application
}

/*
* 添加处理模块
 */
func (ds *Dotweb) RegisterModule(module *HttpModule) {
	ds.Modules = append(ds.Modules, module)
}

/*
设置Debug模式,默认为false
*/
func (ds *Dotweb) SetEnabledDebug(isEnabled bool) {
	ds.HttpServer.ServerConfig.EnabledDebug = isEnabled
}

/*
设置是否启用Session,默认为false
*/
func (ds *Dotweb) SetEnabledSession(isEnabled bool) {
	ds.HttpServer.ServerConfig.EnabledSession = isEnabled
}

/*
设置是否启用gzip,默认为false
*/
func (ds *Dotweb) SetEnabledGzip(isEnabled bool) {
	ds.HttpServer.ServerConfig.EnabledGzip = isEnabled
}

//set session store config
func (ds *Dotweb) SetSessionConfig(config *session.StoreConfig) {
	ds.SessionConfig = config
}

/*
* 设置异常处理函数
 */
func (ds *Dotweb) SetExceptionHandle(handler ExceptionHandle) {
	ds.ExceptionHandler = handler
}

/*
* 启动pprof服务，该端口号请不要与StartServer的端口号一致
 */
func (ds *Dotweb) StartPProfServer(httpport int) error {
	port := ":" + strconv.Itoa(httpport)
	err := http.ListenAndServe(port, nil)
	return err
}

/*
* 设置日志根目录
 */
func (ds *Dotweb) SetLogPath(path string) {
	ds.logpath = path
}

/*启动WebServer
* 需要初始化HttpRoute
* httpPort := 80
 */
func (ds *Dotweb) StartServer(httpport int) error {
	//启动内部日志
	logger.StartLogHandler(ds.logpath)

	//添加框架默认路由规则
	//默认支持pprof信息查看
	ds.HttpServer.GET("/dotweb/debug/pprof/:key", initPProf)
	ds.HttpServer.GET("/dotweb/debug/freemem", freeMemory)
	ds.HttpServer.GET("/dotweb/state", showServerState)
	ds.HttpServer.GET("/dotweb/query/:key", showQuery)

	if ds.ExceptionHandler == nil {
		ds.SetExceptionHandle(ds.DefaultHTTPErrorHandler)
	}

	//init session manager
	if ds.HttpServer.ServerConfig.EnabledSession {
		if ds.SessionConfig == nil {
			panic("no set SessionConfig, but set enabledsession true")
		}
		ds.HttpServer.InitSessionManager(ds.SessionConfig)
	}

	port := ":" + strconv.Itoa(httpport)
	logger.Log("Dotweb:StartServer["+port+"] begin", LogTarget_HttpServer, LogLevel_Debug)
	err := http.ListenAndServe(port, ds.HttpServer)
	return err
}

//默认异常处理
func (ds *Dotweb) DefaultHTTPErrorHandler(ctx *HttpContext, errinfo interface{}) {
	//输出内容
	ctx.Response.WriteHeader(http.StatusInternalServerError)
	ctx.Response.Header().Set(HeaderContentType, CharsetUTF8)
	if ds.HttpServer.ServerConfig.EnabledDebug {
		ctx.WriteString(fmt.Sprintln(errinfo))
	} else {
		ctx.WriteString("Internal Server Error")
	}
}

//query pprof debug info
//key:heap goroutine threadcreate block
func initPProf(ctx *HttpContext) {
	querykey := ctx.GetRouterName("key")
	runtime.GC()
	pprof.Lookup(querykey).WriteTo(ctx.Response.Writer(), 1)
}

func freeMemory(ctx *HttpContext) {
	debug.FreeOSMemory()
}

//显示服务器状态信息
func showServerState(ctx *HttpContext) {
	ctx.WriteString(jsonutil.GetJsonString(GlobalState))
}

//显示服务器状态信息
func showQuery(ctx *HttpContext) {
	querykey := ctx.GetRouterName("key")
	switch querykey {
	case "state":
		ctx.WriteString(jsonutil.GetJsonString(GlobalState))
	case "":
		ctx.WriteString("please input key")
	default:
		ctx.WriteString("not support key => " + querykey)
	}
}
