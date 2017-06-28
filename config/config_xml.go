package config

import (
	"encoding/xml"
)

func fromXml(content []byte,v interface{}) error  {
	return xml.Unmarshal(content, v)
}