package lockx

import (
	"fmt"
	"sync"
	"sync/atomic"
)

//////////////////////////////////
// KeyMutex: Mutex-based locking //
//////////////////////////////////

// Mutex represents a mutex lock for a specific resource with reference counting.
type Mutex struct {
	sync.Mutex       // Mutex for exclusive access to the resource.
	refCount   int32 // Reference count to track how many goroutines hold this lock.
}

// KeyMutex manages multiple mutex locks based on resource keys.
type KeyMutex struct {
	locks sync.Map // Concurrent-safe map for managing multiple resource locks.
}

// NewKeyMutex creates and returns a new instance of KeyMutex.
func NewKeyMutex() *KeyMutex {
	return &KeyMutex{}
}

// Lock locks the specified key's mutex, ensuring exclusive access.
func (km *KeyMutex) Lock(key string) {
	mutex := km.getOrCreateLock(key)
	atomic.AddInt32(&mutex.refCount, 1)
	mutex.Lock()
}

// Unlock releases the specified key's mutex.
func (km *KeyMutex) Unlock(key string) {
	if m, ok := km.getLock(key); ok {
		m.Unlock()
		if atomic.AddInt32(&m.refCount, -1) == 0 {
			km.locks.Delete(key)
		}
		return
	}
	// Optionally, we could panic or log an error here if an unlock is attempted on a non-existent key.
	panic(fmt.Sprintf("Unlock attempted on a key that was not locked: %s", key))

}

// getOrCreateLock retrieves an existing Mutex or creates a new one if it doesn't exist.
func (km *KeyMutex) getLock(key string) (*Mutex, bool) {
	if m, ok := km.locks.Load(key); ok {
		return m.(*Mutex), ok
	}
	return nil, false
}

// getOrCreateLock retrieves an existing Mutex or creates a new one if it doesn't exist.
func (km *KeyMutex) getOrCreateLock(key string) *Mutex {
	actual, _ := km.locks.LoadOrStore(key, &Mutex{})
	return actual.(*Mutex)
}
