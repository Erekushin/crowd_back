package data

import (
	"database/sql/driver"
	"time"
)

type LocalTime time.Time

const localTimeFormat = "2006-01-02 15:04:05"

func (t *LocalTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+localTimeFormat+`"`, string(data), time.Local)
	*t = LocalTime(now)
	return
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(localTimeFormat)+2)
	b = append(b, '"')
	b = append(b, []byte(t.String())...)
	b = append(b, '"')

	return b, nil
}

func (t LocalTime) String() string {
	if time.Time(t).IsZero() {
		return "0000-00-00 00:00:00"
	}

	return time.Time(t).Format(localTimeFormat)
}

func (t LocalTime) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return time.Now(), nil
	}
	return time.Time(t), nil
}

func (t *LocalTime) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = LocalTime(vt)
	case string:
		tTime, _ := time.Parse("2006/01/02 15:04:05", vt)
		*t = LocalTime(tTime)
	default:
		return nil
	}
	return nil
}
