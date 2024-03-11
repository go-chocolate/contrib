package tokenutil

import (
	"context"
	"time"
)

type Storage interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte, expiration ...time.Duration) error
	Del(ctx context.Context, keys ...string) error
}

type memoryStorage struct {
	data map[string][]byte
}

var _ Storage = (*memoryStorage)(nil)

func NewMemoryStorage() Storage {
	return &memoryStorage{data: make(map[string][]byte)}
}

func (s *memoryStorage) Get(ctx context.Context, key string) ([]byte, error) {
	return s.data[key], nil
}

func (s *memoryStorage) Set(ctx context.Context, key string, val []byte, expiration ...time.Duration) error {
	s.data[key] = val
	return nil
}

func (s *memoryStorage) Del(ctx context.Context, keys ...string) error {
	for _, key := range keys {
		delete(s.data, key)
	}
	return nil
}
