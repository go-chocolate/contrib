package goroutine

import (
	"fmt"
	"os"
	"sync"
)

func Go(f func()) {
	go safeGo(f, nil, nil)
}

func GoError(f func(), ch chan error) {
	go safeGo(f, nil, ch)
}

func safeGo(f func(), wait *sync.WaitGroup, ch chan error) {
	if wait != nil {
		defer wait.Done()
	}
	if ch != nil {
		defer close(ch)
	}
	defer func() {
		if recoverError := recover(); recoverError != nil {
			fmt.Fprintf(os.Stderr, "%v", recoverError)
			var err error
			if e, ok := recoverError.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", recoverError)
			}
			if ch != nil {
				ch <- err
			}
		}
	}()
	f()
}
