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

	"context"
	"errors"
	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/logger"
	"github.com/devfeel/dotweb/servers"
	"github.com/devfeel/dotweb/session"
	"reflect"
	"sync"
	"time"
)

type (
	DotWeb struct {
		HttpServer              *HttpServer
		cache                   cache.Cache
		OfflineServer           servers.Server
		Config                  *config.Config
		Middlewares             []Middleware
		ExceptionHandler        ExceptionHandle
		NotFoundHandler         StandardHandle // NotFoundHandler 支持自定义404处理代码能力
		MethodNotAllowedHandler StandardHandle // MethodNotAllowedHandler fixed for #64 增加MethodNotAllowed自定义处理
		Items                   core.ConcurrenceMap
		middlewareMap           map[string]MiddlewareFunc
		middlewareMutex         *sync.RWMutex
		StartMode               string
	}

	// ExceptionHandle 支持自定义异常处理代码能力
	ExceptionHandle func(Context, error)

	// StandardHandle 标准处理函数，需传入Context参数
	StandardHandle func(Context)

	// Handle is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but has a special parameter Context contain all request and response data.
	HttpHandle func(Context) error
)

const (
	// DefaultHTTPPort default http port; fixed for #70 UPDATE default http port 80 to 8080
	DefaultHTTPPort           = 8080

	DefaultLogPath = ""

	// RunMode_Development app runmode in development mode
	RunMode_Development       = "development"
	// RunMode_Production app runmode in production mode
	RunMode_Production        = "production"

	//StartMode_New app startmode in New mode
	StartMode_New     = "New"
	//StartMode_Classic app startmode in Classic mode
	StartMode_Classic = "Classic"
)

//New create and return DotApp instance
//default run mode is RunMode_Production
func New() *DotWeb {
	app := &DotWeb{
		HttpServer:      NewHttpServer(),
		OfflineServer:   servers.NewOfflineServer(),
		Middlewares:     make([]Middleware, 0),
		Items:           core.NewConcurrenceMap(),
		Config:          config.NewConfig(),
		middlewareMap:   make(map[string]MiddlewareFunc),
		middlewareMutex: new(sync.RWMutex),
		StartMode:       StartMode_New,
	}
	//set default run mode = RunMode_Production
	app.Config.App.RunMode = RunMode_Production
	app.HttpServer.setDotApp(app)
	//add default httphandler with middlewares
	//fixed for issue #100
	app.Use(&xMiddleware{})

	//init logger
	logger.InitLog()
	return app
}

// Classic create and return DotApp instance\
// if set logPath = "", it  will use bin-root/logs for log-root
// 1.SetEnabledLog(true)
// 2.use RequestLog Middleware
// 3.print logo
func Classic(logPath string) *DotWeb {
	app := New()
	app.StartMode = StartMode_Classic

	if logPath != ""{
		app.SetLogPath(logPath)
	}
	app.SetEnabledLog(true)
	//print logo
	printDotLogo()
	logger.Logger().Debug("DotWeb Start New AppServer", LogTarget_HttpServer)
	return app
}

// ClassicWithConf create and return DotApp instance
// must set config info
func ClassicWithConf(config *config.Config) *DotWeb{
	app := Classic(config.App.LogPath)
	app.SetConfig(config)
	return app
}

// RegisterMiddlewareFunc register middleware with gived name & middleware
func (app *DotWeb) RegisterMiddlewareFunc(name string, middleFunc MiddlewareFunc) {
	app.middlewareMutex.Lock()
	app.middlewareMap[name] = middleFunc
	app.middlewareMutex.Unlock()
}

// GetMiddlewareFunc get middleware with gived name
func (app *DotWeb) GetMiddlewareFunc(name string) (MiddlewareFunc, bool) {
	app.middlewareMutex.RLock()
	v, exists := app.middlewareMap[name]
	app.middlewareMutex.RUnlock()
	return v, exists
}

// Cache return cache interface
func (app *DotWeb) Cache() cache.Cache {
	return app.cache
}

// SetCache set cache interface
func (app *DotWeb) SetCache(ca cache.Cache) {
	app.cache = ca
}

// RunMode current app run mode, if not set, default set RunMode_Development
func (app *DotWeb) RunMode() string {
	if app.Config.App.RunMode != RunMode_Development && app.Config.App.RunMode != RunMode_Production {
		app.Config.App.RunMode = RunMode_Development
	}
	return app.Config.App.RunMode
}

// IsDevelopmentMode check current run mode is development mode
func (app *DotWeb) IsDevelopmentMode() bool {
	return app.RunMode() == RunMode_Development
}

