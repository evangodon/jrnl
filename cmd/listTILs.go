package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/ui"
	"jrnl/pkg/util"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/urfave/cli/v2"
)

var ListTilsCmd = &cli.Command{
	Name:    "til",
	Aliases: []string{"t"},
	Usage:   "List all TILs entries",
	Action: func(c *cli.Context) error {

		p := tea.NewProgram(initialListTilsModel(c))
		p.EnterAltScreen()

		if err := p.Start(); err != nil {
			fmt.Printf("Error occured: %v", err)
			return err
		}

		return nil
	},
}

type listTilsModel struct {
	list ui.List
	ctx  *cli.Context
}

type tilItem struct {
	itemNum   int
	CreatedAt time.Time
	Content   string
}

func (i tilItem) Title() string { return fmt.Sprintf("TIL #%d", i.itemNum) }

func (i tilItem) Description() string {
	description := i.Content
	maxLength := 40

	if len(description) > maxLength {
		description = fmt.Sprintf("%s...", description[:maxLength])
	}

	return strings.TrimSpace(description)

}

func (i tilItem) GetContent() string {
	out, err := glamour.Render(i.Content, "dark")
	util.CheckError(err)

	return out
}

func (i tilItem) FilterValue() string     { return i.CreatedAt.String() }
func (i tilItem) GetCreatedAt() time.Time { return i.CreatedAt }

func initialListTilsModel(c *cli.Context) listTilsModel {
	return listTilsModel{
		list: ui.CreateList("Journal Entries"),
		ctx:  c,
	}
}

func getTilsEntries(c *cli.Context) tea.Msg {
	var (
		db  sqldb.DB        = sqldb.Connect()
		ctx context.Context = context.Background()
	)

	var tilEntries []tilItem

	err := db.NewSelect().
		Model(&sqldb.TIL{}).
		Column("created_at", "content").
		Order("created_at DESC").
		Scan(ctx, &tilEntries)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal(err)
		}
	}

	var items []ui.ListItem
	for index, entry := range tilEntries {
		var item ui.ListItem = tilItem{
			itemNum:   len(tilEntries) - index,
			CreatedAt: entry.CreatedAt,
			Content:   entry.Content,
		}
		items = append(items, item)
	}

	return items
}

func (m listTilsModel) Init() tea.Cmd {
	teaCmd := func() tea.Msg {
		return getTilsEntries(m.ctx)
	}

	return teaCmd

}

func (m listTilsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.list.HandleMessage(msg)

	return m, cmd
}

func (m listTilsModel) View() string {
	return m.list.View()
}
