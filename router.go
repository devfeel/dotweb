package dotweb

import (
	"github.com/devfeel/dotweb/framework/log"
	"github.com/devfeel/dotweb/routers"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
	"sync"
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
)

type (
	// Router is the interface that wraps the router method.
	Router interface {
		ServeHTTP(w http.ResponseWriter, req *http.Request)
		ServeFiles(path string, root http.FileSystem)
		GET(path string, handle HttpHandle)
		HEAD(path string, handle HttpHandle)
		OPTIONS(path string, handle HttpHandle)
		POST(path string, handle HttpHandle)
		PUT(path string, handle HttpHandle)
		PATCH(path string, handle HttpHandle)
		DELETE(path string, handle HttpHandle)
		HiJack(path string, handle HttpHandle)
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
	HttpMethodMap["GET"] = "GET"
	HttpMethodMap["HEAD"] = "HEAD"
	HttpMethodMap["POST"] = "POST"
	HttpMethodMap["PUT"] = "PUT"
	HttpMethodMap["PATCH"] = "PATCH"
	HttpMethodMap["OPTIONS"] = "OPTIONS"
	HttpMethodMap["DELETE"] = "DELETE"
	HttpMethodMap["HiJack"] = "HiJack"
	HttpMethodMap["WebSocket"] = "WebSocket"

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

func (r *router) ServeFiles(path string, root http.FileSystem) {
	r.ServeFiles(path, root)
}

// GET is a shortcut for router.Handle("GET", path, handle)
func (r *router) GET(path string, handle HttpHandle) {
	r.RegisterRoute(RouteMethod_GET, path, handle)
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
// support GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS\HiJack\WebSocket
func (r *router) RegisterRoute(routeMethod string, path string, handle HttpHandle) {

	routeMethod = strings.ToUpper(routeMethod)

	if _, exists := HttpMethodMap[routeMethod]; !exists {
		logger.Log("Dotweb:Router:RegisterRoute failed [illegal method] ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Warn)
		return
	} else {
		logger.Log("Dotweb:Router:RegisterRoute success ["+routeMethod+"] ["+path+"]", LogTarget_HttpServer, LogLevel_Debug)
	}

	//hijack mode,use get and isHijack = true
	if routeMethod == RouteMethod_HiJack {
		r.router.Handle(RouteMethod_GET, path, r.server.wrapRouterHandle(handle, true))
		return
	}
	//websocket mode,use default httpserver
	if routeMethod == RouteMethod_WebSocket {
		http.Handle(path, websocket.Handler(r.server.wrapWebSocketHandle(handle)))
		return
	}

	//GET\POST\DELETE\PUT\HEAD\PATCH\OPTIONS mode
	r.router.Handle(routeMethod, path, r.server.wrapRouterHandle(handle, false))
	return
}

// ServerFile is a shortcut for router.ServeFiles(path, filepath)
// simple demo:server.ServerFile("/src/*filepath", "/var/www")
func (r *router) ServerFile(urlpath string, filepath string) {
	r.router.ServeFiles(urlpath, http.Dir(filepath))
}
