package cache

import (
	"time"

	"github.com/frozzare/go-cache/store"
)

// Cache represents the cache struct.
type Cache struct {
	store store.Store
}

// New with the given options.
func New(store store.Store) *Cache {
	return &Cache{store}
}

// Flush remove all items from the cache.
func (c *Cache) Flush() error {
	return c.store.Flush()
}

// Get will retrive a item from the cache.
func (c *Cache) Get(key string) (interface{}, error) {
	return c.store.Get(key)
}

// Remove will remove a item from the cache.
func (c *Cache) Remove(key string) error {
	return c.store.Remove(key)
}

// Result will retrieve a item from the cache and stores the
// result in the value pointed to by v.
func (c *Cache) Result(key string, v interface{}) error {
	return c.store.Result(key, v)
}

// Set will store a item in the cache.
func (c *Cache) Set(key string, value interface{}, expiration ...time.Duration) error {
	e := time.Duration(0)

	if len(expiration) > 0 {
		e = expiration[0]
	}

	return c.store.Set(key, value, e)
}
