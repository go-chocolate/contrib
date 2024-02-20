package goroutine

import (
	"context"
	"fmt"
	"sync"
)

type Waiter interface {
	Wait()
}

type Promisor[OUTPUT any] interface {
	Then(ctx context.Context, then func(OUTPUT, error)) Waiter
	Await(ctx context.Context) (OUTPUT, error)
}

type promise[OUTPUT any] struct {
	ch  chan OUTPUT
	ech chan error
}

func (p *promise[OUTPUT]) Then(ctx context.Context, then func(OUTPUT, error)) Waiter {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		then(p.Await(ctx))
	})
	return wg
}

func (p *promise[OUTPUT]) Await(ctx context.Context) (OUTPUT, error) {
	var output OUTPUT
	select {
	case output = <-p.ch:
		return output, nil
	case e := <-p.ech:
		return output, e
	case <-ctx.Done():
		return output, ctx.Err()
	}
}

func Promise[INPUT, OUTPUT any](f func(INPUT) (OUTPUT, error), input INPUT) Promisor[OUTPUT] {
	p := &promise[OUTPUT]{
		ch:  make(chan OUTPUT, 1),
		ech: make(chan error, 1),
	}
	Go(func() {
		result, err := doPromise(f, input)
		if err != nil {
			p.ech <- err
		} else {
			p.ch <- result
		}
	})
	return p
}

func doPromise[INPUT, OUTPUT any](f func(INPUT) (OUTPUT, error), input INPUT) (result OUTPUT, err error) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			if e, ok := panicInfo.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicInfo)
			}
		}
	}()
	result, err = f(input)
	return
}

type AnyPromisor interface {
	Then(ctx context.Context, then func(any, error)) Waiter
	Await(context.Context) (any, error)
}

type anyPromise struct {
	ch  chan any
	ech chan error
}

func (p *anyPromise) Then(ctx context.Context, then func(any, error)) Waiter {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	Go(func() {
		defer wg.Done()
		then(p.Await(ctx))
	})
	return wg
}

func (p *anyPromise) Await(ctx context.Context) (any, error) {
	select {
	case v := <-p.ch:
		return v, nil
	case e := <-p.ech:
		return nil, e
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func AnyPromise(f func(...any) (any, error), args ...any) AnyPromisor {
	p := &anyPromise{
		ch:  make(chan any, 1),
		ech: make(chan error, 1),
	}
	Go(func() {
		result, err := doAnyPromise(f, args...)
		if err != nil {
			p.ech <- err
		} else {
			p.ch <- result
		}
	})
	return p
}

func doAnyPromise(f func(...any) (any, error), args ...any) (result any, err error) {
	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			if e, ok := panicInfo.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicInfo)
			}
		}
	}()
	result, err = f(args...)
	return
}
