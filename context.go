package dotweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"fmt"
	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/routers"
	"github.com/devfeel/dotweb/session"
)

const (
	defaultMemory   = 32 << 20 // 32 MB
	defaultHttpCode = http.StatusOK
)

type HttpContext struct {
	Request      *Request
	RouterNode   *RouterNode
	RouterParams routers.Params
	Response     *Response
	WebSocket    *WebSocket
	HijackConn   *HijackConn
	IsWebSocket  bool
	IsHijack     bool
	isEnd        bool //表示当前处理流程是否需要终止
	HttpServer   *HttpServer
	SessionID    string
	items        *core.ItemContext
	viewData     *core.ItemContext
	Features     *xFeatureTools
}

//reset response attr
func (ctx *HttpContext) Reset(res *Response, r *Request, server *HttpServer, node *RouterNode, params routers.Params) {
	ctx.Request = r
	ctx.Response = res
	ctx.RouterNode = node
	ctx.RouterParams = params
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.HttpServer = server
	ctx.items = nil
	ctx.isEnd = false
	ctx.Features = FeatureTools
}

//release all field
func (ctx *HttpContext) release() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.RouterNode = nil
	ctx.RouterParams = nil
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.HttpServer = nil
	ctx.isEnd = false
	ctx.items = nil
	ctx.viewData = nil
	ctx.SessionID = ""
}

//get application's global appcontext
//issue #3
func (ctx *HttpContext) AppContext() *core.ItemContext {
	if ctx.HttpServer != nil {
		return ctx.HttpServer.DotApp.AppContext
	} else {
		return core.NewItemContext()
	}
}

//get application's global cache
func (ctx *HttpContext) Cache() cache.Cache {
	return ctx.HttpServer.DotApp.Cache()
}

//get request's tem context
//lazy init when first use
func (ctx *HttpContext) Items() *core.ItemContext {
	if ctx.items == nil {
		ctx.items = core.NewItemContext()
	}
	return ctx.items
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
	if ctx.HttpServer == nil {
		//return nil, errors.New("no effective http-server")
		panic("no effective http-server")
	}
	if !ctx.HttpServer.SessionConfig.EnabledSession {
		//return nil, errors.New("http-server not enabled session")
		panic("http-server not enabled session")
	}
	state, _ = ctx.HttpServer.sessionManager.GetSessionState(ctx.SessionID)
	return state
}