// SetDevelopmentMode set run mode on development mode
// 1.SetEnabledLog(true)
// 2.SetEnabledConsole(true)
func (app *DotWeb) SetDevelopmentMode() {
	app.Config.App.RunMode = RunMode_Development
	app.SetEnabledLog(true)
	app.Use(new(RequestLogMiddleware))
	logger.SetEnabledConsole(true)
}

// SetProductionMode set run mode on production mode
func (app *DotWeb) SetProductionMode() {
	app.Config.App.RunMode = RunMode_Production
	logger.SetEnabledConsole(false)
}

// ExcludeUse registers a middleware exclude routers
// like exclude /index or /query/:id
func (app *DotWeb) ExcludeUse(m Middleware, routers ...string) {
	middlewareLen := len(app.Middlewares)
	if m != nil {
		m.Exclude(routers...)
		if middlewareLen > 0 {
			app.Middlewares[middlewareLen-1].SetNext(m)
		}
		app.Middlewares = append(app.Middlewares, m)
	}
}

// Use registers a middleware
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

// UseRequestLog register RequestLogMiddleware
func (app *DotWeb) UseRequestLog() {
	app.Use(&RequestLogMiddleware{})
}

// UseTimeoutHook register TimeoutHookMiddleware
func (app *DotWeb) UseTimeoutHook(handler StandardHandle, timeout time.Duration){
	app.Use(&TimeoutHookMiddleware{
		HookHandle: handler,
		TimeoutDuration:timeout,
	})
}

// SetExceptionHandle set custom error handler
func (app *DotWeb) SetExceptionHandle(handler ExceptionHandle) {
	app.ExceptionHandler = handler
}

// SetNotFoundHandle set custom 404 handler
func (app *DotWeb) SetNotFoundHandle(handler StandardHandle) {
	app.NotFoundHandler = handler
}

// SetMethodNotAllowedHandle set custom 405 handler
func (app *DotWeb) SetMethodNotAllowedHandle(handler StandardHandle) {
	app.MethodNotAllowedHandler = handler
}

// SetPProfConfig set pprofserver config, default is disable
// and don't use same port with StartServer
func (app *DotWeb) SetPProfConfig(enabledPProf bool, httpport int) {
	app.Config.App.EnabledPProf = enabledPProf
	app.Config.App.PProfPort = httpport
	logger.Logger().Debug("DotWeb SetPProfConfig ["+strconv.FormatBool(enabledPProf)+", "+strconv.Itoa(httpport)+"]", LogTarget_HttpServer)
}

// SetLogger set user logger, the logger must implement logger.AppLog interface
func (app *DotWeb) SetLogger(log logger.AppLog) {
	logger.SetLogger(log)
}

// SetLogPath set log root path
func (app *DotWeb) SetLogPath(path string) {
	logger.SetLogPath(path)
	//fixed #74 dotweb.SetEnabledLog 无效
	app.Config.App.LogPath = path
}

// SetEnabledLog set enabled log flag
func (app *DotWeb) SetEnabledLog(enabledLog bool) {
	logger.SetEnabledLog(enabledLog)
	//fixed #74 dotweb.SetEnabledLog 无效
	app.Config.App.EnabledLog = enabledLog
}

// SetConfig set config for app
func (app *DotWeb) SetConfig(config *config.Config) error {
	app.Config = config
	return nil
}

// StartServer start server with http port
// if config the pprof, will be start pprof server
func (app *DotWeb) StartServer(httpPort int) error {
	addr := ":" + strconv.Itoa(httpPort)
	return app.ListenAndServe(addr)
}

// Start start app server with set config
// If an exception occurs, will be return it
// if no set Server.Port, will be use DefaultHttpPort
func (app *DotWeb) Start() error {
	if app.Config == nil {
		return errors.New("no config exists")
	}
	//start server
	port := app.Config.Server.Port
	if port <= 0 {
		port = DefaultHTTPPort
	}
	return app.StartServer(port)
}

// MustStart start app server with set config
// If an exception occurs, will be panic it
// if no set Server.Port, will be use DefaultHttpPort
func (app *DotWeb) MustStart() {
	err := app.Start()
	if err != nil {
		panic(err)
	}
}

// ListenAndServe start server with addr
// not support pprof server auto start
func (app *DotWeb) ListenAndServe(addr string) error {
	app.initAppConfig()
	app.initRegisterConfigMiddleware()
	app.initRegisterConfigRoute()
	app.initRegisterConfigGroup()

	app.initServerEnvironment()

	app.initBindMiddleware()

	if app.StartMode == StartMode_Classic {
		app.IncludeDotwebGroup()
	}

	if app.HttpServer.ServerConfig().EnabledTLS {
		err := app.HttpServer.ListenAndServeTLS(addr, app.HttpServer.ServerConfig().TLSCertFile, app.HttpServer.ServerConfig().TLSKeyFile)
		return err
	}
	err := app.HttpServer.ListenAndServe(addr)
	return err

}

