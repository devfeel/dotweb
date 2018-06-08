package base64x

import "encoding/base64"


// EncodeString encode string use base64 StdEncoding
func EncodeString(source string) string{
	return base64.StdEncoding.EncodeToString([]byte(source))
}

// DecodeString deencode string use base64 StdEncoding
func DecodeString(source string) (string, error){
	dst, err:= base64.StdEncoding.DecodeString(source)
	if err != nil{
		return "", err
	}else{
		return string(dst), nil
	}
}