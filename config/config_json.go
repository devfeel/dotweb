package config

import "encoding/json"

func fromJson(content []byte,v interface{}) error  {
	return json.Unmarshal(content, v)
}
