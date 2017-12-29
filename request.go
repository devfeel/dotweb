package dotweb

import (
	"github.com/devfeel/dotweb/framework/crypto/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	*http.Request
	httpCtx    *HttpContext
	postBody   []byte
	isReadBody bool
	requestID  string
}

//reset response attr
func (req *Request) reset(r *http.Request, ctx *HttpContext) {
	req.Request = r
	req.isReadBody = false
	if ctx.HttpServer().ServerConfig().EnabledRequestID {
		req.requestID = uuid.NewV4().String32()
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

// QueryStrings 返回Get请求方式下查询字符串map表示
func (req *Request) QueryStrings() url.Values {
	return req.URL.Query()
}

// RawQuery 获取原始查询字符串
func (req *Request) RawQuery() string {
	return req.URL.RawQuery
}

// QueryString 根据指定key获取在Get请求中对应参数值
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
func (req *Request) FormFiles()(map[string]*UploadFile, error){
	files := make(map[string]*UploadFile)
	req.parseForm()
	if req.Request.MultipartForm == nil || req.Request.MultipartForm.File == nil {
		return nil, http.ErrMissingFile
	}
	for key, fileMap:=range req.Request.MultipartForm.File{
		if len(fileMap) > 0{
			file, err := fileMap[0].Open()
			if err== nil{
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

//Deprecated: Use the PostFormValue instead
//returns the first value for the named component of the POST
// or PUT request body. URL query parameters are ignored.
func (req *Request) PostString(key string) string {
	return req.PostFormValue(key)
}

/*
* 获取post提交的字节数组
 */
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

//RemoteAddr to an "IP" address
func (req *Request) RemoteIP() string {
	fullIp := req.Request.RemoteAddr
	//special: if run in win10, localIp will be like "[::]:port"
	//fixed for #20 cann't get RemoteIP and RemoteAddr in win10
	lastFlagIndex := strings.LastIndex(fullIp, ":")
	if lastFlagIndex >= 0 {
		return fullIp[:lastFlagIndex]
	}
	return fullIp
}

//RemoteAddr to an "IP:port" address
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
	return req.Header.Get(HeaderXRequestedWith) == "XMLHttpRequest"
}

// Url get request url
func (req *Request) Url() string {
	return req.URL.String()
}
