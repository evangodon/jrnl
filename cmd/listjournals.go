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
	"github.com/urfave/cli/v2"
)

var ListJournalsCmd = &cli.Command{
	Name:    "journals",
	Aliases: []string{"j"},
	Usage:   "List all journal entries",
	Action: func(c *cli.Context) error {

		p := tea.NewProgram(initialListJournalsModel(c))

		if err := p.Start(); err != nil {
			fmt.Printf("Error occured: %v", err)
			return err
		}

		return nil
	},
	Flags: []cli.Flag{
		&cli.BoolFlag{Name: "week", Aliases: []string{"w"}, Usage: "Show the week's entries"},
	},
}

type listJournalsModel struct {
	list ui.List
	ctx  *cli.Context
}

type journalItem struct {
	itemNum   int
	CreatedAt time.Time
	Content   string
}

func (i journalItem) Title() string { return fmt.Sprintf("Journal #%d", i.itemNum) }

func (i journalItem) Description() string {
	return util.FormatToLocalTime(i.CreatedAt, "Monday, January 2, 2006")
}

func (i journalItem) GetContent() string {
	out, err := glamour.Render(i.Content, "dark")
	util.CheckError(err)

	return out
}

func (i journalItem) FilterValue() string { return i.CreatedAt.String() }

func initialListJournalsModel(c *cli.Context) listJournalsModel {
	return listJournalsModel{
		list: ui.CreateList("Journal Entries"),
		ctx:  c,
	}
}

func getJournalEntries(c *cli.Context) tea.Msg {
	var (
		db  sqldb.DB        = sqldb.Connect()
		ctx context.Context = context.Background()
	)

	whereClause := "true"
	if c.Bool("week") {
		whereClause = "created_at >= date('now', 'weekday 0', '-7 days')"
	}

	var journalEntries []journalItem

	err := db.NewSelect().
		Model(&sqldb.Journal{}).
		Column("created_at", "content").
		Order("created_at DESC").
		Where(whereClause).
		Scan(ctx, &journalEntries)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	var items []ui.ListItem
	for index, entry := range journalEntries {
		var item ui.ListItem = journalItem{itemNum: len(journalEntries) - index, CreatedAt: entry.CreatedAt, Content: entry.Content}
		items = append(items, item)
	}

	return items
}

func (m listJournalsModel) Init() tea.Cmd {
	teaCmd := func() tea.Msg {
		return getJournalEntries(m.ctx)
	}

	return teaCmd

}

func (m listJournalsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.list.HandleMessage(msg)

	return m, cmd
}

func (m listJournalsModel) View() string {
	return m.list.View()
}
