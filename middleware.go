package dotweb

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
}

func (x *xMiddleware) Handle(ctx *HttpContext) error {
	if x.next == nil {
		ctx.handle(ctx)
		return nil
	} else {
		return x.Next(ctx)
	}
}
