package bolt

import (
	"fmt"
	"os"
	"strconv"
	"time"

	boltdb "github.com/boltdb/bolt"
	"github.com/frozzare/go-cache/store"
)

var (
	bucket    = []byte("store")
	bucketTTL = []byte("store_ttl")
)

// Options is a alias for bolt options.
type Options = boltdb.Options

// Store represents the redis cache store.
type Store struct {
	db *boltdb.DB
}

// NewStore will create a new redis store with the given options.
func NewStore(name string, permission os.FileMode, opts *Options) (store.Store, error) {
	db, err := boltdb.Open(name, permission, opts)

	return &Store{
		db: db,
	}, err
}

// Close store.
func (s *Store) Close() error {
	return s.db.Close()
}

// Flush remove all items from the cache.
func (s *Store) Flush() error {
	return s.db.Update(func(tx *boltdb.Tx) error {
		if err := tx.DeleteBucket(bucketTTL); err == nil {
			_, err := tx.CreateBucket(bucketTTL)
			return err
		}

		if err := tx.DeleteBucket(bucket); err != nil {
			return err
		}

		_, err := tx.CreateBucket(bucket)
		return err
	})
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
	return s.db.Update(func(tx *boltdb.Tx) error {
		b := tx.Bucket(bucketTTL)
		b.Delete(bucketTTL)
		b = tx.Bucket(bucket)
		return b.Delete([]byte(key))
	})
}

// Result will retrieve a item from the cache and stores the
// result in the value pointed to by value.
func (s *Store) Result(key string, value interface{}) error {
	return s.db.View(func(tx *boltdb.Tx) error {
		b := tx.Bucket(bucketTTL)

		exp := b.Get([]byte(key))
		if len(exp) > 0 {
			i, err := strconv.ParseInt(string(exp), 10, 64)
			if err != nil {
				return err
			}

			if time.Now().UnixNano() > i && i != 0 {
				return fmt.Errorf("Item %s not found", key)
			}
		}

		b = tx.Bucket(bucket)
		buf := b.Get([]byte(key))
		return store.Unmarshal(buf, value)
	})
}

// Set will store a item in the cache.
func (s *Store) Set(key string, value interface{}, expiration time.Duration) error {
	return s.db.Update(func(tx *boltdb.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucket)

		if err != nil {
			return err
		}

		buf, err := store.Marshal(value)

		if err != nil {
			return err
		}

		if err := b.Put([]byte(key), buf); err != nil {
			return err
		}

		b, err = tx.CreateBucketIfNotExists(bucketTTL)

		if err != nil {
			return err
		}

		return b.Put([]byte(key), []byte(fmt.Sprintf("%d", expiration)))
	})
}
