package kv

const (
	REDIS = "redis"
)

type Driver func(c Option) (Storage, error)

var drivers = map[string]Driver{
	REDIS: redisDriver,
}
