package dotweb

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Request struct {
	*http.Request
	postBody   []byte
	isReadBody bool
}

//reset response attr
func (req *Request) reset(r *http.Request) {
	req.Request = r
	req.isReadBody = false
}

func (req *Request) release() {
	req.Request = nil
	req.isReadBody = false
	req.postBody = nil
}

/*
* 返回查询字符串map表示
 */
func (req *Request) QueryStrings() url.Values {
	return req.URL.Query()
}

/*
* 获取原始查询字符串
 */
func (req *Request) RawQuery() string {
	return req.URL.RawQuery
}

/*
* 根据指定key获取对应value
 */
func (req *Request) QueryString(key string) string {
	return req.URL.Query().Get(key)
}

func (req *Request) FormFile(key string) (*UploadFile, error) {
	file, header, err := req.Request.FormFile(key)
	if err != nil {
		return nil, err
	} else {
		return NewUploadFile(file, header), nil
	}
}

/*
* 获取包括post、put和get内的值
 */
func (req *Request) FormValues() map[string][]string {
	req.parseForm()
	return map[string][]string(req.Form)
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

func (req *Request) ContentType() string {
	return req.Header.Get(HeaderContentType)
}

func (req *Request) QueryHeader(key string) string {
	return req.Header.Get(key)
}

/*
* 根据指定key获取包括在post、put内的值
* Obsolete("use PostFormValue replace this")
 */
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
	s := strings.Split(fullIp, ":")
	if len(s) > 1 {
		return s[0]
	} else {
		return fullIp
	}
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

func (req *Request) Url() string {
	return req.URL.String()
}
