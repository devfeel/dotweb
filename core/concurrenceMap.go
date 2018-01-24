package core

import (
	"fmt"
	"sync"
)

type (
	ReadonlyMap interface {
		Get(key string) (value interface{}, exists bool)
		GetString(key string) string
		GetInt(key string) int
		GetUInt64(key string) uint64
		Exists(key string) bool
		Len() int
	}
	ConcurrenceMap interface {
		Get(key string) (value interface{}, exists bool)
		GetString(key string) string
		GetInt(key string) int
		GetUInt64(key string) uint64
		Exists(key string) bool
		GetCurrentMap() map[string]interface{}
		Len() int
		Set(key string, value interface{}) error
		Remove(key string)
		Once(key string) (value interface{}, exists bool)
	}
)

//自带锁，并发安全的Map
type ItemMap struct {
	innerMap map[string]interface{}
	*sync.RWMutex
}

func NewItemMap() *ItemMap {
	return &ItemMap{
		innerMap: make(map[string]interface{}),
		RWMutex:  new(sync.RWMutex),
	}
}

func NewConcurrenceMap() ConcurrenceMap {
	return &ItemMap{
		innerMap: make(map[string]interface{}),
		RWMutex:  new(sync.RWMutex),
	}
}

func NewReadonlyMap() ReadonlyMap {
	return &ItemMap{
		innerMap: make(map[string]interface{}),
		RWMutex:  new(sync.RWMutex),
	}
}

/*
* 以key、value置入AppContext
 */
func (ctx *ItemMap) Set(key string, value interface{}) error {
	ctx.Lock()
	ctx.innerMap[key] = value
	ctx.Unlock()
	return nil
}

/*
* 读取指定key在AppContext中的内容
 */
func (ctx *ItemMap) Get(key string) (value interface{}, exists bool) {
	ctx.RLock()
	value, exists = ctx.innerMap[key]
	ctx.RUnlock()
	return value, exists
}

//remove item by gived key
//if not exists key, do nothing...
func (ctx *ItemMap) Remove(key string) {
	ctx.Lock()
	delete(ctx.innerMap, key)
	ctx.Unlock()
}

//get item by gived key, and remove it
//only can be read once, it will be locked
func (ctx *ItemMap) Once(key string) (value interface{}, exists bool) {
	ctx.Lock()
	defer ctx.Unlock()
	value, exists = ctx.innerMap[key]
	if exists {
		delete(ctx.innerMap, key)
	}
	return value, exists
}

/*
* 读取指定key在AppContext中的内容，以string格式输出
 */
func (ctx *ItemMap) GetString(key string) string {
	value, exists := ctx.Get(key)
	if !exists {
		return ""
	}
	return fmt.Sprint(value)
}

/*
* 读取指定key在AppContext中的内容，以int格式输出
 */
func (ctx *ItemMap) GetInt(key string) int {
	value, exists := ctx.Get(key)
	if !exists {
		return 0
	}
	return value.(int)
}

/*
* 读取指定key在AppContext中的内容，以int格式输出
 */
func (ctx *ItemMap) GetUInt64(key string) uint64 {
	value, exists := ctx.Get(key)
	if !exists {
		return 0
	}
	return value.(uint64)
}

//check exists key
func (ctx *ItemMap) Exists(key string) bool {
	_, exists := ctx.innerMap[key]
	return exists
}

//get current map, returns map[string]interface{}
func (ctx *ItemMap) GetCurrentMap() map[string]interface{} {
	return ctx.innerMap
}

//get context length
func (ctx *ItemMap) Len() int {
	return len(ctx.innerMap)
}
