package base64x

import "testing"

func TestEncodeString(t *testing.T) {
	source := "welcome to dotweb!"
	t.Log(EncodeString(source))
}

func TestDecodeString(t *testing.T) {
	source :=  "welcome to dotweb!"
	encode := EncodeString(source)

	dst, err := DecodeString(encode)
	if err != nil{
		t.Error("TestDecodeString error", err)
	}else{
		t.Log("TestDecodeString success", dst, source)
	}
}
