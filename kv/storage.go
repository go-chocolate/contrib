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
