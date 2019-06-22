package dotweb

import (
	"github.com/devfeel/dotweb/core"
	jsonutil "github.com/devfeel/dotweb/framework/json"
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"strings"
)

// initDotwebGroup init Dotweb route group which start with /dotweb/
func initDotwebGroup(server *HttpServer) {
	gInner := server.Group("/dotweb")
	gInner.GET("/debug/pprof/:key", showPProf)
	gInner.GET("/debug/freemem", freeMemory)
	gInner.GET("/state", showServerState)
	gInner.GET("/state/interval", showIntervalData)
	gInner.GET("/query/:key", showQuery)
	gInner.GET("/routers", showRouters)
}

// query pprof debug info
// key:heap goroutine threadcreate block
func showPProf(ctx Context) error {
	querykey := ctx.GetRouterName("key")
	runtime.GC()
	return pprof.Lookup(querykey).WriteTo(ctx.Response().Writer(), 1)
}

func freeMemory(ctx Context) error {
	debug.FreeOSMemory()
	return nil
}

func showIntervalData(ctx Context) error {
	type data struct {
		Time         string
		RequestCount uint64
		ErrorCount   uint64
	}
	queryKey := ctx.QueryString("querykey")

	d := new(data)
	d.Time = queryKey
	d.RequestCount = core.GlobalState.QueryIntervalRequestData(queryKey)
	d.ErrorCount = core.GlobalState.QueryIntervalErrorData(queryKey)
	return ctx.WriteJson(d)
}

// snow server status
func showServerState(ctx Context) error {
	return ctx.WriteHtml(core.GlobalState.ShowHtmlData(Version, ctx.HttpServer().DotApp.GlobalUniqueID()))
}

// query server information
func showQuery(ctx Context) error {
	querykey := ctx.GetRouterName("key")
	switch querykey {
	case "state":
		return ctx.WriteString(jsonutil.GetJsonString(core.GlobalState))
	case "":
		return ctx.WriteString("please input key")
	default:
		return ctx.WriteString("not support key => " + querykey)
	}
}

func showRouters(ctx Context) error {

	result := ""
	for k, _ := range ctx.HttpServer().router.GetAllRouterExpress() {
		result += strings.Split(k, routerExpressSplit)[1] + "\r\n"
	}
	return ctx.WriteString(result)
}
