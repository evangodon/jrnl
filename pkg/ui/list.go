package ui

import (
	"fmt"
	"jrnl/pkg/util"
	"time"

	bubbleList "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ListItem interface {
	Title() string
	Description() string
	GetContent() string
	FilterValue() string
	GetCreatedAt() time.Time
}

type List struct {
	model   bubbleList.Model
	entries map[int]ListItem
	height  int
}

func CreateList(title string) List {
	list := List{}
	list.model = bubbleList.NewModel([]bubbleList.Item{}, bubbleList.NewDefaultDelegate(), 0, 0)
	list.model.Title = title
	list.model.Styles.Title = lg.NewStyle().
		Foreground(Color.White).
		Background(Color.Primary).
		Padding(0, 1)

	list.entries = map[int]ListItem{}

	return list
}

func (l *List) AddItems(listItems []ListItem) {
	for index, listItem := range listItems {

		l.model.InsertItem(index, listItem)

		l.entries[index] = listItem
	}
}

func (l *List) Cursor() int {
	return (l.model.Paginator.Page * l.model.Paginator.PerPage) + l.model.Cursor()
}

func (l *List) SetHeight(height int) {
	l.height = height
	l.model.SetHeight(height)
}

func (l *List) HandleMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return tea.Quit

		case "up", "k":
			l.model.CursorUp()

		case "down", "j":
			l.model.CursorDown()

		case "?":
			l.model.SetShowHelp(true)
		}

	case tea.WindowSizeMsg:
		l.SetHeight(msg.Height - leftSidemarginTop)

	case []ListItem:
		l.AddItems(msg)
	}

	return nil
}

const leftSidemarginTop = 3

func (l *List) View() string {
	list := l.model.View()

	if len(l.entries) == 0 {
		return ""
	}

	activeEntry := l.entries[l.Cursor()]

	leftSide := lg.NewStyle().
		MarginRight(3).
		MarginTop(leftSidemarginTop).
		Render(list)

	topRightSide := ""
	if activeEntry != nil && len(activeEntry.GetContent()) > 0 {
		topRightSide = lg.NewStyle().
			MarginTop(2).
			BorderStyle(lg.RoundedBorder()).
			BorderForeground(Color.GreyLight).
			Height(l.height - 8).
			Render(activeEntry.GetContent())
	}

	bottomRightSide := createStreakLine(l.entries, activeEntry)
	rightSide := lg.JoinVertical(lg.Left, topRightSide, bottomRightSide)

	return lg.JoinHorizontal(lg.Top, leftSide, rightSide)
}

func createStreakLine(allEntries map[int]ListItem, activeEntry ListItem) string {
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
		p := time.Date(activeEntry.GetCreatedAt().Year(), activeEntry.GetCreatedAt().Month(), i, 0, 0, 0, 0, time.UTC)
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
