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

func Test_GetRandString(t *testing.T) {
	randStr := GetRandString(12)
	rand1 := GetRandString(12)
	rand2 := GetRandString(12)
	rand3 := GetRandString(12)
	if rand1 == rand2 || rand2 == rand3 || rand1 == rand3 {
		t.Error("rand result is same")
	} else {
		t.Log("GetRandString:", randStr)
		test.Equal(t, 12, len(randStr))
	}
}
