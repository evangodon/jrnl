package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/evangodon/jrnl/internal/util"

	lg "github.com/charmbracelet/lipgloss"
)

var (
	dot       = "● "
	activeDot = "▣ "
)

func CreateStreakLine(allEntries []list.Item, activeEntry ListItem) string {
	result := activeEntry.GetCreatedAt().Format("January 2006")

	numberOfdays := util.GetNumberOfDaysInMonth(activeEntry.GetCreatedAt())
	streakLine := ""

	allEntriesOfMonth := map[int]ListItem{}
	for _, e := range allEntries {
		journalItem, ok := e.(JournalItem)
		if !ok {
			panic("Couldn't cast item")
		}

		if journalItem.GetCreatedAt().Month() == activeEntry.GetCreatedAt().Month() {
			allEntriesOfMonth[journalItem.GetCreatedAt().Day()] = journalItem
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

			streakLine += lg.NewStyle().Foreground(color).Render(activeDot)
			continue
		}

		if activeEntry.GetCreatedAt().Day() == i {
			streakLine += lg.NewStyle().Foreground(Color.Primary).Render(dot)
			continue
		}

		entry := allEntriesOfMonth[i]

		if entry != nil {
			streakLine += lg.NewStyle().Foreground(Color.White).Render(dot)
			continue
		}

		streakLine += lg.NewStyle().Foreground(Color.GreyLight).Render(dot)
	}

	result += fmt.Sprintf("\n%s", streakLine)

	return lg.NewStyle().
		MarginTop(2).
		Render(result)
}
