package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	db "github.com/evangodon/jrnl/db"
	ui "github.com/evangodon/jrnl/ui"
	"github.com/evangodon/jrnl/util"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/urfave/cli/v2"
)

var ListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List all journal entries",
	Action: func(c *cli.Context) error {

		p := tea.NewProgram(initialListJournalsModel(c))
		tea.EnterAltScreen()

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

func (i journalItem) FilterValue() string     { return i.CreatedAt.String() }
func (i journalItem) GetCreatedAt() time.Time { return i.CreatedAt }

func initialListJournalsModel(c *cli.Context) listJournalsModel {
	return listJournalsModel{
		list: ui.CreateList("Journal Entries"),
		ctx:  c,
	}
}

func getJournalEntries(c *cli.Context) tea.Msg {
	var (
		dbClient = db.Connect()
		ctx      = context.Background()
	)

	whereCondition := "true"
	if c.Bool("week") {
		whereCondition = "created_at >= date('now', 'weekday 0', '-7 days')"
	}

	var journalEntries []journalItem

	err := dbClient.NewSelect().
		Model(&db.Journal{}).
		Column("created_at", "content").
		Order("created_at DESC").
		Where(whereCondition).
		Scan(ctx, &journalEntries)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	var items []ui.ListItem
	for index, entry := range journalEntries {
		var item ui.ListItem = journalItem{
			itemNum:   len(journalEntries) - index,
			CreatedAt: entry.CreatedAt,
			Content:   entry.Content,
		}
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