// init App Config
func (app *DotWeb) initAppConfig() {
	config := app.Config
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

	//设置启用详细请求数据统计
	if config.Server.EnabledDetailRequestData {
		core.GlobalState.EnabledDetailRequestData = config.Server.EnabledDetailRequestData
	}
}

// init register config's Middleware
func (app *DotWeb) initRegisterConfigMiddleware() {
	config := app.Config
	//register app's middleware
	for _, m := range config.Middlewares {
		if !m.IsUse {
			continue
		}
		if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
			app.Use(mf())
		}
	}
}

// init register config's route
func (app *DotWeb) initRegisterConfigRoute() {
	config := app.Config
	//load router and register
	for _, r := range config.Routers {
		//fmt.Println("config.Routers ", i, " ", config.Routers[i])
		if h, isok := app.HttpServer.Router().GetHandler(r.HandlerName); isok && r.IsUse {
			node := app.HttpServer.Router().RegisterRoute(strings.ToUpper(r.Method), r.Path, h)
			//use middleware
			for _, m := range r.Middlewares {
				if !m.IsUse {
					continue
				}
				if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
					node.Use(mf())
				}
			}
		}
	}
}

// init register config's route
func (app *DotWeb) initRegisterConfigGroup() {
	config := app.Config
	//support group
	for _, v := range config.Groups {
		if !v.IsUse {
			continue
		}
		g := app.HttpServer.Group(v.Path)
		//use middleware
		for _, m := range v.Middlewares {
			if !m.IsUse {
				continue
			}
			if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
				g.Use(mf())
			}
		}
		//init group's router
		for _, r := range v.Routers {
			if h, isok := app.HttpServer.Router().GetHandler(r.HandlerName); isok && r.IsUse {
				node := g.RegisterRoute(strings.ToUpper(r.Method), r.Path, h)
				//use middleware
				for _, m := range r.Middlewares {
					if !m.IsUse {
						continue
					}
					if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
						node.Use(mf())
					}
				}
			}
		}
	}
}

// init bind app's middleware to router node
func (app *DotWeb) initBindMiddleware() {
	router := app.HttpServer.Router().(*router)
	//bind app middlewares
	for fullExpress, _ := range router.allRouterExpress {
		expresses := strings.Split(fullExpress, routerExpressSplit)
		if len(expresses) < 2 {
			continue
		}
		node := router.getNode(expresses[0], expresses[1])
		if node == nil {
			continue
		}

		node.appMiddlewares = app.Middlewares
		for _, m := range node.appMiddlewares {
			if m.HasExclude() && m.ExistsExcludeRouter(node.fullPath) {
				logger.Logger().Debug("DotWeb initBindMiddleware [app] "+fullExpress+" "+reflect.TypeOf(m).String()+" exclude", LogTarget_HttpServer)
				node.hasExcludeMiddleware = true
			} else {
				logger.Logger().Debug("DotWeb initBindMiddleware [app] "+fullExpress+" "+reflect.TypeOf(m).String()+" match", LogTarget_HttpServer)
			}
		}
		if len(node.middlewares) > 0 {
			firstMiddleware := &xMiddleware{}
			firstMiddleware.SetNext(node.middlewares[0])
			node.middlewares = append([]Middleware{firstMiddleware}, node.middlewares...)
		}
	}

	//bind group middlewares
	for _, g := range app.HttpServer.groups {
		xg := g.(*xGroup)
		if len(xg.middlewares) <= 0 {
			continue
		} else {
			firstMiddleware := &xMiddleware{}
			firstMiddleware.SetNext(xg.middlewares[0])
			xg.middlewares = append([]Middleware{firstMiddleware}, xg.middlewares...)
		}
		for fullExpress, _ := range xg.allRouterExpress {
			expresses := strings.Split(fullExpress, "_")
			if len(expresses) < 2 {
				continue
			}
			node := router.getNode(expresses[0], expresses[1])
			if node == nil {
				continue
			}
			node.groupMiddlewares = xg.middlewares
			for _, m := range node.groupMiddlewares {
				if m.HasExclude() && m.ExistsExcludeRouter(node.fullPath) {
					logger.Logger().Debug("DotWeb initBindMiddleware [group] "+fullExpress+" "+reflect.TypeOf(m).String()+" exclude", LogTarget_HttpServer)
					node.hasExcludeMiddleware = true
				} else {
					logger.Logger().Debug("DotWeb initBindMiddleware [group] "+fullExpress+" "+reflect.TypeOf(m).String()+" match", LogTarget_HttpServer)
				}
			}
		}
	}
}

