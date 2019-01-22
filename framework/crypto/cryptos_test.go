package cryptos

import (
	"testing"

	"github.com/devfeel/dotweb/test"
)

//

func Test_GetMd5String_1(t *testing.T) {
	str := "123456789"
	md5str := GetMd5String(str)
	t.Log("GetMd5String:", md5str)
	test.Equal(t, "25f9e794323b453885f5181f1b624d0b", md5str)
}
