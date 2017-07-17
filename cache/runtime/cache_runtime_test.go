package runtime

import (
	"testing"
	"time"
	"github.com/devfeel/dotweb/test"
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
	cache:=NewTestRuntimeCache()
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
	cache:=NewTestRuntimeCache()
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

func NewTestRuntimeCache() *RuntimeCache {
	cache := RuntimeCache{items: make(map[string]*RuntimeItem), gcInterval: DefaultTestGCInterval}
	go cache.gc()
	return &cache
}