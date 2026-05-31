package service

import (
	"sync"
)

type KeyLocker struct {
	mu    sync.Mutex
	locks map[string]*sync.Mutex
}

func NewKeyLocker() *KeyLocker {
	return &KeyLocker{
		locks: make(map[string]*sync.Mutex),
	}
}

func (kl *KeyLocker) Lock(key string) func() {
	kl.mu.Lock()
	l, ok := kl.locks[key]
	if !ok {
		l = &sync.Mutex{}
		kl.locks[key] = l
	}
	kl.mu.Unlock()

	l.Lock()
	return func() {
		l.Unlock()
	}
}
