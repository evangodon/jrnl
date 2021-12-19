package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/ui"
	"jrnl/pkg/util"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
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

type model struct {
	list ui.List
}

type item struct {
	itemNum   int
	createdAt time.Time
	content   string
}

func (i item) Title() string { return fmt.Sprintf("Journal #%d", i.itemNum) }

func (i item) Description() string {
	return util.FormatToLocalTime(i.createdAt, "Monday, January 2, 2006")
}

func (i item) GetContent() string {
	out, err := glamour.Render(i.content, "dark")
	util.CheckError(err)

	return out
}

func (i item) FilterValue() string { return i.createdAt.String() }

func initialModel() model {
	return model{
		list: ui.CreateList("Journal Entries"),
	}
}

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

	var items []ui.ListItem
	for index, entry := range journalEntries {
		var item ui.ListItem = item{itemNum: len(journalEntries) - index, createdAt: entry.CreatedAt, content: entry.Content}
		items = append(items, item)
	}

	return items
}

func (m model) Init() tea.Cmd {
	return getJournalEntries
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.list.HandleMessage(msg)

	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}
