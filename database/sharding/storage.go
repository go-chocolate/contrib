package sharding

import (
	"context"
)

type Item struct {
	PrimaryKey    any
	ShardingTable string
	Order         []Order
	Group         []string
	OrderValues   Record
	GroupValues   Record
}

type Items []*Item

func (c Items) Len() int {
	return len(c)
}
func (c Items) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Items) Less(i, j int) bool {
	left := c[i]
	right := c[j]
	for _, order := range left.Order {
		result := compare(left.OrderValues[order.Column], right.OrderValues[order.Column])
		if order.Desc {
			if result < 0 {
				return false
			} else if result == 0 {
				continue
			} else {
				return true
			}
		} else {
			if result < 0 {
				return true
			} else if result == 0 {
				continue
			} else {
				return false
			}
		}
	}
	return false
}

type Storage interface {
	Count(ctx context.Context, name string, condition *Condition) int64
	Get(ctx context.Context, name string, condition *Condition, offset, limit int) ([]*Item, error)
	Exist(ctx context.Context, name string, condition *Condition) bool
	Put(ctx context.Context, name string, condition *Condition, items []*Item) error
	Del(ctx context.Context, name string, condition *Condition)
}
