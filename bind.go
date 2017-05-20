package dotweb

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/devfeel/dotweb/framework/reflects"
	"strings"
)

const (
	defaultTagName = "form"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(interface{}, Context) error
	}

	binder struct{}
)

func (b *binder) Bind(i interface{}, ctx Context) (err error) {
	req := ctx.Request()
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
	//case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm),
	//	strings.HasPrefix(ctype, MIMETextHTML):
	//	err = reflects.ConvertMapToStruct(defaultTagName, i, ctx.FormValues())
	default:
		//no check content type for fixed issue #6
		err = reflects.ConvertMapToStruct(defaultTagName, i, ctx.Request().FormValues())
	}
	return err
}

func newBinder() *binder {
	return &binder{}
}
