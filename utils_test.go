package dotweb

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
  "net/url"
)

//common init context
func initContext(param *InitContextParam) *HttpContext {
	httpRequest := &http.Request{}
	context := &HttpContext{
		request: &Request{
			Request: httpRequest,
		},
	}
	header := make(map[string][]string)
	header["Accept-Encoding"] = []string{"gzip, deflate"}
	header["Accept-Language"] = []string{"en-us"}
	header["Foo"] = []string{"Bar", "two"}
	//specify json
	header["Content-Type"] = []string{param.contentType}
	context.request.Header = header

	jsonStr := param.convertHandler(param.t, param.v)
	body := format(jsonStr)
	context.request.Request.Body = body

	return context
}

//init response context
func initResponseContext(param *InitContextParam) *HttpContext {
	context := &HttpContext{
		response: &Response{},
	}

	var buf1 bytes.Buffer
	w := io.MultiWriter(&buf1)

	writer := &gzipResponseWriter{
		ResponseWriter: &httpWriter{},
		Writer:         w,
	}

	context.response = NewResponse(writer)

	return context
}

//init request and response context
func initAllContext(param *InitContextParam) *HttpContext {
	context := &HttpContext{
		response: &Response{},
		request: &Request{
			Request: &http.Request{},
		},
	}
	header := make(map[string][]string)
	header["Accept-Encoding"] = []string{"gzip, deflate"}
	header["Accept-Language"] = []string{"en-us"}
	header["Foo"] = []string{"Bar", "two"}
	//specify json
	header["Content-Type"] = []string{param.contentType}
	context.request.Header = header

	u:=&url.URL{
		Path:"/index",
	}

	context.request.URL=u
	context.request.Method="POST"

	jsonStr := param.convertHandler(param.t, param.v)
	body := format(jsonStr)
	context.request.Request.Body = body

	//var buf1 bytes.Buffer
	//w := io.MultiWriter(&buf1)

	w := &httpWriter{}
	//gzip 开关
	/*
		gw, _ := gzip.NewWriterLevel(w, DefaultGzipLevel)
		writer := &gzipResponseWriter{
			ResponseWriter: w,
			Writer:         &gzipResponseWriter{Writer: gw, ResponseWriter: w},
		}
	*/

	context.response = NewResponse(w)

	return context
}

type httpWriter http.Header

func (ho httpWriter) Header() http.Header {
	return http.Header(ho)
}

func (ho httpWriter) Write(byte []byte) (int, error) {
	fmt.Println("string:", string(byte))
	return 0, nil
}

func (ho httpWriter) WriteHeader(code int) {
	fmt.Println("code:", code)
}

func format(b string) io.ReadCloser {
	s := strings.NewReader(b)
	r := ioutil.NopCloser(s)
	r.Close()
	return r
}

type InitContextParam struct {
	t              *testing.T
	v              interface{}
	contentType    string
	convertHandler func(t *testing.T, v interface{}) string
}
