package goroutine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Option func(o *GoPool)

func applyOptions(o *GoPool, options ...Option) {
	for _, opt := range options {
		opt(o)
	}
}

func WithLimit(limit uint32) Option {
	return func(o *GoPool) {
		o.limit = limit
	}
}

type Runner interface {
	Run()
}

type RunFunc func()

func (r RunFunc) Run() {
	r()
}

// GoPool 协程（线程）池
type GoPool struct {
	channel chan Runner        // 待执行队列
	limit   uint32             // 协程数上限限制，0不限制
	count   uint32             // 当前数量
	ctx     context.Context    // 内部context
	cancel  context.CancelFunc // 取消执行

	wg *sync.WaitGroup //全局等待执行完成
}

// NewGoPool 新建协程池
// Deprecated 请使用 goroutine.NewGoPool()
func NewGoPool(options ...Option) *GoPool {
	o := &GoPool{
		limit: 0,
		count: 0,
		wg:    new(sync.WaitGroup),
	}
	applyOptions(o, options...)
	o.channel = make(chan Runner)
	o.ctx, o.cancel = context.WithCancel(context.Background())
	return o
}

// Do 加入协程池异步执行
// WARNING: 当协程池达到上限时，将阻塞等待，直到有空闲资源释放
// 如需要超时，请使用 context.WithTimeout + DoCtx
func (p *GoPool) Do(r Runner) error {
	return p.DoCtx(context.Background(), r)
}

// DoFunc 加入协程池异步执行
// same as Do
func (p *GoPool) DoFunc(f func()) error {
	return p.Do(RunFunc(f))
}

// DoFuncCtx 加入协程池异步执行
// same as DoCtx
func (p *GoPool) DoFuncCtx(ctx context.Context, f func()) error {
	return p.DoCtx(ctx, RunFunc(f))
}

// DoCtx 加入协程池异步执行
// 协程池达到上限时，将阻塞等待，直到成功加入协程池或被context取消
func (p *GoPool) DoCtx(ctx context.Context, r Runner) error {
	select {
	case p.channel <- r:
		return nil
	default:
		return p.scale(ctx, r)
	}
}

func (p *GoPool) scale(ctx context.Context, r Runner) error {
	if p.acquire() {
		p.wg.Add(1)
		go p.process(r, p.channel)
		return nil
	} else {
		select {
		case p.channel <- r:
			return nil
		case <-ctx.Done():
			return errors.New("canceled by input context")
		}
	}
}

// Shutdown 关闭协程池，并等待当前正在执行的任务执行完毕
// WARNING: 请勿在调用 Shutdown 之后再调用 Do
func (p *GoPool) Shutdown() {
	p.cancel()
	p.wg.Wait()
}

// ShutdownAnyway 关闭协程池
// 不等待当前正在执行的任务，也不会取消当前的正在执行的任务。
// 如果进程未退出，当前正在执行的任务仍会继续执行；如果进程已退出，当前正在执行的任务将被强制中断
func (p *GoPool) ShutdownAnyway() {
	p.cancel()
}

// Count 当前正在执行的任务数量
func (p *GoPool) Count() int {
	return int(p.count)
}

// 检查是否达到上限，加“锁”
func (p *GoPool) acquire() bool {
	if p.limit == 0 {
		return true
	}
	//FIXME 两步加起来为非原子操作
	// 但理论上稍微超出限制不会有资源泄露问题
	if count := atomic.LoadUint32(&p.count); count >= p.limit {
		return false
	} else {
		atomic.AddUint32(&p.count, 1)
		return true
	}
}

// 释放“锁”
func (p *GoPool) release() {
	if p.limit == 0 {
		return
	}
	atomic.AddUint32(&p.count, ^uint32(0))
}

// 执行任务
func (p *GoPool) process(current Runner, ch chan Runner) {
	defer p.wg.Done()
	defer p.release()
	if current != nil {
		p.safeRun(current)
	}
	for {
		select {
		case r, ok := <-ch:
			if !ok {
				return
			}
			p.safeRun(r) //执行任务
		case <-p.ctx.Done(): //强制退出
			return
		case <-time.After(time.Second): // 释放空闲协程
			return
		}
	}
}

func (p *GoPool) safeRun(runner Runner) {
	defer func() {
		if recoverError := recover(); recoverError != nil {
			fmt.Fprintf(os.Stderr, "%v", recoverError)
		}
	}()
	runner.Run()
}
