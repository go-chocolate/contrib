package kv

const (
	REDIS  = "redis"
	MEMORY = "memory"
)

type Driver func(c Option) (Storage, error)

var drivers = map[string]Driver{
	REDIS:  redisDriver,
	MEMORY: memoryDriver,
}
