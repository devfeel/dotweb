package dotweb

import (
	"github.com/devfeel/dotweb/core"
	_ "github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/logger"
	"golang.org/x/net/websocket"
	"net/http"
	paths "path"
	"reflect"
	"runtime"
	"strings"
	"sync"
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
	RouteMethod_HiJack    = "HiJack"
	RouteMethod_WebSocket = "WebSocket"
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
		ServeHTTP(w http.ResponseWriter, req *http.Request)
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
		MatchPath(ctx *HttpContext, routePath string) bool
	}

	RouterNode interface {
		Use(m ...Middleware) *Node
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
		Nodes map[string]*Node

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

		// If enabled, the router checks if another method is allowed for the
		// current route, if the current request can not be routed.
		// If this is the case, the request is answered with 'Method Not Allowed'
		// and HTTP status code 405.
		// If no other Method is allowed, the request is delegated to the NotFound
		// handler.
		HandleMethodNotAllowed bool

		// If enabled, the router automatically replies to OPTIONS requests.
		// Custom OPTIONS handlers take priority over automatic replies.
		HandleOPTIONS bool

		// Configurable http.Handler which is called when a request
		// cannot be routed and HandleMethodNotAllowed is true.
		// If it is not set, http.Error with http.StatusMethodNotAllowed is used.
		// The "Allow" header with allowed request methods is set before the handler
		// is called.
		MethodNotAllowed http.Handler
	}

	// Handle is a function that can be registered to a route to handle HTTP
	// requests. Like http.HandlerFunc, but has a third parameter for the values of
	// wildcards (variables).
	RouterHandle func(http.ResponseWriter, *http.Request, *ValueNode)

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
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      true,
		HandleMethodNotAllowed: true,
		HandleOPTIONS:          true,
		server:                 server,
		handlerMap:             make(map[string]HttpHandle),
		handlerMutex:           new(sync.RWMutex),
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

func (r *router) MatchPath(ctx *HttpContext, routePath string) bool {
	if root := r.Nodes[ctx.Method()]; root != nil {
		n := root.getNode(routePath)
		return n == ctx.RouterNode.Node()
	}
	return false
}

// ServeHTTP makes the router implement the http.Handler interface.
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	if root := r.Nodes[req.Method]; root != nil {
		if handle, ps, node, tsr := root.getValue(path); handle != nil {
			vn := valueNodePool.Get().(*ValueNode)
			vn.Params = ps
			vn.Node = node
			vn.Method = req.Method
			//user handle
			handle(w, req, vn)
			vn.Params = nil
			vn.Node = nil
			vn.Method = ""
			valueNodePool.Put(vn)
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
		if r.HandleMethodNotAllowed {
			if allow := r.allowed(path, req.Method); len(allow) > 0 {
				w.Header().Set("Allow", allow)
				if r.MethodNotAllowed != nil {
					r.MethodNotAllowed.ServeHTTP(w, req)
				} else {
					http.Error(w,
						http.StatusText(http.StatusMethodNotAllowed),
						http.StatusMethodNotAllowed,
					)
				}
				return
			}
		}
	}

	// Handle 404
	if r.server.DotApp.NotFoundHandler != nil {
		r.server.DotApp.NotFoundHandler.ServeHTTP(w, req)
	} else {
		http.NotFound(w, req)
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
		logger.Logger().Log("Dotweb:Router:RegisterRoute failed [illegal method] ["+routeMethod+"] ["+path+"] ["+handleName+"]", LogTarget_HttpServer, LogLevel_Warn)
		return nil
	} else {
		logger.Logger().Log("Dotweb:Router:RegisterRoute success ["+routeMethod+"] ["+path+"] ["+handleName+"]", LogTarget_HttpServer, LogLevel_Debug)
	}

	//websocket mode,use default httpserver
	if routeMethod == RouteMethod_WebSocket {
		http.Handle(path, websocket.Handler(r.server.wrapWebSocketHandle(handle)))
		return node
	}

	//hijack mode,use get and isHijack = true
	if routeMethod == RouteMethod_HiJack {
		r.add(RouteMethod_GET, path, r.server.wrapRouterHandle(handle, true))
	} else {
		//GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS mode
		node = r.add(routeMethod, path, r.server.wrapRouterHandle(handle, false))
	}

	//if set auto-head, add head router
	//only enabled in hijack\GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS
	if r.server.ServerConfig.EnabledAutoHEAD {
		if routeMethod == RouteMethod_HiJack {
			r.add(RouteMethod_HEAD, path, r.server.wrapRouterHandle(handle, true))
		} else if routeMethod != RouteMethod_Any {
			r.add(RouteMethod_HEAD, path, r.server.wrapRouterHandle(handle, false))
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
	if !r.server.ServerConfig.EnabledListDir {
		root = &core.HideReaddirFS{root}
	}
	fileServer := http.FileServer(root)
	node = r.add(RouteMethod_GET, path, r.server.wrapFileHandle(fileServer))
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
