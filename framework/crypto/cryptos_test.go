package cryptos

import (
	"testing"
)

//

func Test_GetMd5String_1(t *testing.T) {
	str := "123456789"
	md5str := GetMd5String(str)
	t.Log("GetMd5String:", md5str)
}

func Test_GetUUID_1(t *testing.T) {
	uuid := GetUUID()
	t.Log("GetUUID:", uuid)
}

func Test_GetRandString_1(t *testing.T) {
	for i := 4; i < 9; i++ {
		randStr := GetRandString(i)
		t.Log("GetRandString: length-", i, "randStr-", randStr)
		if len(randStr) != i {
			t.Error("GetRandString: length:", i, "randStr-len:", len(randStr))
		}
	}
}
