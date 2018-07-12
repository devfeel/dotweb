package dotweb

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/devfeel/dotweb/cache"
	"github.com/devfeel/dotweb/core"
	"github.com/devfeel/dotweb/session"
)

const (
	defaultMemory   = 32 << 20 // 32 MB
	defaultHttpCode = http.StatusOK
	// ItemKeyHandleStartTime itemkey name for request handler start time
	ItemKeyHandleStartTime = "dotweb.HttpContext.StartTime"
	// ItemKeyHandleDuration itemkey name for request handler time duration
	ItemKeyHandleDuration = "dotweb.HttpContext.HandleDuration"
)

const (
	innerKeyAddView = "inner_AddView"
)

type (
	Context interface {
		Context() context.Context
		SetTimeoutContext(timeout time.Duration) context.Context
		WithContext(runCtx context.Context)
		HttpServer() *HttpServer
		Response() *Response
		Request() *Request
		WebSocket() *WebSocket
		HijackConn() *HijackConn
		RouterNode() RouterNode
		RouterParams() Params
		Handler() HttpHandle
		AppItems() core.ConcurrenceMap
		Cache() cache.Cache
		Items() core.ConcurrenceMap
		ConfigSet() core.ReadonlyMap
		ViewData() core.ConcurrenceMap
		SessionID() string
		Session() (state *session.SessionState)
		Hijack() (*HijackConn, error)
		IsHijack() bool
		IsWebSocket() bool
		End()
		IsEnd() bool
		Redirect(code int, targetUrl string) error
		QueryString(key string) string
		QueryInt(key string) int
		QueryInt64(key string) int64
		FormValue(key string) string
		PostFormValue(key string) string
		File(file string) (err error)
		Attachment(file string, name string) error
		Inline(file string, name string) error
		Bind(i interface{}) error
		BindJsonBody(i interface{}) error
		// Validate validates provided `i`. It is usually called after `Context#Bind()`.
		Validate(i interface{}) error
		GetRouterName(key string) string
		RemoteIP() string
		SetCookieValue(name, value string, maxAge int)
		SetCookie(cookie *http.Cookie)
		RemoveCookie(name string)
		ReadCookieValue(name string) (string, error)
		ReadCookie(name string) (*http.Cookie, error)
		AddView(name ...string) []string
		View(name string) error
		ViewC(code int, name string) error
		Write(code int, content []byte) (int, error)
		WriteString(contents ...interface{}) error
		WriteStringC(code int, contents ...interface{}) error
		WriteHtml(contents ...interface{}) error
		WriteHtmlC(code int, contents ...interface{}) error
		WriteBlob(contentType string, b []byte) error
		WriteBlobC(code int, contentType string, b []byte) error
		WriteJson(i interface{}) error
		WriteJsonC(code int, i interface{}) error
		WriteJsonBlob(b []byte) error
		WriteJsonBlobC(code int, b []byte) error
		WriteJsonp(callback string, i interface{}) error
		WriteJsonpBlob(callback string, b []byte) error
	}

	HttpContext struct {
		context context.Context
		//暂未启用
		cancle         context.CancelFunc
		middlewareStep string
		request        *Request
		routerNode     RouterNode
		routerParams   Params
		response       *Response
		webSocket      *WebSocket
		hijackConn     *HijackConn
		isWebSocket    bool
		isHijack       bool
		isEnd          bool //表示当前处理流程是否需要终止
		httpServer     *HttpServer
		sessionID      string
		innerItems     core.ConcurrenceMap
		items          core.ConcurrenceMap
		viewData       core.ConcurrenceMap
		features       *xFeatureTools
		handler        HttpHandle
	}
)

//reset response attr
func (ctx *HttpContext) reset(res *Response, r *Request, server *HttpServer, node RouterNode, params Params, handler HttpHandle) {
	ctx.request = r
	ctx.response = res
	ctx.routerNode = node
	ctx.routerParams = params
	ctx.isHijack = false
	ctx.isWebSocket = false
	ctx.httpServer = server
	ctx.innerItems = nil
	ctx.items = nil
	ctx.isEnd = false
	ctx.features = FeatureTools
	ctx.handler = handler
	ctx.Items().Set(ItemKeyHandleStartTime, time.Now())
}

