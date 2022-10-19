package util

import (
	"strings"
	"time"

	"github.com/evangodon/jrnl/internal/db"
)

var dateFormat = "January 2 2006"
var hourFormat = "3:04 PM"

// Checks if daily has existing content, and if not, will format a new one.
func FormatContent(daily db.Journal, now time.Time) string {
	if daily.ID != "" && strings.TrimSpace(daily.Content) != "" {

		content := daily.Content

		if now.Format(dateFormat) == daily.UpdatedAt.Format(dateFormat) {
			fiveMinutesAgo := now.Add(-5 * time.Minute)

			if daily.UpdatedAt.Before(fiveMinutesAgo) {
				return content + "\n-- " + now.Format(hourFormat)
			}
		}

		return daily.Content
	}

	// The entry from not today is being updated.
	if now.Format(dateFormat) != daily.CreatedAt.Format(dateFormat) {
		return "# " + daily.CreatedAt.Format(dateFormat) + "\n"
	}

	return "# " + daily.CreatedAt.Format(dateFormat) + "\n\n-- " + now.Format(hourFormat)
}
