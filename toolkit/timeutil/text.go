package timeutil

import (
	"strconv"
	"strings"
	"time"
)

type (
	Datetime string
	Date     string
	Clock    string
	UnixSec  int64
	UnixMil  int64
	UnixMic  int64
)

func (d Datetime) Time() Time {
	return MustParse(string(d))
}

func (d Date) Time() Time {
	return MustParse(string(d), "2006-01-02")
}

func (c Clock) Time() Time {
	clock := strings.Split(string(c), ":")
	if len(clock) != 3 {
		return Time{}
	}
	h, _ := strconv.Atoi(strings.TrimLeft(clock[0], "0"))
	m, _ := strconv.Atoi(strings.TrimLeft(clock[1], "0"))
	s, _ := strconv.Atoi(strings.TrimLeft(clock[2], "0"))
	Y, M, D := Now().Date()
	return Time{Time: time.Date(Y, M, D, h, m, s, 0, time.Local)}
}

func (u UnixSec) Time() Time {
	return Time{Time: time.Unix(int64(u), 0)}
}

func (u UnixMil) Time() Time {
	return Time{Time: time.UnixMilli(int64(u))}
}

func (u UnixMic) Time() Time {
	return Time{Time: time.UnixMicro(int64(u))}
}
