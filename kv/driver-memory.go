package kv

import (
	"context"
	"errors"
	"sync"
	"time"
)

type memoryItem struct {
	timestamp time.Time
	data      []byte
}

type memoryStorage struct {
	sync.Mutex
	storage map[string]*memoryItem
}

func (s *memoryStorage) Get(ctx context.Context, key string) ([]byte, error) {
	s.Lock()
	defer s.Unlock()
	if item, ok := s.storage[key]; ok {
		if item.timestamp.IsZero() || item.timestamp.After(time.Now()) {
			return item.data, nil
		} else {
			delete(s.storage, key)
		}
	}
	return nil, errors.New("record not found in memory cache")
}

func (s *memoryStorage) Set(ctx context.Context, key string, val []byte, expiration ...time.Duration) error {
	s.Lock()
	defer s.Unlock()
	if len(expiration) > 0 {
		s.storage[key] = &memoryItem{time.Now().Add(expiration[0]), val}
	} else {
		s.storage[key] = &memoryItem{time.Time{}, val}
	}
	return nil
}

func (s *memoryStorage) Del(ctx context.Context, keys ...string) error {
	s.Lock()
	defer s.Unlock()
	for _, key := range keys {
		delete(s.storage, key)
	}
	return nil
}

func memoryDriver(c Option) (Storage, error) { return &memoryStorage{}, nil }
