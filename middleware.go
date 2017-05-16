package dotweb

import (
	"github.com/devfeel/dotweb/framework/convert"
	"github.com/devfeel/dotweb/logger"
	"time"
)

type MiddlewareFunc func() Middleware

// Middleware middleware interface
type Middleware interface {
	Handle(ctx *HttpContext) error
	SetNext(m Middleware)
	Next(ctx *HttpContext) error
}

//middleware 基础类，应用可基于此实现完整Moddleware
type BaseMiddlware struct {
	next Middleware
}

func (bm *BaseMiddlware) SetNext(m Middleware) {
	bm.next = m
}

func (bm *BaseMiddlware) Next(ctx *HttpContext) error {
	return bm.next.Handle(ctx)
}

type xMiddleware struct {
	BaseMiddlware
	IsEnd bool
}

func (x *xMiddleware) Handle(ctx *HttpContext) error {
	if x.IsEnd {
		ctx.handle(ctx)
		return nil
	} else {
		if x.next == nil {
			if len(ctx.RouterNode.Middlewares()) <= 0 {
				ctx.handle(ctx)
			} else {
				ctx.RouterNode.Use(&xMiddleware{IsEnd: true})
				ctx.RouterNode.Middlewares()[0].Handle(ctx)
			}
			return nil
		} else {
			return x.Next(ctx)
		}
	}

}

//请求日志中间件
type RequestLogMiddleware struct {
	BaseMiddlware
}

func (m *RequestLogMiddleware) Handle(ctx *HttpContext) error {
	m.Next(ctx)
	timetaken := int64(time.Now().Sub(ctx.startTime) / time.Millisecond)
	log := ctx.Url() + " " + logContext(ctx, timetaken)
	logger.Logger().Log(log, LogTarget_HttpRequest, LogLevel_Debug)
	return nil
}

//get default log string
func logContext(ctx *HttpContext, timetaken int64) string {
	var reqbytelen, resbytelen, method, proto, status, userip string
	if ctx != nil {
		reqbytelen = convert.Int642String(ctx.Request.ContentLength)
		resbytelen = convert.Int642String(ctx.Response.Size)
		method = ctx.Request.Method
		proto = ctx.Request.Proto
		status = convert.Int2String(ctx.Response.Status)
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
