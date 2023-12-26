package syncutil

import (
	"sync"
	"sync/atomic"
)

// SyncMap
// sync.Map 封装，使其支持计数
type SyncMap struct {
	sync.Map
	count uint32
}

func NewSyncMap() *SyncMap {
	return new(SyncMap)
}

func (s *SyncMap) Swap(key, value any) (any, bool) {
	previous, loaded := s.Map.Swap(key, value)
	if !loaded {
		s.inc()
	}
	return previous, loaded
}

func (s *SyncMap) LoadAndDelete(key any) (any, bool) {
	value, loaded := s.Map.LoadAndDelete(key)
	if loaded {
		s.dec()
	}
	return value, loaded
}

func (s *SyncMap) LoadOrStore(key, value any) (any, bool) {
	actual, loaded := s.Map.LoadOrStore(key, value)
	if !loaded {
		s.inc()
	}
	return actual, loaded
}

func (s *SyncMap) Store(key, value any) {
	if _, loaded := s.Map.Swap(key, value); !loaded {
		s.inc()
	}
}

func (s *SyncMap) Delete(key any) {
	s.LoadAndDelete(key)
}

func (s *SyncMap) CompareAndDelete(key any, old any) bool {
	deleted := s.Map.CompareAndDelete(key, old)
	if deleted {
		s.dec()
	}
	return deleted
}

// 计数+1
func (s *SyncMap) inc() {
	atomic.AddUint32(&s.count, 1)
}

// 计数-1
func (s *SyncMap) dec() {
	atomic.AddUint32(&s.count, ^uint32(0))
}

func (s *SyncMap) Len() int {
	return int(atomic.LoadUint32(&s.count))
}
