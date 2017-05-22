package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/framework/convert"
	"github.com/devfeel/dotweb/framework/exception"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/session"
	"net/http"
	"strings"
	"sync"
	"time"

	"compress/gzip"
	"github.com/devfeel/dotweb/config"
	"github.com/devfeel/dotweb/feature"
	"github.com/devfeel/dotweb/logger"
	"golang.org/x/net/websocket"
	"io"
	"net/url"
)

const (
	DefaultGzipLevel = 9
	gzipScheme       = "gzip"
)

type (
	//HttpModule定义
	HttpModule struct {
		//响应请求时作为 HTTP 执行管线链中的第一个事件发生
		OnBeginRequest func(Context)
		//响应请求时作为 HTTP 执行管线链中的最后一个事件发生。
		OnEndRequest func(Context)
	}

	//HttpServer定义
	HttpServer struct {
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
	return server
}

//ServeHTTP make sure request can be handled correctly
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
			server.Router().ServeHTTP(w, req)
		}
	}
}

//IsOffline check server is set offline state
func (server *HttpServer) IsOffline() bool {
	return server.offline
}

//SetOffline set server offline config
func (server *HttpServer) SetOffline(offline bool, offlineText string, offlineUrl string) {
	server.offline = offline
}

//set session store config
func (server *HttpServer) SetSessionConfig(storeConfig *session.StoreConfig) {
	//sync session config
	server.SessionConfig.Timeout = storeConfig.Maxlifetime
	server.SessionConfig.SessionMode = storeConfig.StoreName
	server.SessionConfig.ServerIP = storeConfig.ServerIP
}

//init session manager
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

/*
* 关联当前HttpServer实例对应的DotServer实例
 */
func (server *HttpServer) setDotApp(dotApp *DotWeb) {
	server.DotApp = dotApp
}

//get session manager in current httpserver
func (server *HttpServer) GetSessionManager() *session.SessionManager {
	if !server.SessionConfig.EnabledSession {
		return nil
	}
	return server.sessionManager
}

//get router interface in server
func (server *HttpServer) Router() Router {
	return server.router
}

//create new group with current HttpServer
func (server *HttpServer) Group(prefix string) Group {
	return NewGroup(prefix, server)
}

//get binder interface in server
func (server *HttpServer) Binder() Binder {
	return server.binder
}

//get renderer interface in server
//if no set, init InnerRenderer
func (server *HttpServer) Renderer() Renderer {
	if server.render == nil {
		server.render = NewInnerRenderer()
	}
	return server.render
}

//set custom renderer in server
func (server *HttpServer) SetRenderer(r Renderer) {
	server.render = r
}

//set EnabledAutoHEAD true or false
func (server *HttpServer) SetEnabledAutoHEAD(autoHEAD bool) {
	server.ServerConfig.EnabledAutoHEAD = autoHEAD
}

/*
设置是否允许目录浏览,默认为false
*/
func (server *HttpServer) SetEnabledListDir(isEnabled bool) {
	server.ServerConfig.EnabledListDir = isEnabled
}

/*
设置是否启用Session,默认为false
*/
func (server *HttpServer) SetEnabledSession(isEnabled bool) {
	server.SessionConfig.EnabledSession = isEnabled
}

/*
设置是否启用gzip,默认为false
*/
func (server *HttpServer) SetEnabledGzip(isEnabled bool) {
	server.ServerConfig.EnabledGzip = isEnabled
}

//do features...
func (server *HttpServer) doFeatures(ctx Context) Context {
	//处理 cros feature
	if server.Features.CROSConfig != nil {
		c := server.Features.CROSConfig
		if c.EnabledCROS {
			FeatureTools.SetCROSConfig(ctx, c)
		}
	}
	return ctx
}

type LogJson struct {
	RequestUrl string
	HttpHeader string
	HttpBody   string
}

