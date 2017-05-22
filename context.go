package dotweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"fmt"
	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/session"
	"time"
)

const (
	defaultMemory   = 32 << 20 // 32 MB
	defaultHttpCode = http.StatusOK
)

type (
	Context interface {
		HttpServer() *HttpServer
		Response() *Response
		Request() *Request
		WebSocket() *WebSocket
		HijackConn() *HijackConn
		RouterNode() RouterNode
		RouterParams() Params
		Handler() HttpHandle
		AppContext() *core.ItemContext
		Cache() cache.Cache
		Items() *core.ItemContext
		AppSetConfig() *core.ItemContext
		ViewData() *core.ItemContext
		SessionID() string
		Session() (state *session.SessionState)
		Hijack() (*HijackConn, error)
		End()
		IsEnd() bool
		Redirect(code int, targetUrl string) error
		QueryString(key string) string
		FormValue(key string) string
		PostFormValue(key string) string
		Bind(i interface{}) error
		GetRouterName(key string) string
		RemoteIP() string
		SetCookieValue(name, value string, maxAge int)
		SetCookie(cookie *http.Cookie)
		RemoveCookie(name string)
		ReadCookieValue(name string) (string, error)
		ReadCookie(name string) (*http.Cookie, error)
		View(name string) error
		ViewC(code int, name string) error
		Write(code int, content []byte) (int, error)
		WriteString(contents ...interface{}) (int, error)
		WriteStringC(code int, contents ...interface{}) (int, error)
		WriteBlob(contentType string, b []byte) (int, error)
		WriteBlobC(code int, contentType string, b []byte) (int, error)
		WriteJson(i interface{}) (int, error)
		WriteJsonC(code int, i interface{}) (int, error)
		WriteJsonBlob(b []byte) (int, error)
		WriteJsonBlobC(code int, b []byte) (int, error)
		WriteJsonp(callback string, i interface{}) (int, error)
		WriteJsonpBlob(callback string, b []byte) (size int, err error)
	}

	HttpContext struct {
		request      *Request
		routerNode   RouterNode
		routerParams Params
		response     *Response
		webSocket    *WebSocket
		hijackConn   *HijackConn
		IsWebSocket  bool
		IsHijack     bool
		isEnd        bool //表示当前处理流程是否需要终止
		httpServer   *HttpServer
		sessionID    string
		items        *core.ItemContext
		viewData     *core.ItemContext
		features     *xFeatureTools
		handler      HttpHandle
		startTime    time.Time
	}
)

//reset response attr
func (ctx *HttpContext) reset(res *Response, r *Request, server *HttpServer, node RouterNode, params Params, handler HttpHandle) {
	ctx.request = r
	ctx.response = res
	ctx.routerNode = node
	ctx.routerParams = params
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.httpServer = server
	ctx.items = nil
	ctx.isEnd = false
	ctx.features = FeatureTools
	ctx.handler = handler
	ctx.startTime = time.Now()
}

//release all field
func (ctx *HttpContext) release() {
	ctx.request = nil
	ctx.response = nil
	ctx.routerNode = nil
	ctx.routerParams = nil
	ctx.webSocket = nil
	ctx.hijackConn = nil
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.httpServer = nil
	ctx.isEnd = false
	ctx.features = nil
	ctx.items = nil
	ctx.viewData = nil
	ctx.sessionID = ""
	ctx.handler = nil
	ctx.startTime = time.Time{}
}

func (ctx *HttpContext) HttpServer() *HttpServer {
	return ctx.httpServer
}

func (ctx *HttpContext) Response() *Response {
	return ctx.response
}

func (ctx *HttpContext) Request() *Request {
	return ctx.request
}

func (ctx *HttpContext) RouterNode() RouterNode {
	return ctx.routerNode
}

func (ctx *HttpContext) WebSocket() *WebSocket {
	return ctx.webSocket
}

func (ctx *HttpContext) HijackConn() *HijackConn {
	return ctx.hijackConn
}

func (ctx *HttpContext) RouterParams() Params {
	return ctx.routerParams
}

func (ctx *HttpContext) Handler() HttpHandle {
	return ctx.handler
}

func (ctx *HttpContext) SessionID() string {
	return ctx.sessionID
}

