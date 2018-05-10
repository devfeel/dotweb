package runtime

import (
	"testing"
	"time"
	"github.com/devfeel/dotweb/test"
	"sync"
)

const (
	DefaultTestGCInterval = 2

	TEST_CACHE_KEY = "joe"
	TEST_CACHE_VALUE = "zou"
	//int value
	TEST_CACHE_INT_VALUE = 1

	//int64 value
	TEST_CACHE_INT64_VALUE = 1
)

func TestRuntimeCache_Get(t *testing.T) {
	cache:=NewRuntimeCache()
	cache.Set(TEST_CACHE_KEY,TEST_CACHE_VALUE,5)
	//check value
	go func(cache *RuntimeCache,t *testing.T) {
		time.Sleep(4*time.Second)
		value,err:=cache.Get(TEST_CACHE_KEY)

		test.Nil(t,err)
		test.Equal(t,TEST_CACHE_VALUE,value)
	}(cache,t)

	//check expired
	go func(cache *RuntimeCache,t *testing.T) {
		time.Sleep(5*time.Second)
		value,err:=cache.Exists(TEST_CACHE_KEY)

		test.Nil(t,err)
		test.Equal(t,true,value)
	}(cache,t)

	time.Sleep(5*time.Second)
}


func TestRuntimeCache_GetInt(t *testing.T) {
	testRuntimeCache(t,TEST_CACHE_INT_VALUE,func(cache *RuntimeCache,key string)(interface{}, error){
		return cache.GetInt(key)
	})
}


func TestRuntimeCache_GetInt64(t *testing.T) {
	testRuntimeCache(t,TEST_CACHE_INT64_VALUE,func(cache *RuntimeCache,key string)(interface{}, error){
		return cache.GetInt64(key)
	})
}

func TestRuntimeCache_GetString(t *testing.T) {
	testRuntimeCache(t,TEST_CACHE_VALUE,func(cache *RuntimeCache,key string)(interface{}, error){
		return cache.GetString(key)
	})
}

func testRuntimeCache(t *testing.T,insertValue interface{},f func(cache *RuntimeCache,key string)(interface{}, error)) {
	cache:=NewRuntimeCache()
	cache.Set(TEST_CACHE_KEY,insertValue,5)
	//check value
	go func(cache *RuntimeCache,t *testing.T) {
		time.Sleep(4*time.Second)
		value,err:=f(cache,TEST_CACHE_KEY)

		test.Nil(t,err)
		test.Equal(t,insertValue,value)
	}(cache,t)

	time.Sleep(5*time.Second)
}

func TestRuntimeCache_Delete(t *testing.T) {
	cache:=NewRuntimeCache()
	cache.Set(TEST_CACHE_KEY,TEST_CACHE_VALUE,5)

	value,e:=cache.Get(TEST_CACHE_KEY)

	test.Nil(t,e)
	test.Equal(t,TEST_CACHE_VALUE,value)

	cache.Delete(TEST_CACHE_KEY)

	value,e=cache.Get(TEST_CACHE_KEY)
	test.Nil(t,e)
	test.Nil(t,value)
}

func TestRuntimeCache_ClearAll(t *testing.T) {
	cache:=NewRuntimeCache()
	cache.Set(TEST_CACHE_KEY,TEST_CACHE_VALUE,5)
	cache.Set("2",TEST_CACHE_VALUE,5)
	cache.Set("3",TEST_CACHE_VALUE,5)

	val2, err := cache.GetString("2")
	if err != nil{
		t.Error(err)
	}
	test.Equal(t,TEST_CACHE_VALUE, val2)

	cache.ClearAll()
	exists2, err := cache.Exists("2")
	if err != nil{
		t.Error(err)
	}
	if exists2{
		t.Error("exists 2 but need not exists")
	}
}

func TestRuntimeCache_Incr(t *testing.T) {
	cache:=NewRuntimeCache()
	var wg sync.WaitGroup
	wg.Add(2)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Incr(TEST_CACHE_KEY)
		}

		wg.Add(-1)
	}(cache)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Incr(TEST_CACHE_KEY)
		}
		wg.Add(-1)
	}(cache)

	wg.Wait()

	value,e:=cache.GetInt(TEST_CACHE_KEY)
	test.Nil(t,e)

	test.Equal(t,100,value)
}

func TestRuntimeCache_Decr(t *testing.T) {
	cache:=NewRuntimeCache()
	var wg sync.WaitGroup
	wg.Add(2)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Decr(TEST_CACHE_KEY)
		}

		wg.Add(-1)
	}(cache)

	go func(cache *RuntimeCache) {
		for i := 0; i < 50; i++ {
			cache.Decr(TEST_CACHE_KEY)
		}
		wg.Add(-1)
	}(cache)

	wg.Wait()

	value,e:=cache.GetInt(TEST_CACHE_KEY)
	test.Nil(t,e)

	test.Equal(t,-100,value)
}