package dotweb

import (
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/session"
	"net/http"
	"strings"
	"sync"

	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/feature"
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/logger"
	"strconv"
)

const (
	DefaultGzipLevel = 9
	gzipScheme       = "gzip"
	DefaultIndexPage = "index.html"
)

type (
	//HttpServer定义
	HttpServer struct {
		stdServer      *http.Server
		router         Router
		groups	 	   []Group
		Modules        []*HttpModule
		DotApp         *DotWeb
		Validator      Validator
		sessionManager *session.SessionManager
		lock_session   *sync.RWMutex
		pool           *pool
		binder         Binder
		render         Renderer
		offline        bool
		Features       *feature.Feature
	}

	//pool定义
	pool struct {
		request  sync.Pool
		response sync.Pool
		context  sync.Pool
	}
)

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		pool: &pool{
			response: sync.Pool{
				New: func() interface{} {
					return &Response{}
				},
			},
			request: sync.Pool{
				New: func() interface{} {
					return &Request{}
				},
			},
			context: sync.Pool{
				New: func() interface{} {
					return &HttpContext{}
				},
			},
		},
		Modules:      make([]*HttpModule, 0),
		lock_session: new(sync.RWMutex),
		binder:       newBinder(),
		Features:     &feature.Feature{},
	}
	//设置router
	server.router = NewRouter(server)
	server.stdServer = &http.Server{Handler: server}
	return server
}

// ServerConfig a shortcut for App.Config.ServerConfig
func (server *HttpServer) ServerConfig() *config.ServerNode {
	return server.DotApp.Config.Server
}

// SessionConfig a shortcut for App.Config.SessionConfig
func (server *HttpServer) SessionConfig() *config.SessionNode {
	return server.DotApp.Config.Session
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
func (server *HttpServer) ListenAndServe(addr string) error {
	server.stdServer.Addr = addr
	logger.Logger().Debug("DotWeb:HttpServer ListenAndServe ["+addr+"]", LogTarget_HttpServer)
	return server.stdServer.ListenAndServe()
}

// ListenAndServeTLS listens on the TCP network address srv.Addr and
// then calls Serve to handle requests on incoming TLS connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// Filenames containing a certificate and matching private key for the
// server must be provided if neither the Server's TLSConfig.Certificates
// nor TLSConfig.GetCertificate are populated. If the certificate is
// signed by a certificate authority, the certFile should be the
// concatenation of the server's certificate, any intermediates, and
// the CA's certificate.
//
// If srv.Addr is blank, ":https" is used.
//
// ListenAndServeTLS always returns a non-nil error.
func (server *HttpServer) ListenAndServeTLS(addr string, certFile, keyFile string) error {
	server.stdServer.Addr = addr
	//check tls config
	if !file.Exist(certFile) {
		logger.Logger().Error("DotWeb:HttpServer ListenAndServeTLS ["+addr+","+certFile+","+keyFile+"] error => Server EnabledTLS is true, but TLSCertFile not exists", LogTarget_HttpServer)
		panic("Server EnabledTLS is true, but TLSCertFile not exists")
	}
	if !file.Exist(keyFile) {
		logger.Logger().Error("DotWeb:HttpServer ListenAndServeTLS ["+addr+","+certFile+","+keyFile+"] error => Server EnabledTLS is true, but TLSKeyFile not exists", LogTarget_HttpServer)
		panic("Server EnabledTLS is true, but TLSKeyFile not exists")
	}
	logger.Logger().Debug("DotWeb:HttpServer ListenAndServeTLS ["+addr+","+certFile+","+keyFile+"]", LogTarget_HttpServer)
	return server.stdServer.ListenAndServeTLS(certFile, keyFile)
}

// ServeHTTP make sure request can be handled correctly
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	core.GlobalState.AddCurrentRequest(1)
	defer core.GlobalState.SubCurrentRequest(1)

	//针对websocket与调试信息特殊处理
	if checkIsWebSocketRequest(req) {
		http.DefaultServeMux.ServeHTTP(w, req)
		//增加状态计数
		core.GlobalState.AddRequestCount(req.URL.Path, defaultHttpCode, 1)
	} else {
		//设置header信息
		w.Header().Set(HeaderServer, DefaultServerName)
		//处理维护
		if server.IsOffline() {
			server.DotApp.OfflineServer.ServeHTTP(w, req)
		} else {
			//get from pool
			response := server.pool.response.Get().(*Response)
			request := server.pool.request.Get().(*Request)
			httpCtx := server.pool.context.Get().(*HttpContext)
			httpCtx.reset(response, request, server, nil, nil, nil)
			response.reset(w)
			request.reset(req, httpCtx)

			//处理前置Module集合
			for _, module := range server.Modules {
				if module.OnBeginRequest != nil {
					module.OnBeginRequest(httpCtx)
				}
			}

			if !httpCtx.IsEnd() {
				server.Router().ServeHTTP(httpCtx)
			}

			//处理后置Module集合
			for _, module := range server.Modules {
				if module.OnEndRequest != nil {
					module.OnEndRequest(httpCtx)
				}
			}

			//增加状态计数
			core.GlobalState.AddRequestCount(httpCtx.Request().Path(), httpCtx.Response().HttpCode(), 1)

			//release response
			response.release()
			server.pool.response.Put(response)
			//release request
			request.release()
			server.pool.request.Put(request)
			//release context
			httpCtx.release()
			server.pool.context.Put(httpCtx)
		}
	}
}