func (ctx *HttpContext) Features() *xFeatureTools {
	return ctx.features
}

//get application's global appcontext
//issue #3
func (ctx *HttpContext) AppContext() *core.ItemContext {
	if ctx.HttpServer != nil {
		return ctx.httpServer.DotApp.AppContext
	} else {
		return core.NewItemContext()
	}
}

//get application's global cache
func (ctx *HttpContext) Cache() cache.Cache {
	return ctx.httpServer.DotApp.Cache()
}

//get request's tem context
//lazy init when first use
func (ctx *HttpContext) Items() *core.ItemContext {
	if ctx.items == nil {
		ctx.items = core.NewItemContext()
	}
	return ctx.items
}

//get appset from config file
func (ctx *HttpContext) AppSetConfig() *core.ItemContext {
	return ctx.HttpServer().DotApp.Config.AppSetConfig
}

//get view data context
//lazy init when first use
func (ctx *HttpContext) ViewData() *core.ItemContext {
	if ctx.viewData == nil {
		ctx.viewData = core.NewItemContext()
	}
	return ctx.viewData
}

//get session state in current context
func (ctx *HttpContext) Session() (state *session.SessionState) {
	if ctx.httpServer == nil {
		//return nil, errors.New("no effective http-server")
		panic("no effective http-server")
	}
	if !ctx.httpServer.SessionConfig.EnabledSession {
		//return nil, errors.New("http-server not enabled session")
		panic("http-server not enabled session")
	}
	state, _ = ctx.httpServer.sessionManager.GetSessionState(ctx.sessionID)
	return state
}

//make current connection to hijack mode
func (ctx *HttpContext) Hijack() (*HijackConn, error) {
	hj, ok := ctx.response.Writer().(http.Hijacker)
	if !ok {
		return nil, errors.New("The Web Server does not support Hijacking! ")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, errors.New("Hijack error:" + err.Error())
	}
	ctx.hijackConn = &HijackConn{Conn: conn, ReadWriter: bufrw, header: "HTTP/1.1 200 OK\r\n"}
	ctx.IsHijack = true
	return ctx.hijackConn, nil
}

//set context user handler process end
//if set HttpContext.End,ignore user handler, but exec all http module  - fixed issue #5
func (ctx *HttpContext) End() {
	ctx.isEnd = true
}

func (ctx *HttpContext) IsEnd() bool {
	return ctx.isEnd
}

//redirect replies to the request with a redirect to url and with httpcode
//default you can use http.StatusFound
func (ctx *HttpContext) Redirect(code int, targetUrl string) error {
	return ctx.response.Redirect(code, targetUrl)
}

/*
* 根据指定key获取对应value
 */
func (ctx *HttpContext) QueryString(key string) string {
	return ctx.request.QueryString(key)
}

/*
* 根据指定key获取包括在post、put和get内的值
 */
func (ctx *HttpContext) FormValue(key string) string {
	return ctx.request.FormValue(key)
}

/*
* 根据指定key获取包括在post、put内的值
 */
func (ctx *HttpContext) PostFormValue(key string) string {
	return ctx.request.PostFormValue(key)
}

/*
* 支持Json、Xml、Form提交的属性绑定
 */
func (ctx *HttpContext) Bind(i interface{}) error {
	return ctx.httpServer.Binder().Bind(i, ctx)
}

func (ctx *HttpContext) GetRouterName(key string) string {
	return ctx.routerParams.ByName(key)
}

//RemoteAddr to an "IP" address
func (ctx *HttpContext) RemoteIP() string {
	return ctx.request.RemoteIP()
}

// write cookie for name & value & maxAge
//
// default path = "/"
// default domain = current domain
// default maxAge = 0 //seconds
// seconds=0 means no 'Max-Age' attribute specified.
// seconds<0 means delete cookie now, equivalently 'Max-Age: 0'
// seconds>0 means Max-Age attribute present and given in seconds
func (ctx *HttpContext) SetCookieValue(name, value string, maxAge int) {
	cookie := &http.Cookie{Name: name, Value: url.QueryEscape(value), MaxAge: maxAge}
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}

// write cookie with cookie-obj
func (ctx *HttpContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.response.Writer(), cookie)
}

