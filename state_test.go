package dotweb

import (
	"testing"
)

// 以下为功能测试

func Test_AddRequestCount_1(t *testing.T) {
	var num uint64 = 0
	for i := 1; i < 100; i++ {
		num = GlobalState.AddRequestCount(uint64(i))
	}
	t.Log("TotalRequestCount:", num)

}

func Test_AddErrorCount_1(t *testing.T) {
	var num uint64 = 0
	for i := 1; i < 100; i++ {
		num = GlobalState.AddErrorCount(uint64(i))
	}
	t.Log("TotalErrorCount:", num)
}
