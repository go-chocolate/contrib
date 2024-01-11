package sharding

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Sharding interface {
	InitializeSharding(ctx context.Context) error
	Insert(ctx context.Context, record Record) error
	BatchInsert(ctx context.Context, records []Record) error
	Find(ctx context.Context, condition *Condition, offset, limit int) ([]Record, int64, error)
	FindOne(ctx context.Context, condition *Condition) (Record, error)
	Count(ctx context.Context, condition *Condition) (int64, error)
}

type Option func(s *sharding)

func applyOptions(s *sharding, options ...Option) {
	for _, opt := range options {
		opt(s)
	}
}

func WithStrategy(strategy Strategy) Option {
	return func(s *sharding) {
		s.strategy = strategy
	}
}

func WithStorage(storage Storage) Option {
	return func(s *sharding) {
		s.storage = storage
	}
}

func WithPrimaryKey(primaryKey string) Option {
	return func(s *sharding) {
		s.primaryKey = primaryKey
	}
}

func WithCreateTable(createTable bool, createTableScript string) Option {
	return func(s *sharding) {
		s.createTable = createTable
		s.createTableScript = createTableScript
	}
}

type sharding struct {
	storage  Storage
	strategy Strategy

	name              string
	primaryKey        string
	shardingColumn    string
	createTable       bool
	createTableScript string

	db *gorm.DB
}

func NewSharding(db *gorm.DB, name string, options ...Option) Sharding {
	sh := &sharding{
		storage:           nil,
		strategy:          nil,
		name:              name,
		primaryKey:        "id",
		shardingColumn:    "id",
		createTableScript: "",
		db:                db,
	}
	applyOptions(sh, options...)
	return sh
}

