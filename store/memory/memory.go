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

// Decrement the value of an item in the cache.
func (s *Store) Decrement(key string, args ...int64) (int64, error) {
	i, err := s.item(key)
	if err != nil {
		i = store.Item{
			Object: int64(0),
		}
	}

	iv := int64(1)
	if len(args) > 0 {
		iv = args[0]
	}

	ev, ok := i.Object.(int64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int64", key)
	}

	nv := ev - iv
	i.Object = nv
	s.setItem(key, i)

	return nv, nil
}

// Flush remove all items from the cache.
func (s *Store) Flush() error {
	s.mu.Lock()
	s.items = make(map[string]store.Item)
	s.mu.Unlock()
	return nil
}

// Result will retrieve a item from the cache and stores the
// result in the value pointed to by v.
func (s *Store) Result(key string, v interface{}) error {
	i, err := s.item(key)
	if err != nil {
		return err
	}

	buf, err := store.Marshal(i.Object)
	if err != nil {
		return err
	}

	if err := store.Unmarshal(buf, &v); err != nil {
		return err
	}

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

// Increment the value of an item in the cache.
func (s *Store) Increment(key string, args ...int64) (int64, error) {
	i, err := s.item(key)
	if err != nil {
		i = store.Item{
			Object: int64(0),
		}
	}

	iv := int64(1)
	if len(args) > 0 {
		iv = args[0]
	}

	ev, ok := i.Object.(int64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int64", key)
	}

	nv := ev + iv
	i.Object = nv
	s.setItem(key, i)

	return nv, nil
}

// Number will retrieve the number set by decrement and increment methods from the cache.
func (s *Store) Number(key string) (int64, error) {
	i, err := s.item(key)
	if err != nil {
		return 0, err
	}

	ev, ok := i.Object.(int64)
	if !ok {
		return 0, fmt.Errorf("The value for %s is not an int64", key)
	}
	return ev, nil
}

// Remember will retrieve a item from the cache, but also store a default value if the requested item
// doesn't exists or is empty.
func (s *Store) Remember(key string, expiration time.Duration, fn store.RememberFunc) (interface{}, error) {
	v, err := s.Get(key)
	if v != nil && err == nil {
		return v, nil
	}

	v = fn()

	if err := s.Set(key, v, expiration); err != nil {
		return nil, err
	}

	return v, nil
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

// Set will store a item in the cache.
func (s *Store) Set(key string, value interface{}, expiration time.Duration) error {
	s.setItem(key, store.Item{
		Object:     value,
		Expiration: expiration,
	})
	return nil
}