package lockx

import (
	"fmt"
	"sync"
)

// RWMutex is a reference-counted read/write lock for a single resource.
type RWMutex struct {
	sync.RWMutex
	refCount int
}

// KeyRWMutex manages multiple RWMutex locks keyed by string.
//
// A single global mutex (km.mu) protects the `locks` map and ensures correct
// refCount increments/decrements. When refCount hits 0, the lock is removed.
type KeyRWMutex struct {
	mu    sync.Mutex
	locks map[string]*RWMutex
}

// NewKeyRWMutex instantiates a KeyRWMutex with an internal map.
func NewKeyRWMutex() *KeyRWMutex {
	return &KeyRWMutex{
		locks: make(map[string]*RWMutex),
	}
}

// ---------------------------
//       Public Methods
// ---------------------------

// Lock obtains an exclusive (write) lock for the given key.
func (km *KeyRWMutex) Lock(key string) {
	km.doLock(key, (*RWMutex).Lock)
}

// RLock obtains a shared (read) lock for the given key.
func (km *KeyRWMutex) RLock(key string) {
	km.doLock(key, (*RWMutex).RLock)
}

// Unlock releases a write lock for the given key.
func (km *KeyRWMutex) Unlock(key string) {
	km.doUnlock(key, func(r *RWMutex) {
		r.Unlock()
	})
}

// RUnlock releases a read lock for the given key.
func (km *KeyRWMutex) RUnlock(key string) {
	km.doUnlock(key, (*RWMutex).RUnlock)
}

// ---------------------------
//       Internal Logic
// ---------------------------

// doLock retrieves or creates an RWMutex for the key, increments its refCount,
// then calls the provided lockMethod (which can be Lock or RLock).
func (km *KeyRWMutex) doLock(key string, lockMethod func(*RWMutex)) {
	km.mu.Lock()
	rw, ok := km.locks[key]
	if !ok {
		rw = &RWMutex{}
		km.locks[key] = rw
	}
	rw.refCount++
	km.mu.Unlock()

	// Actually lock it (either Lock or RLock)
	lockMethod(rw)
}

// doUnlock looks up the RWMutex for the key, calls the provided unlockMethod
// (which can be Unlock or RUnlock), decrements refCount, and if refCount == 0
// removes it from the map.
func (km *KeyRWMutex) doUnlock(key string, unlockMethod func(*RWMutex)) {
	km.mu.Lock()
	rw, ok := km.locks[key]
	if !ok {
		km.mu.Unlock()
		panic(fmt.Sprintf("Unlock called for unknown key %q", key))
	}

	// Perform the actual unlock (either Unlock or RUnlock)
	unlockMethod(rw)

	// Decrement refCount and remove if zero
	rw.refCount--
	if rw.refCount <= 0 {
		delete(km.locks, key)
	}
	km.mu.Unlock()
}
