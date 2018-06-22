package dotweb

import (
	"github.com/devfeel/dotweb/framework/file"
	"html/template"
	"io"
	"path"
	"sync"
)

// Renderer is the interface that wraps the render method.
type Renderer interface {
	SetTemplatePath(path string)
	Render(io.Writer, interface{}, Context, ...string) error
}

type innerRenderer struct {
	templatePath string
	// Template cache (for FromCache())
	enabledCache 	   bool
	templateCache      map[string]*template.Template
	templateCacheMutex sync.RWMutex
}

// Render render view use http/template
func (r *innerRenderer) Render(w io.Writer, data interface{}, ctx Context, tpl ...string) error {
	t, err := r.parseFiles(tpl...)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// SetTemplatePath set default template path
func (r *innerRenderer) SetTemplatePath(path string) {
	r.templatePath = path
}

// 定义函数unescaped
func unescaped(x string) interface{} { return template.HTML(x) }

// return http/template by gived file name
func (r *innerRenderer) parseFiles(fileNames ...string) (*template.Template, error) {
	var realFileNames []string
	var filesCacheKey string
	var err error
	for _, v := range fileNames {
		if !file.Exist(v) {
			v = path.Join(r.templatePath, v)
		}
		realFileNames = append(realFileNames, v)
		filesCacheKey = filesCacheKey + v
	}

	var t *template.Template
	var exists bool
	if r.enabledCache {
		//check from chach
		t, exists = r.parseFilesFromCache(filesCacheKey)
	}
	if !exists{
		t, err = template.ParseFiles(realFileNames...)
		if err != nil {
			return nil, err
		}
		r.templateCacheMutex.Lock()
		defer r.templateCacheMutex.Unlock()
		r.templateCache[filesCacheKey] = t
	}

	t = registeTemplateFunc(t)
	return t, nil
}

func (r *innerRenderer) parseFilesFromCache(filesCacheKey string) (*template.Template, bool){
	r.templateCacheMutex.RLock()
	defer r.templateCacheMutex.RUnlock()
	t, exists:= r.templateCache[filesCacheKey]
	return t, exists
}

// registeTemplateFunc registe default support funcs
func registeTemplateFunc(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{"unescaped": unescaped})
	//TODO:add more func
}

// NewInnerRenderer create a inner renderer instance
func NewInnerRenderer() Renderer {
	r := new(innerRenderer)
	r.enabledCache = true
	r.templateCache = make(map[string]*template.Template)
	return r
}

// NewInnerRendererNoCache create a inner renderer instance with no cache mode
func NewInnerRendererNoCache() Renderer {
	r := new(innerRenderer)
	r.enabledCache = false
	r.templateCache = make(map[string]*template.Template)
	return r
}

