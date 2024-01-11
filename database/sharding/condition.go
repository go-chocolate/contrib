package sharding

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Condition struct {
	PrimaryKey string
	Selects    []string
	Where      clause.Expression
	Order      []Order
	Group      []string
}

func (c *Condition) DB(db *gorm.DB) *gorm.DB {
	if len(c.Selects) > 0 {
		db = db.Select(c.Selects)
	}
	if c.Where != nil {
		db = db.Clauses(c.Where)
	}
	for _, order := range c.Order {
		db = db.Order(order.String())
	}
	for _, group := range c.Group {
		db = db.Group(group)
	}
	return db
}

type Order struct {
	Column string
	Desc   bool
}

func (o *Order) String() string {
	if o.Desc {
		return fmt.Sprintf("`%s` DESC", o.Column)
	}
	return fmt.Sprintf("`%s`", o.Column)
}
