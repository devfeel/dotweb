package cryptos

import (
	"github.com/devfeel/dotweb/test"
	"testing"
)

//

func Test_GetMd5String_1(t *testing.T) {
	str := "123456789"
	md5str := GetMd5String(str)
	t.Log("GetMd5String:", md5str)
	test.Equal(t, "25f9e794323b453885f5181f1b624d0b", md5str)
}

//这个测试用例没按照功能实现所说按照长度生成对应长度字符串？
func Test_GetRandString_1(t *testing.T) {
	for i := 4; i < 9; i++ {
		randStr := GetRandString(i)

		test.Equal(t, i, len(randStr))

		if len(randStr) != i {
			t.Error("GetRandString: length:", i, "randStr-len:", len(randStr))
		} else {
			t.Log("GetRandString: length-", i, "randStr-", randStr)
		}
	}
}
