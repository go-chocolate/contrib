package repository

import (
	"context"
	"fmt"
	"reflect"
)

type Model interface {
	GetID() int64
	GetUpdate() any
	Rebase()
}

type ModelOperator[T any] struct {
	Rep   Repository[T]
	Value *T
}

func NewModelOperator[T any](rep Repository[T], value *T) *ModelOperator[T] {
	_, ok := any(value).(Model)
	if !ok {
		t := reflect.TypeOf(value)
		panic(fmt.Errorf("type %v must implement repository.Model", t))
	}
	return &ModelOperator[T]{
		Rep:   rep,
		Value: value,
	}
}

func (m *ModelOperator[T]) Model() Model {
	return any(m.Value).(Model)
}

func (m *ModelOperator[T]) Load(ctx context.Context) error {
	if data, err := m.Rep.FindOne(ctx, m.Value); err != nil {
		return err
	} else {
		m.Value = data
	}
	m.Model().Rebase()
	return nil
}

func (m *ModelOperator[T]) Update(ctx context.Context) (rowsAffected int64, err error) {
	model := m.Model()
	update := model.GetUpdate()
	rowsAffected, err = m.Rep.Update(ctx, model.GetID(), update)
	if err != nil {
		return
	}
	model.Rebase()
	return
}

func (m *ModelOperator[T]) Create(ctx context.Context) (rowsAffected int64, err error) {
	rowsAffected, err = m.Rep.Insert(ctx, m.Value)
	return
}

func (m *ModelOperator[T]) Delete(ctx context.Context) (rowsAffected int64, err error) {
	model := m.Model()
	if id := model.GetID(); id > 0 {
		return m.Rep.Delete(ctx, id)
	} else {
		return m.Rep.Delete(ctx, model)
	}
}

func (m *ModelOperator[T]) List(ctx context.Context, limit, offset int, order ...any) ([]*T, int64, error) {
	return m.Rep.List(ctx, m.Value, limit, offset, order...)
}
