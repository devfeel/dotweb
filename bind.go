package dotweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/devfeel/dotweb/framework/reflects"
	"strings"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(interface{}, *HttpContext) error
	}

	binder struct{}
)

func (b *binder) Bind(i interface{}, ctx *HttpContext) (err error) {
	req := ctx.Request
	ctype := req.Header.Get(HeaderContentType)
	if req.Body == nil {
		err = errors.New("request body can't be empty")
		return err
	}
	err = errors.New("request unsupported MediaType -> " + ctype)
	switch {
	case strings.HasPrefix(ctype, MIMEApplicationJSON):
		err = json.NewDecoder(req.Body).Decode(i)
	case strings.HasPrefix(ctype, MIMEApplicationXML):
		err = xml.NewDecoder(req.Body).Decode(i)
	case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
		err = reflects.ConvertMapToStruct("form", i, ctx.FormValues())
	}
	return err
}
