package store

import (
	"time"
)

// Store provides a interface to implement cache stores.
type Store interface {
	Flush() error
	Decrement(string, ...int64) (int64, error)
	Get(string) (interface{}, error)
	Number(string) (int64, error)
	Increment(string, ...int64) (int64, error)
	Remember(string, time.Duration, RememberFunc) (interface{}, error)
	Remove(string) error
	Result(string, interface{}) error
	Set(string, interface{}, time.Duration) error
}

// RememberFunc is the function that is used for remember method.
type RememberFunc func() interface{}

// Item represents a item in the cache.
type Item struct {
	Expiration time.Duration
	Object     interface{}
}

// Expired returns true if the item has expired.
func (item Item) Expired() bool {
	if item.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > int64(item.Expiration)
}
