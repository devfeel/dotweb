package cryptos

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/devfeel/dotweb/framework/convert"
)

// GetMd5String compute the md5 sum as string
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetRandString returns randominzed string with given length
func GetRandString(length int) string {
	bytes := []byte("0123456789abcdefghijklmnopqrstuvwxyz")
	result := make([]byte, 0, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return convert.Bytes2String(result)
}
