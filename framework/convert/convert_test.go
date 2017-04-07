package convert

import (
	"testing"
)

//功能测试

func Test_String2Bytes_1(t *testing.T) {
	str := "0123456789"
	b := String2Bytes(str)
	t.Log(str, " to byte: ", b)
}
