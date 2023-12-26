package kv

import (
	"context"
	"fmt"
)

type kvContextKey struct{}

var _kvContextKey = &kvContextKey{}

func WithContext(ctx context.Context, kv Storage) context.Context {
	return context.WithValue(ctx, _kvContextKey, kv)
}

func FromContext(ctx context.Context) Storage {
	if val := ctx.Value(_kvContextKey); val != nil {
		if db, ok := val.(Storage); ok {
			return db
		}
	}
	panic(fmt.Errorf("kv.Storage is not attached to context.Context"))
}