//wrap HttpHandle to Handle
func (server *HttpServer) wrapRouterHandle(handler HttpHandle, isHijack bool) RouterHandle {
	return func(w http.ResponseWriter, r *http.Request, vnode *ValueNode) {
		//get from pool
		res := server.pool.response.Get().(*Response)
		res.reset(w)
		req := server.pool.request.Get().(*Request)
		req.reset(r)
		httpCtx := server.pool.context.Get().(*HttpContext)
		httpCtx.reset(res, req, server, vnode.Node, vnode.Params, handler)

		//gzip
		if server.ServerConfig.EnabledGzip {
			gw, err := gzip.NewWriterLevel(w, DefaultGzipLevel)
			if err != nil {
				panic("use gzip error -> " + err.Error())
			}
			grw := &gzipResponseWriter{Writer: gw, ResponseWriter: w}
			res.reset(grw)
			httpCtx.Response().SetHeader(HeaderContentEncoding, gzipScheme)
		}
		//增加状态计数
		GlobalState.AddRequestCount(1)

		//session
		//if exists client-sessionid, use it
		//if not exists client-sessionid, new one
		if server.SessionConfig.EnabledSession {
			sessionId, err := server.GetSessionManager().GetClientSessionID(r)
			if err == nil && sessionId != "" {
				httpCtx.sessionID = sessionId
			} else {
				httpCtx.sessionID = server.GetSessionManager().NewSessionID()
				cookie := &http.Cookie{
					Name:  server.sessionManager.CookieName,
					Value: url.QueryEscape(httpCtx.sessionID),
					Path:  "/",
				}
				httpCtx.SetCookie(cookie)
			}
		}

		//hijack处理
		if isHijack {
			_, hijack_err := httpCtx.Hijack()
			if hijack_err != nil {
				//输出内容
				httpCtx.Response().WriteHeader(http.StatusInternalServerError)
				httpCtx.Response().Header().Set(HeaderContentType, CharsetUTF8)
				httpCtx.WriteString(hijack_err.Error())
				return
			}
		}

		defer func() {
			var errmsg string
			if err := recover(); err != nil {
				errmsg = exception.CatchError("HttpServer::RouterHandle", LogTarget_HttpServer, err)

				//handler the exception
				if server.DotApp.ExceptionHandler != nil {
					server.DotApp.ExceptionHandler(httpCtx, fmt.Errorf("%v", err))
				}

				//if set enabledLog, take the error log
				if logger.EnabledLog {
					//记录访问日志
					headinfo := fmt.Sprintln(httpCtx.Response().Header)
					logJson := LogJson{
						RequestUrl: httpCtx.Request().RequestURI,
						HttpHeader: headinfo,
						HttpBody:   errmsg,
					}
					logString := jsonutil.GetJsonString(logJson)
					logger.Logger().Log(logString, LogTarget_HttpServer, LogLevel_Error)
				}

				//增加错误计数
				GlobalState.AddErrorCount(1)
			}

			if server.ServerConfig.EnabledGzip {
				var w io.Writer
				w = res.Writer().(*gzipResponseWriter).Writer
				w.(*gzip.Writer).Close()
			}
			//release response
			res.release()
			server.pool.response.Put(res)
			//release request
			req.release()
			server.pool.request.Put(req)
			//release context
			httpCtx.release()
			server.pool.context.Put(httpCtx)
		}()

		//do features
		server.doFeatures(httpCtx)

		//处理前置Module集合
		for _, module := range server.DotApp.Modules {
			if module.OnBeginRequest != nil {
				module.OnBeginRequest(httpCtx)
			}
		}

		//处理用户handle
		//if already set HttpContext.End,ignore user handler - fixed issue #5
		if !httpCtx.IsEnd() {
			var ctxErr error
			if len(server.DotApp.Middlewares) > 0 {
				ctxErr = server.DotApp.Middlewares[0].Handle(httpCtx)
			} else {
				ctxErr = handler(httpCtx)
			}
			if ctxErr != nil {
				//handler the exception
				if server.DotApp.ExceptionHandler != nil {
					server.DotApp.ExceptionHandler(httpCtx, ctxErr)
				}
			}
		}

		//处理后置Module集合
		for _, module := range server.DotApp.Modules {
			if module.OnEndRequest != nil {
				module.OnEndRequest(httpCtx)
			}
		}

	}
}

//wrap fileHandler to httprouter.Handle
func (server *HttpServer) wrapFileHandle(fileHandler http.Handler) RouterHandle {
	return func(w http.ResponseWriter, r *http.Request, vnode *ValueNode) {
		//增加状态计数
		GlobalState.AddRequestCount(1)
		startTime := time.Now()
		r.URL.Path = vnode.ByName("filepath")
		fileHandler.ServeHTTP(w, r)
		timetaken := int64(time.Now().Sub(startTime) / time.Millisecond)
		//HttpServer Logging
		logger.Logger().Log(r.URL.String()+" "+logRequest(r, timetaken), LogTarget_HttpRequest, LogLevel_Debug)
	}
}

//wrap HttpHandle to websocket.Handle
func (server *HttpServer) wrapWebSocketHandle(handler HttpHandle) websocket.Handler {
	return func(ws *websocket.Conn) {
		//get from pool
		req := server.pool.request.Get().(*Request)
		req.reset(ws.Request())
		httpCtx := server.pool.context.Get().(*HttpContext)
		httpCtx.reset(nil, req, server, nil, nil, handler)
		httpCtx.webSocket = &WebSocket{
			Conn: ws,
		}
		httpCtx.IsWebSocket = true

		startTime := time.Now()
		defer func() {
			var errmsg string
			if err := recover(); err != nil {
				errmsg = exception.CatchError("httpserver::WebsocketHandle", LogTarget_HttpServer, err)

				//记录访问日志
				headinfo := fmt.Sprintln(httpCtx.webSocket.Request().Header)
				logJson := LogJson{
					RequestUrl: httpCtx.webSocket.Request().RequestURI,
					HttpHeader: headinfo,
					HttpBody:   errmsg,
				}
				logString := jsonutil.GetJsonString(logJson)
				logger.Logger().Log(logString, LogTarget_HttpServer, LogLevel_Error)

				//增加错误计数
				GlobalState.AddErrorCount(1)
			}
			timetaken := int64(time.Now().Sub(startTime) / time.Millisecond)
			//HttpServer Logging
			logger.Logger().Log(httpCtx.Request().Url()+" "+logWebsocketContext(httpCtx, timetaken), LogTarget_HttpRequest, LogLevel_Debug)

			//release request
			req.release()
			server.pool.request.Put(req)
			//release context
			httpCtx.release()
			server.pool.context.Put(httpCtx)
		}()

		handler(httpCtx)

		//增加状态计数
		GlobalState.AddRequestCount(1)
	}
}

//get default log string
func logWebsocketContext(ctx Context, timetaken int64) string {
	var reqbytelen, resbytelen, method, proto, status, userip string
	if ctx != nil {
		reqbytelen = convert.Int642String(ctx.Request().ContentLength)
		resbytelen = "0"
		method = ctx.Request().Method
		proto = ctx.Request().Proto
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

func logRequest(req *http.Request, timetaken int64) string {
	var reqbytelen, resbytelen, method, proto, status, userip string
	reqbytelen = convert.Int642String(req.ContentLength)
	resbytelen = ""
	method = req.Method
	proto = req.Proto
	status = "200"
	userip = req.RemoteAddr

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