// IsOffline check server is set offline state
func (server *HttpServer) IsOffline() bool {
	return server.offline
}

// SetOffline set server offline config
func (server *HttpServer) SetOffline(offline bool, offlineText string, offlineUrl string) {
	server.offline = offline
}

// IndexPage default index page name
func (server *HttpServer) IndexPage() string {
	if server.ServerConfig().IndexPage == "" {
		return DefaultIndexPage
	} else {
		return server.ServerConfig().IndexPage
	}
}

// SetSessionConfig set session store config
func (server *HttpServer) SetSessionConfig(storeConfig *session.StoreConfig) {
	//sync session config
	server.SessionConfig().Timeout = storeConfig.Maxlifetime
	server.SessionConfig().SessionMode = storeConfig.StoreName
	server.SessionConfig().ServerIP = storeConfig.ServerIP
	server.SessionConfig().StoreKeyPre = storeConfig.StoreKeyPre
	server.SessionConfig().CookieName = storeConfig.CookieName
	logger.Logger().Debug("DotWeb:HttpServer SetSessionConfig ["+jsonutil.GetJsonString(storeConfig)+"]", LogTarget_HttpServer)
}

// InitSessionManager init session manager
func (server *HttpServer) InitSessionManager() {
	storeConfig := new(session.StoreConfig)
	storeConfig.Maxlifetime = server.SessionConfig().Timeout
	storeConfig.StoreName = server.SessionConfig().SessionMode
	storeConfig.ServerIP = server.SessionConfig().ServerIP
	storeConfig.StoreKeyPre = server.SessionConfig().StoreKeyPre
	storeConfig.CookieName = server.SessionConfig().CookieName

	if server.sessionManager == nil {
		//设置Session
		server.lock_session.Lock()
		if manager, err := session.NewDefaultSessionManager(storeConfig); err != nil {
			//panic error with create session manager
			panic(err.Error())
		} else {
			server.sessionManager = manager
		}
		server.lock_session.Unlock()
	}
	logger.Logger().Debug("DotWeb:HttpServer InitSessionManager ["+jsonutil.GetJsonString(storeConfig)+"]", LogTarget_HttpServer)
}

// setDotApp 关联当前HttpServer实例对应的DotServer实例
func (server *HttpServer) setDotApp(dotApp *DotWeb) {
	server.DotApp = dotApp
}

// GetSessionManager get session manager in current httpserver
func (server *HttpServer) GetSessionManager() *session.SessionManager {
	if !server.SessionConfig().EnabledSession {
		return nil
	}
	return server.sessionManager
}

