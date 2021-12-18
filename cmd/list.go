package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/ui"
	"jrnl/pkg/util"
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"
)

var ListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "Create a new journal entry for today",
	Action: func(c *cli.Context) error {

		p := tea.NewProgram(initialModel())
		p.EnterAltScreen()

		if err := p.Start(); err != nil {
			fmt.Printf("Error occured: %v", err)
			return err
		}

		return nil
	},
}

type window struct {
	height int
	width  int
}

type Entry struct {
	content string
}

type model struct {
	list    list.Model
	entries map[int]Entry
	window  window
}

func initialModel() model {
	return model{
		list:    list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0),
		entries: map[int]Entry{},
	}
}

type entries []sqldb.JournalEntry

func getJournalEntries() tea.Msg {
	var (
		db  *bun.DB         = sqldb.Connect()
		ctx context.Context = context.Background()
	)

	var journalEntries []sqldb.JournalEntry

	err := db.NewSelect().
		Model(&sqldb.JournalEntry{}).
		Column("id", "created_at", "content").
		Order("created_at DESC").
		Scan(ctx, &journalEntries)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	return entries(journalEntries)
}

func (m model) Init() tea.Cmd {
	return getJournalEntries
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {

		case "ctrl+c", "q":
			return m, tea.Quit

		case "up", "k":
			m.list.CursorUp()

		case "down", "j":
			m.list.CursorDown()

		case "?":
			m.list.SetShowHelp(true)
		}

	case tea.WindowSizeMsg:
		m.window.width = msg.Width
		m.window.height = msg.Height

	case entries:
		m.list.Title = "Journal Entries"
		m.list.Styles.Title = lipgloss.NewStyle().
			Foreground(ui.ColorWhite).
			Background(ui.ColorPrimary).
			Padding(0, 1)
		m.list.SetHeight(m.window.height)

		entries := msg
		for index, entry := range entries {
			date := util.FormatToLocalTime(entry.CreatedAt, "Monday, January 2, 2006")

			m.list.InsertItem(index, item{
				title: fmt.Sprintf("Journal #%v", len(entries)-index),
				date:  date,
			})
			m.entries[index] = Entry{
				content: entry.Content,
			}
		}
	}

	return m, nil
}

type item struct {
	title string
	date  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.date }
func (i item) FilterValue() string { return i.date }

func (m model) View() string {

	list := m.list.View()
	content := m.entries[m.list.Cursor()].content

	renderedContent, err := glamour.Render(content, "dark")
	util.CheckError(err)

	leftSide := lipgloss.NewStyle().
		MarginRight(3).
		Render(list)
	rightSide := lipgloss.NewStyle().
		Padding(2, 1).
		Height(m.window.height - 2).
		Render(renderedContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftSide, rightSide)
}
