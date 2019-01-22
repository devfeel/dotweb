package des

import (
	"fmt"
	"testing"

	"github.com/devfeel/dotweb/test"
)

//

func Test_ECBEncrypt_1(t *testing.T) {
	key := []byte("01234567")
	origData := []byte("dotweb@devfeel")
	b, e := ECBEncrypt(origData, key)
	if e != nil {
		t.Error(e)
	} else {
		t.Logf("%x\n", b)
	}

	test.Equal(t, "72f9f187eafe43478f9eb3dd49ef7b43", fmt.Sprintf("%x", b))
}

 func Test_ECBDecrypt_1(t *testing.T) {
 	key := []byte("01234567")
	origData := []byte("dotweb@devfeel")
	b1, e1 := ECBEncrypt(origData, key)
	if e1 != nil {
	 t.Error(e1)
	}
 	b, e := ECBDecrypt(b1, key)
 	if e != nil {
 		t.Error(e)
 	} else {
 		t.Logf("%x\n", b)
 	}

 	test.Equal(t, "dotweb@devfeel", string(b))
 }

func Test_PKCS5Padding_1(t *testing.T) {}

func Test_PKCS5UnPadding_1(t *testing.T) {}

func Test_TripleEcbDesDecrypt_1(t *testing.T) {}

func Test_TripleEcbDesEncrypt_1(t *testing.T) {}
