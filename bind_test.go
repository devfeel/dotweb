package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"io"
	"encoding/xml"
)

type Person struct {
	Hair string
	HasGlass bool
	Age int
	Legs []string
}

type InitContextParam struct {
	t *testing.T
	v interface{}
	contentType string
	convertHandler func(t *testing.T,v interface{})string
}

//json
func TestBinder_Bind_json(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}

	//init param
	param:=&InitContextParam{
		t,
		expected,
		"application/json",
		toJson,
	}

	//init param
	context:=initContext(param)
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must nil
	test.Nil(t,err)

	//check expected
	test.Equal(t,expected,person)

	t.Log("person:",person)
	t.Log("expected:",expected)

}

//json
func TestBinder_Bind_json_error(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}

	//init param
	param:=&InitContextParam{
		t,
		expected,
		"application/xml",
		toJson,
	}

	//init param
	context:=initContext(param)
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must not nil
	test.NotNil(t,err)
}

//xml
func TestBinder_Bind_xml(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}
	param:=&InitContextParam{
		t,
		expected,
		"application/xml",
		toXml,
	}

	//init param
	context:=initContext(param)
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must nil
	test.Nil(t,err)

	//check expected
	test.Equal(t,expected,person)

	t.Log("person:",person)
	t.Log("expected:",expected)

}

//xml
func TestBinder_Bind_xml_error(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}
	param:=&InitContextParam{
		t,
		expected,
		"application/json",
		toXml,
	}

	//init param
	context:=initContext(param)
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must not nil
	test.NotNil(t,err)
}

//else
func TestBinder_Bind_default(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}
	param:=&InitContextParam{
		t,
		expected,
		"",
		toDefault,
	}

	//init param
	context:=initContext(param)

	form:=make(map[string][]string)
	form["Hair"]=[]string{"Brown"}
	form["HasGlass"]=[]string{"true"}
	form["Age"]=[]string{"10"}
	form["Legs"]=[]string{"Left", "Right"}

	context.request.Form=form
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must nil
	test.Nil(t,err)

	//check expected
	test.Equal(t,expected,person)

	t.Log("person:",person)
	t.Log("expected:",expected)

}

//else
func TestBinder_Bind_default_error(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}
	param:=&InitContextParam{
		t,
		expected,
		"application/xml",
		toDefault,
	}

	//init param
	context:=initContext(param)

	form:=make(map[string][]string)
	form["Hair"]=[]string{"Brown"}
	form["HasGlass"]=[]string{"true"}
	form["Age"]=[]string{"10"}
	form["Legs"]=[]string{"Left", "Right"}

	context.request.Form=form
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must not nil
	test.NotNil(t,err)

}

//common init context
func initContext(param *InitContextParam) *HttpContext {
	httpRequest:=&http.Request{}
	context:=&HttpContext{
		request:&Request{
			Request:httpRequest,
		},
	}
	header:=make(map[string][]string)
	header["Accept-Encoding"]=[]string{"gzip, deflate"}
	header["Accept-Language"]=[]string{"en-us"}
	header["Foo"]=[]string{"Bar", "two"}
	//specify json
	header["Content-Type"]=[]string{param.contentType}
	context.request.Header=header

	jsonStr:=param.convertHandler(param.t,param.v)
	body:=format(jsonStr)
	context.request.Request.Body=body


	return context
}

//default
//TODO:content type is null but body not null,is it right??
func TestBinder_Bind_ContentTypeNull(t *testing.T) {

	binder:=newBinder()

	if binder==nil {
		t.Error("binder can not be nil!")
	}

	//init DotServer
	app := New()

	if app==nil {
		t.Error("app can not be nil!")
	}

	//expected
	expected:=&Person{
		Hair:"Brown",
		HasGlass:true,
		Age:10,
		Legs:[]string{"Left", "Right"},
	}
	param:=&InitContextParam{
		t,
		expected,
		"",
		toXml,
	}

	//init param
	context:=initContext(param)
	//actual
	person:=&Person{
	}

	err:=binder.Bind(person,context)

	//check error must nil?
	test.Nil(t,err)
}

func toJson(t *testing.T,v interface{}) string{
	b,err:=json.Marshal(v)
	test.Nil(t,err)
	return string(b)
}
func toXml(t *testing.T,v interface{}) string{
	b,err:=xml.Marshal(v)
	test.Nil(t,err)
	//t.Log("xml:",string(b))
	return string(b)
}

func toDefault(t *testing.T,v interface{}) string{
	return ""
}

func format(b string) io.ReadCloser{
	s := strings.NewReader(b)
	r := ioutil.NopCloser(s)
	r.Close()
	return r
}