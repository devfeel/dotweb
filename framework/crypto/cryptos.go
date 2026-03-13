package cryptos

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/devfeel/dotweb/framework/convert"
	"math/big"
)

// GetMd5String compute the md5 sum as string
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetRandString returns randominzed string with given length
func GetRandString(length int) string {
	var container string
	var str = "0123456789abcdefghijklmnopqrstuvwxyz"
	b := bytes.NewBufferString(str)
	len := b.Len()
	bigInt := big.NewInt(int64(len))
	for i := 0; i < length; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
