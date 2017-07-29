package dotweb

import (
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/session"
	"net/http"
	"strings"
	"sync"

	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/feature"
)

const (
	DefaultGzipLevel = 9
	gzipScheme       = "gzip"
	DefaultIndexPage = "index.html"
)

type (
	// Deprecated: Use the Middleware instead
	// HttpModule struct
	HttpModule struct {
		//响应请求时作为 HTTP 执行管线链中的第一个事件发生
		OnBeginRequest func(Context)
		//响应请求时作为 HTTP 执行管线链中的最后一个事件发生。
		OnEndRequest func(Context)
	}

	//HttpServer定义
	HttpServer struct {
		stdServer      *http.Server
		router         Router
		DotApp         *DotWeb
		sessionManager *session.SessionManager
		lock_session   *sync.RWMutex
		pool           *pool
		ServerConfig   *config.ServerNode
		SessionConfig  *config.SessionNode
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
		ServerConfig:  config.NewServerNode(),
		SessionConfig: config.NewSessionNode(),
		lock_session:  new(sync.RWMutex),
		binder:        newBinder(),
		Features:      &feature.Feature{},
	}
	//设置router
	server.router = NewRouter(server)
	server.stdServer = &http.Server{Handler: server}
	return server
}

// ListenAndServe listens on the TCP network address srv.Addr and then
// calls Serve to handle requests on incoming connections.
func (server *HttpServer) ListenAndServe(addr string) error {
	server.stdServer.Addr = addr
	return server.stdServer.ListenAndServe()
}

// ServeHTTP make sure request can be handled correctly
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//针对websocket与调试信息特殊处理
	if checkIsWebSocketRequest(req) {
		http.DefaultServeMux.ServeHTTP(w, req)
	} else {
		//设置header信息
		w.Header().Set(HeaderServer, DefaultServerName)
		//处理维护
		if server.IsOffline() {
			server.DotApp.OfflineServer.ServeHTTP(w, req)
		} else {
			//get from pool
			response := server.pool.response.Get().(*Response)
			response.reset(w)
			request := server.pool.request.Get().(*Request)
			request.reset(req)
			httpCtx := server.pool.context.Get().(*HttpContext)
			httpCtx.reset(response, request, server, nil, nil, nil)

			//增加状态计数
			core.GlobalState.AddRequestCount(1)

			server.Router().ServeHTTP(httpCtx)

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
	if server.ServerConfig.IndexPage == "" {
		return DefaultIndexPage
	} else {
		return server.ServerConfig.IndexPage
	}
}

// SetSessionConfig set session store config
func (server *HttpServer) SetSessionConfig(storeConfig *session.StoreConfig) {
	//sync session config
	server.SessionConfig.Timeout = storeConfig.Maxlifetime
	server.SessionConfig.SessionMode = storeConfig.StoreName
	server.SessionConfig.ServerIP = storeConfig.ServerIP
}

// InitSessionManager init session manager
func (server *HttpServer) InitSessionManager() {
	storeConfig := new(session.StoreConfig)
	storeConfig.Maxlifetime = server.SessionConfig.Timeout
	storeConfig.StoreName = server.SessionConfig.SessionMode
	storeConfig.ServerIP = server.SessionConfig.ServerIP

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
}

// setDotApp 关联当前HttpServer实例对应的DotServer实例
func (server *HttpServer) setDotApp(dotApp *DotWeb) {
	server.DotApp = dotApp
}

// GetSessionManager get session manager in current httpserver
func (server *HttpServer) GetSessionManager() *session.SessionManager {
	if !server.SessionConfig.EnabledSession {
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

func (server *HttpServer) HiJack(path string, handle HttpHandle) {
	server.Router().HiJack(path, handle)
}

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
		server.render = NewInnerRenderer()
	}
	return server.render
}

// SetRenderer set custom renderer in server
func (server *HttpServer) SetRenderer(r Renderer) {
	server.render = r
}

// SetEnabledAutoHEAD set EnabledAutoHEAD true or false
func (server *HttpServer) SetEnabledAutoHEAD(autoHEAD bool) {
	server.ServerConfig.EnabledAutoHEAD = autoHEAD
}

// SetEnabledListDir 设置是否允许目录浏览,默认为false
func (server *HttpServer) SetEnabledListDir(isEnabled bool) {
	server.ServerConfig.EnabledListDir = isEnabled
}

// SetEnabledSession 设置是否启用Session,默认为false
func (server *HttpServer) SetEnabledSession(isEnabled bool) {
	server.SessionConfig.EnabledSession = isEnabled
}

// SetEnabledGzip 设置是否启用gzip,默认为false
func (server *HttpServer) SetEnabledGzip(isEnabled bool) {
	server.ServerConfig.EnabledGzip = isEnabled
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
