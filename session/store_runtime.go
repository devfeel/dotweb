package session

import (
	"container/list"
	"sync"
	"time"
)

// MemProvider Implement the provider interface
type RuntimeStore struct {
	lock        *sync.RWMutex            // locker
	sessions    map[string]*list.Element // map in memory
	list        *list.List               // for gc
	maxlifetime int64
}

func NewRuntimeStore(config *StoreConfig) *RuntimeStore {
	return &RuntimeStore{
		lock:        new(sync.RWMutex),
		sessions:    make(map[string]*list.Element),
		list:        new(list.List),
		maxlifetime: config.Maxlifetime,
	}
}

// SessionRead get session state by sessionId
func (store *RuntimeStore) SessionRead(sessionId string) (*SessionState, error) {
	store.lock.RLock()
	if element, ok := store.sessions[sessionId]; ok {
		go store.SessionAccess(sessionId)
		store.lock.RUnlock()
		return element.Value.(*SessionState), nil
	}
	store.lock.RUnlock()

	//if sessionId of state not exist, create a new state
	state := NewSessionState(store, sessionId, make(map[interface{}]interface{}))
	store.lock.Lock()
	element := store.list.PushFront(state)
	store.sessions[sessionId] = element
	store.lock.Unlock()
	return state, nil
}

// SessionExist check session state exist by sessionId
func (store *RuntimeStore) SessionExist(sessionId string) bool {
	store.lock.RLock()
	defer store.lock.RUnlock()
	if _, ok := store.sessions[sessionId]; ok {
		return true
	}
	return false
}

//SessionUpdate update session state in store
func (store *RuntimeStore) SessionUpdate(state *SessionState) error {
	store.lock.RLock()
	if element, ok := store.sessions[state.sessionId]; ok { //state has exist
		go store.SessionAccess(state.sessionId)
		store.lock.RUnlock()
		element.Value.(*SessionState).values = state.values //only assist update whole session state
		return nil
	}
	store.lock.RUnlock()

	//if sessionId of state not exist, create a new state
	new_state := NewSessionState(store, state.sessionId, state.values)
	store.lock.Lock()
	new_element := store.list.PushFront(new_state)
	store.sessions[state.sessionId] = new_element
	store.lock.Unlock()
	return nil
}

// SessionRemove delete session state in store
func (store *RuntimeStore) SessionRemove(sessionId string) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	if element, ok := store.sessions[sessionId]; ok {
		delete(store.sessions, sessionId)
		store.list.Remove(element)
		return nil
	}
	return nil
}

// SessionGC clean expired session stores in memory session
func (store *RuntimeStore) SessionGC() int {
	num := 0
	store.lock.RLock()
	for {
		element := store.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*SessionState).timeAccessed.Unix() + store.maxlifetime) < time.Now().Unix() {
			store.lock.RUnlock()
			store.lock.Lock()
			store.list.Remove(element)
			delete(store.sessions, element.Value.(*SessionState).SessionID())
			num += 1
			store.lock.Unlock()
			store.lock.RLock()
		} else {
			break
		}
	}
	store.lock.RUnlock()
	return num
}

// SessionAll get count number of memory session
func (store *RuntimeStore) SessionCount() int {
	return store.list.Len()
}

// SessionAccess expand time of session store by id in memory session
func (store *RuntimeStore) SessionAccess(sessionId string) error {
	store.lock.Lock()
	defer store.lock.Unlock()
	if element, ok := store.sessions[sessionId]; ok {
		element.Value.(*SessionState).timeAccessed = time.Now()
		store.list.MoveToFront(element)
		return nil
	}
	return nil
}
