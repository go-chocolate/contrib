package sharding

import (
	"fmt"
	"reflect"
	"time"
)

func compare(a, b any) int {
	if reflect.TypeOf(a).Kind() != reflect.TypeOf(b).Kind() {
		//TODO
		panic("type not match")
	}
	switch val := a.(type) {
	case int8:
		bval := b.(int8)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case int16:
		bval := b.(int16)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case int32:
		bval := b.(int32)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case int64:
		bval := b.(int64)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case int:
		bval := b.(int)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case uint8:
		bval := b.(uint8)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case uint16:
		bval := b.(uint16)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case uint32:
		bval := b.(uint32)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case uint64:
		bval := b.(uint64)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case uint:
		bval := b.(uint)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case float32:
		bval := b.(float32)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case float64:
		bval := b.(float64)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}
	case string:
		bval := b.(string)
		if val < bval {
			return -1
		} else if val == bval {
			return 0
		} else {
			return 1
		}

	case bool:
		bval := b.(bool)
		if val == bval {
			return 0
		} else if !val {
			return -1
		} else {
			return 1
		}
	case time.Time:
		bval := b.(time.Time)
		if val.Before(bval) {
			return -1
		} else if val.Equal(bval) {
			return 0
		} else {
			return 1
		}
	}
	panic(fmt.Errorf("cannot compare for type %v", reflect.TypeOf(a)))
}

func contains(array []string, item string) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}
