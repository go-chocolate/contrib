package gormutil

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/go-chocolate/contrib/database/repository"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) GetDB(ctx context.Context) *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return FromContext(ctx)
}

func (r *Repository[T]) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *Repository[T]) where(cmd *gorm.DB, where any) *gorm.DB {
	var model T
	cmd = cmd.Model(model)
	if where != nil {
		switch condition := where.(type) {
		case clause.Expression:
			cmd = cmd.Clauses(condition)
		case int64:
			cmd = cmd.Where("id = ?", condition)
		default:
			cmd = cmd.Where(condition)
		}
	}
	return cmd
}

func (r *Repository[T]) list(cmd *gorm.DB, offset, limit int, order ...any) ([]*T, int64, error) {
	var count int64
	if err := cmd.Count(&count).Error; err != nil {
		return nil, count, err
	}
	for _, v := range order {
		cmd = cmd.Order(v)
	}
	var dst []*T
	err := cmd.Offset(offset).Limit(limit).Find(dst).Error
	return dst, count, err
}

func (r *Repository[T]) FindOne(ctx context.Context, where any) (dst *T, err error) {
	dst = new(T)
	err = r.where(r.db, where).Take(dst).Error
	return
}

func (r *Repository[T]) List(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error) {
	var cmd = r.where(r.GetDB(ctx), where)
	return r.list(cmd, offset, limit, order)
}

func (r *Repository[T]) Count(ctx context.Context, where any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), where)
	var count int64
	err := cmd.Count(&count).Error
	return count, err
}

func (r *Repository[T]) Update(ctx context.Context, where any, update any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), where).Updates(update)
	return cmd.RowsAffected, cmd.Error
}

func (r *Repository[T]) Insert(ctx context.Context, data any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), nil).Create(data)
	return cmd.RowsAffected, cmd.Error
}

func (r *Repository[T]) Delete(ctx context.Context, where any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), where).Delete(nil)
	return cmd.RowsAffected, cmd.Error
}

func (r *Repository[T]) Iterate(ctx context.Context, column string, where any) (repository.Iterator[T], error) {
	return &columnIterator[T]{
		column: column,
		db:     r.GetDB(ctx),
		where:  where,
	}, nil
}

var _ repository.Repository[any] = (*Repository[any])(nil)
