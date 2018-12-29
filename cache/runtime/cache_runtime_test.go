package runtime

import (
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/devfeel/dotweb/test"
)

const (
	DefaultTestGCInterval = 2

	TEST_Cache_Key   = "joe"
	TEST_Cache_Value = "zou"
	// int value
	TEST_Cache_Int_Value = 1

	// int64 value
	TEST_Cache_Int64_Value = int64(1)
)

func TestRuntimeCache_Get(t *testing.T) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, TEST_Cache_Value, 2)
	var wg sync.WaitGroup

	// check value
	wg.Add(1)
	go func(cache *RuntimeCache, t *testing.T) {
		time.Sleep(1 * time.Second)
		value, err := cache.Get(TEST_Cache_Key)

		test.Nil(t, err)
		test.Equal(t, TEST_Cache_Value, value)
		wg.Done()
	}(cache, t)

	// check expired
	wg.Add(1)
	go func(cache *RuntimeCache, t *testing.T) {
		time.Sleep(2 * time.Second)
		value, err := cache.Exists(TEST_Cache_Key)

		test.Nil(t, err)
		test.Equal(t, false, value)
		wg.Done()
	}(cache, t)

	wg.Wait()
}

func TestRuntimeCache_GetInt(t *testing.T) {
	testRuntimeCache(t, TEST_Cache_Int_Value, func(cache *RuntimeCache, key string) (interface{}, error) {
		return cache.GetInt(key)
	})
}

func TestRuntimeCache_GetInt64(t *testing.T) {
	testRuntimeCache(t, TEST_Cache_Int64_Value, func(cache *RuntimeCache, key string) (interface{}, error) {
		return cache.GetInt64(key)
	})
}

func TestRuntimeCache_GetString(t *testing.T) {
	testRuntimeCache(t, TEST_Cache_Value, func(cache *RuntimeCache, key string) (interface{}, error) {
		return cache.GetString(key)
	})
}

func testRuntimeCache(t *testing.T, insertValue interface{}, f func(cache *RuntimeCache, key string) (interface{}, error)) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, insertValue, 2)
	var wg sync.WaitGroup

	// check value
	wg.Add(1)
	go func(cache *RuntimeCache, t *testing.T) {
		time.Sleep(1 * time.Second)
		value, err := f(cache, TEST_Cache_Key)

		test.Nil(t, err)
		test.Equal(t, insertValue, value)
		wg.Done()
	}(cache, t)
	time.Sleep(2 * time.Second)
	wg.Wait()
}

func TestRuntimeCache_Delete(t *testing.T) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, TEST_Cache_Value, 2)

	value, e := cache.Get(TEST_Cache_Key)

	test.Nil(t, e)
	test.Equal(t, TEST_Cache_Value, value)

	cache.Delete(TEST_Cache_Key)

	value, e = cache.Get(TEST_Cache_Key)
	test.Nil(t, e)
	test.Nil(t, value)
}

func TestRuntimeCache_ClearAll(t *testing.T) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, TEST_Cache_Value, 2)
	cache.Set("2", TEST_Cache_Value, 2)
	cache.Set("3", TEST_Cache_Value, 2)

	val2, err := cache.GetString("2")
	if err != nil {
		t.Error(err)
	}
	test.Equal(t, TEST_Cache_Value, val2)

	cache.ClearAll()
	exists2, err := cache.Exists("2")
	if err != nil {
		t.Error(err)
	}
	if exists2 {
		t.Error("exists 2 but need not exists")
	}
}

func TestRuntimeCache_Incr(t *testing.T) {
	cache := NewRuntimeCache()
	var wg sync.WaitGroup
	wg.Add(2)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Incr(TEST_Cache_Key)
		}

		wg.Add(-1)
	}(cache)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Incr(TEST_Cache_Key)
		}
		wg.Add(-1)
	}(cache)

	wg.Wait()

	value, e := cache.GetInt(TEST_Cache_Key)
	test.Nil(t, e)

	test.Equal(t, 100, value)
}

func TestRuntimeCache_Decr(t *testing.T) {
	cache := NewRuntimeCache()
	var wg sync.WaitGroup
	wg.Add(2)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Decr(TEST_Cache_Key)
		}

		wg.Add(-1)
	}(cache)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Decr(TEST_Cache_Key)
		}
		wg.Add(-1)
	}(cache)

	wg.Wait()

	value, e := cache.GetInt(TEST_Cache_Key)
	test.Nil(t, e)

	test.Equal(t, -100, value)
}

func BenchmarkTestRuntimeCache_Get(b *testing.B) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, TEST_Cache_Value, 200000)
	for i := 0; i < b.N; i++ {
		cache.Get(TEST_Cache_Key)
	}
}

func BenchmarkTestRuntimeCache_Set(b *testing.B) {
	cache := NewRuntimeCache()
	for i := 0; i < b.N; i++ {
		cache.Set(TEST_Cache_Key + strconv.Itoa(i), TEST_Cache_Value, 0)
	}
}

func TestRuntimeCache_ConcurrentGetSetError(t *testing.T) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, TEST_Cache_Value, 200000)
	for i := 0; i < 10000; i++ {
		go cache.Get(TEST_Cache_Key)
	}

	for i := 0; i < 10000; i++ {
		go cache.Set(TEST_Cache_Key + strconv.Itoa(i), TEST_Cache_Value, 0)
	}
	time.Sleep(time.Minute)
}

func TestRuntimeCache_ConcurrentIncrDecrError(t *testing.T) {
	cache := NewRuntimeCache()
	cache.Set(TEST_Cache_Key, TEST_Cache_Value, 200000)
	for i := 0; i < 10000; i++ {
		go cache.Incr(TEST_Cache_Key + strconv.Itoa(i))
	}

	for i := 0; i < 10000; i++ {
		go cache.Decr(TEST_Cache_Key + strconv.Itoa(i))
	}
	time.Sleep(time.Minute)
}

