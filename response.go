package dotweb

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"net"
	"net/http"
)

type (
	Response struct {
		writer    http.ResponseWriter
		Status    int
		Size      int64
		body      []byte
		committed bool
		header    http.Header
		isEnd     bool
	}

	gzipResponseWriter struct {
		io.Writer
		http.ResponseWriter
	}
)

func NewResponse(w http.ResponseWriter) (r *Response) {
	return &Response{writer: w,
		header: w.Header()}
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) QueryHeader(key string) string {
	return r.Header().Get(key)
}

func (r *Response) Redirect(code int, targetUrl string) error {
	r.Header().Set(HeaderCacheControl, "no-cache")
	r.Header().Set(HeaderLocation, targetUrl)
	return r.WriteHeader(code)
}

func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

func (r *Response) SetWriter(w http.ResponseWriter) *Response {
	r.writer = w
	return r
}

//HttpCode return http code format int
func (r *Response) HttpCode() int {
	return r.Status
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) BodyString() string {
	return string(r.body)
}

func (r *Response) SetHeader(key, val string) {
	r.Header().Set(key, val)
}

func (r *Response) SetContentType(contenttype string) {
	r.SetHeader(HeaderContentType, contenttype)
}

func (r *Response) SetStatusCode(code int) error {
	return r.WriteHeader(code)
}

// WriteHeader sends an HTTP response header with status code. If WriteHeader is
// not called explicitly, the first call to Write will trigger an implicit
// WriteHeader(http.StatusOK). Thus explicit calls to WriteHeader are mainly
// used to send error codes.
func (r *Response) WriteHeader(code int) error {
	if r.committed {
		return errors.New("response already set status")
	}
	r.Status = code
	r.writer.WriteHeader(code)
	r.committed = true
	return nil
}

// Write writes the data to the connection as part of an HTTP reply.
func (r *Response) Write(code int, b []byte) (n int, err error) {
	if !r.committed {
		r.WriteHeader(code)
	}
	n, err = r.writer.Write(b)
	r.Size += int64(n)
	r.body = append(r.body, b[0:]...)
	return
}

//stop current response
func (r *Response) End() {
	r.isEnd = true
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (r *Response) Flush() {
	r.writer.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See https://golang.org/pkg/net/http/#Hijacker
func (r *Response) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return r.writer.(http.Hijacker).Hijack()
}

//reset response attr
func (r *Response) reset(w http.ResponseWriter) {
	r.writer = w
	r.header = w.Header()
	r.Status = http.StatusOK
	r.Size = 0
	r.body = nil
	r.committed = false
}

//reset response attr
func (r *Response) release() {
	r.writer = nil
	r.header = nil
	r.Status = http.StatusOK
	r.Size = 0
	r.body = nil
	r.committed = false
}

/*gzipResponseWriter*/
func (w *gzipResponseWriter) WriteHeader(code int) {
	if code == http.StatusNoContent { // Issue #489
		w.ResponseWriter.Header().Del(HeaderContentEncoding)
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	if w.Header().Get(HeaderContentType) == "" {
		w.Header().Set(HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Flush() {
	w.Writer.(*gzip.Writer).Flush()
}

func (w *gzipResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.(http.Hijacker).Hijack()
}
