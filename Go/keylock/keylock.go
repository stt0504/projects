package keylock

import (
	"sort"
)

type Mutex struct {
	ch chan int
}

func (m *Mutex) Lock() {
	m.ch <- 0
}

func (m *Mutex) Unlock() {
	<-m.ch
}

func NewMutex() *Mutex {
	return &Mutex{make(chan int, 1)}
}

type KeyLock struct {
	mu    Mutex
	locks map[string]chan struct{}
}

func New() *KeyLock {
	return &KeyLock{
		locks: make(map[string]chan struct{}),
		mu:    *NewMutex(),
	}
}

func (kl *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	keysCopy := append([]string(nil), keys...)
	sort.Strings(keysCopy)

	kl.mu.Lock()
	for _, key := range keysCopy {
		if _, exists := kl.locks[key]; !exists {
			kl.locks[key] = make(chan struct{}, 1)
		}
	}
	kl.mu.Unlock()

	for i, key := range keysCopy {
		kl.mu.Lock()
		lock := kl.locks[key]
		kl.mu.Unlock()
		select {
		case lock <- struct{}{}:
		case <-cancel:
			kl.mu.Lock()
			for _, k := range keysCopy[:i] {
				<-kl.locks[k]
			}
			kl.mu.Unlock()
			return true, nil
		}
	}

	unlock = func() {
		kl.mu.Lock()
		for _, key := range keysCopy {
			<-kl.locks[key]
		}
		kl.mu.Unlock()
	}

	return false, unlock
}