// Router get router interface in server
func (server *HttpServer) Router() Router {
	return server.router
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (server *HttpServer) GET(path string, handle HttpHandle) RouterNode {
	return server.Router().GET(path, handle)
}

// ANY is a shortcut for router.Handle("Any", path, handle)
// it support GET\HEAD\POST\PUT\PATCH\OPTIONS\DELETE
func (server *HttpServer) Any(path string, handle HttpHandle) {
	server.Router().Any(path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (server *HttpServer) HEAD(path string, handle HttpHandle) RouterNode {
	return server.Router().HEAD(path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (server *HttpServer) OPTIONS(path string, handle HttpHandle) RouterNode {
	return server.Router().OPTIONS(path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (server *HttpServer) POST(path string, handle HttpHandle) RouterNode {
	return server.Router().POST(path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (server *HttpServer) PUT(path string, handle HttpHandle) RouterNode {
	return server.Router().PUT(path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (server *HttpServer) PATCH(path string, handle HttpHandle) RouterNode {
	return server.Router().PATCH(path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (server *HttpServer) DELETE(path string, handle HttpHandle) RouterNode {
	return server.Router().DELETE(path, handle)
}

// ServerFile is a shortcut for router.ServeFiles(path, filepath)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (server *HttpServer) ServerFile(path string, fileroot string) RouterNode {
	return server.Router().ServerFile(path, fileroot)
}

// HiJack is a shortcut for router.HiJack(path, handle)
func (server *HttpServer) HiJack(path string, handle HttpHandle) {
	server.Router().HiJack(path, handle)
}

// WebSocket is a shortcut for router.WebSocket(path, handle)
func (server *HttpServer) WebSocket(path string, handle HttpHandle) {
	server.Router().WebSocket(path, handle)
}

// Group create new group with current HttpServer
func (server *HttpServer) Group(prefix string) Group {
	return NewGroup(prefix, server)
}

// Binder get binder interface in server
func (server *HttpServer) Binder() Binder {
	return server.binder
}

// Renderer get renderer interface in server
// if no set, init InnerRenderer
func (server *HttpServer) Renderer() Renderer {
	if server.render == nil {
		if server.DotApp.RunMode() == RunMode_Development{
			server.SetRenderer(NewInnerRendererNoCache())
		}else{
			server.SetRenderer(NewInnerRenderer())
		}
	}
	return server.render
}

// SetRenderer set custom renderer in server
func (server *HttpServer) SetRenderer(r Renderer) {
	server.render = r
}

// SetEnabledAutoHEAD set route use auto head
// set EnabledAutoHEAD true or false
// default is false
func (server *HttpServer) SetEnabledAutoHEAD(isEnabled bool) {
	server.ServerConfig().EnabledAutoHEAD = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledAutoHEAD ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledRequestID set create unique request id per request
// set EnabledRequestID true or false
// default is false
func (server *HttpServer) SetEnabledRequestID(isEnabled bool) {
	server.ServerConfig().EnabledRequestID = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledRequestID ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledListDir 设置是否允许目录浏览,默认为false
func (server *HttpServer) SetEnabledListDir(isEnabled bool) {
	server.ServerConfig().EnabledListDir = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledListDir ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledSession 设置是否启用Session,默认为false
func (server *HttpServer) SetEnabledSession(isEnabled bool) {
	server.SessionConfig().EnabledSession = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledSession ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledGzip 设置是否启用gzip,默认为false
func (server *HttpServer) SetEnabledGzip(isEnabled bool) {
	server.ServerConfig().EnabledGzip = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledGzip ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledBindUseJsonTag 设置bind是否启用json标签,默认为false, fixed for issue #91
func (server *HttpServer) SetEnabledBindUseJsonTag(isEnabled bool) {
	server.ServerConfig().EnabledBindUseJsonTag = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledBindUseJsonTag ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}


// SetEnabledIgnoreFavicon set IgnoreFavicon Enabled
// default is false
func (server *HttpServer) SetEnabledIgnoreFavicon(isEnabled bool) {
	server.ServerConfig().EnabledIgnoreFavicon = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledIgnoreFavicon ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
	server.RegisterModule(getIgnoreFaviconModule())
}

// SetEnabledTLS set tls enabled
// default is false
// if it's true, must input certificate\private key fileName
func (server *HttpServer) SetEnabledTLS(isEnabled bool, certFile, keyFile string) {
	server.ServerConfig().EnabledTLS = isEnabled
	server.ServerConfig().TLSCertFile = certFile
	server.ServerConfig().TLSKeyFile = keyFile
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledTLS ["+strconv.FormatBool(isEnabled)+","+certFile+","+keyFile+"]", LogTarget_HttpServer)
}

// SetEnabledDetailRequestData 设置是否启用详细请求数据统计,默认为false
func (server *HttpServer) SetEnabledDetailRequestData(isEnabled bool) {
	server.ServerConfig().EnabledDetailRequestData = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledDetailRequestData ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// RegisterModule 添加处理模块
func (server *HttpServer) RegisterModule(module *HttpModule) {
	server.Modules = append(server.Modules, module)
	logger.Logger().Debug("DotWeb:HttpServer RegisterModule ["+module.Name+"]", LogTarget_HttpServer)
}

type LogJson struct {
	RequestUrl string
	HttpHeader string
	HttpBody   string
}

//check request is the websocket request
//check Connection contains upgrade
func checkIsWebSocketRequest(req *http.Request) bool {
	if strings.Index(strings.ToLower(req.Header.Get("Connection")), "upgrade") >= 0 {
		return true
	}
	return false
}

//check request is startwith /debug/
func checkIsDebugRequest(req *http.Request) bool {
	if strings.Index(req.RequestURI, "/debug/") == 0 {
		return true
	}
	return false
}
