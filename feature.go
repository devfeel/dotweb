package dotweb

import (
	"compress/gzip"
	"github.com/devfeel/dotweb/feature"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type xFeatureTools struct{}

var FeatureTools *xFeatureTools

func init() {
	FeatureTools = new(xFeatureTools)
}

//set CROS config on HttpContext
func (f *xFeatureTools) SetCROSConfig(ctx *HttpContext, c *feature.CROSConfig) {
	ctx.Response().SetHeader(HeaderAccessControlAllowOrigin, c.AllowedOrigins)
	ctx.Response().SetHeader(HeaderAccessControlAllowMethods, c.AllowedMethods)
	ctx.Response().SetHeader(HeaderAccessControlAllowHeaders, c.AllowedHeaders)
	ctx.Response().SetHeader(HeaderAccessControlAllowCredentials, strconv.FormatBool(c.AllowCredentials))
	ctx.Response().SetHeader(HeaderP3P, c.AllowedP3P)
}

//set CROS config on HttpContext
func (f *xFeatureTools) SetSession(httpCtx *HttpContext) {
	sessionId, err := httpCtx.HttpServer().GetSessionManager().GetClientSessionID(httpCtx.Request().Request)
	if err == nil && sessionId != "" {
		httpCtx.sessionID = sessionId
	} else {
		httpCtx.sessionID = httpCtx.HttpServer().GetSessionManager().NewSessionID()
		cookie := &http.Cookie{
			Name:  httpCtx.HttpServer().sessionManager.StoreConfig().CookieName,
			Value: url.QueryEscape(httpCtx.SessionID()),
			Path:  "/",
		}
		httpCtx.SetCookie(cookie)
	}
}

func (f *xFeatureTools) SetGzip(httpCtx *HttpContext) {
	gw, err := gzip.NewWriterLevel(httpCtx.Response().Writer(), DefaultGzipLevel)
	if err != nil {
		panic("use gzip error -> " + err.Error())
	}
	grw := &gzipResponseWriter{Writer: gw, ResponseWriter: httpCtx.Response().Writer()}
	httpCtx.Response().reset(grw)
	httpCtx.Response().SetHeader(HeaderContentEncoding, gzipScheme)
}

// doFeatures do features...
func (f *xFeatureTools) InitFeatures(server *HttpServer, httpCtx *HttpContext) {

	//gzip
	if server.ServerConfig().EnabledGzip {
		FeatureTools.SetGzip(httpCtx)
	}

	//session
	//if exists client-sessionid, use it
	//if not exists client-sessionid, new one
	if server.SessionConfig().EnabledSession {
		FeatureTools.SetSession(httpCtx)
	}

	//处理 cros feature
	if server.Features.CROSConfig != nil {
		c := server.Features.CROSConfig
		if c.EnabledCROS {
			FeatureTools.SetCROSConfig(httpCtx, c)
		}
	}

}

func (f *xFeatureTools) ReleaseFeatures(server *HttpServer, httpCtx *HttpContext) {
	if server.ServerConfig().EnabledGzip {
		var w io.Writer
		w = httpCtx.Response().Writer().(*gzipResponseWriter).Writer
		w.(*gzip.Writer).Close()
	}
}
