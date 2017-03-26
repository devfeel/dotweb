package dotweb

import (
	"github.com/devfeel/dotweb/framework/file"
	"github.com/devfeel/dotweb/framework/log"
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

type (
	// Router is the interface that wraps the router method.
	Router interface {
		ServeHTTP(w http.ResponseWriter, req *http.Request)
		ServerFile(path string, fileRoot string)
		GET(path string, handle HttpHandle)
		HEAD(path string, handle HttpHandle)
		OPTIONS(path string, handle HttpHandle)
		POST(path string, handle HttpHandle)
		PUT(path string, handle HttpHandle)
		PATCH(path string, handle HttpHandle)
		DELETE(path string, handle HttpHandle)
		HiJack(path string, handle HttpHandle)
		Any(path string, handle HttpHandle)
		RegisterRoute(routeMethod string, path string, handle HttpHandle)
		RegisterHandler(name string, handler HttpHandle)
		GetHandler(name string) (HttpHandle, bool)
	}
	router struct {
		router       *routers.Router
		server       *HttpServer
		handlerMap   map[string]HttpHandle
		handlerMutex *sync.RWMutex
	}
)

var (
	HttpMethodMap map[string]string
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

}

func NewRouter(server *HttpServer) Router {
	r := new(router)
	r.router = routers.New()
	r.server = server
	r.handlerMap = make(map[string]HttpHandle)
	r.handlerMutex = new(sync.RWMutex)
	return r
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

//use router ServerHTTP
func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *router) GET(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_GET, path, handle)
}

// ANY is a shortcut for router.Handle("Any", path, handle)
// it support GET\HEAD\POST\PUT\PATCH\OPTIONS\DELETE
func (r *router) Any(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_Any, path, handle)
}

// HEAD is a shortcut for router.Handle("HEAD", path, handle)
func (r *router) HEAD(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_HEAD, path, handle)
}

// OPTIONS is a shortcut for router.Handle("OPTIONS", path, handle)
func (r *router) OPTIONS(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_OPTIONS, path, handle)
}

// POST is a shortcut for router.Handle("POST", path, handle)
func (r *router) POST(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_POST, path, handle)
}

// PUT is a shortcut for router.Handle("PUT", path, handle)
func (r *router) PUT(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_PUT, path, handle)
}

// PATCH is a shortcut for router.Handle("PATCH", path, handle)
func (r *router) PATCH(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_PATCH, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *router) DELETE(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_DELETE, path, handle)
}

// DELETE is a shortcut for router.Handle("DELETE", path, handle)
func (r *router) HiJack(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_GET, path, handle)
}

// shortcut for router.Handle(httpmethod, path, handle)
// support GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS\HiJack\WebSocket\ANY
func (r *router) RegisterRoute(routeMethod string, path string, handle HttpHandle) {

	routeMethod = strings.ToUpper(routeMethod)

	if _, exists := HttpMethodMap[routeMethod]; !exists {
		logger.Log("Dotweb:Router:RegisterRoute failed [illegal method] ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Warn)
		return
	} else {
		logger.Log("Dotweb:Router:RegisterRoute success ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Debug)
	}

	//websocket mode,use default httpserver
	if routeMethod == RouteMethod_WebSocket {
		http.Handle(path, websocket.Handler(r.server.wrapWebSocketHandle(handle)))
		return
	}

	//hijack mode,use get and isHijack = true
	if routeMethod == RouteMethod_HiJack {
		r.router.Handle(RouteMethod_GET, path, r.server.wrapRouterHandle(handle, true))
	} else if routeMethod == RouteMethod_Any {
		r.router.ANY(path, r.server.wrapRouterHandle(handle, false))
	} else {
		//GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS mode
		r.router.Handle(routeMethod, path, r.server.wrapRouterHandle(handle, false))
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
	return
}

// ServerFile is a shortcut for router.ServeFiles(path, filepath)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (r *router) ServerFile(path string, fileroot string) {
	if len(path) < 10 || path[len(path)-10:] != "/*filepath" {
		panic("path must end with /*filepath in path '" + path + "'")
	}
	var root http.FileSystem
	root = http.Dir(fileroot)
	if !r.server.ServerConfig.EnabledListDir {
		root = &file.HideReaddirFS{root}
	}
	fileServer := http.FileServer(root)
	r.router.Handle(RouteMethod_GET, path, r.server.wrapFileHandle(fileServer))
	return
}
