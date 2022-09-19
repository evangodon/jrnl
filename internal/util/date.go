package util

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const LocalTimeZone = "America/New_York"

// FormatToLocalTime formats a time.Time to a string using the given format and the local time zone.
func FormatToLocalTime(date time.Time, format string) string {
	location, err := time.LoadLocation(LocalTimeZone)
	CheckError(err)

	return date.In(location).Format(format)
}

// CreateTimeDate creates a time.Time from a string using the YYYY-MM-DD format.
func CreateTimeDate(date string) (time.Time, error) {
	r, _ := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}")

	if !r.MatchString(date) {
		return time.Now(), errors.New("invalid date format used.  Use YYYY-MM-DD")
	}

	dateArgs := strings.Split(date, "-")
	year, _ := strconv.Atoi(dateArgs[0])
	month, _ := strconv.Atoi(dateArgs[1])
	day, _ := strconv.Atoi(dateArgs[2])

	return time.
		Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local), nil
}

// GetNumberOfDaysInMonth calculates the number of days in a given month.
func GetNumberOfDaysInMonth(date time.Time) int {
	year := date.Year()
	month := date.Month()

	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
}

// IsToday checks if the given date is today.
func IsToday(date time.Time) bool {
	n := time.Now()
	return date.Year() == n.Year() && date.Month() == n.Month() && date.Day() == n.Day()
}
