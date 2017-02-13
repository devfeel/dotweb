package dotweb

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type Response struct {
	writer    http.ResponseWriter
	Status    int
	Size      int64
	body      []byte
	committed bool
	header    http.Header
}

func NewResponse(w http.ResponseWriter) (r *Response) {
	return &Response{writer: w,
		header: w.Header()}
}

func (r *Response) Header() http.Header {
	return r.header
}

func (r *Response) Writer() http.ResponseWriter {
	return r.writer
}

func (r *Response) Body() []byte {
	return r.body
}

func (r *Response) BodyString() string {
	return string(r.body)
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
func (r *Response) Write(b []byte) (n int, err error) {
	if !r.committed {
		r.WriteHeader(http.StatusOK)
	}
	n, err = r.writer.Write(b)
	r.Size += int64(n)
	r.body = append(r.body, b[0:]...)
	return
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
func (r *Response) Reset(w http.ResponseWriter) {
	r.writer = w
	r.header = w.Header()
	r.Status = http.StatusOK
	r.Size = 0
	r.committed = false
	r.writer = w
}
