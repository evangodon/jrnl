package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
	"github.com/evangodon/jrnl/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
)

var ListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List all journal dailies",
	Action: func(c *cli.Context) error {
		p := tea.NewProgram(initialModel(c), tea.WithAltScreen())

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

type model struct {
	perPage int
	page    int
	total   int
	uiList  ui.List
	client  api.Client
}

func initialModel(c *cli.Context) model {
	return model{
		perPage: 20,
		page:    1,
		uiList:  ui.NewList("Journal Entries"),
		client: api.Client{
			Config: cfg.GetConfig(),
		},
	}
}

func (m model) getJournalEntries() tea.Msg {
	path := fmt.Sprintf("/daily/list?perPage=%d&page=%d", m.perPage, m.page)
	res, err := m.client.MakeRequest(http.MethodGet, path, nil)
	if err != nil {
		return err
	}

	payload := struct {
		DailyEntries []db.Journal `json:"daily_entries"`
		Total        int          `json:"total"`
	}{}

	err = json.Unmarshal(res.Body, &payload)
	if err != nil {
		return err
	}

	entries := payload.DailyEntries

	var listItems []ui.ListItem
	for index, entry := range entries {
		var listItem ui.ListItem = ui.JournalItem{
			ItemNum:   payload.Total - (len(m.uiList.Model.Items()) + index),
			CreatedAt: entry.CreatedAt,
			Content:   entry.Content,
		}
		listItems = append(listItems, listItem)
	}

	return ui.JournalEntriesRes{
		ListItems: listItems,
		Total:     payload.Total,
	}
}

func (m model) Init() tea.Cmd {
	teaCmd := func() tea.Msg {
		return m.getJournalEntries()
	}

	return teaCmd
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmd := m.uiList.HandleMessage(msg)
	localItemsNum := len(m.uiList.Model.Items())
	getMore := localItemsNum > 0 && m.uiList.Model.Paginator.OnLastPage() && m.total > localItemsNum

	if getMore {
		m.page += 1
		return m, func() tea.Msg {
			return m.getJournalEntries()
		}
	}

	switch msg := msg.(type) {
	case ui.JournalEntriesRes:
		m.total = msg.Total
	}

	return m, cmd
}

func (m model) View() string {
	return m.uiList.View()
}
