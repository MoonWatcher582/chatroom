package common

import "sync"

type AtomicMap struct {
	atmap map[string]*AtomicConn
	lock  *sync.Mutex
}

func NewAtomicMap() *AtomicMap {
	return &AtomicMap{atmap: make(map[string]*AtomicConn), lock: &sync.Mutex{}}
}

func (a *AtomicMap) Get(key string) *AtomicConn {
	a.lock.Lock()
	val := a.atmap[key]
	a.lock.Unlock()
	return val
}

func (a *AtomicMap) Set(key string, val *AtomicConn) {
	a.lock.Lock()
	a.atmap[key] = val
	a.lock.Unlock()
}

func (a *AtomicMap) Remove(key string) {
	a.lock.Lock()
	delete(a.atmap, key)
	a.lock.Unlock()
}

func (a *AtomicMap) Keys() []string {
	keys := make([]string, len(a.atmap))

	i := 0
	for k := range a.atmap {
		keys[i] = k
		i++
	}
	return keys
}
