package dotweb

import (
	"github.com/devfeel/dotweb/framework/convert"
	"github.com/devfeel/dotweb/logger"
	"time"
)

const (
	middleware_App    = "app"
	middleware_Group  = "group"
	middleware_Router = "router"
)

type MiddlewareFunc func() Middleware

//middleware执行优先级：
//优先级1：app级别middleware
//优先级2：group级别middleware
//优先级3：router级别middleware

// Middleware middleware interface
type Middleware interface {
	Handle(ctx Context) error
	SetNext(m Middleware)
	Next(ctx Context) error
	Exclude(routers ...string)
	HasExclude() bool
	ExistsExcludeRouter(router string) bool
}

//middleware 基础类，应用可基于此实现完整Moddleware
type BaseMiddlware struct {
	next           Middleware
	excludeRouters map[string]struct{}
}

func (bm *BaseMiddlware) SetNext(m Middleware) {
	bm.next = m
}

func (bm *BaseMiddlware) Next(ctx Context) error {
	httpCtx := ctx.(*HttpContext)
	if httpCtx.middlewareStep == "" {
		httpCtx.middlewareStep = middleware_App
	}
	if bm.next == nil {
		if httpCtx.middlewareStep == middleware_App {
			httpCtx.middlewareStep = middleware_Group
			if len(httpCtx.RouterNode().GroupMiddlewares()) > 0 {
				return httpCtx.RouterNode().GroupMiddlewares()[0].Handle(ctx)
			}
		}
		if httpCtx.middlewareStep == middleware_Group {
			httpCtx.middlewareStep = middleware_Router
			if len(httpCtx.RouterNode().Middlewares()) > 0 {
				return httpCtx.RouterNode().Middlewares()[0].Handle(ctx)
			}
		}

		if httpCtx.middlewareStep == middleware_Router {
			return httpCtx.Handler()(ctx)
		}
	} else {
		//check exclude config
		if ctx.RouterNode().Node().hasExcludeMiddleware && bm.next.HasExclude() {
			if bm.next.ExistsExcludeRouter(ctx.RouterNode().Node().fullPath) {
				return bm.next.Next(ctx)
			}
		}
		return bm.next.Handle(ctx)
	}
	return nil
}

// Exclude Exclude this middleware with router
func (bm *BaseMiddlware) Exclude(routers ...string) {
	if bm.excludeRouters == nil {
		bm.excludeRouters = make(map[string]struct{})
	}
	for _, v := range routers {
		bm.excludeRouters[v] = struct{}{}
	}
}

// HasExclude check has set exclude router
func (bm *BaseMiddlware) HasExclude() bool {
	if bm.excludeRouters == nil {
		return false
	}
	if len(bm.excludeRouters) > 0 {
		return true
	} else {
		return false
	}
}

// ExistsExcludeRouter check is exists router in exclude map
func (bm *BaseMiddlware) ExistsExcludeRouter(router string) bool {
	if bm.excludeRouters == nil {
		return false
	}
	_, exists := bm.excludeRouters[router]
	return exists
}

type xMiddleware struct {
	BaseMiddlware
	IsEnd bool
}

func (x *xMiddleware) Handle(ctx Context) error {
	httpCtx := ctx.(*HttpContext)
	if httpCtx.middlewareStep == "" {
		httpCtx.middlewareStep = middleware_App
	}
	if x.IsEnd {
		return httpCtx.Handler()(ctx)
	}
	return x.Next(ctx)
}

//请求日志中间件
type RequestLogMiddleware struct {
	BaseMiddlware
}

func (m *RequestLogMiddleware) Handle(ctx Context) error {
	var timeDuration time.Duration
	var timeTaken uint64
	var err error
	m.Next(ctx)
	if ctx.Items().Exists(ItemKeyHandleDuration){
		timeDuration, err = time.ParseDuration(ctx.Items().GetString(ItemKeyHandleDuration))
		if err != nil{
			timeTaken = 0
		}else{
			timeTaken = uint64(timeDuration/time.Millisecond)
		}
	}else{
		var begin time.Time
		beginVal, exists := ctx.Items().Get(ItemKeyHandleStartTime)
		if !exists{
			begin  = time.Now()
		}else{
			begin = beginVal.(time.Time)
		}
		timeTaken = uint64(time.Now().Sub(begin) / time.Millisecond)
	}
	log := ctx.Request().Url() + " " + logContext(ctx, timeTaken)
	logger.Logger().Debug(log, LogTarget_HttpRequest)
	return nil
}

//get default log string
func logContext(ctx Context, timetaken uint64) string {
	var reqbytelen, resbytelen, method, proto, status, userip string
	if ctx != nil {
		reqbytelen = convert.Int642String(ctx.Request().ContentLength)
		resbytelen = convert.Int642String(ctx.Response().Size)
		method = ctx.Request().Method
		proto = ctx.Request().Proto
		status = convert.Int2String(ctx.Response().Status)
		userip = ctx.RemoteIP()
	}

	log := method + " "
	log += userip + " "
	log += proto + " "
	log += status + " "
	log += reqbytelen + " "
	log += resbytelen + " "
	log += convert.UInt642String(timetaken)

	return log
}

// TimeoutHookMiddleware 超时钩子中间件
type TimeoutHookMiddleware struct {
	BaseMiddlware
	HookHandle StandardHandle
	TimeoutDuration time.Duration
}

func (m *TimeoutHookMiddleware) Handle(ctx Context) error {
	var begin time.Time
	if m.HookHandle != nil{
		beginVal, exists := ctx.Items().Get(ItemKeyHandleStartTime)
		if !exists{
			begin  = time.Now()
		}else{
			begin = beginVal.(time.Time)
		}
	}
	//Do next
	m.Next(ctx)
	if m.HookHandle != nil{
		realDuration := time.Now().Sub(begin)
		ctx.Items().Set(ItemKeyHandleDuration, realDuration)
		if realDuration > m.TimeoutDuration{
			m.HookHandle(ctx)
		}
	}
	return nil
}