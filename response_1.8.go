// +build go1.8

package dotweb

import "net/http"

func (r *Response) Push(target string, opts *http.PushOptions) error {
	return r.writer.(http.Pusher).Push(target, opts)
}
