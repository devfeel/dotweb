package config

import (
	"encoding/xml"
)

// UnmarshalXML parses the XML-encoded data and stores the result in
// the value pointed to by v, which must be an arbitrary struct,
// slice, or string. Well-formed data that does not fit into v is
// discarded.
func UnmarshalXML(content []byte, v interface{}) error {
	return xml.Unmarshal(content, v)
}

// MarshalXML returns the XML encoding of v.
func MarshalXML(v interface{}) (out []byte, err error) {
	return xml.Marshal(v)
}

// MarshalXMLString returns the XML encoding string format of v.
func MarshalXMLString(v interface{}) (out string) {
	marshal, err := xml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}
