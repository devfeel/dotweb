package dotweb

import (
	"errors"
	"github.com/devfeel/dotweb/framework/file"
	"html/template"
	"io"
	"io/ioutil"
	"path/filepath"
)

// Renderer is the interface that wraps the render method.
type Renderer interface {
	//set default template path, support multi path
	//模板查找顺序从最后一个插入的元素开始往前找
	//默认添加base、base/templates、base/views
	SetTemplatePath(path ...string)
	Render(io.Writer, string, interface{}, Context) error
}

type innerRenderer struct {
	templatePaths []string
}

// Render render view use http/template
func (r *innerRenderer) Render(w io.Writer, tpl string, data interface{}, ctx Context) error {
	t, err := r.parseFile(tpl)
	if err != nil {
		return err
	}
	return t.Execute(w, data)
}

// SetTemplatePath set default template paths, support multi path
// 模板查找顺序从最后一个插入的元素开始往前找
// 默认添加base、base/templates、base/views
func (r *innerRenderer) SetTemplatePath(path ...string) {
	r.templatePaths = append(r.templatePaths, path...)
}

// 定义函数unescaped
func unescaped(x string) interface{} { return template.HTML(x) }

// return http/template by gived file name
func (r *innerRenderer) parseFile(filename string) (*template.Template, error) {
	var t *template.Template
	findlog := filename
	realTempFile := filename
	isExist := false
	//多级检查
	if !file.Exist(filename) {
		tmpFileName := filename
		for i := len(r.templatePaths) - 1; i >= 0; i-- {
			tmpFileName = r.templatePaths[i] + "/" + filename
			findlog += "\r\n" + tmpFileName
			if file.Exist(tmpFileName) {
				realTempFile = tmpFileName
				isExist = true
				break
			}
		}
	} else {
		isExist = true
	}

	if !isExist {
		return nil, errors.New("not found template file=>\r\n " + findlog)
	}

	b, err := ioutil.ReadFile(realTempFile)
	if err != nil {
		return nil, err
	}
	s := string(b)
	name := filepath.Base(realTempFile)
	t = template.New(name)
	t = registeTemplateFunc(t)
	_, err = t.Parse(s)
	if err != nil {
		return nil, err
	}
	return t, nil
}

// registeTemplateFunc registe default support funcs
func registeTemplateFunc(t *template.Template) *template.Template {
	return t.Funcs(template.FuncMap{"unescaped": unescaped})
	//TODO:add more func
}

// NewInnerRenderer create a inner renderer instance
func NewInnerRenderer() *innerRenderer {
	r := new(innerRenderer)
	r.templatePaths = make([]string, 3)
	//添加基础路径
	//base、base/templates、base/views
	r.templatePaths[0] = file.GetCurrentDirectory()
	r.templatePaths[1] = file.GetCurrentDirectory() + "/templates"
	r.templatePaths[2] = file.GetCurrentDirectory() + "/views"
	return r
}
