package util

import "time"

func FormatToLocalTime(date time.Time, format string) string {
	location, err := time.LoadLocation("America/New_York")
	CheckError(err)

	return date.In(location).Format(format)
}
