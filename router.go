package dotweb

import (
	"fmt"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/framework/convert"
	"github.com/devfeel/dotweb/framework/exception"
	_ "github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/framework/json"
	"github.com/devfeel/dotweb/logger"
	"golang.org/x/net/websocket"
	"net/http"
	paths "path"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (
	RouteMethod_Any       = "ANY"
	RouteMethod_GET       = "GET"
	RouteMethod_HEAD      = "HEAD"
	RouteMethod_OPTIONS   = "OPTIONS"
	RouteMethod_POST      = "POST"
	RouteMethod_PUT       = "PUT"
	RouteMethod_PATCH     = "PATCH"
	RouteMethod_DELETE    = "DELETE"
	RouteMethod_HiJack    = "HIJACK"
	RouteMethod_WebSocket = "WEBSOCKET"
)

const(
	routerExpressSplit = "^$^"
)

var (
	HttpMethodMap map[string]string
	valueNodePool sync.Pool
)

func init() {
	HttpMethodMap = make(map[string]string)
	HttpMethodMap["ANY"] = RouteMethod_Any
	HttpMethodMap["GET"] = RouteMethod_GET
	HttpMethodMap["HEAD"] = RouteMethod_HEAD
	HttpMethodMap["POST"] = RouteMethod_POST
	HttpMethodMap["PUT"] = RouteMethod_PUT
	HttpMethodMap["PATCH"] = RouteMethod_PATCH
	HttpMethodMap["OPTIONS"] = RouteMethod_OPTIONS
	HttpMethodMap["DELETE"] = RouteMethod_DELETE
	HttpMethodMap["HIJACK"] = RouteMethod_HiJack
	HttpMethodMap["WEBSOCKET"] = RouteMethod_WebSocket

	valueNodePool = sync.Pool{
		New: func() interface{} {
			return &ValueNode{}
		},
	}

}

type (
	// Router is the interface that wraps the router method.
	Router interface {
		ServeHTTP(ctx *HttpContext)
		ServerFile(path string, fileRoot string) RouterNode
		GET(path string, handle HttpHandle) RouterNode
		HEAD(path string, handle HttpHandle) RouterNode
		OPTIONS(path string, handle HttpHandle) RouterNode
		POST(path string, handle HttpHandle) RouterNode
		PUT(path string, handle HttpHandle) RouterNode
		PATCH(path string, handle HttpHandle) RouterNode
		DELETE(path string, handle HttpHandle) RouterNode
		HiJack(path string, handle HttpHandle)
		WebSocket(path string, handle HttpHandle)
		Any(path string, handle HttpHandle)
		RegisterRoute(routeMethod string, path string, handle HttpHandle) RouterNode
		RegisterHandler(name string, handler HttpHandle)
		GetHandler(name string) (HttpHandle, bool)
		MatchPath(ctx Context, routePath string) bool
	}

	RouterNode interface {
		Use(m ...Middleware) *Node
		AppMiddlewares() []Middleware
		GroupMiddlewares() []Middleware
		Middlewares() []Middleware
		Node() *Node
	}

	ValueNode struct {
		Params
		Method string
		Node   *Node
	}

	// router is a http.Handler which can be used to dispatch requests to different
	// handler functions via configurable routes
	router struct {
		Nodes        map[string]*Node
		allRouterExpress   map[string]struct{}
		server       *HttpServer
		handlerMap   map[string]HttpHandle
		handlerMutex *sync.RWMutex

		// Enables automatic redirection if the current route can't be matched but a
		// handler for the path with (without) the trailing slash exists.
		// For example if /foo/ is requested but a route only exists for /foo, the
		// client is redirected to /foo with http status code 301 for GET requests
		// and 307 for all other request methods.
		RedirectTrailingSlash bool

		// If enabled, the router tries to fix the current request path, if no
		// handle is registered for it.
		// First superfluous path elements like ../ or // are removed.
		// Afterwards the router does a case-insensitive lookup of the cleaned path.
		// If a handle can be found for this route, the router makes a redirection
		// to the corrected path with status code 301 for GET requests and 307 for
		// all other request methods.
		// For example /FOO and /..//Foo could be redirected to /foo.
		// RedirectTrailingSlash is independent of this option.
		RedirectFixedPath bool

		// If enabled, the router automatically replies to OPTIONS requests.
		// Custom OPTIONS handlers take priority over automatic replies.
		HandleOPTIONS bool
	}

	// Handle is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but has a third parameter for the values of
	// wildcards (variables).
	RouterHandle func(ctx *HttpContext)

	// Param is a single URL parameter, consisting of a key and a value.
	Param struct {
		Key   string
		Value string
	}

	// Params is a Param-slice, as returned by the router.
	// The slice is ordered, the first URL parameter is also the first slice value.
	// It is therefore safe to read values by the index.
	Params []Param
)

// ByName returns the value of the first Param which key matches the given name.
// If no matching Param is found, an empty string is returned.
func (ps Params) ByName(name string) string {
	for i := range ps {
		if ps[i].Key == name {
			return ps[i].Value
		}
	}
	return ""
}

// New returns a new initialized Router.
// Path auto-correction, including trailing slashes, is enabled by default.
func NewRouter(server *HttpServer) *router {
	return &router{
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
		HandleOPTIONS:         true,
		allRouterExpress:      make(map[string]struct{}),
		server:                server,
		handlerMap:            make(map[string]HttpHandle),
		handlerMutex:          new(sync.RWMutex),
	}
}

func (r *router) RegisterHandler(name string, handler HttpHandle) {
	r.handlerMutex.Lock()
	r.handlerMap[name] = handler
	r.handlerMutex.Unlock()
}

func (r *router) GetHandler(name string) (HttpHandle, bool) {
	r.handlerMutex.RLock()
	v, exists := r.handlerMap[name]
	r.handlerMutex.RUnlock()
	return v, exists
}

func (r *router) MatchPath(ctx Context, routePath string) bool {
	if root := r.Nodes[ctx.Request().Method]; root != nil {
		n := root.getNode(routePath)
		return n == ctx.RouterNode().Node()
	}
	return false
}

func (r *router) getNode(httpMethod string, routePath string) *Node{
	if root := r.Nodes[httpMethod]; root != nil {
		n := root.getNode(routePath)
		return n
	}
	return nil
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *router) ServeHTTP(ctx *HttpContext) {
	req := ctx.Request().Request
	w := ctx.Response().Writer()
	path := req.URL.Path
	if root := r.Nodes[req.Method]; root != nil {
		if handle, ps, node, tsr := root.getValue(path); handle != nil {
			ctx.routerParams = ps
			ctx.routerNode = node
			handle(ctx)
			return
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301 // Permanent redirect, request with GET method
			if req.Method != "GET" {
				// Temporary redirect, request with same method
				// As of Go 1.3, Go does not support status code 308.
				code = 307
			}

			if tsr && r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(w, req, req.URL.String(), code)
				return
			}

			// Try to fix the request path
			if r.RedirectFixedPath {
				fixedPath, found := root.findCaseInsensitivePath(
					//file.CleanPath(path),
					paths.Clean(path),
					r.RedirectTrailingSlash,
				)
				if found {
					req.URL.Path = string(fixedPath)
					http.Redirect(w, req, req.URL.String(), code)
					return
				}
			}
		}
	}

	if req.Method == "OPTIONS" {
		// Handle OPTIONS requests
		if r.HandleOPTIONS {
			if allow := r.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set("Allow", allow)
				return
			}
		}
	} else {
		// Handle 405
		if allow := r.allowed(path, req.Method); len(allow) > 0 {
			w.Header().Set("Allow", allow)
			ctx.Response().SetStatusCode(http.StatusMethodNotAllowed)
			r.server.DotApp.MethodNotAllowedHandler(ctx)
			return
		}
	}

	// Handle 404
	ctx.Response().SetStatusCode(http.StatusNotFound)
	r.server.DotApp.NotFoundHandler(ctx)
}

