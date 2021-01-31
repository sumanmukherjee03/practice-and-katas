package date_utils

import "time"

const (
	dateFormat string = "2006-01-02T15:04:05Z" // This is a special string denoting the format and needs to be kept this way. Dont use any other random date or time here.
)

func GetNow() time.Time {
	return time.Now().UTC()
}

func GetNowString() string {
	t := GetNow()
	return t.Format(dateFormat)
}
