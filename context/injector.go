package context

import "context"

type Injector interface {
	WithContext(ctx context.Context) context.Context
}

func WithContext(ctx context.Context, injector ...Injector) context.Context {
	for _, v := range injector {
		ctx = v.WithContext(ctx)
	}
	return ctx
}