//wrap HttpHandle to Handle
func (r *router) wrapRouterHandle(handler HttpHandle, isHijack bool) RouterHandle {
	return func(httpCtx *HttpContext) {
		httpCtx.handler = handler

		//do features
		FeatureTools.InitFeatures(r.server, httpCtx)

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
				if r.server.DotApp.ExceptionHandler != nil {
					r.server.DotApp.ExceptionHandler(httpCtx, fmt.Errorf("%v", err))
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
					logger.Logger().Error(logString, LogTarget_HttpServer)
				}

				//增加错误计数
				core.GlobalState.AddErrorCount(httpCtx.Request().Path(), fmt.Errorf("%v", err), 1)
			}

			FeatureTools.ReleaseFeatures(r.server, httpCtx)

			//cancle Context
			if httpCtx.cancle != nil {
				httpCtx.cancle()
			}
		}()

		//处理用户handle
		var ctxErr error
		//if len(r.server.DotApp.Middlewares) > 0 {
		//	ctxErr = r.server.DotApp.Middlewares[0].Handle(httpCtx)
		//} else {
		//	ctxErr = handler(httpCtx)
		//}

		if len(httpCtx.routerNode.AppMiddlewares()) > 0 {
			ctxErr = httpCtx.routerNode.AppMiddlewares()[0].Handle(httpCtx)
		} else {
			ctxErr = handler(httpCtx)
		}

		if ctxErr != nil {
			//handler the exception
			if r.server.DotApp.ExceptionHandler != nil {
				r.server.DotApp.ExceptionHandler(httpCtx, ctxErr)
				//增加错误计数
				core.GlobalState.AddErrorCount(httpCtx.Request().Path(), ctxErr, 1)
			}
		}

	}
}

