package repository

import "context"

type Repository[T any] interface {
	FindOne(ctx context.Context, where any) (*T, error)
	FindOneById(ctx context.Context, id int64) (*T, error)
	FindList(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error)
	Count(ctx context.Context, where any) (int64, error)
	Update(ctx context.Context, where any, update any) (int64, error)
	Delete(ctx context.Context, where any) (int64, error)
	Insert(ctx context.Context, data any) (int64, error)
}
