package config

import "encoding/json"

// UnmarshalJSON parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func UnmarshalJSON(content []byte, v interface{}) error {
	return json.Unmarshal(content, v)
}

// MarshalJSON returns the JSON encoding of v.
func MarshalJSON(v interface{}) (out []byte, err error) {
	return json.Marshal(v)
}

// MarshalJSONString returns the JSON encoding string format of v.
func MarshalJSONString(v interface{}) (out string) {
	marshal, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}
