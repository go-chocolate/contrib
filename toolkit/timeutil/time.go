package timeutil

import "time"

type Time struct {
	time.Time
}

func (t Time) Format(layout ...string) string {
	if t.IsZero() {
		return ""
	}
	if len(layout) > 0 {
		return t.Time.Format(layout[0])
	}
	return t.Time.Format(time.DateTime)
}

func (t Time) FormatWithDefault(def string, layout ...string) string {
	if t.IsZero() {
		return def
	}
	return t.Format(layout...)
}

func (t Time) DateFormat() string {
	return t.Format("2006-01-02")
}

func (t Time) ClockFormat() string {
	return t.Format("15:04:05")
}

func (t Time) TodayBegin() time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

func (t Time) TodayEnd() time.Time {
	today := t.TodayBegin()
	return today.Add(24*time.Hour - 1)
}

func (t Time) TodayRange() [2]time.Time {
	return [2]time.Time{t.TodayBegin(), t.TodayEnd()}
}

func Parse(text string, layout ...string) (Time, error) {
	var lay string
	if len(layout) > 0 {
		lay = layout[0]
	} else {
		lay = time.DateTime
	}
	ti, err := time.ParseInLocation(lay, text, time.Local)
	if err != nil {
		return Time{}, err
	}
	return Time{Time: ti}, nil
}

func MustParse(text string, layout ...string) Time {
	t, _ := Parse(text, layout...)
	return t
}

func Now() Time {
	return Time{Time: time.Now()}
}
