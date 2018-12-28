package dotweb

// HttpModule global module in http server
// it will be no effect when websocket request or use offline mode
type HttpModule struct {
	Name string
	// OnBeginRequest is the first event in the execution chain
	OnBeginRequest func(Context)
	// OnEndRequest is the last event in the execution chain
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