// IncludeDotwebGroup init inner routers
func (app *DotWeb) IncludeDotwebGroup() {
	//默认支持pprof信息查看
	gInner := app.HttpServer.Group("/dotweb")
	gInner.GET("/debug/pprof/:key", initPProf)
	gInner.GET("/debug/freemem", freeMemory)
	gInner.GET("/state", showServerState)
	gInner.GET("/state/interval", showIntervalData)
	gInner.GET("/query/:key", showQuery)
}

// init Server Environment
func (app *DotWeb) initServerEnvironment() {
	if app.ExceptionHandler == nil {
		app.SetExceptionHandle(app.DefaultHTTPErrorHandler)
	}

	if app.NotFoundHandler == nil {
		app.SetNotFoundHandle(app.DefaultNotFoundHandler)
	}

	if app.MethodNotAllowedHandler == nil {
		app.SetMethodNotAllowedHandle(app.DefaultMethodNotAllowedHandler)
	}

	//init session manager
	if app.HttpServer.SessionConfig().EnabledSession {
		if app.HttpServer.SessionConfig().SessionMode == "" {
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
	//if is develop mode, it will use nocache mode
	if app.HttpServer.Renderer() == nil {
		if app.RunMode() == RunMode_Development{
			app.HttpServer.SetRenderer(NewInnerRendererNoCache())
		}else{
			app.HttpServer.SetRenderer(NewInnerRenderer())
		}
	}

	//start pprof server
	if app.Config.App.EnabledPProf {
		logger.Logger().Debug("DotWeb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] Begin", LogTarget_HttpServer)
		go func() {
			err := http.ListenAndServe(":"+strconv.Itoa(app.Config.App.PProfPort), nil)
			if err != nil {
				logger.Logger().Error("DotWeb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] error: "+err.Error(), LogTarget_HttpServer)
				//panic the error
				panic(err)
			}
		}()
	}
}

// DefaultHTTPErrorHandler default exception handler
func (app *DotWeb) DefaultHTTPErrorHandler(ctx Context, err error) {
	//输出内容
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	//if in development mode, output the error info
	if app.IsDevelopmentMode() {
		stack := string(debug.Stack())
		ctx.WriteStringC(http.StatusInternalServerError, fmt.Sprintln(err)+stack)
	} else {
		ctx.WriteStringC(http.StatusInternalServerError, "Internal Server Error")
	}
}

// DefaultNotFoundHandler default exception handler
func (app *DotWeb) DefaultNotFoundHandler(ctx Context) {
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	ctx.WriteStringC(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// DefaultMethodNotAllowedHandler default exception handler
func (app *DotWeb) DefaultMethodNotAllowedHandler(ctx Context) {
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	ctx.WriteStringC(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}

// Close immediately stops the server.
// It internally calls `http.Server#Close()`.
func (app *DotWeb) Close() error {
	return app.HttpServer.stdServer.Close()
}

// Shutdown stops server the gracefully.
// It internally calls `http.Server#Shutdown()`.
func (app *DotWeb) Shutdown(ctx context.Context) error {
	return app.HttpServer.stdServer.Shutdown(ctx)
}

// HTTPNotFound simple notfound function for Context
func HTTPNotFound(ctx Context) {
	http.NotFound(ctx.Response().Writer(), ctx.Request().Request)
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

func showIntervalData(ctx Context) error {
	type data struct {
		Time         string
		RequestCount uint64
		ErrorCount   uint64
	}
	queryKey := ctx.QueryString("querykey")

	d := new(data)
	d.Time = queryKey
	d.RequestCount = core.GlobalState.QueryIntervalRequestData(queryKey)
	d.ErrorCount = core.GlobalState.QueryIntervalErrorData(queryKey)
	ctx.WriteJson(d)
	return nil
}

//显示服务器状态信息
func showServerState(ctx Context) error {
	ctx.WriteHtml(core.GlobalState.ShowHtmlData(Version))
	return nil
}

//显示服务器状态信息
func showQuery(ctx Context) error {
	querykey := ctx.GetRouterName("key")
	switch querykey {
	case "state":
		ctx.WriteString(jsonutil.GetJsonString(core.GlobalState))
	case "":
		ctx.WriteString("please input key")
	default:
		ctx.WriteString("not support key => " + querykey)
	}
	return nil
}

func printDotLogo() {
	logger.Logger().Print(`    ____           __                     __`, LogTarget_HttpServer)
	logger.Logger().Print(`   / __ \  ____   / /_ _      __  ___    / /_`, LogTarget_HttpServer)
	logger.Logger().Print(`  / / / / / __ \ / __/| | /| / / / _ \  / __ \`, LogTarget_HttpServer)
	logger.Logger().Print(` / /_/ / / /_/ // /_  | |/ |/ / /  __/ / /_/ /`, LogTarget_HttpServer)
	logger.Logger().Print(`/_____/  \____/ \__/  |__/|__/  \___/ /_.___/`, LogTarget_HttpServer)
}
