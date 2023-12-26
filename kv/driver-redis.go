package kv

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisStorage struct {
	prefix string
	client redis.UniversalClient
}

func redisDriver(c Option) (Storage, error) {
	var options = redis.UniversalOptions{
		Addrs:            c.Strings("Addrs"),
		ClientName:       c.String("ClientName"),
		DB:               int(c.Int64("DB")),
		Username:         c.String("Username"),
		Password:         c.String("Password"),
		SentinelUsername: c.String("SentinelUsername"),
		SentinelPassword: c.String("SentinelPassword"),
		MasterName:       c.String("MasterName"),
		//MaxRetries:            0,
		//MinRetryBackoff:       0,
		//MaxRetryBackoff:       0,
		//DialTimeout:           0,
		//ReadTimeout:           0,
		//WriteTimeout:          0,
		//ContextTimeoutEnabled: false,
		//PoolFIFO:              false,
		//PoolSize:              0,
		//PoolTimeout:           0,
		//MinIdleConns:          0,
		//MaxIdleConns:          0,
		//MaxActiveConns:        0,
		//ConnMaxIdleTime:       0,
		//ConnMaxLifetime:       0,
		//TLSConfig:             nil,
		//MaxRedirects:          0,
		//ReadOnly:              false,
		//RouteByLatency:        false,
		//RouteRandomly:         false,
		//DisableIndentity:      false,
	}
	client := redis.NewUniversalClient(&options)
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &redisStorage{client: client, prefix: c.String("Prefix")}, nil
}

func (s *redisStorage) key(key string) string {
	if s.prefix == "" {
		return key
	}
	return fmt.Sprintf("%s:%s", s.prefix, key)
}

func (s *redisStorage) Get(ctx context.Context, key string) ([]byte, error) {
	return s.client.Get(ctx, s.key(key)).Bytes()
}

func (s *redisStorage) Set(ctx context.Context, key string, val []byte, expiration ...time.Duration) error {
	var exp time.Duration
	if len(expiration) > 0 {
		exp = expiration[0]
	}
	return s.client.Set(ctx, s.key(key), val, exp).Err()
}

func (s *redisStorage) Del(ctx context.Context, keys ...string) error {
	return s.client.Del(ctx, keys...).Err()
}
