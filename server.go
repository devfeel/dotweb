package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/framework/convert"
	"github.com/devfeel/dotweb/framework/exception"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/framework/log"
	"github.com/devfeel/dotweb/session"
	"net/http"
	"strings"
	"sync"
	"time"

	"compress/gzip"
	"github.com/devfeel/dotweb/router"
	"golang.org/x/net/websocket"
	"io"
	"net/url"
)

const (
	RouteMethod_GET       = "GET"
	RouteMethod_HEAD      = "HEAD"
	RouteMethod_OPTIONS   = "OPTIONS"
	RouteMethod_POST      = "POST"
	RouteMethod_PUT       = "PUT"
	RouteMethod_PATCH     = "PATCH"
	RouteMethod_DELETE    = "DELETE"
	RouteMethod_HiJack    = "HiJack"
	RouteMethod_WebSocket = "WebSocket"

	DefaultGzipLevel = 9
	gzipScheme       = "gzip"
)

type (
	//HttpModule定义
	HttpModule struct {
		//响应请求时作为 HTTP 执行管线链中的第一个事件发生
		OnBeginRequest func(*HttpContext)
		//响应请求时作为 HTTP 执行管线链中的最后一个事件发生。
		OnEndRequest func(*HttpContext)
	}

	//HttpServer定义
	HttpServer struct {
		router         *router.Router
		DotApp         *DotWeb
		sessionManager *session.SessionManager
		lock_session   *sync.RWMutex
		pool           *pool
		ServerConfig   *ServerConfig
		binder         Binder
	}

	//server global config
	ServerConfig struct {
		EnabledDebug   bool //is enabled debug mode
		EnabledSession bool //is enabled seesion
		EnabledGzip    bool //is enabled gzip
	}

	//pool定义
	pool struct {
		response sync.Pool
		context  sync.Pool
	}
)

func NewServerConfig() *ServerConfig {
	config := &ServerConfig{
		EnabledDebug:   false,
		EnabledSession: false,
	}
	return config
}

// Handle is a function that can be registered to a route to handle HTTP
// requests. Like http.HandlerFunc, but has a third parameter for the values of
// wildcards (variables).
type HttpHandle func(*HttpContext)

func NewHttpServer() *HttpServer {
	server := &HttpServer{
		router: router.New(),
		pool: &pool{
			response: sync.Pool{
				New: func() interface{} {
					return &Response{}
				},
			},
			context: sync.Pool{
				New: func() interface{} {
					return &HttpContext{}
				},
			},
		},
		ServerConfig: NewServerConfig(),
		lock_session: new(sync.RWMutex),
		binder:       &binder{},
	}

	return server
}

//ServeHTTP makes the httprouter implement the http.Handler interface.
func (server *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//针对websocket与调试信息特殊处理
	if checkIsWebSocketRequest(req) {
		http.DefaultServeMux.ServeHTTP(w, req)
	} else {
		//设置header信息
		w.Header().Set(HeaderServer, DefaultServerName)
		server.router.ServeHTTP(w, req)
	}
}

//init session manager
func (server *HttpServer) InitSessionManager(config *session.StoreConfig) {
	if server.sessionManager == nil {
		//设置Session
		server.lock_session.Lock()
		if manager, err := session.NewDefaultSessionManager(config); err != nil {
			//panic error with create session manager
			panic(err.Error())
		} else {
			server.sessionManager = manager
		}
		server.lock_session.Unlock()
	}
}

/*
* 关联当前HttpServer实例对应的DotServer实例
 */
func (server *HttpServer) setDotApp(dotApp *DotWeb) {
	server.DotApp = dotApp
}

