package util

import "time"

const LOCAL_TIME_ZONE = "America/New_York"

// FormatToLocalTime formats a time.Time to a string using the given format and the local time zone.
func FormatToLocalTime(date time.Time, format string) string {
	location, err := time.LoadLocation(LOCAL_TIME_ZONE)
	CheckError(err)

	return date.In(location).Format(format)
}
