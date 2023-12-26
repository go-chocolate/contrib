package timeutil

import (
	"database/sql/driver"
	"fmt"
	"reflect"
	"time"
)

func (t *Time) Scan(v any) error {
	var s string
	switch val := v.(type) {
	case string:
		s = val
	case []byte:
		s = string(val)
	case time.Time:
		t.Time = val
		return nil
	default:
		return fmt.Errorf("cannot unmarshal field from type %v to time.Time", reflect.TypeOf(v))
	}
	var err error
	*t, err = Parse(s)
	return err
}

func (t Time) Value() (driver.Value, error) {
	return t.Time, nil
}
