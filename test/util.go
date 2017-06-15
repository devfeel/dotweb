package test

import (
	"encoding/json"
	"encoding/xml"
	"testing"
)

// ToJson
func ToJson(t *testing.T, v interface{}) string {
	b, err := json.Marshal(v)
	Nil(t, err)
	return string(b)
}

// ToXML
func ToXML(t *testing.T, v interface{}) string {
	b, err := xml.Marshal(v)
	Nil(t, err)
	//t.Log("xml:",string(b))
	return string(b)
}

//ToDefault
func ToDefault(t *testing.T, v interface{}) string {
	return ""
}
