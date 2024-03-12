package repository

import "context"

type Repository[T any] interface {
	FindOne(ctx context.Context, where any) (*T, error)
	List(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error)
	Count(ctx context.Context, where any) (int64, error)
	Update(ctx context.Context, where any, update any) (int64, error)
	Delete(ctx context.Context, where any) (int64, error)
	Insert(ctx context.Context, data any) (int64, error)
	Iterate(ctx context.Context, column string, where any) (Iterator[T], error)
}

type Iterator[T any] interface {
	Count() (int64, error)
	Next() ([]*T, error)
}
