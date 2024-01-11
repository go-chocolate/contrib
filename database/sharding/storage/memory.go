package storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm/clause"

	"github.com/go-chocolate/contrib/database/sharding"
)

type ClauseBuilder struct {
	bytes.Buffer
}

func (b *ClauseBuilder) WriteQuoted(field interface{}) {
	b.WriteByte('`')
	defer b.WriteByte('`')
	switch val := field.(type) {
	case string:
		b.WriteString(val)
	case fmt.Stringer:
		b.WriteString(val.String())
	default:
		fmt.Fprintf(b, "%v", field)
	}

}

func (b *ClauseBuilder) AddVar(w clause.Writer, values ...interface{}) {
	if len(values) == 0 {
		return
	}
	for _, v := range values {
		switch val := v.(type) {
		case string:
			b.WriteString("'" + val + "',")
		case time.Time:
			b.WriteString("'" + val.Format(time.DateTime) + "'")
		case fmt.Stringer:
			b.WriteString("'" + val.String() + "',")
		case json.Number:
			b.WriteString(val.String())
		default:
			fmt.Fprintf(b, "%v,", v)
		}
	}
	b.Truncate(b.Len() - 1)
}

func (b *ClauseBuilder) AddError(error) error {
	return nil
}

func NewClauseBuilder() *ClauseBuilder {
	return &ClauseBuilder{}
}

type MemoryStorage struct {
	sync.RWMutex
	store map[string][]*sharding.Item
}

func (s *MemoryStorage) key(name string, condition *sharding.Condition) string {
	b := NewClauseBuilder()
	b.WriteString(name)
	if condition.Where != nil {
		condition.Where.Build(b)
	}
	for _, order := range condition.Order {
		fmt.Fprintf(b, "%s,", order.String())
	}
	b.WriteString(strings.Join(condition.Group, ","))
	hash := md5.Sum(b.Bytes())
	return hex.EncodeToString(hash[:])
}

func (s *MemoryStorage) Count(ctx context.Context, name string, condition *sharding.Condition) int64 {
	key := s.key(name, condition)
	s.RLock()
	defer s.RUnlock()
	return int64(len(s.store[key]))
}

func (s *MemoryStorage) Get(ctx context.Context, name string, condition *sharding.Condition, offset, limit int) ([]*sharding.Item, error) {
	key := s.key(name, condition)
	s.Lock()
	defer s.Unlock()

	items := s.store[key]
	if len(condition.Order) > 0 {
		sort.Sort(sharding.Items(items))
	}
	if len(condition.Group) > 0 {
		//TODO
		return nil, fmt.Errorf("memory storage not support for group")
	}
	if offset > len(items) {
		return nil, nil
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end], nil
}

func (s *MemoryStorage) Exist(ctx context.Context, name string, condition *sharding.Condition) bool {
	key := s.key(name, condition)
	s.RLock()
	defer s.RUnlock()
	_, ok := s.store[key]
	return ok
}

func (s *MemoryStorage) Put(ctx context.Context, name string, condition *sharding.Condition, items []*sharding.Item) error {
	key := s.key(name, condition)
	s.Lock()
	defer s.Unlock()
	if _, ok := s.store[key]; ok {
		s.store[key] = append(s.store[key], items...)
	} else {
		s.store[key] = items
	}
	return nil
}

func (s *MemoryStorage) Del(ctx context.Context, name string, condition *sharding.Condition) {
	key := s.key(name, condition)
	s.Lock()
	defer s.Unlock()
	delete(s.store, key)
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{store: map[string][]*sharding.Item{}}
}
