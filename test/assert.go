package test

import (
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func Contains(t *testing.T, expected, actual interface{}) {
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		formattedLog(t, "value type not match %T (expected)\n\n\t != %T (actual)", expected, actual)
	}
	switch e := expected.(type) {
	case string:
		a := actual.(string)
		if !strings.Contains(a, e) {
			formattedLog(t, "   %v (expected)\n\n\t!= %v (actual)",
				expected, actual)
			t.FailNow()
		}
	default:
		formattedLog(t, "unsupported type %T(expected)", expected)
		t.FailNow()
	}
}

func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		formattedLog(t, "   %v (expected)\n\n\t!= %v (actual)",
			expected, actual)
		t.FailNow()
	}
}

func NotEqual(t *testing.T, expected, actual interface{}) {
	if reflect.DeepEqual(expected, actual) {
		formattedLog(t, "value should not equal %#v", actual)
		t.FailNow()
	}
}

func Nil(t *testing.T, object interface{}) {
	if !isNil(object) {
		formattedLog(t, "   <nil> (expected)\n\n\t!= %#v (actual)", object)
		t.FailNow()
	}
}

func NotNil(t *testing.T, object interface{}) {
	if isNil(object) {
		formattedLog(t, "Expected value not to be <nil>", object)
		t.FailNow()
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	kind := value.Kind()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}

	return false
}

func formattedLog(t *testing.T, fmt string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	file = filepath.Base(file)
	targs := make([]interface{}, len(args)+2)
	targs[0] = file
	targs[1] = line
	copy(targs[2:], args)
	t.Logf("\033[31m%s:%d:\n\n\t"+fmt+"\033[39m\n\n", targs...)
}