//wrap fileHandler to httprouter.Handle
func (r *router) wrapFileHandle(fileHandler http.Handler) RouterHandle {
	return func(httpCtx *HttpContext) {
		startTime := time.Now()
		httpCtx.Request().URL.Path = httpCtx.RouterParams().ByName("filepath")
		fileHandler.ServeHTTP(httpCtx.Response().Writer(), httpCtx.Request().Request)
		timetaken := int64(time.Now().Sub(startTime) / time.Millisecond)
		//HttpServer Logging
		logger.Logger().Debug(httpCtx.Request().Url()+" "+logRequest(httpCtx.Request().Request, timetaken), LogTarget_HttpRequest)
	}
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *router) GET(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_GET, path, handle)
}

// ANY is a shortcut for router.Handle("Any", path, handle)
// it support GET\HEAD\POST\PUT\PATCH\OPTIONS\DELETE
func (r *router) Any(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_HEAD, path, handle)
	r.RegisterRoute(RouteMethod_GET, path, handle)
	r.RegisterRoute(RouteMethod_POST, path, handle)
	r.RegisterRoute(RouteMethod_PUT, path, handle)
	r.RegisterRoute(RouteMethod_DELETE, path, handle)
	r.RegisterRoute(RouteMethod_PATCH, path, handle)
	r.RegisterRoute(RouteMethod_OPTIONS, path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *router) HEAD(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_HEAD, path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *router) OPTIONS(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_OPTIONS, path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *router) POST(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_POST, path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *router) PUT(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_PUT, path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *router) PATCH(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_PATCH, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *router) DELETE(path string, handle HttpHandle) RouterNode {
	return r.RegisterRoute(RouteMethod_DELETE, path, handle)
}

func (r *router) HiJack(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_HiJack, path, handle)
}

func (r *router) WebSocket(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_WebSocket, path, handle)
}

