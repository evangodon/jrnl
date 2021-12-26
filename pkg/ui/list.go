package ui

import (
	bubbleList "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	lg "github.com/charmbracelet/lipgloss"
)

type ListItem interface {
	Title() string
	Description() string
	GetContent() string
	FilterValue() string
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
	return l.model.Cursor()
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

	content := ""
	if len(l.entries) > 0 {
		content = l.entries[l.Cursor()].GetContent()
	}

	leftSide := lg.NewStyle().
		MarginRight(3).
		MarginTop(leftSidemarginTop).
		Render(list)

	rightSide := ""
	if len(content) > 0 {
		rightSide = lg.NewStyle().
			MarginTop(2).
			BorderStyle(lg.RoundedBorder()).
			BorderForeground(ColorGreyLight).
			Height(l.height - 8).
			Render(content)
	}

	return lg.JoinHorizontal(lg.Top, leftSide, rightSide)
}
