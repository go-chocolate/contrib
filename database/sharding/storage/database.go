package storage

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/go-chocolate/contrib/database/sharding"
)

type DatabaseStorage struct {
	db          *gorm.DB
	indexTable  string
	readOnly    bool
	primaryKey  string
	tableColumn string
}

func NewDatabaseStorage(db *gorm.DB, table string, readOnly bool) *DatabaseStorage {
	return &DatabaseStorage{
		db:         db,
		indexTable: table,
		readOnly:   readOnly,
	}
}

func (s *DatabaseStorage) Count(ctx context.Context, name string, condition *sharding.Condition) int64 {
	db := s.db.WithContext(ctx).Table(s.indexTable)

	var count int64
	db.Clauses(condition.Where).Count(&count)
	return count
}

func (s *DatabaseStorage) Get(ctx context.Context, name string, condition *sharding.Condition, offset, limit int) ([]*sharding.Item, error) {
	exec := s.build(nil, condition, offset, limit)
	var results []sharding.Record
	err := s.db.WithContext(ctx).Raw(exec).Find(&results).Error
	if err != nil {
		return nil, err
	}
	var items []*sharding.Item
	for _, v := range results {
		item := &sharding.Item{
			PrimaryKey:    condition.PrimaryKey,
			ShardingTable: v.String(s.tableColumn),
			Order:         condition.Order,
			Group:         condition.Group,
			OrderValues:   sharding.Record{},
			GroupValues:   sharding.Record{},
		}
		for _, order := range condition.Order {
			item.OrderValues[order.Column] = v[order.Column]
		}
		for _, group := range condition.Group {
			item.GroupValues[group] = v[group]
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *DatabaseStorage) Exist(ctx context.Context, name string, condition *sharding.Condition) bool {
	if s.readOnly {
		return true
	}
	if !s.db.Migrator().HasTable(s.indexTable) {
		return false
	}
	var count int64
	s.db.Table(s.indexTable).Count(&count)
	return count > 0
}

func (s *DatabaseStorage) Put(ctx context.Context, name string, condition *sharding.Condition, items []*sharding.Item) error {
	if s.readOnly {
		return nil
	}
	return fmt.Errorf("unimplemented")
}

func (s *DatabaseStorage) Del(ctx context.Context, name string, condition *sharding.Condition) {
	if s.readOnly {
		return
	}
}

func (s *DatabaseStorage) build(selects []string, condition *sharding.Condition, offset, limit int) string {
	b := NewClauseBuilder()
	b.WriteString("SELECT ")
	if len(selects) > 0 {
		b.WriteString(strings.Join(selects, ","))
	} else {
		b.WriteString("*")
	}
	b.WriteString(" FROM ")
	b.WriteString(s.indexTable)
	if condition.Where != nil {
		b.WriteString(" WHERE ")
		condition.Where.Build(b)
	}
	if len(condition.Group) > 0 {
		b.WriteString(" GROUP BY ")
		b.WriteString("`" + strings.Join(condition.Group, "`, `") + "`")
	}
	if len(condition.Order) > 0 {
		b.WriteString(" ORDER BY ")
		for _, order := range condition.Order {
			b.WriteString(order.String() + ",")
		}
		b.Truncate(b.Len() - 1)
	}

	if limit > 0 {
		fmt.Fprintf(b, " LIMIT %d,%d", offset, limit)
	}
	return b.String()
}
