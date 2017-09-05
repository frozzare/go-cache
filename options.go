package cache

import "github.com/frozzare/go-cache/store/redis"

// Option function.
type Option func(*Cache)

// Redis option will create a new Redis store with the given options.
func Redis(r *redis.Options) Option {
	return func(c *Cache) {
		c.store = redis.NewStore(r)
	}
}
