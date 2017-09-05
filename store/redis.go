package store

import (
	"time"

	goredis "github.com/go-redis/redis"
)

// RedisOptions is a alias for the options for the redis client.
type RedisOptions = goredis.Options

// RedisStore represents the redis cache store.
type RedisStore struct {
	client *goredis.Client
}

// NewRedisStore will create a new redis store with the given options.
func NewRedisStore(o *RedisOptions) *RedisStore {
	if o == nil {
		o = &RedisOptions{
			Addr: "localhost:6379",
		}
	}

	return &RedisStore{
		client: goredis.NewClient(o),
	}
}

// Decrement the value of an item in the cache.
func (s *RedisStore) Decrement(key string, args ...int64) (int64, error) {
	value := int64(1)
	if len(args) > 0 {
		value = args[0]
	}

	return s.client.DecrBy(key, value).Result()
}

// Flush remove all items from the cache.
func (s *RedisStore) Flush() error {
	return s.client.FlushDB().Err()
}

// Result will retrieve a item from the cache and stores the
// result in the value pointed to by v.
func (s *RedisStore) Result(key string, v interface{}) error {
	b, err := s.client.Get(key).Bytes()
	if len(b) == 0 {
		return nil
	}

	if err != nil {
		return err
	}

	if err := Unmarshal(b, v); err != nil {
		return err
	}

	return nil
}

// Get will retrieve a item from the cache.
func (s *RedisStore) Get(key string) (interface{}, error) {
	var i *Item

	if err := s.Result(key, &i); err != nil {
		return nil, err
	}

	return i.Object, nil
}

// Increment the value of an item in the cache.
func (s *RedisStore) Increment(key string, args ...int64) (int64, error) {
	value := int64(1)
	if len(args) > 0 {
		value = args[0]
	}

	return s.client.IncrBy(key, value).Result()
}

// Number will retrieve the number set by decrement and increment methods from the cache.
func (s *RedisStore) Number(key string) (int64, error) {
	return s.client.Get(key).Int64()
}

// Remember will retrieve a item from the cache, but also store a default value if the requested item
// doesn't exists or is empty.
func (s *RedisStore) Remember(key string, expiration time.Duration, fn RememberFunc) (interface{}, error) {
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
func (s *RedisStore) Remove(key string) error {
	return s.client.Del(key).Err()
}

// Set will store a item in the cache.
func (s *RedisStore) Set(key string, value interface{}, expiration time.Duration) error {
	var b []byte
	var err error

	switch value.(type) {
	case []byte:
		b = value.([]byte)
	default:
		b, err = Marshal(value)
		if err != nil {
			return err
		}
	}

	return s.client.Set(key, b, expiration).Err()
}
