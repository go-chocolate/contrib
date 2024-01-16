package context

import (
	"context"
	"sync"
	"time"
)

type Context = context.Context

type Setter interface {
	Set(key, val any)
}

type valueContext struct {
	context.Context
	sync.Map
}

func (c *valueContext) Value(key any) any {
	if val, ok := c.Load(key); ok {
		return val
	}
	return c.Context.Value(key)
}

func (c *valueContext) Set(key, val any) {
	c.Store(key, val)
}

var (
	Background      = context.Background
	WithCancel      = context.WithCancel
	TODO            = context.TODO
	WithCancelCause = context.WithCancelCause
	WithDeadline    = context.WithDeadline
	WithTimeout     = context.WithTimeout
)

// WithoutCancel returns a copy of parent that is not canceled when parent is canceled.
// The returned context returns no Deadline or Err, and its Done channel is nil.
// Calling [Cause] on the returned context returns nil.
// WARNING: this function is supported in go 1.21+
func WithoutCancel(parent context.Context) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}
	return withoutCancelCtx{parent}
}

type withoutCancelCtx struct {
	context.Context
}

func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (withoutCancelCtx) Err() error {
	return nil
}