//release all field
func (ctx *HttpContext) release() {
	ctx.request = nil
	ctx.response = nil
	ctx.middlewareStep = ""
	ctx.routerNode = nil
	ctx.routerParams = nil
	ctx.webSocket = nil
	ctx.hijackConn = nil
	ctx.isHijack = false
	ctx.isWebSocket = false
	ctx.httpServer = nil
	ctx.isEnd = false
	ctx.features = nil
	ctx.innerItems = nil
	ctx.items = nil
	ctx.viewData = nil
	ctx.sessionID = ""
	ctx.handler = nil
	ctx.Items().Remove(ItemKeyHandleStartTime)
	ctx.Items().Remove(ItemKeyHandleDuration)
}

// Context return context.Context
func (ctx *HttpContext) Context() context.Context {
	return ctx.context
}

// SetTimeoutContext set new Timeout Context
// set Context & cancle
// withvalue RequestID
func (ctx *HttpContext) SetTimeoutContext(timeout time.Duration) context.Context {
	ctx.context, ctx.cancle = context.WithTimeout(context.Background(), timeout)
	ctx.context = context.WithValue(ctx.context, "RequestID", ctx.Request().RequestID())
	return ctx.context
}

// WithContext set Context with RequestID
func (ctx *HttpContext) WithContext(runCtx context.Context) {
	if runCtx == nil {
		panic("nil context")
	}
	ctx.context = runCtx
	ctx.context = context.WithValue(ctx.context, "RequestID", ctx.Request().RequestID())
}

// HttpServer return HttpServer
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

func (ctx *HttpContext) IsWebSocket() bool {
	return ctx.isWebSocket
}

func (ctx *HttpContext) IsHijack() bool {
	return ctx.isHijack
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

// AppContext get application's global appcontext
// issue #3
func (ctx *HttpContext) AppItems() core.ConcurrenceMap {
	if ctx.httpServer != nil {
		return ctx.httpServer.DotApp.Items
	} else {
		return core.NewConcurrenceMap()
	}
}

// Cache get application's global cache
func (ctx *HttpContext) Cache() cache.Cache {
	return ctx.httpServer.DotApp.Cache()
}

// getInnerItems get request's inner item context
// lazy init when first use
func (ctx *HttpContext) getInnerItems() core.ConcurrenceMap {
	if ctx.innerItems == nil {
		ctx.innerItems = core.NewConcurrenceMap()
	}
	return ctx.innerItems
}

// Items get request's item context
// lazy init when first use
func (ctx *HttpContext) Items() core.ConcurrenceMap {
	if ctx.items == nil {
		ctx.items = core.NewConcurrenceMap()
	}
	return ctx.items
}

// AppSetConfig get appset from config file
// update for issue #16 配置文件
func (ctx *HttpContext) ConfigSet() core.ReadonlyMap {
	return ctx.HttpServer().DotApp.Config.ConfigSet
}

// ViewData get view data context
// lazy init when first use
func (ctx *HttpContext) ViewData() core.ConcurrenceMap {
	if ctx.viewData == nil {
		ctx.viewData = core.NewConcurrenceMap()
	}
	return ctx.viewData
}

// Session get session state in current context
func (ctx *HttpContext) Session() (state *session.SessionState) {
	if ctx.httpServer == nil {
		//return nil, errors.New("no effective http-server")
		panic("no effective http-server")
	}
	if !ctx.httpServer.SessionConfig().EnabledSession {
		//return nil, errors.New("http-server not enabled session")
		panic("http-server not enabled session")
	}
	state, _ = ctx.httpServer.sessionManager.GetSessionState(ctx.sessionID)
	return state
}

// Hijack make current connection to hijack mode
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
	ctx.isHijack = true
	return ctx.hijackConn, nil
}

// End set context user handler process end
// if set HttpContext.End,ignore user handler, but exec all http module  - fixed issue #5
func (ctx *HttpContext) End() {
	ctx.isEnd = true
}

func (ctx *HttpContext) IsEnd() bool {
	return ctx.isEnd
}

// Redirect redirect replies to the request with a redirect to url and with httpcode
// default you can use http.StatusFound
func (ctx *HttpContext) Redirect(code int, targetUrl string) error {
	return ctx.response.Redirect(code, targetUrl)
}

/*
* 根据指定key获取在Get请求中对应参数值
 */
func (ctx *HttpContext) QueryString(key string) string {
	return ctx.request.QueryString(key)
}

