package strategy

import (
	"fmt"

	"github.com/go-chocolate/contrib/database/sharding"
)

type ReminderStrategy struct {
	Num int
}

func (s *ReminderStrategy) TableSuffix(record sharding.Record, column string) (string, error) {
	val := record.Int(column) % s.Num
	if val == 0 {
		val = s.Num
	}
	return fmt.Sprintf("%04d", val), nil
}

func NewReminderStrategy(num int) *ReminderStrategy {
	return &ReminderStrategy{Num: num}
}
