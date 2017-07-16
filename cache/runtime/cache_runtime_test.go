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
		test.Equal(t,false,value)
	}(cache,t)

	time.Sleep(5*time.Second)
}

func NewTestRuntimeCache() *RuntimeCache {
	cache := RuntimeCache{items: make(map[string]*RuntimeItem), gcInterval: DefaultTestGCInterval}
	go cache.gc()
	return &cache
}