// remove cookie for path&name
func (ctx *HttpContext) RemoveCookie(name string) {
	cookie := &http.Cookie{Name: name, MaxAge: -1}
	ctx.SetCookie(cookie)
}

// read cookie value for name
func (ctx *HttpContext) ReadCookieValue(name string) (string, error) {
	cookieobj, err := ctx.request.Cookie(name)
	if err != nil {
		return "", err
	} else {
		return url.QueryUnescape(cookieobj.Value)
	}
}

// read cookie object for name
func (ctx *HttpContext) ReadCookie(name string) (*http.Cookie, error) {
	return ctx.request.Cookie(name)
}

// write view content to response
func (ctx *HttpContext) View(name string) error {
	return ctx.ViewC(defaultHttpCode, name)
}

// write (httpCode, view content) to response
func (ctx *HttpContext) ViewC(code int, name string) error {
	ctx.response.SetStatusCode(code)
	err := ctx.httpServer.Renderer().Render(ctx.response.Writer(), name, ctx.ViewData().GetCurrentMap(), ctx)
	return err
}

// write code and content content to response
func (ctx *HttpContext) Write(code int, content []byte) (int, error) {
	if ctx.IsHijack {
		//TODO:hijack mode, status-code set default 200
		return ctx.hijackConn.WriteBlob(content)
	} else {
		return ctx.response.Write(code, content)
	}
}

// write string content to response
func (ctx *HttpContext) WriteString(contents ...interface{}) (int, error) {
	return ctx.WriteStringC(defaultHttpCode, contents...)
}

// write (httpCode, string) to response
func (ctx *HttpContext) WriteStringC(code int, contents ...interface{}) (int, error) {
	content := fmt.Sprint(contents...)
	if ctx.IsHijack {
		return ctx.hijackConn.WriteString(content)
	} else {
		return ctx.response.Write(code, []byte(content))
	}
}

// write []byte content to response
func (ctx *HttpContext) WriteBlob(contentType string, b []byte) (int, error) {
	return ctx.WriteBlobC(defaultHttpCode, contentType, b)
}

// write (httpCode, []byte) to response
func (ctx *HttpContext) WriteBlobC(code int, contentType string, b []byte) (int, error) {
	if contentType != "" {
		ctx.response.SetContentType(contentType)
	}
	if ctx.IsHijack {
		return ctx.hijackConn.WriteBlob(b)
	} else {
		return ctx.response.Write(code, b)
	}
}

// write (httpCode, json string) to response
// auto convert interface{} to json string
func (ctx *HttpContext) WriteJson(i interface{}) (int, error) {
	return ctx.WriteJsonC(defaultHttpCode, i)
}

// write (httpCode, json string) to response
// auto convert interface{} to json string
func (ctx *HttpContext) WriteJsonC(code int, i interface{}) (int, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}
	return ctx.WriteJsonBlobC(code, b)
}

// write json []byte to response
func (ctx *HttpContext) WriteJsonBlob(b []byte) (int, error) {
	return ctx.WriteJsonBlobC(defaultHttpCode, b)
}

// write (httpCode, json []byte) to response
func (ctx *HttpContext) WriteJsonBlobC(code int, b []byte) (int, error) {
	return ctx.WriteBlobC(code, MIMEApplicationJSONCharsetUTF8, b)
}

// write jsonp string to response
func (ctx *HttpContext) WriteJsonp(callback string, i interface{}) (int, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}
	return ctx.WriteJsonpBlob(callback, b)
}

// write jsonp string as []byte to response
func (ctx *HttpContext) WriteJsonpBlob(callback string, b []byte) (size int, err error) {
	ctx.response.SetContentType(MIMEApplicationJavaScriptCharsetUTF8)
	//特殊处理，如果为hijack，需要先行WriteBlob头部
	if ctx.IsHijack {
		if size, err = ctx.hijackConn.WriteBlob([]byte(ctx.hijackConn.header + "\r\n")); err != nil {
			return
		}
	}
	if size, err = ctx.WriteBlob("", []byte(callback+"(")); err != nil {
		return
	}
	if size, err = ctx.WriteBlob("", b); err != nil {
		return
	}
	size, err = ctx.WriteBlob("", []byte(");"))
	return
}
