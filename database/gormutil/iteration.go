package gormutil

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var ErrIterationEoF = errors.New("iteration eof")

type limitIterator[T any] struct {
	offset, limit int
	query         *gorm.DB
}

func (i *limitIterator[T]) Count() (count int64, err error) {
	err = i.query.Count(&count).Error
	return
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

func (i *columnIterator[T]) Count() (count int64, err error) {
	query := i.db
	if i.where != nil {
		switch cond := i.where.(type) {
		case clause.Expression:
			query = query.Clauses(cond)
		default:
			query = query.Where(cond)
		}
	}
	err = query.Count(&count).Error
	return
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
