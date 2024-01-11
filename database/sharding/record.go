package sharding

import (
	"fmt"
	"strconv"
	"time"
)

type Record map[string]any

func (r Record) Int(key string) int {
	val := r[key]
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		i, _ := strconv.Atoi(v)
		return i
	}
	return 0
}

func (r Record) Uint(key string) uint {
	val := r[key]
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return uint(v)
	case int8:
		return uint(v)
	case int16:
		return uint(v)
	case int32:
		return uint(v)
	case int64:
		return uint(v)
	case uint:
		return uint(v)
	case uint8:
		return uint(v)
	case uint16:
		return uint(v)
	case uint32:
		return uint(v)
	case uint64:
		return uint(v)
	case float32:
		return uint(v)
	case float64:
		return uint(v)
	case string:
		i, _ := strconv.ParseUint(v, 10, 64)
		return uint(i)
	}
	return 0
}

func (r Record) String(key string) string {
	val := r[key]
	if val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", val)
	}
}

func (r Record) Float(key string) float64 {
	val := r[key]
	if val == nil {
		return 0
	}
	switch v := val.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case string:
		i, _ := strconv.ParseFloat(v, 64)
		return i
	}
	return 0
}

func (r Record) Boolean(key string) bool {
	val := r[key]
	if val == nil {
		return false
	}
	switch v := val.(type) {
	case bool:
		return v
	case string:
		return v == "true"
	}
	return false
}

func (r Record) Time(key string) time.Time {
	val := r[key]
	if val == nil {
		return time.Time{}
	}
	switch v := val.(type) {
	case time.Time:
		return v
	case string:
		t, _ := time.ParseInLocation(time.DateTime, v, time.Local)
		return t
	case fmt.Stringer:
		t, _ := time.ParseInLocation(time.DateTime, v.String(), time.Local)
		return t
	}
	return time.Time{}
}
