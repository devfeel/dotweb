package dotweb

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/devfeel/dotweb/router"
	"github.com/devfeel/dotweb/session"
)

type HttpContext struct {
	Request      *http.Request
	RouterParams router.Params
	Response     *Response
	WebSocket    *WebSocket
	HijackConn   *HijackConn
	IsWebSocket  bool
	IsHijack     bool
	isEnd        bool //表示当前处理流程是否需要终止
	HttpServer   *HttpServer
	SessionID    string
}

//set context process end
func (ctx *HttpContext) End() {
	ctx.isEnd = true
}

func (ctx *HttpContext) IsEnd() bool {
	return ctx.isEnd
}

//reset response attr
func (ctx *HttpContext) Reset(res *Response, r *http.Request, server *HttpServer, params router.Params) {
	ctx.Request = r
	ctx.Response = res
	ctx.RouterParams = params
	ctx.IsHijack = false
	ctx.IsWebSocket = false
	ctx.HttpServer = server
	ctx.isEnd = false
}

//get session state in current context
func (ctx *HttpContext) Session() (session *session.SessionState) {
	if ctx.HttpServer == nil {
		//return nil, errors.New("no effective http-server")
		panic("no effective http-server")
	}
	if !ctx.HttpServer.ServerConfig.EnabledSession {
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

func (ctx *HttpContext) Server() string {
	return ctx.QueryHeader(HeaderServer)
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

// write cookie for domain&name&liveseconds
//
// default path = "/"
// default domain = current domain
// default seconds = 0
func (ctx *HttpContext) WriteCookie(name, value string, seconds int) {
	cookie := http.Cookie{Name: name, Value: value, MaxAge: seconds}
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

// write string content to response
func (ctx *HttpContext) WriteString(content string) (int, error) {
	if ctx.IsHijack {
		return ctx.HijackConn.WriteString(content)
	} else {
		return ctx.Response.Write([]byte(content))
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
		return ctx.Response.Write(b)
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