//make current connection to hijack mode
func (ctx *HttpContext) Hijack() (*HijackConn, error) {
	hj, ok := ctx.Response.Writer().(http.Hijacker)
	if !ok {
		return nil, errors.New("The Web Server does not support Hijacking! ")
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		return nil, errors.New("Hijack error:" + err.Error())
	}
	ctx.HijackConn = &HijackConn{Conn: conn, ReadWriter: bufrw, header: "HTTP/1.1 200 OK\r\n"}
	ctx.IsHijack = true
	return ctx.HijackConn, nil
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
func (ctx *HttpContext) Redirect(code int, targetUrl string) {
	ctx.Response.Redirect(code, targetUrl)
}

/*
* 返回查询字符串map表示
 */
func (ctx *HttpContext) QueryStrings() url.Values {
	return ctx.Request.QueryStrings()
}

/*
* 获取原始查询字符串
 */
func (ctx *HttpContext) RawQuery() string {
	return ctx.Request.RawQuery()
}

/*
* 根据指定key获取对应value
 */
func (ctx *HttpContext) QueryString(key string) string {
	return ctx.Request.QueryString(key)
}

/*
* 根据指定key获取包括在post、put和get内的值
 */
func (ctx *HttpContext) FormValue(key string) string {
	return ctx.Request.FormValue(key)
}

func (ctx *HttpContext) FormFile(key string) (*UploadFile, error) {
	return ctx.Request.FormFile(key)
}

/*
* 获取包括post、put和get内的值
 */
func (ctx *HttpContext) FormValues() map[string][]string {
	return ctx.Request.FormValues()
}

/*
* 根据指定key获取包括在post、put内的值
 */
func (ctx *HttpContext) PostFormValue(key string) string {
	return ctx.Request.PostFormValue(key)
}

/*
* 根据指定key获取包括在post、put内的值
* Obsolete("use PostFormValue replace this")
 */
func (ctx *HttpContext) PostString(key string) string {
	return ctx.Request.PostFormValue(key)
}

/*
* 获取post提交的字节数组
 */
func (ctx *HttpContext) PostBody() []byte {
	return ctx.Request.PostBody()
}

/*
* 支持Json、Xml、Form提交的属性绑定
 */
func (ctx *HttpContext) Bind(i interface{}) error {
	return ctx.HttpServer.Binder().Bind(i, ctx)
}

func (ctx *HttpContext) QueryHeader(key string) string {
	return ctx.Request.Header.Get(key)
}

func (ctx *HttpContext) DelHeader(key string) {
	ctx.Response.Header().Del(key)
}

//set response header kv info
func (ctx *HttpContext) SetHeader(key, value string) {
	if ctx.IsHijack {
		ctx.HijackConn.SetHeader(key, value)
	} else {
		ctx.Response.Header().Set(key, value)
	}
}

func (ctx *HttpContext) Url() string {
	return ctx.Request.URL.String()
}

func (ctx *HttpContext) ContentType() string {
	return ctx.Request.Header.Get(HeaderContentType)
}

func (ctx *HttpContext) GetRouterName(key string) string {
	return ctx.RouterParams.ByName(key)
}

//RemoteAddr to an "IP" address
func (ctx *HttpContext) RemoteIP() string {
	return ctx.Request.RemoteIP()
}

func (ctx *HttpContext) SetContentType(contenttype string) {
	ctx.SetHeader(HeaderContentType, contenttype)
}

func (ctx *HttpContext) SetStatusCode(code int) error {
	return ctx.Response.WriteHeader(code)
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
	cookie := http.Cookie{Name: name, Value: url.QueryEscape(value), MaxAge: maxAge}
	cookie.Path = "/"
	ctx.SetCookie(cookie)
}

// write cookie with cookie-obj
func (ctx *HttpContext) SetCookie(cookie http.Cookie) {
	http.SetCookie(ctx.Response.Writer(), &cookie)
}

// remove cookie for path&name
func (ctx *HttpContext) RemoveCookie(name string) {
	cookie := http.Cookie{Name: name, MaxAge: -1}
	ctx.SetCookie(cookie)
}

// read cookie value for name
func (ctx *HttpContext) ReadCookieValue(name string) (string, error) {
	cookieobj, err := ctx.Request.Cookie(name)
	if err != nil {
		return "", err
	} else {
		return url.QueryUnescape(cookieobj.Value)
	}
}

// read cookie object for name
func (ctx *HttpContext) ReadCookie(name string) (*http.Cookie, error) {
	return ctx.Request.Cookie(name)
}

// write view content to response
func (ctx *HttpContext) View(name string) error {
	err := ctx.HttpServer.Renderer().Render(ctx.Response.Writer(), name, ctx.ViewData().GetCurrentMap(), ctx)
	if err != nil {
		panic(err.Error())
	}
	return err
}

// write code and content content to response
func (ctx *HttpContext) Write(code int, content []byte) (int, error) {
	if ctx.IsHijack {
		//TODO:hijack mode, status-code set default 200
		return ctx.HijackConn.WriteBlob(content)
	} else {
		return ctx.Response.Write(code, content)
	}
}

// write string content to response
func (ctx *HttpContext) WriteString(contents ...interface{}) (int, error) {
	content := fmt.Sprint(contents...)
	if ctx.IsHijack {
		return ctx.HijackConn.WriteString(content)
	} else {
		return ctx.Response.Write(defaultHttpCode, []byte(content))
	}
}

// write []byte content to response
func (ctx *HttpContext) WriteBlob(contentType string, b []byte) (int, error) {
	if contentType != "" {
		ctx.SetContentType(contentType)
	}
	if ctx.IsHijack {
		return ctx.HijackConn.WriteBlob(b)
	} else {
		return ctx.Response.Write(defaultHttpCode, b)
	}
}

// write json string to response
//
// auto convert interface{} to json string
func (ctx *HttpContext) WriteJson(i interface{}) (int, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return 0, err
	}
	return ctx.WriteJsonBlob(b)
}

// write json string as []byte to response
func (ctx *HttpContext) WriteJsonBlob(b []byte) (int, error) {
	return ctx.WriteBlob(MIMEApplicationJSONCharsetUTF8, b)
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
	ctx.SetContentType(MIMEApplicationJavaScriptCharsetUTF8)
	//特殊处理，如果为hijack，需要先行WriteBlob头部
	if ctx.IsHijack {
		if size, err = ctx.HijackConn.WriteBlob([]byte(ctx.HijackConn.header + "\r\n")); err != nil {
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
