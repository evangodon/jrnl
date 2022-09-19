package ui

import (
	"fmt"
	"time"

	"github.com/evangodon/jrnl/internal/util"

	lg "github.com/charmbracelet/lipgloss"
)

func CreateStreakLine(allEntries map[int]ListItem, activeEntry ListItem) string {
	result := activeEntry.GetCreatedAt().Format("January 2006")

	numberOfdays := util.GetNumberOfDaysInMonth(activeEntry.GetCreatedAt())
	streakLine := ""

	allEntriesOfMonth := map[int]ListItem{}
	for _, e := range allEntries {
		if e.GetCreatedAt().Month() == activeEntry.GetCreatedAt().Month() {
			allEntriesOfMonth[e.GetCreatedAt().Day()] = e
		}
	}

	for i := 1; i <= numberOfdays+1; i++ {
		p := time.Date(
			activeEntry.GetCreatedAt().Year(),
			activeEntry.GetCreatedAt().Month(),
			i,
			0,
			0,
			0,
			0,
			time.UTC,
		)
		if util.IsToday(p) {
			var color lg.TerminalColor
			if activeEntry.GetCreatedAt().Day() == i {
				color = Color.Primary
			} else {
				color = Color.Secondary
			}

			streakLine += lg.NewStyle().Foreground(color).Render("▣ ")
			continue
		}

		if activeEntry.GetCreatedAt().Day() == i {
			streakLine += lg.NewStyle().Foreground(Color.Primary).Render("● ")
			continue
		}

		entry := allEntriesOfMonth[i]

		if entry != nil {
			streakLine += lg.NewStyle().Foreground(Color.White).Render("● ")
			continue
		}

		streakLine += lg.NewStyle().Foreground(Color.GreyLight).Render("○ ")
	}

	result += fmt.Sprintf("\n%s", streakLine)

	return lg.NewStyle().
		MarginTop(2).
		Render(result)
}