// shortcut for router.Handle(httpmethod, path, handle)
// support GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS\HiJack\WebSocket\ANY
func (r *router) RegisterRoute(routeMethod string, path string, handle HttpHandle) RouterNode {
	var node *Node
	handleName := handlerName(handle)
	routeMethod = strings.ToUpper(routeMethod)
	if _, exists := HttpMethodMap[routeMethod]; !exists {
		logger.Logger().Warn("DotWeb:Router:RegisterRoute failed [illegal method] ["+routeMethod+"] ["+path+"] ["+handleName+"]", LogTarget_HttpServer)
		return nil
	} else {
		logger.Logger().Debug("DotWeb:Router:RegisterRoute success ["+routeMethod+"] ["+path+"] ["+handleName+"]", LogTarget_HttpServer)
	}

	//websocket mode,use default httpserver
	if routeMethod == RouteMethod_WebSocket {
		http.Handle(path, websocket.Handler(r.wrapWebSocketHandle(handle)))
		return node
	}

	//hijack mode,use get and isHijack = true
	if routeMethod == RouteMethod_HiJack {
		r.add(RouteMethod_GET, path, r.wrapRouterHandle(handle, true))
	} else {
		//GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS mode
		node = r.add(routeMethod, path, r.wrapRouterHandle(handle, false))
	}

	//if set auto-head, add head router
	//only enabled in hijack\GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS
	if r.server.ServerConfig().EnabledAutoHEAD {
		if routeMethod == RouteMethod_HiJack {
			r.add(RouteMethod_HEAD, path, r.wrapRouterHandle(handle, true))
		} else if routeMethod != RouteMethod_Any {
			r.add(RouteMethod_HEAD, path, r.wrapRouterHandle(handle, false))
		}
	}
	return node
}

// ServerFile is a shortcut for router.ServeFiles(path, filepath)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (r *router) ServerFile(path string, fileroot string) RouterNode {
	node := &Node{}
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}
	var root http.FileSystem
	root = http.Dir(fileroot)
	if !r.server.ServerConfig().EnabledListDir {
		root = &core.HideReaddirFS{root}
	}
	fileServer := http.FileServer(root)
	node = r.add(RouteMethod_GET, path, r.wrapFileHandle(fileServer))
	return node
}

func handlerName(h HttpHandle) string {
	t := reflect.ValueOf(h).Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
	}
	return t.String()
}

// Handle registers a new request handle with the given path and method.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (r *router) add(method, path string, handle RouterHandle, m ...Middleware) (outnode *Node) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.Nodes == nil {
		r.Nodes = make(map[string]*Node)
	}

	root := r.Nodes[method]
	if root == nil {
		root = new(Node)
		r.Nodes[method] = root
	}
	//fmt.Println("Handle => ", method, " - ", *root, " - ", path)
	outnode = root.addRoute(path, handle, m...)
	outnode.fullPath = path
	r.allRouterExpress[method + routerExpressSplit + path] = struct{}{}
	return
}

func (r *router) allowed(path, reqMethod string) (allow string) {
	if path == "*" { // server-wide
		for method := range r.Nodes {
			if method == "OPTIONS" {
				continue
			}

			// add request method to list of allowed methods
			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else { // specific path
		for method := range r.Nodes {
			// Skip the requested method - we already tried this one
			if method == reqMethod || method == "OPTIONS" {
				continue
			}

			handle, _, _, _ := r.Nodes[method].getValue(path)
			if handle != nil {
				// add request method to list of allowed methods
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}
	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

//wrap HttpHandle to websocket.Handle
func (r *router) wrapWebSocketHandle(handler HttpHandle) websocket.Handler {
	return func(ws *websocket.Conn) {
		//get from pool
		req := r.server.pool.request.Get().(*Request)
		httpCtx := r.server.pool.context.Get().(*HttpContext)
		httpCtx.reset(nil, req, r.server, nil, nil, handler)
		req.reset(ws.Request(), httpCtx)
		httpCtx.webSocket = &WebSocket{
			Conn: ws,
		}
		httpCtx.isWebSocket = true

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
				logger.Logger().Error(logString, LogTarget_HttpServer)

				//增加错误计数
				core.GlobalState.AddErrorCount(httpCtx.Request().Path(), fmt.Errorf("%v", err), 1)
			}
			timetaken := int64(time.Now().Sub(startTime) / time.Millisecond)
			//HttpServer Logging
			logger.Logger().Debug(httpCtx.Request().Url()+" "+logWebsocketContext(httpCtx, timetaken), LogTarget_HttpRequest)

			//release request
			req.release()
			r.server.pool.request.Put(req)
			//release context
			httpCtx.release()
			r.server.pool.context.Put(httpCtx)
		}()

		handler(httpCtx)
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
