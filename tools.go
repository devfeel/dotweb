package dotweb

import (
	"encoding/json"
	"fmt"
)

type Tools struct {
}

func (t *Tools) PrettyJson(data interface{}) string {
	by, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return fmt.Sprint(data)
	}
	return string(by)
}
