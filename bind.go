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
	jsonTagName = "json"
)

type (
	// Binder is the interface that wraps the Bind method.
	Binder interface {
		Bind(interface{}, Context) error
		BindJsonBody(interface{}, Context) error
	}

	binder struct{}
)

//Bind decode req.Body or form-value to struct
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
		err = json.Unmarshal(ctx.Request().PostBody(), i)
	case strings.HasPrefix(ctype, MIMEApplicationXML):
		err = xml.Unmarshal(ctx.Request().PostBody(), i)
	//case strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm),
	//	strings.HasPrefix(ctype, MIMETextHTML):
	//	err = reflects.ConvertMapToStruct(defaultTagName, i, ctx.FormValues())
	default:
		//check is use json tag, fixed for issue #91
		tagName := defaultTagName
		if ctx.HttpServer().ServerConfig().EnabledBindUseJsonTag{
			tagName = jsonTagName
		}
		//no check content type for fixed issue #6
		err = reflects.ConvertMapToStruct(tagName, i, ctx.Request().FormValues())
	}
	return err
}

//BindJsonBody default use json decode req.Body to struct
func (b *binder) BindJsonBody(i interface{}, ctx Context) (err error) {
	if ctx.Request().PostBody() == nil {
		err = errors.New("request body can't be empty")
		return err
	}
	err = json.Unmarshal(ctx.Request().PostBody(), i)
	return err
}

func newBinder() *binder {
	return &binder{}
}
