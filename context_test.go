package dotweb

import (
	"testing"
	"github.com/devfeel/dotweb/test"
)

type Animal struct{
	Hair     string
	HasMouth bool
}

func TestWriteJsonpBlob(t *testing.T) {
	param := &InitContextParam{
		t,
		&Animal{},
		"",
		test.ToDefault,
	}

	//init param
	context := initResponseContext(param)

	excepted:=&Animal{
		"Black",
		true,
	}

	context.WriteJsonp("jsonCallBack",excepted)

	t.Log(string(context.response.body))
}