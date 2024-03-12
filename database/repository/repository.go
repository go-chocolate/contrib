package repository

import (
	"context"
	"errors"
)

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

var ErrUnimplemented = errors.New("not implemented")

type UnimplementedRepository[T any] struct{}

func (rep *UnimplementedRepository[T]) FindOne(ctx context.Context, where any) (*T, error) {
	return nil, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) List(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error) {
	return nil, 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Count(ctx context.Context, where any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Update(ctx context.Context, where any, update any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Delete(ctx context.Context, where any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Insert(ctx context.Context, data any) (int64, error) {
	return 0, ErrUnimplemented
}
func (rep *UnimplementedRepository[T]) Iterate(ctx context.Context, column string, where any) (Iterator[T], error) {
	return nil, ErrUnimplemented
}
