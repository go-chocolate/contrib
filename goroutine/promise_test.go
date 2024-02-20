package goroutine

import (
	"strconv"
	"testing"
)

func Atoi(s string) (int, error) {
	return strconv.Atoi(s)
}

func TestPromise(t *testing.T) {
	t.Log(
		Promise(Atoi, "1").
			Then(func(result int, err error) { t.Log(result, err) }).
			Await(),
	)
}
