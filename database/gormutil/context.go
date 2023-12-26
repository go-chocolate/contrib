package gormutil

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type gormContextKey struct{}

var _gormContextKey = &gormContextKey{}

func WithContext(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, _gormContextKey, db)
}

func FromContext(ctx context.Context) *gorm.DB {
	if val := ctx.Value(_gormContextKey); val != nil {
		if db, ok := val.(*gorm.DB); ok {
			return db
		}
	}
	panic(fmt.Errorf("gorm.DB is not attached to context.Context"))
}
