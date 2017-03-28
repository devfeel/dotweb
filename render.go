package dotweb

import (
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

//Renderer is the interface that wraps the render method.
type Renderer interface {
	Render(io.Writer, string, interface{}, *HttpContext) error
}

type innerRenderer struct {
}

// Render render view use http/template
func (r *innerRenderer) Render(w io.Writer, tpl string, data interface{}, ctx *HttpContext) error {
	t, err := parseFile(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// 定义函数unescaped
func unescaped(x string) interface{} { return template.HTML(x) }

// return http/template by gived file name
func parseFile(filename string) (*template.Template, error) {
	var t *template.Template
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := string(b)
	name := filepath.Base(filename)
	t = template.New(name)
	t = registeTemplateFunc(t)
	_, err = t.Parse(s)
	if err != nil {
		return nil, err
	}
	return t, nil
}

//registe default support funcs
func registeTemplateFunc(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{"unescaped": unescaped})
	//TODO:add more func
}

// NewInnerRenderer create a inner renderer instance
func NewInnerRenderer() *innerRenderer {
	r := new(innerRenderer)
	return r
}
