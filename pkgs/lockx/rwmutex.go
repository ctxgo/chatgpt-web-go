package lockx

import (
	"fmt"
	"sync"
	"sync/atomic"
)

////////////////////////////////////
// KeyRWMutex: RWMutex-based lock  //
////////////////////////////////////

// RWMutex represents a read/write lock for a specific resource with reference counting.
type RWMutex struct {
	sync.RWMutex       // RWMutex for shared (read) or exclusive (write) access.
	refCount     int32 // Reference count to track how many goroutines hold this lock.
}

// KeyRWMutex manages multiple read/write locks based on resource keys.
type KeyRWMutex struct {
	locks sync.Map // Concurrent-safe map for managing multiple resource locks.
}

// NewKeyRWMutex creates and returns a new instance of KeyRWMutex.
func NewKeyRWMutex() *KeyRWMutex {
	return &KeyRWMutex{}
}

func (km *KeyRWMutex) lock(key string, lockFunc func(*RWMutex)) {
	rwMutex := km.getOrCreateLock(key)
	atomic.AddInt32(&rwMutex.refCount, 1)
	lockFunc(rwMutex)
}

// Lock locks the specified key's mutex for writing, ensuring exclusive access.
func (km *KeyRWMutex) Lock(key string) {
	km.lock(key, func(r *RWMutex) {
		r.Lock()
	})
}

// RLock locks the specified key's mutex for reading, allowing shared access.
func (km *KeyRWMutex) RLock(key string) {
	km.lock(key, func(r *RWMutex) {
		r.RLock()
	})
}

func (km *KeyRWMutex) unlock(key string, unlockFunc func(k string)) {
	if m, ok := km.getLock(key); ok {
		unlockFunc(key)
		if atomic.AddInt32(&m.refCount, -1) == 0 {
			km.locks.Delete(key)
		}
		return
	}
	// Optionally, we could panic or log an error here if an unlock is attempted on a non-existent key.
	panic(fmt.Sprintf("Unlock attempted on a key that was not locked: %s", key))

}

// Unlock releases the specified key's write lock.
func (km *KeyRWMutex) Unlock(key string) {
	km.unlock(key, km.Unlock)

}

// RUnlock releases the specified key's read lock.
func (km *KeyRWMutex) RUnlock(key string) {
	km.unlock(key, km.RUnlock)
}

// getOrCreateLock retrieves an existing RWMutex or creates a new one if it doesn't exist.
func (km *KeyRWMutex) getLock(key string) (*RWMutex, bool) {
	actual, ok := km.locks.Load(key)
	return actual.(*RWMutex), ok
}

// getOrCreateLock retrieves an existing RWMutex or creates a new one if it doesn't exist.
func (km *KeyRWMutex) getOrCreateLock(key string) *RWMutex {
	actual, _ := km.locks.LoadOrStore(key, &RWMutex{})
	return actual.(*RWMutex)
}
