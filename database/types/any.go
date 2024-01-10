package types

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type Any struct {
	val any
}

func (a *Any) Scan(v any) error {
	switch val := v.(type) {
	case int:
		a.val = int64(val)
	case int8:
		a.val = int64(val)
	case int16:
		a.val = int64(val)
	case int32:
		a.val = int64(val)
	case uint:
		a.val = uint64(val)
	case uint8:
		a.val = uint64(val)
	case uint16:
		a.val = uint64(val)
	case uint32:
		a.val = uint64(val)
	case float32:
		a.val = float64(val)
	case string, int64, uint64, float64, bool:
		a.val = val
	default:
		return fmt.Errorf("unsupported type for scan to any: %v", reflect.TypeOf(v))
	}
	return nil
}

func (a Any) Value() (driver.Value, error) {
	return a.val, nil
}

func (a *Any) Set(val any) *Any {
	if err := a.Scan(val); err != nil {
		panic(err)
	}
	return a
}

func (a *Any) Any() any {
	return a.val
}

func (a *Any) Int() int64 {
	switch val := a.val.(type) {
	case int64:
		return int64(val)
	case uint64:
		return int64(val)
	case float64:
		return int64(val)
	case string:
		v, _ := strconv.ParseInt(val, 10, 64)
		return v
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func (a *Any) Uint() uint64 {
	switch val := a.val.(type) {
	case int64:
		return uint64(val)
	case uint64:
		return uint64(val)
	case float64:
		return uint64(val)
	case string:
		v, _ := strconv.ParseUint(val, 10, 64)
		return v
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func (a *Any) String() string {
	switch val := a.val.(type) {
	case string:
		return val
	case fmt.Stringer:
		return val.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (a *Any) Float() float64 {
	switch val := a.val.(type) {
	case int64:
		return float64(val)
	case uint64:
		return float64(val)
	case float64:
		return val
	case string:
		v, _ := strconv.ParseFloat(val, 64)
		return v
	case bool:
		if val {
			return 1
		}
	}
	return 0
}

func (a *Any) Boolean() bool {
	switch val := a.val.(type) {
	case bool:
		return val
	case string:
		return val == "true" || val == "TRUE" || val == "True"
	}
	return false
}

func (a *Any) Time() time.Time {
	switch val := a.val.(type) {
	case time.Time:
		return val
	case string:
		t, _ := time.ParseInLocation(time.DateTime, val, time.Local)
		return t
	}
	return time.Time{}
}
