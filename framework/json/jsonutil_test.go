package jsonutil

import (
	"testing"
)

// 以下是功能测试

type (
	colorGroup struct {
		ID     int
		Name   string
		Colors []string
	}
	mymap map[string]interface{}
)

func Test_GetJsonString_1(t *testing.T) {
	g := newGroup()
	json := GetJsonString(g)
	t.Log(json)
}

func Test_Marshal_1(t *testing.T) {
	g := newGroup()
	json, err := Marshal(g)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(json)
	}

}

func Test_Unmarshal_1(t *testing.T) {
	var group colorGroup
	g := newGroup()
	json := GetJsonString(g)
	err := Unmarshal(json, &group)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(group)
	}

}

func newGroup() colorGroup {
	group := colorGroup{
		ID:     1,
		Name:   "Reds",
		Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	}
	return group
}
