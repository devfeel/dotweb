package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
	"encoding/json"
	"fmt"
	"net/http"
)

type Animal struct{
	Hair     string
	HasMouth bool
}

//normal write
func TestWrite(t *testing.T) {
	param := &InitContextParam{
		t,
		&Animal{},
		"",
		test.ToDefault,
	}

	//init param
	context := initResponseContext(param)

	exceptedObject:=&Animal{
		"Black",
		true,
	}

	animalJson,err:=json.Marshal(exceptedObject)
	test.Nil(t,err)

	//call function
	status:=http.StatusNotFound
	_,contextErr:=context.Write(status,animalJson)
	test.Nil(t,contextErr)

	//check result

	//header
	contentType:=context.response.header.Get(HeaderContentType)
	//因writer中的header方法调用过http.Header默认设置
	test.Equal(t,CharsetUTF8,contentType)
	test.Equal(t,status,context.response.Status)

	//body
	body:=string(context.response.body)

	test.Equal(t,body,string(animalJson))
}

//normal jsonp
func TestWriteJsonp(t *testing.T) {
	param := &InitContextParam{
		t,
		&Animal{},
		"",
		test.ToDefault,
	}

	//init param
	context := initResponseContext(param)

	exceptedObject:=&Animal{
		"Black",
		true,
	}

	callback:="jsonCallBack"

	//call function
	_,err:=context.WriteJsonp(callback,exceptedObject)
	test.Nil(t,err)

	//check result

	//header
	contentType:=context.response.header.Get(HeaderContentType)
	test.Equal(t,MIMEApplicationJavaScriptCharsetUTF8,contentType)
	test.Equal(t,defaultHttpCode,context.response.Status)

	//body
	body:=string(context.response.body)

	animalJson,err:=json.Marshal(exceptedObject)
	test.Nil(t,err)

	excepted:=fmt.Sprint(callback,"(",string(animalJson),");")

	test.Equal(t,body,excepted)
}