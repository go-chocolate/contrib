package types

import (
	"database/sql/driver"
	"fmt"
	"io"
	"reflect"
)

type Binary []byte

func (a *Binary) Scan(v any) error {
	switch val := v.(type) {
	case string:
		*a = Binary(val)
		return nil
	case []byte:
		*a = Binary(val)
		return nil
	case fmt.Stringer:
		*a = Binary(val.String())
		return nil
	case io.Reader:
		if buf, err := io.ReadAll(val); err != nil {
			return err
		} else {
			*a = Binary(buf)
		}
	}
	return fmt.Errorf("cannot decode type %v to string", reflect.TypeOf(v))
}

func (a Binary) Value() (driver.Value, error) {
	return []byte(a), nil
}
