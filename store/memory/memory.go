package memory

import (
	"fmt"
	"sync"
	"time"

	"github.com/frozzare/go-cache/store"
)

// Store represents the redis cache store.
type Store struct {
	items map[string]store.Item
	mu    sync.RWMutex
}

// NewStore will create a new redis store with the given options.
func NewStore() store.Store {
	return &Store{
		items: make(map[string]store.Item),
		mu:    sync.RWMutex{},
	}
}

func (s *Store) setItem(key string, i store.Item) {
	s.mu.Lock()
	s.items[key] = i
	s.mu.Unlock()
}

func (s *Store) item(key string) (store.Item, error) {
	s.mu.Lock()
	i, ok := s.items[key]
	if !ok || i.Expired() {
		s.mu.Unlock()
		return store.Item{}, fmt.Errorf("Item %s not found", key)
	}
	s.mu.Unlock()
	return i, nil
}

// Close store.
func (s *Store) Close() error {
	return nil
}

// Flush remove all items from the cache.
func (s *Store) Flush() error {
	s.mu.Lock()
	s.items = make(map[string]store.Item)
	s.mu.Unlock()
	return nil
}

// Get will retrieve a item from the cache.
func (s *Store) Get(key string) (interface{}, error) {
	i, err := s.item(key)
	if err != nil {
		return nil, err
	}

	return i.Object, nil
}

// Remove will remove a item from the cache.
func (s *Store) Remove(key string) error {
	s.mu.Lock()

	i, ok := s.items[key]
	if !ok || i.Expired() {
		s.mu.Unlock()
		return fmt.Errorf("Item %s not found", key)
	}

	delete(s.items, key)
	s.mu.Unlock()

	return nil
}

// Result will retrieve a item from the cache and stores the
// result in the value pointed to by value.
func (s *Store) Result(key string, value interface{}) error {
	i, err := s.item(key)
	if err != nil {
		return err
	}

	buf, err := store.Marshal(i.Object)
	if err != nil {
		return err
	}

	return store.Unmarshal(buf, &value)
}

// Set will store a item in the cache.
func (s *Store) Set(key string, value interface{}, expiration time.Duration) error {
	s.setItem(key, store.Item{
		Object:     value,
		Expiration: expiration,
	})
	return nil
}
