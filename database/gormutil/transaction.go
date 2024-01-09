package gormutil

import (
	"context"

	"gorm.io/gorm"
)

func Tx(ctx context.Context, f func(ctx context.Context) error) error {
	return FromContext(ctx).Transaction(func(tx *gorm.DB) error {
		return f(WithContext(ctx, tx))
	})
}
