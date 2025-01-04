package lockx

import (
	"fmt"
	"sync"
)

// Mutex represents a mutex lock for a specific resource with reference counting.
type Mutex struct {
	sync.Mutex
	refCount int32
}

// KeyMutex manages multiple mutex locks based on resource keys.
type KeyMutex struct {
	mu    sync.Mutex
	locks map[string]*Mutex
}

// NewKeyMutex creates and returns a new instance of KeyMutex.
func NewKeyMutex() *KeyMutex {
	return &KeyMutex{
		locks: make(map[string]*Mutex),
	}
}

// Lock locks the specified key's mutex, ensuring exclusive access.
// 1) Acquire (or create) a *Mutex with refCount++
// 2) Then call m.Lock() on that mutex
func (km *KeyMutex) Lock(key string) {
	m := km.acquireLock(key)
	m.Lock()
}

// Unlock releases the specified key's mutex.
// 1) Actually unlock the underlying mutex
// 2) Decrement refCount, if <= 0 then remove it from the map
func (km *KeyMutex) Unlock(key string) {
	km.releaseLock(key)
}

// acquireLock finds (or creates) the Mutex for `key` and increments its refCount.
func (km *KeyMutex) acquireLock(key string) *Mutex {
	km.mu.Lock()
	defer km.mu.Unlock()

	m, ok := km.locks[key]
	if !ok {
		m = &Mutex{}
		km.locks[key] = m
	}
	m.refCount++ // protected by km.mu, so no atomic needed
	return m
}

// releaseLock looks up the Mutex for `key`, unlocks it, decrements refCount,
// and removes it from the map if no one else is using it. Panics if not found.
func (km *KeyMutex) releaseLock(key string) {
	km.mu.Lock()
	m, ok := km.locks[key]
	if !ok {
		km.mu.Unlock()
		panic(fmt.Sprintf("Unlock attempted on a key that was not locked: %s", key))
	}

	// Step1: Unlock the actual mutex
	m.Unlock()

	// Step2: Decrement refCount, remove if 0
	m.refCount--
	if m.refCount <= 0 {
		delete(km.locks, key)
	}
	km.mu.Unlock()
}
