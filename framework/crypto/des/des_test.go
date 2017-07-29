package des

import (
	"testing"
	"github.com/devfeel/dotweb/test"
	"fmt"
)

//

func Test_ECBEncrypt_1(t *testing.T) {
	key := []byte("01234567")
	origData := []byte("cphpbb@hotmail.com")
	b, e := ECBEncrypt(origData, key)
	if e != nil {
		t.Error(e)
	} else {
		t.Logf("%x\n", b)
	}

	test.Equal(t,"a5296e4c525693a3892bbe31e1ed630121f26338ce9aa280",fmt.Sprintf("%x",b))
}

//ECBDecrypt方法有bug，这个方法会报空指针
func Test_ECBDecrypt_1(t *testing.T) {
	hextext := []byte("a5296e4c525693a3892bbe31e1ed630121f26338ce9aa280")
	key := []byte("01234567")
	b, e := ECBDecrypt(hextext, key)
	if e != nil {
		t.Error(e)
	} else {
		t.Logf("%x\n", b)
	}

	//test.Equal(t,"a5296e4c525693a3892bbe31e1ed630121f26338ce9aa280",fmt.Sprintf("%x",b))
}

func Test_PKCS5Padding_1(t *testing.T) {}

func Test_PKCS5UnPadding_1(t *testing.T) {}

func Test_TripleEcbDesDecrypt_1(t *testing.T) {}

func Test_TripleEcbDesEncrypt_1(t *testing.T) {}
