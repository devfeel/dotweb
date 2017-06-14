package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"strings"
	"io"
)

type Person struct {
	Hair string
	HasGlass bool
	Age int
	Legs []string
}

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
	context:=initContext(t,expected)
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

//common init context
func initContext(t *testing.T,v interface{}) *HttpContext {
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
	header["Content-Type"]=[]string{"application/json"}
	context.request.Header=header

	body:=format(t,v)
	context.request.Request.Body=body


	return context
}

func format(t *testing.T,v interface{}) io.ReadCloser{
	b,err:=json.Marshal(v)
	test.Nil(t,err)

	s := strings.NewReader(string(b))
	r := ioutil.NopCloser(s)
	r.Close()
	return r
}