package jsonutil

import (
	"encoding/json"
)

// GetJsonString marshals the object as string
func GetJsonString(obj interface{}) string {
	resByte, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(resByte)
}

// Marshal marshals the value as string
func Marshal(v interface{}) (string, error) {
	resByte, err := json.Marshal(v)
	if err != nil {
		return "", err
	} else {
		return string(resByte), nil
	}
}

// Unmarshal converts the jsonstring into value
func Unmarshal(jsonstring string, v interface{}) error {
	return json.Unmarshal([]byte(jsonstring), v)
}
