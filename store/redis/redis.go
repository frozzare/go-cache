package redis

import (
	"time"

	"github.com/frozzare/go-cache/store"
	goredis "github.com/go-redis/redis"
)

// Options is a alias for the options for the redis client.
type Options = goredis.Options

// Store represents the redis cache store.
type Store struct {
	client *goredis.Client
}

// NewStore will create a new redis store with the given options.
func NewStore(o *Options) store.Store {
	if o == nil {
		o = &Options{}
	}

	if len(o.Addr) == 0 {
		o.Addr = "localhost:6379"
	}

	return &Store{
		client: goredis.NewClient(o),
	}
}

// Close store.
func (s *Store) Close() error {
	return s.client.Close()
}

// Flush remove all items from the cache.
func (s *Store) Flush() error {
	return s.client.FlushDB().Err()
}

// Result will retrieve a item from the cache and stores the
// result in the value pointed to by value.
func (s *Store) Result(key string, value interface{}) error {
	b, err := s.client.Get(key).Bytes()
	if len(b) == 0 {
		return nil
	}

	if err != nil {
		return err
	}

	return store.Unmarshal(b, value)
}

// Get will retrieve a item from the cache.
func (s *Store) Get(key string) (interface{}, error) {
	var i *store.Item

	if err := s.Result(key, &i); err != nil {
		return nil, err
	}

	if i.Object != nil {
		return i.Object, nil
	}

	var i2 interface{}

	if err := s.Result(key, &i2); err != nil {
		return nil, err
	}

	return i2, nil
}

// Remove will remove a item from the cache.
func (s *Store) Remove(key string) error {
	return s.client.Del(key).Err()
}

// Set will store a item in the cache.
func (s *Store) Set(key string, value interface{}, expiration time.Duration) error {
	var b []byte
	var err error

	switch value.(type) {
	case []byte:
		b = value.([]byte)
	default:
		b, err = store.Marshal(value)
		if err != nil {
			return err
		}
	}

	return s.client.Set(key, b, expiration).Err()
}
