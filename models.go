package main

import sync "sync"

type consumer struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	Name       string `json:"name"`
}

type safeMap[K comparable, V any] struct {
	data map[K]V
	mu   sync.RWMutex
}

func (s *safeMap[K, V]) get(k K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, e := s.data[k]
	return v, e
}

func (s *safeMap[K, V]) set(k K, v V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[k] = v
}

func (m *safeMap[K, V]) delete(k K) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, k)
}

func (s *safeMap[K, V]) forEach(block func(k K, v V) bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for k, v := range s.data {
		if block(k, v) {
			break
		}
	}
}

func newSafeMap[K comparable, V any]() *safeMap[K, V] {
	return &safeMap[K, V]{data: make(map[K]V)}
}
