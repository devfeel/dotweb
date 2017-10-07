package dotweb

import (
	"github.com/devfeel/dotweb/framework/convert"
	"github.com/devfeel/dotweb/logger"
	"reflect"
	"time"
)

type MiddlewareFunc func() Middleware

// Middleware middleware interface
type Middleware interface {
	Handle(ctx Context) error
	SetNext(m Middleware)
	Next(ctx Context) error
}

//middleware 基础类，应用可基于此实现完整Moddleware
type BaseMiddlware struct {
	next Middleware
}

func (bm *BaseMiddlware) SetNext(m Middleware) {
	bm.next = m
}

func (bm *BaseMiddlware) Next(ctx Context) error {
	return bm.next.Handle(ctx)
}

type xMiddleware struct {
	BaseMiddlware
	IsEnd bool
}

func (x *xMiddleware) Handle(ctx Context) error {
	len := len(ctx.RouterNode().Middlewares())
	if x.IsEnd {
		return ctx.Handler()(ctx)
	} else {
		if x.next == nil {
			if len <= 0 {
				return ctx.Handler()(ctx)
			} else {
				if reflect.TypeOf(ctx.RouterNode().Middlewares()[len-1]).String() != "*dotweb.xMiddleware" {
					ctx.RouterNode().Use(&xMiddleware{IsEnd: true})
				}
				return ctx.RouterNode().Middlewares()[0].Handle(ctx)
			}
		} else {
			return x.Next(ctx)
		}
	}

}

//请求日志中间件
type RequestLogMiddleware struct {
	BaseMiddlware
}

func (m *RequestLogMiddleware) Handle(ctx Context) error {
	m.Next(ctx)
	timetaken := int64(time.Now().Sub(ctx.(*HttpContext).startTime) / time.Millisecond)
	log := ctx.Request().Url() + " " + logContext(ctx, timetaken)
	logger.Logger().Debug(log, LogTarget_HttpRequest)
	return nil
}

//get default log string
func logContext(ctx Context, timetaken int64) string {
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
	log += convert.Int642String(timetaken)

	return log
}
