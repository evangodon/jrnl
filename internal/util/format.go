package util

import (
	"strings"
	"time"

	"github.com/evangodon/jrnl/internal/db"
)

var dateFormat = "January 2 2006"
var hourFormat = "3:04 AM"

// Checks if daily has existing content, and if not, will format a new one.
func FormatContent(daily db.Journal, now time.Time) string {
	if daily.ID != "" && strings.TrimSpace(daily.Content) != "" {

		content := daily.Content

		if now.Format(dateFormat) == daily.UpdatedAt.Format(dateFormat) {
			fiveMinutesAgo := now.Add(-5 * time.Minute)

			println(
				"five minutes ago",
				daily.UpdatedAt.Format(hourFormat),
				fiveMinutesAgo.Format(hourFormat),
			)

			if daily.UpdatedAt.Before(fiveMinutesAgo) {
				return content + "\n-- " + now.Format(hourFormat)
			}
		}

		return daily.Content
	}

	return "# " + now.Format(dateFormat) + "\n\n" + now.Format(hourFormat)
}
