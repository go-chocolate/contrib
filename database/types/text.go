package types

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
)

type Text string

func (a *Text) Scan(v any) error {
	switch val := v.(type) {
	case string:
		*a = Text(val)
		return nil
	case []byte:
		*a = Text(val)
		return nil
	case fmt.Stringer:
		*a = Text(val.String())
		return nil
	case io.Reader:
		if buf, err := io.ReadAll(val); err != nil {
			return err
		} else {
			*a = Text(buf)
		}
	}
	return fmt.Errorf("cannot decode type %v to string", reflect.TypeOf(v))
}

func (a Text) Value() (driver.Value, error) {
	return string(a), nil
}
