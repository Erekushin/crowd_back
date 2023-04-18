package data

import (
	"database/sql/driver"
	"time"
)

type LocalDate time.Time

const localDateFormat = "2006-01-02"

func (t *LocalDate) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+localDateFormat+`"`, string(data), time.Local)
	*t = LocalDate(now)
	return
}

func (t LocalDate) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(localDateFormat)+2)
	b = append(b, '"')
	b = append(b, []byte(t.String())...)
	b = append(b, '"')

	return b, nil
}

func (t LocalDate) String() string {
	if time.Time(t).IsZero() {
		return "0000-00-00"
	}

	return time.Time(t).Format(localDateFormat)
}

func (t LocalDate) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return time.Now(), nil
	}
	return time.Time(t), nil
}

func (t *LocalDate) Scan(v interface{}) error {
	switch vt := v.(type) {
	case time.Time:
		*t = LocalDate(vt)
	case string:
		tTime, _ := time.Parse("2006/01/02 15:04:05", vt)
		*t = LocalDate(tTime)
	default:
		return nil
	}
	return nil
}
