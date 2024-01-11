package sharding_test

import (
	"context"
	"testing"
	"time"

	"gorm.io/gorm/clause"

	"github.com/go-chocolate/contrib/database/gormutil"
	"github.com/go-chocolate/contrib/database/sharding"
	"github.com/go-chocolate/contrib/database/sharding/storage"
	"github.com/go-chocolate/contrib/database/sharding/strategy"
)

type ExampleSharding struct {
	sharding.Information
}

func (e *ExampleSharding) TableName() string {
	return "example_sharding"
}

func TestSharding(t *testing.T) {
	db, err := gormutil.OpenMemory(gormutil.WithStdLogger())
	if err != nil {
		t.Error(err)
	}

	{
		db.Migrator().CreateTable(&ExampleSharding{})
	}

	sh := sharding.NewSharding(
		db,
		"example",
		sharding.WithStrategy(strategy.NewReminderStrategy(4)),
		sharding.WithStorage(storage.NewMemoryStorage()),
		sharding.WithCreateTable(true, "(`id` bigint primary key, name varchar(255), created_time datetime)"),
	)
	{
		records := []sharding.Record{
			{"id": 1, "name": "Zhangsan", "created_time": time.Now().AddDate(0, 0, 0)},
			{"id": 2, "name": "Lisi", "created_time": time.Now().AddDate(0, 0, -1)},
			{"id": 3, "name": "Wangwu", "created_time": time.Now().AddDate(0, 0, -2)},
			{"id": 4, "name": "Zhaoliu", "created_time": time.Now().AddDate(0, 0, -3)},
			{"id": 5, "name": "Tianqi", "created_time": time.Now().AddDate(0, 0, -4)},
			{"id": 6, "name": "Zhangsan", "created_time": time.Now().AddDate(0, 0, -5)},
			{"id": 7, "name": "Lisi", "created_time": time.Now().AddDate(0, 0, -6)},
			{"id": 8, "name": "Wangwu", "created_time": time.Now().AddDate(0, 0, -7)},
			{"id": 9, "name": "Zhaoliu", "created_time": time.Now().AddDate(0, 0, -8)},
			{"id": 10, "name": "Tianqi", "created_time": time.Now().AddDate(0, 0, -9)},
			{"id": 11, "name": "Zhangsan", "created_time": time.Now().AddDate(0, 0, -10)},
			{"id": 12, "name": "Lisi", "created_time": time.Now().AddDate(0, 0, -11)},
			{"id": 13, "name": "Wangwu", "created_time": time.Now().AddDate(0, 0, -12)},
			{"id": 14, "name": "Zhaoliu", "created_time": time.Now().AddDate(0, 0, -13)},
			{"id": 15, "name": "Tianqi", "created_time": time.Now().AddDate(0, 0, -14)},
		}
		if err := sh.BatchInsert(context.Background(), records); err != nil {
			t.Error(err)
		}
	}

	{
		one, err := sh.FindOne(context.Background(), &sharding.Condition{Where: clause.Eq{Column: "name", Value: "Zhangsan"}})
		if err != nil {
			t.Error(err)
		} else {
			t.Log(one)
		}
	}

	{
		condition := &sharding.Condition{
			//Where: clause.Gt{Column: "created_time", Value: "2024-01-01 00:00:00"},
			Order: []sharding.Order{{Column: "id", Desc: true}},
		}
		data, count, err := sh.Find(context.Background(), condition, 10, 10)
		if err != nil {
			t.Error(err)
		} else {
			for _, v := range data {
				t.Log(v)
			}
			t.Log(len(data), count)
		}
	}
}