func (s *sharding) InitializeSharding(ctx context.Context) error {
	info := &Information{Name: s.name}
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, exec := range info.exec() {
			if err := tx.Exec(exec).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *sharding) Insert(ctx context.Context, record Record) error {
	db := s.db.WithContext(ctx)
	tableName, err := s.table(ctx, record)
	if err != nil {
		return err
	}
	return db.Table(tableName).Create(map[string]any(record)).Error
}

func (s *sharding) BatchInsert(ctx context.Context, records []Record) error {
	batches, err := s.splitBatch(ctx, records)
	if err != nil {
		return err
	}
	db := s.db.WithContext(ctx)
	err = db.Transaction(func(tx *gorm.DB) error {
		for table, batch := range batches {
			if err := tx.Table(table).Create(batch).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func (s *sharding) Find(ctx context.Context, condition *Condition, offset, limit int) ([]Record, int64, error) {
	condition.PrimaryKey = s.primaryKey
	if !s.storage.Exist(ctx, s.name, condition) {
		if err := s.buildIndex(ctx, condition); err != nil {
			return nil, 0, err
		}
	}
	count := s.storage.Count(ctx, s.name, condition)
	items, err := s.storage.Get(ctx, s.name, condition, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	results, err := s.findByItems(ctx, condition.Selects, items)
	return results, count, err
}

func (s *sharding) FindOne(ctx context.Context, condition *Condition) (Record, error) {
	tables, err := s.getSharding(ctx)
	if err != nil {
		return nil, err
	}
	db := s.db.WithContext(ctx)
	for _, info := range tables {
		var item = map[string]any{}
		err = condition.DB(db.Table(info.Table)).Take(&item).Error
		switch err {
		case gorm.ErrRecordNotFound:
			continue
		case nil:
			return item, nil
		default:
			return nil, err
		}
	}
	return nil, err
}

func (s *sharding) Count(ctx context.Context, condition *Condition) (int64, error) {
	tables, err := s.getSharding(ctx)
	if err != nil {
		return 0, err
	}
	var total int64
	db := s.db.WithContext(ctx)
	for _, info := range tables {
		var count int64
		if err := condition.DB(db.Table(info.Table)).Count(&count).Error; err != nil {
			return 0, err
		}
		total += count
	}
	return total, nil
}

func (s *sharding) splitBatch(ctx context.Context, data []Record) (map[string][]map[string]any, error) {
	var results = make(map[string][]map[string]any)
	for _, item := range data {
		tableName, err := s.table(ctx, item)
		if err != nil {
			return nil, err
		}
		results[tableName] = append(results[tableName], item)
	}
	return results, nil
}

func (s *sharding) buildIndex(ctx context.Context, condition *Condition) error {
	var columns = make([]string, 1, len(condition.Order)+len(condition.Group)+1)
	columns[0] = s.primaryKey
	for _, order := range condition.Order {
		columns = append(columns, order.Column)
	}
	for _, group := range condition.Group {
		columns = append(columns, group)
	}
	columns = s.unique(columns)

	tables, err := s.getSharding(ctx)
	if err != nil {
		return err
	}
	for _, info := range tables {
		if err := s.buildIndexIter(ctx, info.Table, columns, condition); err != nil {
			return err
		}
	}
	return nil
}

func (s *sharding) buildIndexIter(ctx context.Context, table string, columns []string, condition *Condition) error {
	db := s.db.WithContext(ctx)
	var lastId any

	query := db.Table(table).Select(columns)
	for {
		var where = condition.Where
		if lastId != nil {
			lastIdGt := clause.Gt{Column: s.primaryKey, Value: lastId}
			if where == nil {
				where = lastIdGt
			} else {
				where = clause.And(where, lastIdGt)
			}
		}
		var records []map[string]any
		if where != nil {
			query = query.Clauses(where)
		}
		if err := query.Order(s.primaryKey).Limit(2000).Find(&records).Error; err != nil {
			return err
		}
		if len(records) == 0 {
			break
		}
		lastId = records[len(records)-1][s.primaryKey]
		if lastId == nil {
			return fmt.Errorf("primary key '%s' not found for table '%s'", s.primaryKey, table)
		}
		var items []*Item
		for _, v := range records {
			item := &Item{
				PrimaryKey:    v[s.primaryKey],
				ShardingTable: table,
				Order:         condition.Order,
				Group:         condition.Group,
				OrderValues:   Record{},
				GroupValues:   Record{},
			}
			for _, order := range condition.Order {
				item.OrderValues[order.Column] = v[order.Column]
			}
			for _, group := range condition.Group {
				item.GroupValues[group] = v[group]
			}
			items = append(items, item)
		}
		// build index
		if err := s.storage.Put(ctx, s.name, condition, items); err != nil {
			return err
		}
	}
	return nil
}

func (s *sharding) findByItems(ctx context.Context, selects []string, items []*Item) ([]Record, error) {
	db := s.db.WithContext(ctx)

	if len(items) == 0 {
		return nil, nil
	}
	var primaryKeys = make(map[string][]any)
	var order = make([]any, 0, len(items))
	for _, v := range items {
		primaryKeys[v.ShardingTable] = append(primaryKeys[v.ShardingTable], v.PrimaryKey)
		order = append(order, v.PrimaryKey)
	}
	if len(selects) > 0 && !contains(selects, s.primaryKey) {
		selects = append([]string{s.primaryKey}, selects...)
	}

	var cache = make(map[any]map[string]any)

	for table, pks := range primaryKeys {
		var records []map[string]any
		query := db.Table(table)
		if len(selects) > 0 {
			query = query.Select(selects)
		}
		if err := query.Where("`"+s.primaryKey+"`"+" IN (?)", pks).Find(&records).Error; err != nil {
			return nil, err
		}
		for _, record := range records {
			cache[record[s.primaryKey]] = record
		}
	}

	var results = make([]Record, 0, len(cache))
	for _, pk := range order {
		results = append(results, cache[pk])
	}
	return results, nil
}

func (s *sharding) getSharding(ctx context.Context) ([]*Information, error) {
	db := s.db.WithContext(ctx)

	var tables []*Information
	if err := db.Table((&Information{Name: s.name}).tableName()).Where(&Information{Name: s.name}).Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (s *sharding) table(ctx context.Context, record Record) (string, error) {
	suffix, err := s.strategy.TableSuffix(record, s.shardingColumn)
	if err != nil {
		return "", err
	}
	tableName := fmt.Sprintf("%s_%s", s.name, suffix)

	if !s.createTable {
		return tableName, nil
	}

	db := s.db.WithContext(ctx)
	where := &Information{Name: s.name, Table: tableName}
	if err := db.Table(where.tableName()).Where("`name` = ? AND `table` = ?", s.name, tableName).Take(where).Error; err == gorm.ErrRecordNotFound {
		exec := fmt.Sprintf("CREATE TABLE `%s` %s", tableName, s.createTableScript)
		err = s.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(exec).Error; err != nil {
				return err
			}
			if err := tx.Table(where.tableName()).Create(&Information{
				ID:          int64(uuid.New().ID()),
				Name:        s.name,
				Table:       tableName,
				CreatedTime: time.Now(),
				UpdatedTime: time.Now(),
			}).Error; err != nil {
				return err
			}
			return nil
		})
		return tableName, err
	}
	return tableName, nil
}

func (s *sharding) unique(keys []string) []string {
	var unique = make([]string, 0, len(keys))
	var cache = make(map[string]struct{})
	for _, key := range keys {
		if _, ok := cache[key]; !ok {
			unique = append(unique, key)
			cache[key] = struct{}{}
		}
	}
	return unique
}
