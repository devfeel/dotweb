package session

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

var sessionStatePool sync.Pool

func init() {
	sessionStatePool = sync.Pool{
		New: func() interface{} {
			return &SessionState{}
		},
	}
}

//session state
type SessionState struct {
	sessionId    string                      //session id
	timeAccessed time.Time                   //last access time
	values       map[interface{}]interface{} //session store
	lock         *sync.RWMutex
	store        SessionStore
}

func NewSessionState(store SessionStore, sessionId string, values map[interface{}]interface{}) *SessionState {
	state := sessionStatePool.Get().(*SessionState)
	state.reset(store, sessionId, values, time.Now())
	return state
}

// Set set key-value to current state
func (state *SessionState) reset(store SessionStore, sessionId string, values map[interface{}]interface{}, accessTime time.Time) {
	state.values = values
	state.sessionId = sessionId
	state.timeAccessed = accessTime
	state.store = store
	state.lock = new(sync.RWMutex)
}

// Set set key-value to current state
func (state *SessionState) Set(key, value interface{}) error {
	state.lock.Lock()
	defer state.lock.Unlock()
	state.values[key] = value
	return state.store.SessionUpdate(state)

}

// Get get value by key in current state
func (state *SessionState) Get(key interface{}) interface{} {
	state.lock.RLock()
	defer state.lock.RUnlock()
	if v, ok := state.values[key]; ok {
		return v
	}
	return nil
}

// Get get value as string by key in current state
func (state *SessionState) GetString(key interface{}) string {
	v := state.Get(key)
	return fmt.Sprint(v)
}

// Get get value as int by key in current state
func (state *SessionState) GetInt(key interface{}) int {
	v, _ := strconv.Atoi(state.GetString(key))
	return v
}

// Get get value as int64 by key in current state
func (state *SessionState) GetInt64(key interface{}) int64 {
	v, _ := strconv.ParseInt(state.GetString(key), 10, 64)
	return v
}

// Remove remove value by key in current state
func (state *SessionState) Remove(key interface{}) error {
	state.lock.Lock()
	defer state.lock.Unlock()
	delete(state.values, key)
	return nil
}

// Clear clear all values in current store
func (state *SessionState) Clear() error {
	state.lock.Lock()
	defer state.lock.Unlock()
	state.values = make(map[interface{}]interface{})
	return nil
}

// SessionID get this id in current state
func (state *SessionState) SessionID() string {
	return state.sessionId
}

// Count get all item's count in current state
func (state *SessionState) Count() int {
	return len(state.values)
}
