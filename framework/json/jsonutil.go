package jsonutil

import (
	"encoding/json"
)

//将传入对象转换为json字符串
func GetJsonString(obj interface{}) string {
	resByte, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(resByte)
}

//将传入对象转换为json字符串
func Marshal(v interface{}) (string, error) {
	resByte, err := json.Marshal(v)
	if err != nil {
		return "", err
	} else {
		return string(resByte), nil
	}
}

//将传入的json字符串转换为对象
func Unmarshal(jsonstring string, v interface{}) error {
	return json.Unmarshal([]byte(jsonstring), v)
}
