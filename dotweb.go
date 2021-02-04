package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/framework/crypto/uuid"
	"github.com/devfeel/dotweb/framework/exception"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"
	"strconv"
	"strings"

	"context"
	"errors"
	"reflect"
	"sync"
	"time"

	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/logger"
	"github.com/devfeel/dotweb/session"
)

var (
	// ErrValidatorNotRegistered error for not register Validator
	ErrValidatorNotRegistered = errors.New("validator not registered")

	// ErrNotFound error for not found file
	ErrNotFound = errors.New("not found file")
)

type (
	DotWeb struct {
		HttpServer              *HttpServer
		cache                   cache.Cache
		Config                  *config.Config
		Mock                    Mock
		Middlewares             []Middleware
		ExceptionHandler        ExceptionHandle
		NotFoundHandler         StandardHandle // NotFoundHandler supports user defined 404 handler
		MethodNotAllowedHandler StandardHandle // MethodNotAllowedHandler fixed for #64 supports user defined MethodNotAllowed handler
		Items                   core.ConcurrenceMap
		middlewareMap           map[string]MiddlewareFunc
		middlewareMutex         *sync.RWMutex
		pluginMap               map[string]Plugin
		pluginMutex             *sync.RWMutex
		StartMode               string
		IDGenerater             IdGenerate
		globalUniqueID          string
		appLog                  logger.AppLog
		serverStateInfo         *core.ServerStateInfo
		isRun                   bool
	}

	// ExceptionHandle supports exception handling
	ExceptionHandle func(Context, error)

	// StandardHandle for standard request handling
	StandardHandle func(Context)

	// HttpHandle is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but has a special parameter Context contain all request and response data.
	HttpHandle func(Context) error

	// IdGenerater the handler for create Unique Id
	// default is use dotweb.
	IdGenerate func() string

	// Validator is the interface that wraps the Validate function.
	Validator interface {
		Validate(i interface{}) error
	}
)

const (
	// DefaultHTTPPort default http port; fixed for #70 UPDATE default http port 80 to 8080
	DefaultHTTPPort = 8080

	DefaultLogPath = ""

	// RunMode_Development app runmode in development mode
	RunMode_Development = "development"
	// RunMode_Production app runmode in production mode
	RunMode_Production = "production"

	// StartMode_New app startmode in New mode
	StartMode_New = "New"
	// StartMode_Classic app startmode in Classic mode
	StartMode_Classic = "Classic"
)

// New create and return DotApp instance
// default run mode is RunMode_Production
func New() *DotWeb {
	app := &DotWeb{
		HttpServer:      NewHttpServer(),
		Middlewares:     make([]Middleware, 0),
		Items:           core.NewConcurrenceMap(),
		Config:          config.NewConfig(),
		middlewareMap:   make(map[string]MiddlewareFunc),
		middlewareMutex: new(sync.RWMutex),
		pluginMap:       make(map[string]Plugin),
		pluginMutex:     new(sync.RWMutex),
		StartMode:       StartMode_New,
		serverStateInfo: core.NewServerStateInfo(),
	}
	// set default run mode = RunMode_Production
	app.Config.App.RunMode = RunMode_Production
	app.HttpServer.setDotApp(app)
	// add default httphandler with middlewares
	// fixed for issue #100
	app.Use(&xMiddleware{})

	// init logger
	app.appLog = logger.NewAppLog()

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

	if logPath != "" {
		app.SetLogPath(logPath)
	}
	app.SetEnabledLog(true)

	// print logo
	app.printDotLogo()

	app.Logger().Debug("DotWeb Start New AppServer", LogTarget_HttpServer)
	return app
}

// ClassicWithConf create and return DotApp instance
// must set config info
func ClassicWithConf(config *config.Config) *DotWeb {
	app := Classic(config.App.LogPath)
	app.SetConfig(config)
	return app
}

// Logger return app's logger
func (app *DotWeb) Logger() logger.AppLog {
	return app.appLog
}

// StateInfo return app's ServerStateInfo
func (app *DotWeb) StateInfo() *core.ServerStateInfo {
	return app.serverStateInfo
}

// RegisterMiddlewareFunc register middleware with given name & middleware
func (app *DotWeb) RegisterMiddlewareFunc(name string, middleFunc MiddlewareFunc) {
	app.middlewareMutex.Lock()
	app.middlewareMap[name] = middleFunc
	app.middlewareMutex.Unlock()
}

// GetMiddlewareFunc get middleware with given name
func (app *DotWeb) GetMiddlewareFunc(name string) (MiddlewareFunc, bool) {
	app.middlewareMutex.RLock()
	v, exists := app.middlewareMap[name]
	app.middlewareMutex.RUnlock()
	return v, exists
}