// QueryInt get query key with int format
// if not exists or not int type, return 0
func (ctx *HttpContext) QueryInt(key string) int {
	param := ctx.request.QueryString(key)
	if param == "" {
		return 0
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return 0
	}
	return val
}

// QueryInt64 get query key with int64 format
// if not exists or not int64 type, return 0
func (ctx *HttpContext) QueryInt64(key string) int64 {
	param := ctx.request.QueryString(key)
	if param == "" {
		return 0
	}
	val, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0
	}
	return val
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

// File sends a response with the content of the file
// if file not exists, response 404
// for issue #39
func (ctx *HttpContext) File(file string) (err error) {
	f, err := os.Open(file)
	if err != nil {
		HTTPNotFound(ctx)
		return nil
	}
	defer f.Close()

	fi, _ := f.Stat()
	if fi.IsDir() {
		file = filepath.Join(file, ctx.HttpServer().IndexPage())
		f, err = os.Open(file)
		if err != nil {
			HTTPNotFound(ctx)
			return nil
		}
		defer f.Close()
		if fi, err = f.Stat(); err != nil {
			return err
		}
	}
	http.ServeContent(ctx.Response().Writer(), ctx.Request().Request, fi.Name(), fi.ModTime(), f)
	return nil
}

// Attachment sends a response as attachment, prompting client to save the file.
// for issue #39
func (ctx *HttpContext) Attachment(file, name string) (err error) {
	return ctx.contentDisposition(file, name, "attachment")
}

// Inline sends a response as inline, opening the file in the browser.
// if file not exists, response 404
// for issue #39
func (ctx *HttpContext) Inline(file, name string) (err error) {
	return ctx.contentDisposition(file, name, "inline")
}

// contentDisposition set Content-disposition and response file
func (ctx *HttpContext) contentDisposition(file, name, dispositionType string) (err error) {
	ctx.Response().SetHeader(HeaderContentDisposition, fmt.Sprintf("%s; filename=%s", dispositionType, name))
	ctx.File(file)
	return
}

// Bind decode req.Body or form-value to struct
func (ctx *HttpContext) Bind(i interface{}) error {
	return ctx.httpServer.Binder().Bind(i, ctx)
}

// BindJsonBody default use json decode req.Body to struct
func (ctx *HttpContext) BindJsonBody(i interface{}) error {
	return ctx.httpServer.Binder().BindJsonBody(i, ctx)
}
// Validate validates data with HttpServer::Validator
// We will implementing inner validator on next version
func (ctx *HttpContext) Validate(i interface{}) error {
	if ctx.httpServer.Validator == nil {
		return ErrValidatorNotRegistered
	}
	return ctx.httpServer.Validator.Validate(i)
}

// GetRouterName get router name
func (ctx *HttpContext) GetRouterName(key string) string {
	return ctx.routerParams.ByName(key)
}

// RemoteIP return user IP address
func (ctx *HttpContext) RemoteIP() string {
	return ctx.request.RemoteIP()
}

// SetCookieValue write cookie for name & value & maxAge
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

// SetCookie write cookie with cookie-obj
func (ctx *HttpContext) SetCookie(cookie *http.Cookie) {
	http.SetCookie(ctx.response.Writer(), cookie)
}

// RemoveCookie remove cookie for path&name
func (ctx *HttpContext) RemoveCookie(name string) {
	cookie := &http.Cookie{Name: name, MaxAge: -1}
	ctx.SetCookie(cookie)
}

// ReadCookieValue read cookie value for name
func (ctx *HttpContext) ReadCookieValue(name string) (string, error) {
	cookieobj, err := ctx.request.Cookie(name)
	if err != nil {
		return "", err
	} else {
		return url.QueryUnescape(cookieobj.Value)
	}
}

// ReadCookie read cookie object for name
func (ctx *HttpContext) ReadCookie(name string) (*http.Cookie, error) {
	return ctx.request.Cookie(name)
}

// AddView add need parse views before View()
func (ctx *HttpContext) AddView(names ...string) []string {
	var views []string
	item, exists := ctx.getInnerItems().Get(innerKeyAddView)
	if exists {
		views = item.([]string)
	}
	views = append(views, names...)
	ctx.getInnerItems().Set(innerKeyAddView, views)
	return views
}

