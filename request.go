package dotweb

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
)


type Request struct {
	*http.Request
	httpCtx    *HttpContext
	postBody   []byte
	realUrl    string
	isReadBody bool
	requestID  string
}

// reset response attr
func (req *Request) reset(r *http.Request, ctx *HttpContext) {
	req.Request = r
	req.isReadBody = false
	if ctx.HttpServer().ServerConfig().EnabledRequestID {
		req.requestID = ctx.HttpServer().DotApp.IDGenerater()
		ctx.response.SetHeader(HeaderRequestID, req.requestID)
	} else {
		req.requestID = ""
	}
}

func (req *Request) release() {
	req.Request = nil
	req.isReadBody = false
	req.postBody = nil
	req.requestID = ""
}

// RequestID get unique ID with current request
// must HttpServer.SetEnabledRequestID(true)
// default is empty string
func (req *Request) RequestID() string {
	return req.requestID
}

// QueryStrings parses RawQuery and returns the corresponding values.
func (req *Request) QueryStrings() url.Values {
	return req.URL.Query()
}

// RawQuery returns the original query string
func (req *Request) RawQuery() string {
	return req.URL.RawQuery
}

// QueryString returns the first value associated with the given key.
func (req *Request) QueryString(key string) string {
	return req.URL.Query().Get(key)
}

// FormFile get file by form key
func (req *Request) FormFile(key string) (*UploadFile, error) {
	file, header, err := req.Request.FormFile(key)
	if err != nil {
		return nil, err
	} else {
		return NewUploadFile(file, header), nil
	}
}

// FormFiles get multi files
// fixed #92
func (req *Request) FormFiles() (map[string]*UploadFile, error) {
	files := make(map[string]*UploadFile)
	req.parseForm()
	if req.Request.MultipartForm == nil || req.Request.MultipartForm.File == nil {
		return nil, http.ErrMissingFile
	}
	for key, fileMap := range req.Request.MultipartForm.File {
		if len(fileMap) > 0 {
			file, err := fileMap[0].Open()
			if err == nil {
				files[key] = NewUploadFile(file, fileMap[0])
			}
		}
	}
	return files, nil
}

// FormValues including both the URL field's query parameters and the POST or PUT form data
func (req *Request) FormValues() map[string][]string {
	req.parseForm()
	return map[string][]string(req.Form)
}

// PostValues contains the parsed form data from POST, PATCH, or PUT body parameters
func (req *Request) PostValues() map[string][]string {
	req.parseForm()
	return map[string][]string(req.PostForm)
}

func (req *Request) parseForm() error {
	if strings.HasPrefix(req.QueryHeader(HeaderContentType), MIMEMultipartForm) {
		if err := req.ParseMultipartForm(defaultMemory); err != nil {
			return err
		}
	} else {
		if err := req.ParseForm(); err != nil {
			return err
		}
	}
	return nil
}

// ContentType get ContentType
func (req *Request) ContentType() string {
	return req.Header.Get(HeaderContentType)
}

// QueryHeader query header value by key
func (req *Request) QueryHeader(key string) string {
	return req.Header.Get(key)
}

// PostString returns the first value for the named component of the POST
// or PUT request body. URL query parameters are ignored.
// Deprecated: Use the PostFormValue instead
func (req *Request) PostString(key string) string {
	return req.PostFormValue(key)
}

// PostBody returns data from the POST or PUT request body
func (req *Request) PostBody() []byte {
	if !req.isReadBody {
		bts, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return []byte{}
		} else {
			req.isReadBody = true
			req.postBody = bts
		}
	}
	return req.postBody
}

// RemoteIP RemoteAddr to an "IP" address
func (req *Request) RemoteIP() string {
	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	return host
}

// RealIP returns the first ip from 'X-Forwarded-For' or 'X-Real-IP' header key
// if not exists data, returns request.RemoteAddr
// fixed for #164
func (req *Request) RealIP() string {
	if ip := req.Header.Get(HeaderXForwardedFor); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := req.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	host, _, _ := net.SplitHostPort(req.RemoteAddr)
	return host
}

// FullRemoteIP RemoteAddr to an "IP:port" address
func (req *Request) FullRemoteIP() string {
	fullIp := req.Request.RemoteAddr
	return fullIp
}

// Path returns requested path.
//
// The path is valid until returning from RequestHandler.
func (req *Request) Path() string {
	return req.URL.Path
}

// IsAJAX returns if it is a ajax request
func (req *Request) IsAJAX() bool {
	return strings.Contains(req.Header.Get(HeaderXRequestedWith), "XMLHttpRequest")
}

// Url get request url
func (req *Request) Url() string {
	if req.realUrl != "" {
		return req.realUrl
	} else {
		return req.URL.String()
	}
}
