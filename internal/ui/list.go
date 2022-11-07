package ui

import (
	"fmt"
	"strconv"
	"time"

	bubbleList "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	lg "github.com/charmbracelet/lipgloss"

	"github.com/evangodon/jrnl/internal/util"
)

type JournalItem struct {
	ItemNum   int
	CreatedAt time.Time
	Content   string
}

func (i JournalItem) Title() string { return fmt.Sprintf("Journal #%d", i.ItemNum) }

func (i JournalItem) Description() string {
	return util.FormatToLocalTime(i.CreatedAt, "Monday, January 2, 2006")
}

func (i JournalItem) GetContent() string {
	out, err := glamour.Render(i.Content, "dark")
	util.CheckError(err)

	return out
}

func (i JournalItem) FilterValue() string     { return i.CreatedAt.String() }
func (i JournalItem) GetCreatedAt() time.Time { return i.CreatedAt }

func (i JournalItem) GetItemIndex() int {
	date := i.CreatedAt.Format("02012006")

	index, err := strconv.Atoi(date)
	if err != nil {
		panic(err)
	}
	return index
}

type ListItem interface {
	Title() string
	Description() string
	GetContent() string
	FilterValue() string
	GetCreatedAt() time.Time
	GetItemIndex() int
}

type List struct {
	Model  bubbleList.Model
	height int
}

func NewList(title string) List {
	list := List{}
	list.Model = bubbleList.NewModel([]bubbleList.Item{}, bubbleList.NewDefaultDelegate(), 0, 0)
	list.Model.Title = title
	list.Model.Styles.Title = lg.NewStyle().
		Foreground(Color.White).
		Background(Color.Primary).
		Padding(0, 1)

	return list
}

type JournalEntriesRes struct {
	ListItems []ListItem
	Total     int `json:"total"`
}

func (l *List) AddItems(listItems []ListItem) {
	for _, listItem := range listItems {
		l.Model.InsertItem(listItem.GetItemIndex(), listItem)
	}
}

func (l *List) SetHeight(height int) {
	l.height = height
	l.Model.SetHeight(height)
}

func (l *List) HandleMessage(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return tea.Quit

		case "up", "k":
			l.Model.CursorUp()

		case "down", "j":
			l.Model.CursorDown()

		case "?":
			l.Model.SetShowHelp(true)
		}

	case tea.WindowSizeMsg:
		l.SetHeight(msg.Height - leftSidemarginTop)

	case JournalEntriesRes:
		l.AddItems(msg.ListItems)
	}

	return nil
}

const leftSidemarginTop = 3

func (l *List) View() string {
	list := l.Model.View()
	items := l.Model.Items()

	if len(items) == 0 {
		return "No entries created yet"
	}

	listItem := items[l.Model.Cursor()]
	activeEntry, ok := listItem.(JournalItem)
	if !ok {
		panic("Couldn't cast item")
	}

	leftSide := lg.NewStyle().
		MarginRight(3).
		MarginTop(leftSidemarginTop).
		Render(list)

	topRightSide := ""
	if len(activeEntry.GetContent()) > 0 {
		topRightSide = lg.NewStyle().
			MarginTop(2).
			BorderStyle(lg.RoundedBorder()).
			BorderForeground(Color.GreyLight).
			Height(l.height - 8).
			Render(activeEntry.GetContent())
	}

	bottomRightSide := CreateStreakLine(l.Model.Items(), activeEntry)
	rightSide := lg.JoinVertical(lg.Left, topRightSide, bottomRightSide)

	return lg.JoinHorizontal(lg.Top, leftSide, rightSide)
}
