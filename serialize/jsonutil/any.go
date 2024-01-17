package jsonutil

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Type int8

const (
	Unknown Type = iota
	Null
	Text
	Number
	Boolean
	Object
	Array
)

func (t Type) String() string {
	switch t {
	case Unknown:
		return "unknown"
	case Null:
		return "null"
	case Text:
		return "text"
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case Object:
		return "object"
	case Array:
		return "array"
	}
	return ""
}

var ErrUnknownType = errors.New("unknown field type")

type Any struct {
	value string
	Type  Type
}

type AnyArray []*Any
type AnyMap map[string]*Any

func (a *Any) UnmarshalJSON(b []byte) error {
	val := string(b)
	if len(val) == 0 {
		return fmt.Errorf("unexpected json input")
	}
	if len(val) == 1 {
		if b[0] >= '0' && b[0] <= '9' {
			a.Type = Number
			a.value = string(b)
			return nil
		} else {
			return fmt.Errorf("unexpected json input")
		}
	}
	if bytes.Equal(b, []byte("null")) {
		a.Type = Null
		a.value = string(b)
		return nil
	}

	if b[0] == '[' && b[len(b)-1] == ']' {
		a.Type = Array
		a.value = string(b)
		return nil
	}

	if b[0] == '{' && b[len(b)-1] == '}' {
		a.Type = Object
		a.value = string(b)
		return nil
	}

	if b[0] == '"' && b[len(b)-1] == '"' {
		a.Type = Text
		a.value = string(b[1 : len(b)-1])
		return nil
	}
	if isBoolean(b) {
		a.Type = Boolean
		a.value = string(b)
		return nil
	}
	if isNumber(b) {
		a.Type = Number
		a.value = string(b)
		return nil
	}
	a.Type = Unknown
	return fmt.Errorf("unexpected json input")
}

func isNumber(b []byte) bool {
	dot := false
	for _, v := range b {
		if v == '.' {
			if !dot {
				dot = true
				continue
			} else {
				return false
			}
		}
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

func isBoolean(b []byte) bool {
	return bytes.Equal(b, []byte("true")) || bytes.Equal(b, []byte("false"))
}

func (a Any) MarshalJSON() ([]byte, error) {
	data := []byte(a.value)
	if a.Type == Text {
		buf := make([]byte, len(a.value)+2)
		buf = append(buf, '"')
		buf = append(buf, data...)
		buf = append(buf, '"')
		return buf, nil
	}
	return data, nil
}

func (a *Any) Value() (any, error) {
	switch a.Type {
	case Null:
		return nil, nil
	case Text:
		return a.value, nil
	case Number:
		if strings.Contains(a.value, ".") {
			v, _ := strconv.ParseFloat(a.value, 64)
			return v, nil
		} else {
			v, _ := strconv.ParseInt(a.value, 10, 64)
			return v, nil
		}
	case Boolean:
		return a.value == "true", nil
	case Object:
		var dst = make(map[string]*Any)
		if err := Unmarshal([]byte(a.value), &dst); err != nil {
			return nil, err
		} else {
			return dst, nil
		}
	case Array:
		var dst []*Any
		if err := Unmarshal([]byte(a.value), &dst); err != nil {
			return nil, err
		} else {
			return dst, nil
		}
	}
	return nil, ErrUnknownType
}

func (a *Any) Object() (map[string]*Any, error) {
	if a.Type != Object {
		return nil, fmt.Errorf("content type is not object, but %v", a.Type)
	}
	var dst = make(map[string]*Any)
	if err := Unmarshal([]byte(a.value), &dst); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

func (a *Any) Array() ([]*Any, error) {
	if a.Type != Array {
		return nil, fmt.Errorf("content type is not array, but %v", a.Type)
	}
	var dst []*Any
	if err := Unmarshal([]byte(a.value), &dst); err != nil {
		return nil, err
	} else {
		return dst, nil
	}
}

func (a *Any) Int64() (int64, error) {
	if a.Type != Number {
		return 0, fmt.Errorf("content type is not number, but %v", a.Type)
	}
	if strings.Contains(a.value, ".") {
		val, err := a.Float64()
		return int64(val), err
	}
	return strconv.ParseInt(a.value, 10, 64)
}

func (a *Any) Float64() (float64, error) {
	if a.Type != Number {
		return 0, fmt.Errorf("content type is not number, but %v", a.Type)
	}
	return strconv.ParseFloat(a.value, 64)
}

func (a *Any) Boolean() (bool, error) {
	if a.Type != Number {
		return false, fmt.Errorf("content type is not boolean, but %v", a.Type)
	}
	return a.value == "true", nil
}

func (a *Any) String() string {
	return a.value
}

func (a *Any) IsNull() bool {
	return a.Type == Null
}

func (a *Any) MustValue() any {
	val, err := a.Value()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustObject() map[string]*Any {
	val, err := a.Object()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustArray() []*Any {
	val, err := a.Array()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustInt64() int64 {
	val, err := a.Int64()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustFloat64() float64 {
	val, err := a.Float64()
	if err != nil {
		panic(err)
	}
	return val
}

func (a *Any) MustBoolean() bool {
	val, err := a.Boolean()
	if err != nil {
		panic(err)
	}
	return val
}
