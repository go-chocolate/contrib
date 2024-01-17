package gormutil

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

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

func (r *RepositoryT[T]) IterationByLimit(ctx context.Context, where any, limit int, order ...any) Iterator[T] {
	var model T
	query := r.where(r.GetDB(ctx), model, where)
	for _, v := range order {
		query = query.Order(v)
	}
	return &limitIterator[T]{
		limit: limit,
		query: query,
	}
}

func (r *RepositoryT[T]) IterationByColumn(ctx context.Context, where any, column ...string) Iterator[T] {
	var model T
	iter := &columnIterator[T]{db: r.GetDB(ctx).Model(model), where: where}
	if len(column) > 0 {
		iter.column = column[0]
	} else {
		iter.column = "id"
	}
	return iter
}

var ErrIterationEoF = errors.New("iteration eof")

type Iterator[T any] interface {
	Next() ([]*T, error)
}

type limitIterator[T any] struct {
	offset, limit int
	query         *gorm.DB
}

func (i *limitIterator[T]) Next() ([]*T, error) {
	var results []*T
	err := i.query.Offset(i.offset).Limit(i.limit).Find(&results).Error
	if err == nil && len(results) == 0 {
		return results, ErrIterationEoF
	}
	i.offset += i.limit
	return results, err
}

type columnIterator[T any] struct {
	lastID       any
	lastIDColumn string
	column       string
	db           *gorm.DB
	where        any
}

func (i *columnIterator[T]) Next() ([]*T, error) {
	query := i.db
	if i.where != nil {
		switch cond := i.where.(type) {
		case clause.Expression:
			if i.lastID != nil {
				query = query.Clauses(clause.And(cond, clause.Gt{Column: "`" + i.column + "`", Value: i.lastID}))
			} else {
				query = query.Clauses(cond)
			}
		default:
			query = query.Where(cond)
			if i.lastID != nil {
				query = query.Where("`"+i.column+"` > ?", i.lastID)
			}
		}
	} else if i.lastID != nil {
		query = query.Where("`"+i.column+"` > ?", i.lastID)
	}
	query = query.Order("`" + i.column + "`")

	var result []*T
	err := query.Find(&result).Error
	if err != nil {
		return result, err
	}
	if len(result) == 0 {
		return nil, ErrIterationEoF
	}
	return result, i.extractLastID(result[len(result)-1])
}

func (i *columnIterator[T]) extractLastID(item any) error {
	val := reflect.ValueOf(item)

	if i.lastIDColumn == "" {
		typ := reflect.TypeOf(item)
		for n := 0; n < val.NumField(); n++ {
			field := typ.Field(n)
			//decode gorm tag and find column name
			if tag := field.Tag.Get("gorm"); tag != "" {
				var column string
				for _, v := range strings.Split(tag, ";") {
					if len(v) > 7 && strings.ToLower(v[:7]) == "column:" {
						column = v[7:]
						break
					}
				}
				if column != "" && column == i.column {
					i.lastIDColumn = field.Name
					break
				}
			}

			// decode table column name to struct field name
			var column []byte
			var toUpper = true //make first character to upper
			for _, v := range column {
				if v == '_' { // skip underline and make next character to upper
					toUpper = true
					continue
				}
				if toUpper && (v >= 'a' && v <= 'z') {
					v = v - 32
				}
				column = append(column, v)
			}

			if string(column) == field.Name {
				i.lastIDColumn = field.Name
				break
			}
		}
	}
	if i.lastIDColumn == "" {
		typ := reflect.TypeOf(item)
		return fmt.Errorf("struct %s does not contain field %s", typ.String(), i.column)
	}
	i.lastID = val.FieldByName(i.lastIDColumn).Interface()
	return nil
}
