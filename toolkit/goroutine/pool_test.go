package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func Job(i int) RunFunc {
	return func() {
		time.Sleep(time.Second)
		fmt.Println(i)
	}
}

func TestGoroutinePool(t *testing.T) {
	pool := NewGoPool(WithLimit(5))
	for i := 0; i < 10; i++ {
		pool.Do(Job(i))
	}
	pool.Shutdown()
}
