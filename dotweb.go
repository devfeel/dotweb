package dotweb

import (
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/framework/log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
)

type Dotweb struct {
	HttpServer       *HttpServer
	Modules          []*HttpModule
	logpath          string
	ExceptionHandler ExceptionHandle
}

type ExceptionHandle func(*HttpContext, interface{})

/*
* 创建DotServer实例，返回指针
 */
func New() *Dotweb {
	dotweb := &Dotweb{
		HttpServer: NewHttpServer(),
		Modules:    make([]*HttpModule, 0, 10),
	}
	dotweb.HttpServer.setDotweb(dotweb)

	return dotweb
}

/*
* 添加处理模块
 */
func (ds *Dotweb) RegisterModule(module *HttpModule) {
	ds.Modules = append(ds.Modules, module)
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
	port := ":" + strconv.Itoa(httpport)
	logger.Log("Dotweb:StartServer["+port+"] begin", LogTarget_HttpServer, LogLevel_Debug)

	//添加框架默认路由规则
	//默认支持pprof信息查看
	ds.HttpServer.GET("/dotweb/debug/pprof/:key", initPProf)
	ds.HttpServer.GET("/dotweb/debug/freemem", freeMemory)
	ds.HttpServer.GET("/dotweb/state", showServerState)
	ds.HttpServer.GET("/dotweb/query/:key", showQuery)

	err := http.ListenAndServe(port, ds.HttpServer)
	return err
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
