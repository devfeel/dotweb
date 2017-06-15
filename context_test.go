package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
	"encoding/json"
	"fmt"
)

type Animal struct{
	Hair     string
	HasMouth bool
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
	context.WriteJsonp(callback,exceptedObject)

	//check result

	//header
	contentType:=context.response.header.Get(HeaderContentType)
	test.Equal(t,contentType,MIMEApplicationJavaScriptCharsetUTF8)

	//body
	body:=string(context.response.body)

	animalJson,err:=json.Marshal(exceptedObject)
	test.Nil(t,err)

	excepted:=fmt.Sprint(callback,"(",string(animalJson),");")

	test.Equal(t,body,excepted)
}