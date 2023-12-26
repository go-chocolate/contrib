package goroutine

import "sync"

type Promisor[INPUT, OUTPUT any] interface {
	Then(f func(OUTPUT, error)) Promisor[INPUT, OUTPUT]
	Await() (OUTPUT, error)
}

type promise[INPUT, OUTPUT any] struct {
	f      func(INPUT) (OUTPUT, error)
	output OUTPUT
	err    error
	wait   *sync.WaitGroup
}

func (p *promise[INPUT, OUTPUT]) do(input INPUT) {
	p.wait = new(sync.WaitGroup)
	p.wait.Add(1)
	ech := make(chan error, 1)
	safeGo(func() { p.output, p.err = p.f(input) }, p.wait, ech)
	if err, ok := <-ech; ok {
		p.err = err
	}
}

func (p *promise[INPUT, OUTPUT]) Then(f func(OUTPUT, error)) Promisor[INPUT, OUTPUT] {
	safeGo(func() { f(p.output, p.err) }, nil, nil)
	return p
}

func (p *promise[INPUT, OUTPUT]) Await() (OUTPUT, error) {
	p.wait.Wait()
	return p.output, p.err
}

func Promise[INPUT, OUTPUT any](f func(INPUT) (OUTPUT, error), input INPUT) Promisor[INPUT, OUTPUT] {
	p := &promise[INPUT, OUTPUT]{f: f}
	p.do(input)
	return p
}
