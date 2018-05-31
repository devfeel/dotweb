package cache

import "testing"

var runtimeCache Cache
var key string
var val []byte

func init(){
	runtimeCache = NewRuntimeCache()
	key = "abc"
	val = []byte("def")
}


func DoSet(cache Cache){
	expire := 60 // expire in 60 seconds
	cache.Set(key, val, int64(expire))
}

func DoGet(cache Cache){
	cache.Get(key)
}


func BenchmarkTestSet(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoSet(runtimeCache)
	}
}

func BenchmarkTestGet(b *testing.B) {
	DoSet(runtimeCache)
	for i := 0; i < b.N; i++ {
		DoGet(runtimeCache)
	}
}