// View write view content to response
func (ctx *HttpContext) View(name string) error {
	return ctx.ViewC(defaultHttpCode, name)
}

// ViewC write (httpCode, view content) to response
func (ctx *HttpContext) ViewC(code int, name string) error {
	ctx.response.SetStatusCode(code)
	ctx.AddView(name)
	var views []string
	item, exists := ctx.getInnerItems().Get(innerKeyAddView)
	if exists {
		views = item.([]string)
	} else {
		return errors.New("get view info error")
	}
	err := ctx.httpServer.Renderer().Render(ctx.response.Writer(), ctx.ViewData().GetCurrentMap(), ctx, views...)
	return err
}

// Write write code and content content to response
func (ctx *HttpContext) Write(code int, content []byte) (int, error) {
	if ctx.IsHijack() {
		//TODO:hijack mode, status-code set default 200
		return ctx.hijackConn.WriteBlob(content)
	} else {
		return ctx.response.Write(code, content)
	}
}

// WriteString write (200, string, text/plain) to response
func (ctx *HttpContext) WriteString(contents ...interface{}) error {
	return ctx.WriteStringC(defaultHttpCode, contents...)
}

// WriteStringC write (httpCode, string, text/plain) to response
func (ctx *HttpContext) WriteStringC(code int, contents ...interface{}) error {
	content := fmt.Sprint(contents...)
	return ctx.WriteBlobC(code, MIMETextPlainCharsetUTF8, []byte(content))
}

// WriteString write (200, string, text/html) to response
func (ctx *HttpContext) WriteHtml(contents ...interface{}) error {
	return ctx.WriteHtmlC(defaultHttpCode, contents...)
}

// WriteHtmlC write (httpCode, string, text/html) to response
func (ctx *HttpContext) WriteHtmlC(code int, contents ...interface{}) error {
	content := fmt.Sprint(contents...)
	return ctx.WriteBlobC(code, MIMETextHTMLCharsetUTF8, []byte(content))
}

// WriteBlob write []byte content to response
func (ctx *HttpContext) WriteBlob(contentType string, b []byte) error {
	return ctx.WriteBlobC(defaultHttpCode, contentType, b)
}

// WriteBlobC write (httpCode, []byte) to response
func (ctx *HttpContext) WriteBlobC(code int, contentType string, b []byte) error {
	if contentType != "" {
		ctx.response.SetContentType(contentType)
	}
	if ctx.IsHijack() {
		_, err := ctx.hijackConn.WriteBlob(b)
		return err
	} else {
		_, err := ctx.response.Write(code, b)
		return err
	}
}

// WriteJson write (httpCode, json string) to response
// auto convert interface{} to json string
func (ctx *HttpContext) WriteJson(i interface{}) error {
	return ctx.WriteJsonC(defaultHttpCode, i)
}

// WriteJsonC write (httpCode, json string) to response
// auto convert interface{} to json string
func (ctx *HttpContext) WriteJsonC(code int, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return ctx.WriteJsonBlobC(code, b)
}

// WriteJsonBlob write json []byte to response
func (ctx *HttpContext) WriteJsonBlob(b []byte) error {
	return ctx.WriteJsonBlobC(defaultHttpCode, b)
}

// WriteJsonBlobC write (httpCode, json []byte) to response
func (ctx *HttpContext) WriteJsonBlobC(code int, b []byte) error {
	return ctx.WriteBlobC(code, MIMEApplicationJSONCharsetUTF8, b)
}

// WriteJsonp write jsonp string to response
func (ctx *HttpContext) WriteJsonp(callback string, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return ctx.WriteJsonpBlob(callback, b)
}

// WriteJsonpBlob write jsonp string as []byte to response
func (ctx *HttpContext) WriteJsonpBlob(callback string, b []byte) error {
	var err error
	ctx.response.SetContentType(MIMEApplicationJavaScriptCharsetUTF8)
	//特殊处理，如果为hijack，需要先行WriteBlob头部
	if ctx.IsHijack() {
		if _, err = ctx.hijackConn.WriteBlob([]byte(ctx.hijackConn.header + "\r\n")); err != nil {
			return err
		}
	}
	if err = ctx.WriteBlob("", []byte(callback+"(")); err != nil {
		return err
	}
	if err = ctx.WriteBlob("", b); err != nil {
		return err
	}
	err = ctx.WriteBlob("", []byte(");"))
	return err
}