// GlobalUniqueID return app's GlobalUniqueID
// it will be Initializationed when StartServer
func (app *DotWeb) GlobalUniqueID() string {
	return app.globalUniqueID
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

	// enabled auto OPTIONS
	app.HttpServer.SetEnabledAutoOPTIONS(true)
	// enabled auto HEAD
	app.HttpServer.SetEnabledAutoHEAD(true)

	app.SetEnabledLog(true)
	app.Use(new(RequestLogMiddleware))
	app.Logger().SetEnabledConsole(true)
}

// SetProductionMode set run mode on production mode
func (app *DotWeb) SetProductionMode() {
	app.Config.App.RunMode = RunMode_Production
	app.appLog.SetEnabledConsole(true)
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

// UsePlugin registers plugins
func (app *DotWeb) UsePlugin(plugins ...Plugin) {
	app.pluginMutex.Lock()
	defer app.pluginMutex.Unlock()
	for _, p := range plugins {
		app.pluginMap[p.Name()] = p
	}
}

// Use registers middlewares
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
func (app *DotWeb) UseTimeoutHook(handler StandardHandle, timeout time.Duration) {
	app.Use(&TimeoutHookMiddleware{
		HookHandle:      handler,
		TimeoutDuration: timeout,
	})
}

// SetMock set mock logic
func (app *DotWeb) SetMock(mock Mock) {
	app.Mock = mock
	app.Logger().Debug("DotWeb Mock SetMock", LogTarget_HttpServer)
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
	app.Logger().Debug("DotWeb SetPProfConfig ["+strconv.FormatBool(enabledPProf)+", "+strconv.Itoa(httpport)+"]", LogTarget_HttpServer)
}

// SetLogger set user logger, the logger must implement logger.AppLog interface
func (app *DotWeb) SetLogger(log logger.AppLog) {
	app.appLog = log
}

// SetLogPath set log root path
func (app *DotWeb) SetLogPath(path string) {
	app.Logger().SetLogPath(path)
	// fixed #74 dotweb.SetEnabledLog 无效
	app.Config.App.LogPath = path
}

// SetEnabledLog set enabled log flag
func (app *DotWeb) SetEnabledLog(enabledLog bool) {
	app.Logger().SetEnabledLog(enabledLog)
	// fixed #74 dotweb.SetEnabledLog 无效
	app.Config.App.EnabledLog = enabledLog
}

// SetConfig set config for app
func (app *DotWeb) SetConfig(config *config.Config) {
	app.Config = config
}

// ReSetConfig reset config for app
// only apply when app is running
// Port can not be modify
// if EnabledPProf, EnabledPProf flag and PProfPort can not be modify
func (app *DotWeb) ReSetConfig(config *config.Config) {
	if !app.isRun {
		app.Logger().Debug("DotWeb is not running, ReSetConfig can not be call", LogTarget_HttpServer)
		return
	}

	config.Server.Port = app.Config.Server.Port
	if app.Config.App.EnabledPProf {
		config.App.PProfPort = app.Config.App.PProfPort
		config.App.EnabledPProf = app.Config.App.EnabledPProf
	}
	app.Config = config
	app.appLog = logger.NewAppLog()
	app.initAppConfig()
	app.Logger().Debug("DotWeb ReSetConfig is done.", LogTarget_HttpServer)
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
	// start server
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

	// create unique id for dotweb app
	app.globalUniqueID = app.IDGenerater()

	if app.StartMode == StartMode_Classic {
		app.IncludeDotwebGroup()
	}

	// special, if run mode is not develop, auto stop mock
	if app.RunMode() != RunMode_Development {
		if app.Mock != nil {
			app.Logger().Debug("DotWeb Mock RunMode is not DevelopMode, Auto stop mock", LogTarget_HttpServer)
		}
		app.Mock = nil
	}
	// output run mode
	app.Logger().Debug("DotWeb RunMode is "+app.RunMode(), LogTarget_HttpServer)

	// start plugins
	app.initPlugins()

	if app.HttpServer.ServerConfig().EnabledTLS {
		err := app.HttpServer.ListenAndServeTLS(addr, app.HttpServer.ServerConfig().TLSCertFile, app.HttpServer.ServerConfig().TLSKeyFile)
		return err
	}
	app.isRun = true
	err := app.HttpServer.ListenAndServe(addr)
	return err

}

// init App Config
func (app *DotWeb) initAppConfig() {
	config := app.Config
	// log config
	if config.App.LogPath != "" {
		app.SetLogPath(config.App.LogPath)
	}
	app.SetEnabledLog(config.App.EnabledLog)

	// run mode config
	if app.Config.App.RunMode != RunMode_Development && app.Config.App.RunMode != RunMode_Production {
		app.Config.App.RunMode = RunMode_Development
	}

	app.HttpServer.initConfig()

	// detailed request metrics
	if config.Server.EnabledDetailRequestData {
		app.StateInfo().EnabledDetailRequestData = config.Server.EnabledDetailRequestData
	}
}

// init register config's Middleware
func (app *DotWeb) initRegisterConfigMiddleware() {
	config := app.Config
	// register app's middleware
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
	// load router and register
	for _, r := range config.Routers {
		// fmt.Println("config.Routers ", i, " ", config.Routers[i])
		if h, isok := app.HttpServer.Router().GetHandler(r.HandlerName); isok && r.IsUse {
			node := app.HttpServer.Router().RegisterRoute(strings.ToUpper(r.Method), r.Path, h)
			// use middleware
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
	// support group
	for _, v := range config.Groups {
		if !v.IsUse {
			continue
		}
		g := app.HttpServer.Group(v.Path)
		// use middleware
		for _, m := range v.Middlewares {
			if !m.IsUse {
				continue
			}
			if mf, isok := app.GetMiddlewareFunc(m.Name); isok {
				g.Use(mf())
			}
		}
		// init group's router
		for _, r := range v.Routers {
			if h, isok := app.HttpServer.Router().GetHandler(r.HandlerName); isok && r.IsUse {
				node := g.RegisterRoute(strings.ToUpper(r.Method), r.Path, h)
				// use middleware
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

// initPlugins init and run plugins
func (app *DotWeb) initPlugins() {
	for _, p := range app.pluginMap {
		if p.IsValidate() {
			go func(p Plugin) {
				defer func() {
					if err := recover(); err != nil {
						app.Logger().Error(exception.CatchError("DotWeb::initPlugins run error plugin - "+p.Name(), "", err), LogTarget_HttpServer)
					}
				}()
				p.Run()
			}(p)
			app.Logger().Debug("DotWeb initPlugins start run plugin - "+p.Name(), LogTarget_HttpServer)
		} else {
			app.Logger().Debug("DotWeb initPlugins not validate plugin - "+p.Name(), LogTarget_HttpServer)
		}
	}
}

// init bind app's middleware to router node
func (app *DotWeb) initBindMiddleware() {
	router := app.HttpServer.Router().(*router)
	// bind app middlewares
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
				app.Logger().Debug("DotWeb initBindMiddleware [app] "+fullExpress+" "+reflect.TypeOf(m).String()+" exclude", LogTarget_HttpServer)
				node.hasExcludeMiddleware = true
			} else {
				app.Logger().Debug("DotWeb initBindMiddleware [app] "+fullExpress+" "+reflect.TypeOf(m).String()+" match", LogTarget_HttpServer)
			}
		}
		if len(node.middlewares) > 0 {
			firstMiddleware := &xMiddleware{}
			firstMiddleware.SetNext(node.middlewares[0])
			node.middlewares = append([]Middleware{firstMiddleware}, node.middlewares...)
		}
	}

	// bind group middlewares
	for _, g := range app.HttpServer.groups {
		xg := g.(*xGroup)
		if len(xg.middlewares) <= 0 {
			continue
		}
		for fullExpress, _ := range xg.allRouterExpress {
			expresses := strings.Split(fullExpress, routerExpressSplit)
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
					app.Logger().Debug("DotWeb initBindMiddleware [group] "+fullExpress+" "+reflect.TypeOf(m).String()+" exclude", LogTarget_HttpServer)
					node.hasExcludeMiddleware = true
				} else {
					app.Logger().Debug("DotWeb initBindMiddleware [group] "+fullExpress+" "+reflect.TypeOf(m).String()+" match", LogTarget_HttpServer)
				}
			}
		}
	}
}

// IncludeDotwebGroup init inner routers which start with /dotweb/
func (app *DotWeb) IncludeDotwebGroup() {
	initDotwebGroup(app.HttpServer)
}

// init Server Environment
func (app *DotWeb) initServerEnvironment() {
	if app.ExceptionHandler == nil {
		app.SetExceptionHandle(app.DefaultHTTPErrorHandler)
	}

	if app.NotFoundHandler == nil {
		app.SetNotFoundHandle(DefaultNotFoundHandler)
	}

	if app.MethodNotAllowedHandler == nil {
		app.SetMethodNotAllowedHandle(DefaultMethodNotAllowedHandler)
	}

	// set default unique id generater
	if app.IDGenerater == nil {
		app.IDGenerater = DefaultUniqueIDGenerater
	}

	// init session manager
	if app.HttpServer.SessionConfig().EnabledSession {
		if app.HttpServer.SessionConfig().SessionMode == "" {
			// panic("no set SessionConfig, but set enabledsession true")
			app.Logger().Warn("not set SessionMode, but set enabledsession true, now will use default runtime session", LogTarget_HttpServer)
			app.HttpServer.SetSessionConfig(session.NewDefaultRuntimeConfig())
		}
		app.HttpServer.InitSessionManager()
	}

	// if cache not set, create default runtime cache
	if app.Cache() == nil {
		app.cache = cache.NewRuntimeCache()
	}

	// if renderer not set, create inner renderer
	// if is develop mode, it will use nocache mode
	if app.HttpServer.Renderer() == nil {
		if app.RunMode() == RunMode_Development {
			app.HttpServer.SetRenderer(NewInnerRendererNoCache())
		} else {
			app.HttpServer.SetRenderer(NewInnerRenderer())
		}
	}

	// start pprof server
	if app.Config.App.EnabledPProf {
		app.Logger().Debug("DotWeb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] Begin", LogTarget_HttpServer)
		go func() {
			err := http.ListenAndServe(":"+strconv.Itoa(app.Config.App.PProfPort), nil)
			if err != nil {
				app.Logger().Error("DotWeb:StartPProfServer["+strconv.Itoa(app.Config.App.PProfPort)+"] error: "+err.Error(), LogTarget_HttpServer)
				// panic the error
				panic(err)
			}
		}()
	}
}

// DefaultHTTPErrorHandler default exception handler
func (app *DotWeb) DefaultHTTPErrorHandler(ctx Context, err error) {
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	// if in development mode, output the error info
	if app.IsDevelopmentMode() {
		stack := string(debug.Stack())
		ctx.WriteStringC(http.StatusInternalServerError, fmt.Sprintln(err)+stack)
	} else {
		ctx.WriteStringC(http.StatusInternalServerError, "Internal Server Error")
	}
}

func (app *DotWeb) printDotLogo() {
	app.Logger().Print(`    ____           __                     __`, LogTarget_HttpServer)
	app.Logger().Print(`   / __ \  ____   / /_ _      __  ___    / /_`, LogTarget_HttpServer)
	app.Logger().Print(`  / / / / / __ \ / __/| | /| / / / _ \  / __ \`, LogTarget_HttpServer)
	app.Logger().Print(` / /_/ / / /_/ // /_  | |/ |/ / /  __/ / /_/ /`, LogTarget_HttpServer)
	app.Logger().Print(`/_____/  \____/ \__/  |__/|__/  \___/ /_.___/`, LogTarget_HttpServer)
	app.Logger().Print(`                             Version `+Version, LogTarget_HttpServer)
}

// Close immediately stops the server.
// It internally calls `http.Server#Close()`.
func (app *DotWeb) Close() error {
	return app.HttpServer.stdServer.Close()
}

// Shutdown stops server gracefully.
// It internally calls `http.Server#Shutdown()`.
func (app *DotWeb) Shutdown(ctx context.Context) error {
	return app.HttpServer.stdServer.Shutdown(ctx)
}

// HTTPNotFound simple notfound function for Context
func HTTPNotFound(ctx Context) {
	http.NotFound(ctx.Response().Writer(), ctx.Request().Request)
}

// DefaultNotFoundHandler default exception handler
func DefaultNotFoundHandler(ctx Context) {
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	ctx.WriteStringC(http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

// DefaultMethodNotAllowedHandler default exception handler
func DefaultMethodNotAllowedHandler(ctx Context) {
	ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)
	ctx.WriteStringC(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))
}

// DefaultAutoOPTIONSHandler default handler for options request
// if set HttpServer.EnabledAutoOPTIONS, auto bind this handler
func DefaultAutoOPTIONSHandler(ctx Context) error {
	return ctx.WriteStringC(http.StatusNoContent, "")
}

// DefaultUniqueIDGenerater default generater used to create Unique Id
func DefaultUniqueIDGenerater() string {
	return uuid.NewV1().String32()
}

func DefaultTimeoutHookHandler(ctx Context) {
	realDration := ctx.Items().GetTimeDuration(ItemKeyHandleDuration)
	logs := fmt.Sprintf("req %v, cost %v", ctx.Request().Url(), realDration.Seconds())
	ctx.HttpServer().DotApp.Logger().Warn(logs, LogTarget_RequestTimeout)
}
