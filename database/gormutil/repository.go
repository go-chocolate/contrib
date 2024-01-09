package gormutil

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository struct {
	db *gorm.DB
}

func (r *Repository) GetDB(ctx context.Context) *gorm.DB {
	if r.db != nil {
		return r.db
	}
	return FromContext(ctx)
}

func (r *Repository) SetDB(db *gorm.DB) {
	r.db = db
}

func (r *Repository) FindOne(ctx context.Context, dst any, where any) error {
	switch condition := where.(type) {
	case clause.Expression:
		return r.GetDB(ctx).Clauses(condition).Take(dst).Error
	default:
		return r.GetDB(ctx).Where(condition).Take(dst).Error
	}
}

func (r *Repository) FindOneById(ctx context.Context, dst any, id int64) error {
	return r.GetDB(ctx).Take(dst, "id = ?", id).Error
}

func (r *Repository) FindOneBy(ctx context.Context, dst any) error {
	return r.GetDB(ctx).Take(dst).Error
}

func (r *Repository) where(cmd *gorm.DB, model any, where any) *gorm.DB {
	if model != nil {
		if table, ok := model.(string); ok {
			cmd = cmd.Table(table)
		} else {
			cmd = cmd.Model(model)
		}
	}
	if where != nil {
		switch condition := where.(type) {
		case clause.Expression:
			cmd = cmd.Clauses(condition)
		default:
			cmd = cmd.Where(condition)
		}
	}
	return cmd
}

func (r *Repository) list(dst any, cmd *gorm.DB, offset, limit int, order ...any) (int64, error) {
	var count int64
	if err := cmd.Count(&count).Error; err != nil {
		return count, err
	}
	for _, v := range order {
		cmd = cmd.Order(v)
	}
	err := cmd.Offset(offset).Limit(limit).Find(dst).Error
	return count, err
}

func (r *Repository) FindList(ctx context.Context, dst any, model, where any, offset, limit int, order ...any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), model, where)
	return r.list(dst, cmd, offset, limit, order)
}

func (r *Repository) Count(ctx context.Context, model, where any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), model, where)
	var count int64
	err := cmd.Count(&count).Error
	return count, err
}

func (r *Repository) Update(ctx context.Context, model, where any, update any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), model, where).Updates(update)
	return cmd.RowsAffected, cmd.Error
}

func (r *Repository) Insert(ctx context.Context, model, data any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), model, nil).Create(data)
	return cmd.RowsAffected, cmd.Error
}

func (r *Repository) Delete(ctx context.Context, model, where any) (int64, error) {
	var cmd = r.where(r.GetDB(ctx), model, where).Delete(model)
	return cmd.RowsAffected, cmd.Error
}

type RepositoryT[T any] struct {
	Repository
}

func (r *RepositoryT[T]) FindOne(ctx context.Context, where any) (*T, error) {
	var result = new(T)
	err := r.Repository.FindOne(ctx, result, where)
	return result, err

}
func (r *RepositoryT[T]) FindOneById(ctx context.Context, id int64) (*T, error) {
	var result = new(T)
	err := r.Repository.FindOneById(ctx, result, id)
	return result, err
}

func (r *RepositoryT[T]) FindList(ctx context.Context, where any, offset, limit int, order ...any) ([]*T, int64, error) {
	var results []*T
	var model T
	count, err := r.Repository.FindList(ctx, &results, &model, where, offset, limit, order...)
	return results, count, err
}

func (r *RepositoryT[T]) Count(ctx context.Context, where any) (int64, error) {
	var model T
	return r.Repository.Count(ctx, &model, where)
}

func (r *RepositoryT[T]) Update(ctx context.Context, where any, update any) (int64, error) {
	var model T
	return r.Repository.Update(ctx, model, where, update)
}

func (r *RepositoryT[T]) Delete(ctx context.Context, where any) (int64, error) {
	var model T
	return r.Repository.Delete(ctx, model, where)
}

func (r *RepositoryT[T]) Insert(ctx context.Context, data any) (int64, error) {
	var model T
	return r.Repository.Insert(ctx, &model, data)
}
