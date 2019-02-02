package dotweb

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/session"

	"strconv"

	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/logger"
)

const (
	DefaultGzipLevel = 9
	gzipScheme       = "gzip"
	DefaultIndexPage = "index.html"
)

type (
	HttpServer struct {
		stdServer      *http.Server
		router         Router
		groups         []Group
		Modules        []*HttpModule
		DotApp         *DotWeb
		Validator      Validator
		sessionManager *session.SessionManager
		lock_session   *sync.RWMutex
		pool           *pool
		binder         Binder
		render         Renderer
		offline        bool
		virtualPath    string // virtual path when deploy on no root path
	}

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
	}
	// setup router
	server.router = NewRouter(server)
	server.stdServer = &http.Server{Handler: server}
	return server
}

// initConfig init config from app config
func (server *HttpServer) initConfig() {
	server.SetEnabledGzip(server.ServerConfig().EnabledGzip)

	// VirtualPath config
	if server.virtualPath == "" {
		server.virtualPath = server.ServerConfig().VirtualPath
	}
}

// ServerConfig a shortcut for App.Config.ServerConfig
func (server *HttpServer) ServerConfig() *config.ServerNode {
	return server.DotApp.Config.Server
}

// SessionConfig a shortcut for App.Config.SessionConfig
func (server *HttpServer) SessionConfig() *config.SessionNode {
	return server.DotApp.Config.Session
}

// SetBinder set custom Binder on HttpServer
func (server *HttpServer) SetBinder(binder Binder) {
	server.binder = binder
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
	// check tls config
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

	// special handling for websocket and debugging
	if checkIsWebSocketRequest(req) {
		http.DefaultServeMux.ServeHTTP(w, req)
		core.GlobalState.AddRequestCount(req.URL.Path, defaultHttpCode, 1)
	} else {
		// setup header
		w.Header().Set(HeaderServer, DefaultServerName)
		// maintenance mode
		if server.IsOffline() {
			server.DotApp.OfflineServer.ServeHTTP(w, req)
		} else {
			httpCtx := prepareHttpContext(server, w, req)

			// process OnBeginRequest of modules
			for _, module := range server.Modules {
				if module.OnBeginRequest != nil {
					module.OnBeginRequest(httpCtx)
				}
			}

			if !httpCtx.IsEnd() {
				server.Router().ServeHTTP(httpCtx)
			}

			// process OnEndRequest of modules
			for _, module := range server.Modules {
				if module.OnEndRequest != nil {
					module.OnEndRequest(httpCtx)
				}
			}
			core.GlobalState.AddRequestCount(httpCtx.Request().Path(), httpCtx.Response().HttpCode(), 1)

			releaseHttpContext(server, httpCtx)
		}
	}
}

// IsOffline check server is set offline state
func (server *HttpServer) IsOffline() bool {
	return server.offline
}

// SetVirtualPath set current server's VirtualPath
func (server *HttpServer) SetVirtualPath(path string) {
	server.virtualPath = path
	logger.Logger().Debug("DotWeb:HttpServer SetVirtualPath ["+path+"]", LogTarget_HttpServer)

}

// VirtualPath return current server's VirtualPath
func (server *HttpServer) VirtualPath() string {
	return server.virtualPath
}

// SetOffline set server offline config
func (server *HttpServer) SetOffline(offline bool, offlineText string, offlineUrl string) {
	server.offline = offline
	logger.Logger().Debug("DotWeb:HttpServer SetOffline ["+strconv.FormatBool(offline)+", "+offlineText+", "+offlineUrl+"]", LogTarget_HttpServer)

}

// IndexPage default index page name
func (server *HttpServer) IndexPage() string {
	if server.ServerConfig().IndexPage == "" {
		return DefaultIndexPage
	} else {
		return server.ServerConfig().IndexPage
	}
}

// SetIndexPage set default index page name
func (server *HttpServer) SetIndexPage(indexPage string){
	server.ServerConfig().IndexPage = indexPage
	logger.Logger().Debug("DotWeb:HttpServer SetIndexPage ["+indexPage+"]", LogTarget_HttpServer)
}

// SetSessionConfig set session store config
func (server *HttpServer) SetSessionConfig(storeConfig *session.StoreConfig) {
	// sync session config
	server.SessionConfig().Timeout = storeConfig.Maxlifetime
	server.SessionConfig().SessionMode = storeConfig.StoreName
	server.SessionConfig().ServerIP = storeConfig.ServerIP
	server.SessionConfig().BackupServerUrl = storeConfig.BackupServerUrl
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
	storeConfig.BackupServerUrl = server.SessionConfig().BackupServerUrl
	storeConfig.StoreKeyPre = server.SessionConfig().StoreKeyPre
	storeConfig.CookieName = server.SessionConfig().CookieName

	if server.sessionManager == nil {
		// setup session
		server.lock_session.Lock()
		if manager, err := session.NewDefaultSessionManager(storeConfig); err != nil {
			// panic error with create session manager
			panic(err.Error())
		} else {
			server.sessionManager = manager
		}
		server.lock_session.Unlock()
	}
	logger.Logger().Debug("DotWeb:HttpServer InitSessionManager ["+jsonutil.GetJsonString(storeConfig)+"]", LogTarget_HttpServer)
}

// setDotApp associates the dotApp to the current HttpServer
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

// ServerFile a shortcut for router.ServeFiles(path, fileRoot)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (server *HttpServer) ServerFile(path string, fileRoot string) RouterNode {
	return server.Router().ServerFile(path, fileRoot)
}

