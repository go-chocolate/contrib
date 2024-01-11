package sharding

type Strategy interface {
	TableSuffix(record Record, column string) (string, error)
}

type StrategyFunc func(record Record, column string) (string, error)

func (f StrategyFunc) TableSuffix(record Record, column string) (string, error) {
	return f(record, column)
}
