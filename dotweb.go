package dotweb

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strconv"
	"strings"

	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/logger"
	"github.com/devfeel/dotweb/servers"
	"github.com/devfeel/dotweb/session"
)

type (
	DotWeb struct {
		HttpServer       *HttpServer
		cache            cache.Cache
		OfflineServer    servers.Server
		Config           *config.Config
		Modules          []*HttpModule
		ExceptionHandler ExceptionHandle
		AppContext       *core.ItemContext
	}

	ExceptionHandle func(*HttpContext, interface{})
	NotFoundHandle  func(*HttpContext)
)

const (
	DefaultHttpPort     = 80 //default http port
	RunMode_Development = "development"
	RunMode_Production  = "production"
)

/*
* 创建DotServer实例，返回指针
 */
func New() *DotWeb {
	app := &DotWeb{
		HttpServer:    NewHttpServer(),
		OfflineServer: servers.NewOfflineServer(),
		Modules:       make([]*HttpModule, 0, 10),
		AppContext:    core.NewItemContext(),
		Config:        config.NewConfig(),
	}
	app.HttpServer.setDotApp(app)
	return app
}

/*
* return cache interface
 */
func (app *DotWeb) Cache() cache.Cache {
	return app.cache
}

/*
* set cache interface
 */
func (app *DotWeb) SetCache(ca cache.Cache) {
	app.cache = ca
}

//current app run mode, if not set, default set RunMode_Development
func (app *DotWeb) RunMode() string {
	if app.Config.App.RunMode != RunMode_Development && app.Config.App.RunMode != RunMode_Production {
		app.Config.App.RunMode = RunMode_Development
	}
	return app.Config.App.RunMode
}

//check current run mode is development mode
func (app *DotWeb) IsDevelopmentMode() bool {
	return app.RunMode() == RunMode_Development
}

//set run mode on development mode
func (app *DotWeb) SetDevelopmentMode() {
	app.Config.App.RunMode = RunMode_Development
}

//set run mode on production mode
func (app *DotWeb) SetProductionMode() {
	app.Config.App.RunMode = RunMode_Production
}

/*
* 添加处理模块
 */
func (app *DotWeb) RegisterModule(module *HttpModule) {
	app.Modules = append(app.Modules, module)
	module.Server = app.HttpServer
}

/*
* 设置异常处理函数
 */
func (app *DotWeb) SetExceptionHandle(handler ExceptionHandle) {
	app.ExceptionHandler = handler
}

//设置pprofserver启动配置，默认不启动，且该端口号请不要与StartServer的端口号一致
func (app *DotWeb) SetPProfConfig(enabledPProf bool, httpport int) {
	app.Config.App.EnabledPProf = enabledPProf
	app.Config.App.PProfPort = httpport
}

//set log root path
func (app *DotWeb) SetLogPath(path string) {
	logger.Logger().SetLogPath(path)
}

//set enabled log flag
func (app *DotWeb) SetEnabledLog(enabledLog bool) {
	logger.Logger().SetEnabledLog(enabledLog)
}

/*启动WebServer
* 需要初始化HttpRoute
* httpPort := 80
 */
func (app *DotWeb) StartServer(httpport int) error {
	//添加框架默认路由规则
	//默认支持pprof信息查看
	app.HttpServer.Router().GET("/dotweb/debug/pprof/:key", initPProf)
	app.HttpServer.Router().GET("/dotweb/debug/freemem", freeMemory)
	app.HttpServer.Router().GET("/dotweb/state", showServerState)
	app.HttpServer.Router().GET("/dotweb/query/:key", showQuery)

	if app.ExceptionHandler == nil {
		app.SetExceptionHandle(app.DefaultHTTPErrorHandler)
	}

	//init session manager
	if app.HttpServer.SessionConfig.EnabledSession {
		if app.HttpServer.SessionConfig.SessionMode == "" {
			//panic("no set SessionConfig, but set enabledsession true")
			logger.Logger().Warn("not set SessionMode, but set enabledsession true, now will use default runtime session", LogTarget_HttpServer)
			app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
		}
		app.HttpServer.InitSessionManager()
	}

	//if cache not set, create default runtime cache
	if app.Cache() == nil {
		app.cache = cache.NewRuntimeCache()
	}

	//if renderer not set, create inner renderer
	if app.HttpServer.Renderer() == nil {
		app.HttpServer.SetRenderer(NewInnerRenderer())
	}

	//start pprof server
	if app.Config.App.EnabledPProf {
		if app.Config.App.PProfPort == httpport {
			errStr := "PProf Server and HttpServer have the same port"
			logger.Logger().Warn("Dotweb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] failed: "+errStr, LogTarget_HttpServer)
		} else {
			logger.Logger().Debug("Dotweb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] Begin", LogTarget_HttpServer)
			go func() {
				err := http.ListenAndServe(":"+strconv.Itoa(app.Config.App.PProfPort), nil)
				if err != nil {
					logger.Logger().Warn("Dotweb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] error: "+err.Error(), LogTarget_HttpServer)
				}
			}()
		}
	}

	port := ":" + strconv.Itoa(httpport)
	if app.Config.App.EnabledLog {
		logger.Logger().Log("Dotweb:StartServer["+port+"] begin", LogTarget_HttpServer, LogLevel_Debug)
	}
	err := http.ListenAndServe(port, app.HttpServer)
	return err
}

//start server with appconfig
func (app *DotWeb) StartServerWithConfig(config *config.Config) error {
	app.Config = config

	//log config
	if config.App.LogPath != "" {
		logger.Logger().SetLogPath(config.App.LogPath)
	}
	logger.Logger().SetEnabledLog(config.App.EnabledLog)

	//run mode config
	if app.Config.App.RunMode != RunMode_Development && app.Config.App.RunMode != RunMode_Production {
		app.Config.App.RunMode = RunMode_Development
	} else {
		app.Config.App.RunMode = RunMode_Development
	}

	//CROS Config
	if config.Server.EnabledAutoCORS {
		app.HttpServer.Features.SetEnabledCROS()
	}

	app.HttpServer.SetEnabledGzip(config.Server.EnabledGzip)

	//设置维护
	if config.Offline.Offline {
		app.HttpServer.SetOffline(config.Offline.Offline, config.Offline.OfflineText, config.Offline.OfflineUrl)
		app.OfflineServer.SetOffline(config.Offline.Offline, config.Offline.OfflineText, config.Offline.OfflineUrl)
	}

	//设置session
	if config.Session.EnabledSession {
		app.HttpServer.SetEnabledSession(config.Session.EnabledSession)
		app.HttpServer.SetSessionConfig(session.NewStoreConfig(config.Session.SessionMode, config.Session.Timeout, config.Session.ServerIP))
	}

	//load router and register
	for _, v := range config.Routers {
		if h, isok := app.HttpServer.Router().GetHandler(v.HandlerName); isok && v.IsUse {
			app.HttpServer.Router().RegisterRoute(strings.ToUpper(v.Method), v.Path, h)
		}
	}

	//start server
	port := config.Server.Port
	if port <= 0 {
		port = DefaultHttpPort
	}
	return app.StartServer(port)

}

//default exception handler
func (ds *DotWeb) DefaultHTTPErrorHandler(ctx *HttpContext, errinfo interface{}) {
	//输出内容
	ctx.Response.WriteHeader(http.StatusInternalServerError)
	ctx.Response.Header().Set(HeaderContentType, CharsetUTF8)
	//if in development mode, output the error info
	if ds.IsDevelopmentMode() {
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
