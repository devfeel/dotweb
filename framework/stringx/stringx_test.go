package stringx

import (
	"github.com/devfeel/dotweb/test"
	"testing"
)

func TestCompletionRight(t *testing.T) {
	content := "ab"
	flag := "cbc"
	length := 6
	wantResult := "abcbcc"
	test.Equal(t, wantResult, CompletionRight(content, flag, length))
}

func TestCompletionLeft(t *testing.T) {
	content := "ab"
	flag := "cbc"
	length := 6
	wantResult := "cbccab"
	test.Equal(t, wantResult, CompletionLeft(content, flag, length))
}
