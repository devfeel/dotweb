package core

import (
	"testing"
)

// 以下为功能测试

func Test_AddRequestCount_1(t *testing.T) {
	var num uint64 = 1
	var count uint64
	for i := 0; i < 100; i++ {
		count = GlobalState.AddRequestCount(num)
	}
	t.Log("TotalRequestCount:", count)
}

func Test_AddRequestCount_2(t *testing.T) {
	var num uint64 = 1
	var count uint64
	for i := 0; i < 100; i++ {
		count = GlobalState.AddRequestCount(num)
		num++
	}
	t.Log("TotalRequestCount:", count)
}

func Test_AddErrorCount_1(t *testing.T) {
	var num, count uint64
	for i := 0; i < 100; i++ {
		num = 1
		count = GlobalState.AddErrorCount(num)
	}
	t.Log("TotalErrorCount:", count)
}

func Test_AddErrorCount_2(t *testing.T) {
	var num, count uint64
	for i := 0; i < 100; i++ {
		count = GlobalState.AddErrorCount(num)
		num++
	}
	t.Log("TotalErrorCount:", count)
}

// 以下是性能测试

//基准测试
func Benchmark_AddErrorCount_1(b *testing.B) {
	var num uint64 = 1
	for i := 0; i < b.N; i++ {
		GlobalState.AddErrorCount(num)
	}
}

// 测试并发效率
func Benchmark_AddErrorCount_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var num uint64 = 1
		for pb.Next() {
			GlobalState.AddErrorCount(num)
		}
	})
}

//基准测试
func Benchmark_AddRequestCount_1(b *testing.B) {
	var num uint64 = 1
	for i := 0; i < b.N; i++ {
		GlobalState.AddRequestCount(num)
	}
}

// 测试并发效率
func Benchmark_AddRequestCount_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var num uint64 = 1
		for pb.Next() {
			GlobalState.AddRequestCount(num)
		}
	})
}
