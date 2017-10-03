package dotweb

// HttpModule global module in http server
// it will be no effect when websocket request or use offline mode
type HttpModule struct {
	Name string
	//响应请求时作为 HTTP 执行管线链中的第一个事件发生
	OnBeginRequest func(Context)
	//响应请求时作为 HTTP 执行管线链中的最后一个事件发生。
	OnEndRequest func(Context)
}

func getIgnoreFaviconModule() *HttpModule {
	return &HttpModule{
		Name: "IgnoreFavicon",
		OnBeginRequest: func(ctx Context) {
			if ctx.Request().Path() == "/favicon.ico" {
				ctx.End()
			}
		},
	}
}