//get session manager in current httpserver
func (server *HttpServer) GetSessionManager() *session.SessionManager {
	if !server.ServerConfig.EnabledSession {
		return nil
	}
	return server.sessionManager
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (server *HttpServer) GET(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_GET, path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (server *HttpServer) HEAD(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_HEAD, path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (server *HttpServer) OPTIONS(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_OPTIONS, path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (server *HttpServer) POST(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_POST, path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (server *HttpServer) PUT(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_PUT, path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (server *HttpServer) PATCH(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_PATCH, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (server *HttpServer) DELETE(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_DELETE, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (server *HttpServer) HiJack(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_GET, path, handle)
}

// shortcut for router.Handle(httpmethod, path, handle)
// support GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS\HiJack\WebSocket
func (server *HttpServer) RegisterRoute(routeMethod string, path string, handle HttpHandle) {

	logger.Log("Dotweb:RegisterRoute ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Debug)

	//hijack mode,use get and isHijack = true
	if routeMethod == RouteMethod_HiJack {
		server.router.Handle(RouteMethod_GET, path, server.wrapRouterHandle(handle, true))
		return
	}
	//websocket mode,use default httpserver
	if routeMethod == RouteMethod_WebSocket {
		http.Handle(path, websocket.Handler(server.wrapWebSocketHandle(handle)))
		return
	}

	//GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS mode
	server.router.Handle(routeMethod, path, server.wrapRouterHandle(handle, false))
	return
}

// ServerFile is a shortcut for router.ServeFiles(path, filepath)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (server *HttpServer) ServerFile(urlpath string, filepath string) {
	server.router.ServeFiles(urlpath, http.Dir(filepath))
}

// WebSocket is a shortcut for websocket.Handler
func (server *HttpServer) WebSocket(path string, handle HttpHandle) {
	server.RegisterRoute(RouteMethod_WebSocket, path, handle)
}

func (server *HttpServer) Binder() Binder {
	return server.binder
}

type LogJson struct {
	RequestUrl string
	HttpHeader string
	HttpBody   string
}

//wrap HttpHandle to httprouter.Handle
func (server *HttpServer) wrapRouterHandle(handle HttpHandle, isHijack bool) router.Handle {
	return func(w http.ResponseWriter, r *http.Request, params router.Params) {
		//get from pool
		res := server.pool.response.Get().(*Response)
		res.Reset(w)
		httpCtx := server.pool.context.Get().(*HttpContext)
		httpCtx.Reset(res, r, server, params)

		//gzip
		if server.ServerConfig.EnabledGzip {
			gw, err := gzip.NewWriterLevel(w, DefaultGzipLevel)
			if err != nil {
				panic("use gzip error -> " + err.Error())
			}
			grw := &gzipResponseWriter{Writer: gw, ResponseWriter: w}
			res.Reset(grw)
			httpCtx.SetHeader(HeaderContentEncoding, gzipScheme)
		}
		//增加状态计数
		GlobalState.AddRequestCount(1)

		//session
		//if exists client-sessionid, use it
		//if not exists client-sessionid, new one
		if server.ServerConfig.EnabledSession {
			sessionId, err := server.GetSessionManager().GetClientSessionID(r)
			if err == nil && sessionId != "" {
				httpCtx.SessionID = sessionId
			} else {
				httpCtx.SessionID = server.GetSessionManager().NewSessionID()
				cookie := http.Cookie{
					Name:  server.sessionManager.CookieName,
					Value: url.QueryEscape(httpCtx.SessionID),
					Path:  "/",
				}
				httpCtx.WriteCookieObj(cookie)
			}
		}

		//hijack处理
		if isHijack {
			_, hijack_err := httpCtx.Hijack()
			if hijack_err != nil {
				//输出内容
				httpCtx.Response.WriteHeader(http.StatusInternalServerError)
				httpCtx.Response.Header().Set(HeaderContentType, CharsetUTF8)
				httpCtx.WriteString(hijack_err.Error())
			}
		}

		startTime := time.Now()
		defer func() {
			var errmsg string
			if err := recover(); err != nil {
				errmsg = exception.CatchError("httpserver::RouterHandle", LogTarget_HttpServer, err)

				//默认异常处理
				if server.DotApp.ExceptionHandler != nil {
					server.DotApp.ExceptionHandler(httpCtx, err)
				}

				//记录访问日志
				headinfo := fmt.Sprintln(httpCtx.Response.Header)
				logJson := LogJson{
					RequestUrl: httpCtx.Request.RequestURI,
					HttpHeader: headinfo,
					HttpBody:   errmsg,
				}
				logString := jsonutil.GetJsonString(logJson)
				logger.Log(logString, LogTarget_HttpServer, LogLevel_Error)

				//增加错误计数
				GlobalState.AddErrorCount(1)
			}
			timetaken := int64(time.Now().Sub(startTime) / time.Millisecond)
			//HttpServer Logging
			logger.Log(httpCtx.Url()+" "+logString(httpCtx, timetaken), LogTarget_HttpRequest, LogLevel_Debug)

			if server.ServerConfig.EnabledGzip {
				var w io.Writer
				w = res.Writer().(*gzipResponseWriter).Writer
				w.(*gzip.Writer).Close()
			}
			// Return to pool
			server.pool.response.Put(res)
			//release context
			httpCtx.release()
			server.pool.context.Put(httpCtx)
		}()

		//处理前置Module集合
		for _, module := range server.DotApp.Modules {
			if module.OnBeginRequest != nil {
				module.OnBeginRequest(httpCtx)
			}
		}

		//处理用户handle
		//if already set HttpContext.End,ignore user handler - fixed issue #5
		if !httpCtx.IsEnd() {
			handle(httpCtx)
		}

		//处理后置Module集合
		for _, module := range server.DotApp.Modules {
			if module.OnEndRequest != nil {
				module.OnEndRequest(httpCtx)
			}
		}

	}
}

//wrap HttpHandle to websocket.Handle
func (server *HttpServer) wrapWebSocketHandle(handle HttpHandle) websocket.Handler {
	return func(ws *websocket.Conn) {
		//get from pool
		httpCtx := server.pool.context.Get().(*HttpContext)
		httpCtx.Reset(nil, ws.Request(), server, nil)
		httpCtx.WebSocket = &WebSocket{
			Conn: ws,
		}
		httpCtx.IsWebSocket = true

		startTime := time.Now()
		defer func() {
			var errmsg string
			if err := recover(); err != nil {
				errmsg = exception.CatchError("httpserver::WebsocketHandle", LogTarget_HttpServer, err)

				//记录访问日志
				headinfo := fmt.Sprintln(httpCtx.WebSocket.Conn.Request().Header)
				logJson := LogJson{
					RequestUrl: httpCtx.WebSocket.Conn.Request().RequestURI,
					HttpHeader: headinfo,
					HttpBody:   errmsg,
				}
				logString := jsonutil.GetJsonString(logJson)
				logger.Log(logString, LogTarget_HttpServer, LogLevel_Error)

				//增加错误计数
				GlobalState.AddErrorCount(1)
			}
			timetaken := int64(time.Now().Sub(startTime) / time.Millisecond)
			//HttpServer Logging
			logger.Log(httpCtx.Url()+" "+logString(httpCtx, timetaken), LogTarget_HttpRequest, LogLevel_Debug)

			// Return to pool
			server.pool.context.Put(httpCtx)
		}()

		handle(httpCtx)

		//增加状态计数
		GlobalState.AddRequestCount(1)
	}
}

//get default log string
func logString(ctx *HttpContext, timetaken int64) string {

	var reqbytelen, resbytelen, method, proto, status, userip string
	if !ctx.IsWebSocket {
		reqbytelen = convert.Int642String(ctx.Request.ContentLength)
		resbytelen = convert.Int642String(ctx.Response.Size)
		method = ctx.Request.Method
		proto = ctx.Request.Proto
		status = convert.Int2String(ctx.Response.Status)
		userip = ctx.RemoteIP()
	} else {
		reqbytelen = convert.Int642String(ctx.Request.ContentLength)
		resbytelen = "0"
		method = ctx.Request.Method
		proto = ctx.Request.Proto
		status = "0"
		userip = ctx.RemoteIP()
	}

	log := method + " "
	log += userip + " "
	log += proto + " "
	log += status + " "
	log += reqbytelen + " "
	log += resbytelen + " "
	log += convert.Int642String(timetaken)

	return log
}

//check request is the websocket request
func checkIsWebSocketRequest(req *http.Request) bool {
	if req.Header.Get("Connection") == "Upgrade" {
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
