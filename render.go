package dotweb

import (
	"github.com/devfeel/dotweb/framework/file"
	"html/template"
	"io"
	"path"
	"sync"
	"path/filepath"
	"errors"
)

// Renderer is the interface that wraps the render method.
type Renderer interface {
	SetTemplatePath(path string)
	Render(io.Writer, interface{}, Context, ...string) error
	RegisterTemplateFunc(string, interface{})
}

type innerRenderer struct {
	templatePath string
	// Template cache (for FromCache())
	enabledCache 	   bool
	templateCache      map[string]*template.Template
	templateCacheMutex sync.RWMutex

	// used to manager template func
	templateFuncs 	   map[string]interface{}
	templateFuncsMutex *sync.RWMutex

}

// Render render view use http/template
func (r *innerRenderer) Render(w io.Writer, data interface{}, ctx Context, tpl ...string) error {
	if len(tpl) <= 0{
		return errors.New("no enough render template files")
	}
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

// RegisterTemplateFunc used to register template func in renderer
func (r *innerRenderer) RegisterTemplateFunc(funcName string, funcHandler interface{}){
	r.templateFuncsMutex.Lock()
	r.templateFuncs[funcName] = funcHandler
	r.templateFuncsMutex.Unlock()
}


// unescaped inner template func used to encapsulates a known safe HTML document fragment
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
		name := filepath.Base(fileNames[0])
		t = template.New(name)
		if len(r.templateFuncs) > 0{
			t = t.Funcs(r.templateFuncs)
		}
		t, err = t.ParseFiles(realFileNames...)
		if err != nil {
			return nil, err
		}
		r.templateCacheMutex.Lock()
		defer r.templateCacheMutex.Unlock()
		r.templateCache[filesCacheKey] = t
	}

	return t, nil
}

func (r *innerRenderer) parseFilesFromCache(filesCacheKey string) (*template.Template, bool){
	r.templateCacheMutex.RLock()
	defer r.templateCacheMutex.RUnlock()
	t, exists:= r.templateCache[filesCacheKey]
	return t, exists
}

// registeInnerTemplateFunc registe default support funcs
func registeInnerTemplateFunc(funcMap map[string]interface{}){
	funcMap["unescaped"] = unescaped
}

// NewInnerRenderer create a inner renderer instance
func NewInnerRenderer() Renderer {
	r := new(innerRenderer)
	r.enabledCache = true
	r.templateCache = make(map[string]*template.Template)
	r.templateFuncs = make(map[string]interface{})
	r.templateFuncsMutex = new(sync.RWMutex)
	registeInnerTemplateFunc(r.templateFuncs)
	return r
}

// NewInnerRendererNoCache create a inner renderer instance with no cache mode
func NewInnerRendererNoCache() Renderer {
	r := NewInnerRenderer().(*innerRenderer)
	r.enabledCache = false
	return r
}

