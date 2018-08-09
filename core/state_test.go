package core

import (
	"errors"
	"github.com/devfeel/dotweb/test"
	"sync"
	"testing"
)

// 以下为功能测试

func Test_AddRequestCount_1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go addRequestCount(&wg, 50)

	go addRequestCount(&wg, 60)

	wg.Wait()

	test.Equal(t, uint64(110), GlobalState.TotalRequestCount)
}

func addRequestCount(wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		GlobalState.AddRequestCount("test", 200, 1)
	}
	wg.Add(-1)
}

func Test_CurrentRequestCount(t *testing.T) {
	//var num uint64 = 1
	GlobalState.AddCurrentRequest(1000465)
	t.Log(GlobalState.CurrentRequestCount)
	GlobalState.SubCurrentRequest(2561)
	t.Log(GlobalState.CurrentRequestCount)
}

func Test_AddRequestCount_2(t *testing.T) {
	var num uint64 = 1
	for i := 0; i < 100; i++ {
		GlobalState.AddRequestCount("test", 200, num)
		num++
	}
}

func Test_AddErrorCount_1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	go addErrorCount(&wg, 50)

	go addErrorCount(&wg, 60)

	wg.Wait()

	test.Equal(t, uint64(110), GlobalState.TotalErrorCount)
}

func Test_AddErrorCount_2(t *testing.T) {
	var num, count uint64
	for i := 0; i < 100; i++ {
		count = GlobalState.AddErrorCount("test", errors.New("test error"), num)
		num++
	}
	t.Log("TotalErrorCount:", count)
}

func addErrorCount(wg *sync.WaitGroup, count int) {
	for i := 0; i < count; i++ {
		GlobalState.AddErrorCount("test", errors.New("test error"), 1)
	}
	wg.Add(-1)
}

// 以下是性能测试

//基准测试
func Benchmark_AddErrorCount_1(b *testing.B) {
	var num uint64 = 1
	for i := 0; i < b.N; i++ {
		GlobalState.AddErrorCount("test", errors.New("test error"), num)
	}
}

// 测试并发效率
func Benchmark_AddErrorCount_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var num uint64 = 1
		for pb.Next() {
			GlobalState.AddErrorCount("test", errors.New("test error"), num)
		}
	})
}

//基准测试
func Benchmark_AddRequestCount_1(b *testing.B) {
	var num uint64 = 1
	for i := 0; i < b.N; i++ {
		GlobalState.AddRequestCount("test", 200, num)
	}
}


//基准测试
func Benchmark_AddCurrentRequestCount_1(b *testing.B) {
	var num uint64 = 1
	for i := 0; i < b.N; i++ {
		GlobalState.AddCurrentRequest(num)
	}
}


// 测试并发效率
func Benchmark_AddRequestCount_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var num uint64 = 1
		for pb.Next() {
			GlobalState.AddRequestCount("test", 200, num)
		}
	})
}
