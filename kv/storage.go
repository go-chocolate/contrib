package kv

import (
	"context"
	"fmt"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, expiration ...time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

func New(c Config) (Storage, error) {
	if c.Driver == "" {
		c.Driver = REDIS
	}
	driver := drivers[c.Driver]
	if driver == nil {
		return nil, fmt.Errorf("unknown kv storage driver: %s", c.Driver)
	}
	return driver(c.Option)
}

func MustNew(c Config) Storage {
	storage, err := New(c)
	if err != nil {
		panic(err)
	}
	return storage
}

type prefixStorage struct {
	prefix  string
	storage Storage
}

func Prefix(prefix string, storage Storage) Storage {
	return &prefixStorage{prefix: prefix, storage: storage}
}

func (s *prefixStorage) Get(ctx context.Context, key string) ([]byte, error) {
	return s.storage.Get(ctx, s.prefix+key)
}

func (s *prefixStorage) Set(ctx context.Context, key string, val []byte, expiration ...time.Duration) error {
	return s.storage.Set(ctx, s.prefix+key, val, expiration...)
}

func (s *prefixStorage) Del(ctx context.Context, keys ...string) error {
	var ks = make([]string, 0, len(keys))
	for i := range keys {
		ks = append(ks, s.prefix+keys[i])
	}
	return s.storage.Del(ctx, ks...)
}
