package config

import (
	"gopkg.in/yaml.v2"
)

// UnmarshalYaml decodes the first document found within the in byte slice
// and assigns decoded values into the out value.
// For example:
//
//     type T struct {
//         F int `yaml:"a,omitempty"`
//         B int
//     }
//     var t T
//     yaml.Unmarshal([]byte("a: 1\nb: 2"), &t)
func UnmarshalYaml(content []byte, v interface{}) error {
	return yaml.Unmarshal(content, v)
}

// MarshalYaml Marshal serializes the value provided into a YAML document.
// For example:
//
//     type T struct {
//         F int "a,omitempty"
//         B int
//     }
//     yaml.Marshal(&T{B: 2}) // Returns "b: 2\n"
//     yaml.Marshal(&T{F: 1}} // Returns "a: 1\nb: 0\n"
func MarshalYaml(v interface{}) (out []byte, err error) {
	return yaml.Marshal(v)
}

// MarshalYamlString returns the Ymal encoding string format of v.
func MarshalYamlString(v interface{}) (out string) {
	marshal, err := yaml.Marshal(v)
	if err != nil {
		return ""
	}
	return string(marshal)
}
