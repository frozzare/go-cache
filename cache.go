package cache

import (
	"time"

	"github.com/frozzare/go-cache/store"
)

// Cache represents the cache struct.
type Cache struct {
	store store.Store
}

// Option function.
type Option func(*Cache)

// New with the given options.
func New(options ...Option) *Cache {
	var v Cache
	for _, o := range options {
		o(&v)
	}
	return &v
}

// Redis option will create a new Redis store with the given options.
func Redis(r *store.RedisOptions) Option {
	return func(c *Cache) {
		c.store = store.NewRedisStore(r)
	}
}

// Decrement the value of an item in the cache.
func (c *Cache) Decrement(key string, args ...int64) (int64, error) {
	return c.store.Decrement(key, args...)
}

// Flush remove all items from the cache.
func (c *Cache) Flush() error {
	return c.store.Flush()
}

// Get will retrive a item from the cache.
func (c *Cache) Get(key string) (interface{}, error) {
	return c.store.Get(key)
}

// Increment the value of an item in the cache.
func (c *Cache) Increment(key string, args ...int64) (int64, error) {
	return c.store.Increment(key, args...)
}

// Number will retrieve the number set by decrement and increment methods from the cache.
func (c *Cache) Number(key string) (int64, error) {
	return c.store.Number(key)
}

// Remember will retrive a item from the cache, but also store a default value if the requested item
// doesn't exists or is empty.
func (c *Cache) Remember(key string, exp time.Duration, fn store.RememberFunc) (interface{}, error) {
	return c.store.Remember(key, exp, fn)
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
