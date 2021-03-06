package core

import (
	"fmt"
	"sync"
	"time"
)

type (
	// ReadonlyMap only support readonly method for map
	ReadonlyMap interface {
		Get(key string) (value interface{}, exists bool)
		GetString(key string) string
		GetTimeDuration(key string) time.Duration
		GetInt(key string) int
		GetUInt64(key string) uint64
		Exists(key string) bool
		Len() int
	}

	// ReadonlyMap support concurrence for map
	ConcurrenceMap interface {
		Get(key string) (value interface{}, exists bool)
		GetString(key string) string
		GetTimeDuration(key string) time.Duration
		GetInt(key string) int
		GetUInt64(key string) uint64
		Exists(key string) bool
		GetCurrentMap() map[string]interface{}
		Len() int
		Set(key string, value interface{})
		Remove(key string)
		Once(key string) (value interface{}, exists bool)
	}
)

// ItemMap concurrence map
type ItemMap struct {
	innerMap map[string]interface{}
	*sync.RWMutex
}

// NewItemMap create new ItemMap
func NewItemMap() *ItemMap {
	return &ItemMap{
		innerMap: make(map[string]interface{}),
		RWMutex:  new(sync.RWMutex),
	}
}

// NewConcurrenceMap create new ConcurrenceMap
func NewConcurrenceMap() ConcurrenceMap {
	return &ItemMap{
		innerMap: make(map[string]interface{}),
		RWMutex:  new(sync.RWMutex),
	}
}

// NewReadonlyMap create new ReadonlyMap
func NewReadonlyMap() ReadonlyMap {
	return &ItemMap{
		innerMap: make(map[string]interface{}),
		RWMutex:  new(sync.RWMutex),
	}
}

// Set put key, value into ItemMap
func (ctx *ItemMap) Set(key string, value interface{}) {
	ctx.Lock()
	ctx.innerMap[key] = value
	ctx.Unlock()
}

// Get returns value of specified key
func (ctx *ItemMap) Get(key string) (value interface{}, exists bool) {
	ctx.RLock()
	value, exists = ctx.innerMap[key]
	ctx.RUnlock()
	return value, exists
}

// Remove remove item by gived key
// if not exists key, do nothing...
func (ctx *ItemMap) Remove(key string) {
	ctx.Lock()
	delete(ctx.innerMap, key)
	ctx.Unlock()
}

// Once get item by gived key, and remove it
// only can be read once, it will be locked
func (ctx *ItemMap) Once(key string) (value interface{}, exists bool) {
	ctx.Lock()
	defer ctx.Unlock()
	value, exists = ctx.innerMap[key]
	if exists {
		delete(ctx.innerMap, key)
	}
	return value, exists
}

// GetString returns value as string specified by key
// return empty string if key not exists
func (ctx *ItemMap) GetString(key string) string {
	value, exists := ctx.Get(key)
	if !exists {
		return ""
	}
	return fmt.Sprint(value)
}

// GetInt returns value as int specified by key
// return 0 if key not exists
func (ctx *ItemMap) GetInt(key string) int {
	value, exists := ctx.Get(key)
	if !exists {
		return 0
	}
	return value.(int)
}

// GetUInt64 returns value as uint64 specified by key
// return 0 if key not exists or value cannot be converted to int64
func (ctx *ItemMap) GetUInt64(key string) uint64 {
	value, exists := ctx.Get(key)
	if !exists {
		return 0
	}
	return value.(uint64)
}

// GetTimeDuration returns value as time.Duration specified by key
// return 0 if key not exists or value cannot be converted to time.Duration
func (ctx *ItemMap) GetTimeDuration(key string) time.Duration {
	timeDuration, err := time.ParseDuration(ctx.GetString(key))
	if err != nil {
		return 0
	}
	return timeDuration
}

// Exists check exists key
func (ctx *ItemMap) Exists(key string) bool {
	_, exists := ctx.innerMap[key]
	return exists
}

// GetCurrentMap get current map, returns map[string]interface{}
func (ctx *ItemMap) GetCurrentMap() map[string]interface{} {
	return ctx.innerMap
}

// Len get context length
func (ctx *ItemMap) Len() int {
	return len(ctx.innerMap)
}
