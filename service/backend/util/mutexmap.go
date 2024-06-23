package util

import (
	"fmt"
	"sync"
)

type MutexMap struct {
	mutexInner sync.Mutex
	inner      map[interface{}]*mutexMapEntry
}

type mutexMapEntry struct {
	parent *MutexMap
	mutex  sync.Mutex
	count  int
	key    interface{}
}

type Unlocker interface {
	Unlock()
}

func MutexMapNew() *MutexMap {
	return &MutexMap{inner: make(map[interface{}]*mutexMapEntry)}
}

func (m *MutexMap) Lock(key interface{}) Unlocker {
	m.mutexInner.Lock()
	e, ok := m.inner[key]
	if !ok {
		e = &mutexMapEntry{parent: m, key: key}
		m.inner[key] = e
	}
	e.count++
	m.mutexInner.Unlock()

	e.mutex.Lock()

	return e
}

func (me *mutexMapEntry) Unlock() {
	m := me.parent

	m.mutexInner.Lock()
	e, ok := m.inner[me.key]
	if !ok {
		m.mutexInner.Unlock()
		panic(fmt.Errorf("Unlock requested for key=%v but no entry found", me.key))
	}
	e.count--
	if e.count < 1 {
		delete(m.inner, me.key)
	}
	m.mutexInner.Unlock()

	e.mutex.Unlock()
}
