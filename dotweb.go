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
	"sync"
)

type (
	DotWeb struct {
		HttpServer       *HttpServer
		cache            cache.Cache
		OfflineServer    servers.Server
		Config           *config.Config
		Modules          []*HttpModule
		Middlewares      []Middleware
		ExceptionHandler ExceptionHandle
		NotFoundHandler  NotFoundHandle
		AppContext       *core.ItemContext
		middlewareMap    map[string]MiddlewareFunc
		middlewareMutex  *sync.RWMutex
	}

	ExceptionHandle func(Context, error)
	NotFoundHandle  http.Handler

	// Handle is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but has a special parameter *HttpContext contain all request and response data.
	HttpHandle func(Context) error
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
		HttpServer:      NewHttpServer(),
		OfflineServer:   servers.NewOfflineServer(),
		Modules:         make([]*HttpModule, 0),
		Middlewares:     make([]Middleware, 0),
		AppContext:      core.NewItemContext(),
		Config:          config.NewConfig(),
		middlewareMap:   make(map[string]MiddlewareFunc),
		middlewareMutex: new(sync.RWMutex),
	}
	app.HttpServer.setDotApp(app)

	//init logger
	logger.InitLog()
	return app
}

//register middleware with gived name & middleware
func (app *DotWeb) RegisterMiddlewareFunc(name string, middleFunc MiddlewareFunc) {
	app.middlewareMutex.Lock()
	app.middlewareMap[name] = middleFunc
	app.middlewareMutex.Unlock()
}

//get middleware with gived name
func (app *DotWeb) GetMiddlewareFunc(name string) (MiddlewareFunc, bool) {
	app.middlewareMutex.RLock()
	v, exists := app.middlewareMap[name]
	app.middlewareMutex.RUnlock()
	return v, exists
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

//Use registers a middleware
func (app *DotWeb) Use(m ...Middleware) {
	step := len(app.Middlewares) - 1
	for i := range m {
		if m[i] != nil {
			if step >= 0 {
				app.Middlewares[step].SetNext(m[i])
			}
			app.Middlewares = append(app.Middlewares, m[i])
			step++
		}
	}
}

//UseRequestLog register RequestLog middleware
func (app *DotWeb) UseRequestLog() {
	app.Use(&RequestLogMiddleware{})
}

/*
* 添加处理模块
 */
func (app *DotWeb) RegisterModule(module *HttpModule) {
	app.Modules = append(app.Modules, module)
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

//set user logger, the logger must implement logger.AppLog interface
func (app *DotWeb) SetLogger(log logger.AppLog) {
	logger.SetLogger(log)
}

//set log root path
func (app *DotWeb) SetLogPath(path string) {
	logger.SetLogPath(path)
}

//set enabled log flag
func (app *DotWeb) SetEnabledLog(enabledLog bool) {
	logger.SetEnabledLog(enabledLog)
}

//set config for app
func (app *DotWeb) SetConfig(config *config.Config) error {
	app.Config = config

	//log config
	if config.App.LogPath != "" {
		logger.SetLogPath(config.App.LogPath)
	}
	logger.SetEnabledLog(config.App.EnabledLog)

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

	//register app's middleware
	for _, m := range config.Middlewares {
		if m.IsUse {
			if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
				app.Use(mf())
			}
		}
	}

	//load router and register
	for _, r := range config.Routers {
		//fmt.Println("config.Routers ", i, " ", config.Routers[i])
		if h, isok := app.HttpServer.Router().GetHandler(r.HandlerName); isok && r.IsUse {
			node := app.HttpServer.Router().RegisterRoute(strings.ToUpper(r.Method), r.Path, h)
			//use middleware
			for _, m := range r.Middlewares {
				if m.IsUse {
					if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
						node.Use(mf())
					}
				}
			}
		}
	}

	//support group
	for _, v := range config.Groups {
		if v.IsUse {
			g := app.HttpServer.Group(v.Path)
			//use middleware
			for _, m := range v.Middlewares {
				if m.IsUse {
					if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
						g.Use(mf())
					}
				}
			}
			//init group's router
			for _, r := range v.Routers {
				if h, isok := app.HttpServer.Router().GetHandler(r.HandlerName); isok && r.IsUse {
					node := g.RegisterRoute(strings.ToUpper(r.Method), r.Path, h)
					//use middleware
					for _, m := range r.Middlewares {
						if m.IsUse {
							if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
								node.Use(mf())
							}
						}
					}
				}
			}
		}
	}
	return nil
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

	//add default httphandler with middlewares
	app.Use(&xMiddleware{})

	port := ":" + strconv.Itoa(httpport)
	logger.Logger().Log("Dotweb:StartServer["+port+"] begin", LogTarget_HttpServer, LogLevel_Debug)
	err := http.ListenAndServe(port, app.HttpServer)
	return err
}

//start server with appconfig
func (app *DotWeb) StartServerWithConfig(config *config.Config) error {

	err := app.SetConfig(config)
	if err != nil {
		return err
	}
	//start server
	port := config.Server.Port
	if port <= 0 {
		port = DefaultHttpPort
	}
	return app.StartServer(port)

}

//default exception handler
func (ds *DotWeb) DefaultHTTPErrorHandler(ctx Context, err error) {
	//输出内容
	ctx.Response().WriteHeader(http.StatusInternalServerError)
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	//if in development mode, output the error info
	if ds.IsDevelopmentMode() {
		stack := string(debug.Stack())
		ctx.WriteString(fmt.Sprintln(err) + stack)
	} else {
		ctx.WriteString("Internal Server Error")
	}
}

//query pprof debug info
//key:heap goroutine threadcreate block
func initPProf(ctx Context) error {
	querykey := ctx.GetRouterName("key")
	runtime.GC()
	pprof.Lookup(querykey).WriteTo(ctx.Response().Writer(), 1)
	return nil
}

func freeMemory(ctx Context) error {
	debug.FreeOSMemory()
	return nil
}

//显示服务器状态信息
func showServerState(ctx Context) error {
	ctx.WriteString(jsonutil.GetJsonString(GlobalState))
	return nil
}

//显示服务器状态信息
func showQuery(ctx Context) error {
	querykey := ctx.GetRouterName("key")
	switch querykey {
	case "state":
		ctx.WriteString(jsonutil.GetJsonString(GlobalState))
	case "":
		ctx.WriteString("please input key")
	default:
		ctx.WriteString("not support key => " + querykey)
	}
	return nil
}
