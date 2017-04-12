package dotweb

import (
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/feature"
	"github.com/devfeel/dotweb/logger"
	"github.com/devfeel/dotweb/routers"
	"golang.org/x/net/websocket"
	"net/http"
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
	featuresMap   map[interface{}]*feature.Feature
	lock_feature  *sync.RWMutex
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
	HttpMethodMap["HiJack"] = RouteMethod_HiJack
	HttpMethodMap["WebSocket"] = RouteMethod_WebSocket

	featuresMap = make(map[interface{}]*feature.Feature)
	lock_feature = new(sync.RWMutex)
}

type (
	// Router is the interface that wraps the router method.
	Router interface {
		ServeHTTP(w http.ResponseWriter, req *http.Request)
		ServerFile(path string, fileRoot string) *RouterNode
		GET(path string, handle HttpHandle) *RouterNode
		HEAD(path string, handle HttpHandle) *RouterNode
		OPTIONS(path string, handle HttpHandle) *RouterNode
		POST(path string, handle HttpHandle) *RouterNode
		PUT(path string, handle HttpHandle) *RouterNode
		PATCH(path string, handle HttpHandle) *RouterNode
		DELETE(path string, handle HttpHandle) *RouterNode
		HiJack(path string, handle HttpHandle)
		Any(path string, handle HttpHandle)
		RegisterRoute(routeMethod string, path string, handle HttpHandle) *RouterNode
		RegisterHandler(name string, handler HttpHandle)
		GetHandler(name string) (HttpHandle, bool)
		MatchPath(ctx *HttpContext, routePath string) bool
	}
	xRouter struct {
		router       *routers.Router
		server       *HttpServer
		handlerMap   map[string]HttpHandle
		handlerMutex *sync.RWMutex
	}

	RouterNode struct {
		Node   interface{}
		Method string
	}
)

func NewRouterNode(n interface{}, method string) *RouterNode {
	return &RouterNode{
		Node:   n,
		Method: method,
	}
}

func (n *RouterNode) SetEnabledCROS() *feature.CROSConfig {
	var f *feature.Feature
	var isok bool
	lock_feature.RLock()
	f, isok = featuresMap[n.Node]
	lock_feature.RUnlock()
	if !isok {
		f = feature.NewFeature()
		lock_feature.Lock()
		featuresMap[n.Node] = f
		lock_feature.Unlock()
	}
	f.CROSConfig.EnabledCROS = true
	f.CROSConfig.UseDefault()
	//special set method use current router's http method
	f.CROSConfig.SetMethod(n.Method)
	return f.CROSConfig
}

//do features...
func (n *RouterNode) doFeatures(ctx *HttpContext) *HttpContext {
	//处理 cros feature
	lock_feature.RLock()
	f, isok := featuresMap[n.Node]
	lock_feature.RUnlock()
	if isok && f != nil {
		c := f.CROSConfig
		if c.EnabledCROS {
			FeatureTools.SetCROSConfig(ctx, c)
		}
	}
	return ctx
}

func NewRouter(server *HttpServer) Router {
	r := new(xRouter)
	r.router = routers.New()
	r.server = server
	r.handlerMap = make(map[string]HttpHandle)
	r.handlerMutex = new(sync.RWMutex)
	return r
}

func (r *xRouter) RegisterHandler(name string, handler HttpHandle) {
	r.handlerMutex.Lock()
	r.handlerMap[name] = handler
	r.handlerMutex.Unlock()
}

func (r *xRouter) GetHandler(name string) (HttpHandle, bool) {
	r.handlerMutex.RLock()
	v, exists := r.handlerMap[name]
	r.handlerMutex.RUnlock()
	return v, exists
}

//use router ServerHTTP
func (r *xRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *xRouter) MatchPath(ctx *HttpContext, routePath string) bool {
	return r.router.MatchPath(ctx.Request.Method, ctx.RouterNode.Node.(*routers.Node), routePath)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *xRouter) GET(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_GET, path, handle)
}

// ANY is a shortcut for router.Handle("Any", path, handle)
// it support GET\HEAD\POST\PUT\PATCH\OPTIONS\DELETE
func (r *xRouter) Any(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_Any, path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *xRouter) HEAD(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_HEAD, path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *xRouter) OPTIONS(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_OPTIONS, path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *xRouter) POST(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_POST, path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *xRouter) PUT(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_PUT, path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *xRouter) PATCH(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_PATCH, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *xRouter) DELETE(path string, handle HttpHandle) *RouterNode {
	return r.RegisterRoute(RouteMethod_DELETE, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *xRouter) HiJack(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_GET, path, handle)
}

// shortcut for router.Handle(httpmethod, path, handle)
// support GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS\HiJack\WebSocket\ANY
func (r *xRouter) RegisterRoute(routeMethod string, path string, handle HttpHandle) *RouterNode {

	routeMethod = strings.ToUpper(routeMethod)
	rn := &RouterNode{Node: new(routers.Node), Method: routeMethod}
	if _, exists := HttpMethodMap[routeMethod]; !exists {
		logger.Logger().Log("Dotweb:Router:RegisterRoute failed [illegal method] ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Warn)
		return rn
	} else {
		logger.Logger().Log("Dotweb:Router:RegisterRoute success ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Debug)
	}

	//websocket mode,use default httpserver
	if routeMethod == RouteMethod_WebSocket {
		http.Handle(path, websocket.Handler(r.server.wrapWebSocketHandle(handle)))
		return rn
	}

	//hijack mode,use get and isHijack = true
	if routeMethod == RouteMethod_HiJack {
		r.router.Handle(RouteMethod_GET, path, r.server.wrapRouterHandle(handle, true))
	} else if routeMethod == RouteMethod_Any {
		r.router.ANY(path, r.server.wrapRouterHandle(handle, false))
	} else {
		//GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS mode
		rn.Node = r.router.Handle(routeMethod, path, r.server.wrapRouterHandle(handle, false))
	}

	//if set auto-head, add head router
	//only enabled in hijack\GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS
	if r.server.ServerConfig.EnabledAutoHEAD {
		if routeMethod == RouteMethod_HiJack {
			r.router.Handle(RouteMethod_HEAD, path, r.server.wrapRouterHandle(handle, true))
		} else if routeMethod != RouteMethod_Any {
			r.router.Handle(RouteMethod_HEAD, path, r.server.wrapRouterHandle(handle, false))
		}
	}
	return rn
}

// ServerFile is a shortcut for router.ServeFiles(path, filepath)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (r *xRouter) ServerFile(path string, fileroot string) *RouterNode {
	rn := &RouterNode{Node: new(routers.Node)}
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}
	var root http.FileSystem
	root = http.Dir(fileroot)
	if !r.server.ServerConfig.EnabledListDir {
		root = &core.HideReaddirFS{root}
	}
	fileServer := http.FileServer(root)
	rn.Node = r.router.Handle(RouteMethod_GET, path, r.server.wrapFileHandle(fileServer))
	return rn
}
