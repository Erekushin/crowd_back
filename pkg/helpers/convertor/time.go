package convertor

import "time"

//TimeToDateString
func TimeToDateString(t time.Time) string {
	return t.Format("2006-01-02")
}

//TimeDateTimeString
func TimeToDateTimeString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

//TimeTimeString
func TimeToTimeString(t time.Time) string {
	return t.Format("15:04:05")
}
