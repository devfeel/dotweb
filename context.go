package dotweb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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
	Request      *http.Request
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
}

//reset response attr
func (ctx *HttpContext) Reset(res *Response, r *http.Request, server *HttpServer, params routers.Params) {
	ctx.Request = r
	ctx.Response = res
	ctx.RouterParams = params
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.HttpServer = server
	ctx.isEnd = false
}

//release all field
func (ctx *HttpContext) release() {
	ctx.Request = nil
	ctx.Response = nil
	ctx.RouterParams = nil
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.HttpServer = nil
	ctx.isEnd = false
	ctx.items = nil
	ctx.viewData = nil
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
func (ctx *HttpContext) Session() (session *session.SessionState) {
	if ctx.HttpServer == nil {
		//return nil, errors.New("no effective http-server")
		panic("no effective http-server")
	}
	if !ctx.HttpServer.SessionConfig.EnabledSession {
		//return nil, errors.New("http-server not enabled session")
		panic("http-server not enabled session")
	}
	state, err := ctx.HttpServer.sessionManager.GetSessionState(ctx.SessionID)
	if err != nil {
		panic(err.Error())
	}
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
	return ctx.Request.URL.Query()
}

/*
* 获取原始查询字符串
 */
func (ctx *HttpContext) RawQuery() string {
	return ctx.Request.URL.RawQuery
}

/*
* 根据指定key获取对应value
 */
func (ctx *HttpContext) QueryString(key string) string {
	return ctx.Request.URL.Query().Get(key)
}

/*
* 根据指定key获取包括在post、put和get内的值
 */
func (ctx *HttpContext) FormValue(key string) string {
	return ctx.Request.FormValue(key)
}

func (ctx *HttpContext) FormFile(key string) (*UploadFile, error) {
	file, header, err := ctx.Request.FormFile(key)
	if err != nil {
		return nil, err
	} else {
		return NewUploadFile(file, header), nil
	}
}

/*
* 获取包括post、put和get内的值
 */
func (ctx *HttpContext) FormValues() map[string][]string {
	ctx.parseForm()
	return map[string][]string(ctx.Request.Form)
}

func (ctx *HttpContext) parseForm() error {
	if strings.HasPrefix(ctx.QueryHeader(HeaderContentType), MIMEMultipartForm) {
		if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil {
			return err
		}
	} else {
		if err := ctx.Request.ParseForm(); err != nil {
			return err
		}
	}
	return nil
}

/*
* 根据指定key获取包括在post、put内的值
 */
func (ctx *HttpContext) PostFormValue(key string) string {
	return ctx.Request.PostFormValue(key)
}

/*
* 根据指定key获取包括在post、put内的值
 */
func (ctx *HttpContext) PostString(key string) string {
	return ctx.Request.PostFormValue(key)
}

/*
* 获取post提交的字节数组
 */
func (ctx *HttpContext) PostBody() []byte {
	bts, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return []byte{}
	} else {
		return bts
	}
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

// IsAJAX returns if it is a ajax request
func (ctx *HttpContext) IsAJAX() bool {
	return ctx.Request.Header.Get(HeaderXRequestedWith) == "XMLHttpRequest"
}

func (ctx *HttpContext) Proto() string {
	return ctx.Request.Proto
}

func (ctx *HttpContext) Method() string {
	return ctx.Request.Method
}

//RemoteAddr to an "IP" address
func (ctx *HttpContext) RemoteIP() string {
	fullIp := ctx.Request.RemoteAddr
	s := strings.Split(fullIp, ":")
	if len(s) > 1 {
		return s[0]
	} else {
		return fullIp
	}
}

//RemoteAddr to an "IP:port" address
func (ctx *HttpContext) FullRemoteIP() string {
	fullIp := ctx.Request.RemoteAddr
	return fullIp
}

// Referer returns request referer.
//
// The referer is valid until returning from RequestHandler.
func (ctx *HttpContext) Referer() string {
	return ctx.Request.Referer()
}

// UserAgent returns User-Agent header value from the request.
func (ctx *HttpContext) UserAgent() string {
	return ctx.Request.UserAgent()
}

// Path returns requested path.
//
// The path is valid until returning from RequestHandler.
func (ctx *HttpContext) Path() string {
	return ctx.Request.URL.Path
}

// Host returns requested host.
//
// The host is valid until returning from RequestHandler.
func (ctx *HttpContext) Host() string {
	return ctx.Request.Host
}

func (ctx *HttpContext) SetContentType(contenttype string) {
	ctx.SetHeader(HeaderContentType, contenttype)
}

func (ctx *HttpContext) SetStatusCode(code int) error {
	return ctx.Response.WriteHeader(code)
}

// write cookie for domain & name & maxAge
//
// default path = "/"
// default domain = current domain
// default maxAge = 0 //seconds
// seconds=0 means no 'Max-Age' attribute specified.
// seconds<0 means delete cookie now, equivalently 'Max-Age: 0'
// seconds>0 means Max-Age attribute present and given in seconds
func (ctx *HttpContext) WriteCookie(name, value string, maxAge int) {
	cookie := http.Cookie{Name: name, Value: value, MaxAge: maxAge}
	http.SetCookie(ctx.Response.Writer(), &cookie)
}

// write cookie with cookie-obj
func (ctx *HttpContext) WriteCookieObj(cookie http.Cookie) {
	http.SetCookie(ctx.Response.Writer(), &cookie)
}

// remove cookie for path&name
func (ctx *HttpContext) RemoveCookie(name string) {
	cookie := http.Cookie{Name: name, MaxAge: -1}
	http.SetCookie(ctx.Response.Writer(), &cookie)
}

// read cookie value for name
func (ctx *HttpContext) ReadCookie(name string) (string, error) {
	cookieobj, err := ctx.Request.Cookie(name)
	if err != nil {
		return "", err
	} else {
		return cookieobj.Value, nil
	}
}

// read cookie object for name
func (ctx *HttpContext) ReadCookieObj(name string) (*http.Cookie, error) {
	return ctx.Request.Cookie(name)
}

// write view content to response
func (ctx *HttpContext) View(name string) error {
	err := ctx.HttpServer.Renderer().Render(ctx.Response.Writer(), name, ctx.ViewData().GetCurrentMap(), ctx)
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
