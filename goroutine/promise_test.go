package goroutine

import (
	"context"
	"strconv"
	"testing"
)

func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

func TestPromise(t *testing.T) {
	Promise(Atoi, "1").Then(context.Background(), func(result int, err error) { t.Log(result, err) }).Wait()

	t.Log(Promise(Atoi, "2").Await(context.Background()))
}

func TestAnyPromise(t *testing.T) {
	f := func(a ...any) (any, error) {
		return Atoi(a[0].(string))
	}
	AnyPromise(f, "111").Then(context.Background(), func(result any, err error) { t.Log(result, err) }).Wait()

	t.Log(AnyPromise(f, "222").Await(context.Background()))
}