// RegisterServerFile a shortcut for router.RegisterServerFile(routeMethod, path, fileRoot)
// simple demo:server.RegisterServerFile(RouteMethod_GET, "/src/*filepath", "/var/www")
func (server *HttpServer) RegisterServerFile(routeMethod string, path string, fileRoot string) RouterNode {
	return server.Router().RegisterServerFile(routeMethod, path, fileRoot)
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
		if server.DotApp.RunMode() == RunMode_Development {
			server.SetRenderer(NewInnerRendererNoCache())
		} else {
			server.SetRenderer(NewInnerRenderer())
		}
	}
	return server.render
}

// SetRenderer set custom renderer in server
func (server *HttpServer) SetRenderer(r Renderer) {
	server.render = r
	logger.Logger().Debug("DotWeb:HttpServer SetRenderer", LogTarget_HttpServer)
}

// SetEnabledAutoHEAD set route use auto head
// set EnabledAutoHEAD true or false
// default is false
func (server *HttpServer) SetEnabledAutoHEAD(isEnabled bool) {
	server.ServerConfig().EnabledAutoHEAD = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledAutoHEAD ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledAutoOPTIONS set route use auto options
// set SetEnabledAutoOPTIONS true or false
// default is false
func (server *HttpServer) SetEnabledAutoOPTIONS(isEnabled bool) {
	server.ServerConfig().EnabledAutoOPTIONS = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledAutoOPTIONS ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledRequestID set create unique request id per request
// set EnabledRequestID true or false
// default is false
func (server *HttpServer) SetEnabledRequestID(isEnabled bool) {
	server.ServerConfig().EnabledRequestID = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledRequestID ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledListDir set whether to allow listing of directories, default is false
func (server *HttpServer) SetEnabledListDir(isEnabled bool) {
	server.ServerConfig().EnabledListDir = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledListDir ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledSession set whether to enable session, default is false
func (server *HttpServer) SetEnabledSession(isEnabled bool) {
	server.SessionConfig().EnabledSession = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledSession ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledGzip set whether to enable gzip, default is false
func (server *HttpServer) SetEnabledGzip(isEnabled bool) {
	server.ServerConfig().EnabledGzip = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledGzip ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// SetEnabledBindUseJsonTag set whethr to enable json tab on Bind, default is false
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

// SetEnabledStaticFileMiddleware set flag which enabled or disabled middleware for static-file route
func (server *HttpServer) SetEnabledStaticFileMiddleware(isEnabled bool) {
	server.ServerConfig().EnabledStaticFileMiddleware = isEnabled
	logger.Logger().Debug("DotWeb:HttpServer SetEnabledStaticFileMiddleware ["+strconv.FormatBool(isEnabled)+"]", LogTarget_HttpServer)
}

// RegisterModule add HttpModule
func (server *HttpServer) RegisterModule(module *HttpModule) {
	server.Modules = append(server.Modules, module)
	logger.Logger().Debug("DotWeb:HttpServer RegisterModule ["+module.Name+"]", LogTarget_HttpServer)
}

type LogJson struct {
	RequestUrl string
	HttpHeader string
	HttpBody   string
}

// check request is the websocket request
// check Connection contains upgrade
func checkIsWebSocketRequest(req *http.Request) bool {
	if strings.Index(strings.ToLower(req.Header.Get("Connection")), "upgrade") >= 0 {
		return true
	}
	return false
}

// check request is startwith /debug/
func checkIsDebugRequest(req *http.Request) bool {
	if strings.Index(req.RequestURI, "/debug/") == 0 {
		return true
	}
	return false
}

// prepareHttpContext init HttpContext, init session & gzip config on HttpContext
func prepareHttpContext(server *HttpServer, w http.ResponseWriter, req *http.Request) *HttpContext {
	// get from pool
	response := server.pool.response.Get().(*Response)
	request := server.pool.request.Get().(*Request)
	httpCtx := server.pool.context.Get().(*HttpContext)
	httpCtx.reset(response, request, server, nil, nil, nil)
	response.reset(w)
	request.reset(req, httpCtx)

	// session
	// if exists client-sessionid, use it
	// if not exists client-sessionid, new one
	if httpCtx.HttpServer().SessionConfig().EnabledSession {
		sessionId, err := httpCtx.HttpServer().GetSessionManager().GetClientSessionID(httpCtx.Request().Request)
		if err == nil && sessionId != "" {
			httpCtx.sessionID = sessionId
		} else {
			httpCtx.sessionID = httpCtx.HttpServer().GetSessionManager().NewSessionID()
			cookie := &http.Cookie{
				Name:  httpCtx.HttpServer().sessionManager.StoreConfig().CookieName,
				Value: url.QueryEscape(httpCtx.SessionID()),
				Path:  "/",
			}
			httpCtx.SetCookie(cookie)
		}
	}
	// init gzip
	if httpCtx.HttpServer().ServerConfig().EnabledGzip {
		gw, err := gzip.NewWriterLevel(httpCtx.Response().Writer(), DefaultGzipLevel)
		if err != nil {
			panic("use gzip error -> " + err.Error())
		}
		grw := &gzipResponseWriter{Writer: gw, ResponseWriter: httpCtx.Response().Writer()}
		httpCtx.Response().reset(grw)
		httpCtx.Response().SetHeader(HeaderContentEncoding, gzipScheme)
	}

	return httpCtx
}

// releaseHttpContext release HttpContext, release gzip writer
func releaseHttpContext(server *HttpServer, httpCtx *HttpContext){
	// release response
	httpCtx.Response().release()
	server.pool.response.Put(httpCtx.Response())
	// release request
	httpCtx.Request().release()
	server.pool.request.Put(httpCtx.Request())
	// release context
	httpCtx.release()
	server.pool.context.Put(httpCtx)

	if server.ServerConfig().EnabledGzip {
		var w io.Writer
		w = httpCtx.Response().Writer().(*gzipResponseWriter).Writer
		w.(*gzip.Writer).Close()
	}
}