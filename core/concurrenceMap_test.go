package core

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

var ic ConcurrenceMap
var keys []string

func init() {
	ic = NewConcurrenceMap()
	for i := 0; i < 10000000; i++ {
		keys = append(keys, time.Now().String())
	}
	fmt.Println("len of keys ", len(keys))

}

func TestItemContext_Get_Set(t *testing.T) {

	t.Log(ic.Set("foo", "bar"))
	t.Log(ic.Get("foo"))
	t.Log(ic.Exists("foo"))

	t.Log(ic.Get("none"))
	t.Log(ic.Exists("none"))
}

func TestItemContext_Get_Once(t *testing.T) {
	ic.Set("foo", "bar")
	t.Log(ic.Once("foo"))
	t.Log(ic.Get("foo"))
}

func TestItemContext_Remove(t *testing.T) {

	ic.Set("foo", "bar")
	ic.Set("foo1", "bar1")
	t.Log(len(ic.GetCurrentMap()))
	ic.Remove("foo")
	t.Log(ic.GetString("foo"))
}

func TestItemContext_Current(t *testing.T) {
	lock := &sync.Mutex{}
	j := 0
	for i := 0; i < 9; i++ {
		go func() {
			lock.Lock()
			fmt.Println("go", j)
			j++
			v := "bar" + strconv.Itoa(j)
			fmt.Println(v)
			ic.Set(strconv.Itoa(j), v)
			lock.Unlock()
		}()
	}

	time.Sleep(3 * time.Second)

	t.Log(ic.GetCurrentMap())

}

//性能测试

//基准测试
func BenchmarkItemContext_Set_1(b *testing.B) {
	var num uint64 = 1
	for i := 0; i < b.N; i++ {
		ic.Set(string(num), num)
	}
}

//并发效率
func BenchmarkItemContext_Set_Parallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		var num uint64 = 1
		for pb.Next() {
			ic.Set(string(num), num)
		}
	})
}

//基准测试
func BenchmarkItemContext_Get_1(b *testing.B) {
	ic.Set("foo", "bar")
	for i := 0; i < b.N; i++ {
		ic.Get("foo")
	}
}

//并发效率
func BenchmarkItemContext_Get_Parallel(b *testing.B) {
	ic.Set("foo", "bar")
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ic.Get("foo")
		}
	})

}
