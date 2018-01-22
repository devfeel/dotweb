package dotweb

import (
	"github.com/devfeel/dotweb/logger"
)

type (
	Group interface {
		Use(m ...Middleware) Group
		Group(prefix string, m ...Middleware) Group
		DELETE(path string, h HttpHandle) RouterNode
		GET(path string, h HttpHandle) RouterNode
		HEAD(path string, h HttpHandle) RouterNode
		OPTIONS(path string, h HttpHandle) RouterNode
		PATCH(path string, h HttpHandle) RouterNode
		POST(path string, h HttpHandle) RouterNode
		PUT(path string, h HttpHandle) RouterNode
		RegisterRoute(method, path string, h HttpHandle) RouterNode
	}
	xGroup struct {
		prefix      string
		middlewares []Middleware
		allRouterExpress   map[string]struct{}
		server      *HttpServer
	}
)

func NewGroup(prefix string, server *HttpServer) Group {
	g := &xGroup{prefix: prefix, server: server, allRouterExpress:make(map[string]struct{})}
	server.groups = append(server.groups, g)
	logger.Logger().Debug("DotWeb:Group NewGroup ["+prefix+"]", LogTarget_HttpServer)
	return g
}

// Use implements `Router#Use()` for sub-routes within the Group.
func (g *xGroup) Use(m ...Middleware) Group {
	if len(m) <= 0 {
		return g
	}
	step := len(g.middlewares) - 1
	for i := range m {
		if m[i] != nil {
			if step >= 0 {
				g.middlewares[step].SetNext(m[i])
			}
			g.middlewares = append(g.middlewares, m[i])
			step++
		}
	}
	return g
}

// DELETE implements `Router#DELETE()` for sub-routes within the Group.
func (g *xGroup) DELETE(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_DELETE, path, h)
}

// GET implements `Router#GET()` for sub-routes within the Group.
func (g *xGroup) GET(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_GET, path, h)
}

// HEAD implements `Router#HEAD()` for sub-routes within the Group.
func (g *xGroup) HEAD(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_HEAD, path, h)
}

// OPTIONS implements `Router#OPTIONS()` for sub-routes within the Group.
func (g *xGroup) OPTIONS(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_OPTIONS, path, h)
}

// PATCH implements `Router#PATCH()` for sub-routes within the Group.
func (g *xGroup) PATCH(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_PATCH, path, h)
}

// POST implements `Router#POST()` for sub-routes within the Group.
func (g *xGroup) POST(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_POST, path, h)
}

// PUT implements `Router#PUT()` for sub-routes within the Group.
func (g *xGroup) PUT(path string, h HttpHandle) RouterNode {
	return g.add(RouteMethod_PUT, path, h)
}

// Group creates a new sub-group with prefix and optional sub-group-level middleware.
func (g *xGroup) Group(prefix string, m ...Middleware) Group {
	return NewGroup(g.prefix+prefix, g.server).Use(g.middlewares...).Use(m...)
}

func (g *xGroup) RegisterRoute(method, path string, handler HttpHandle) RouterNode {
	return g.add(method, path, handler)
}

func (g *xGroup) add(method, path string, handler HttpHandle) RouterNode {
	node := g.server.Router().RegisterRoute(method, g.prefix+path, handler)
	g.allRouterExpress[method+"_"+g.prefix+path] = struct{}{}
	node.Node().groupMiddlewares = g.middlewares
	return node
}
