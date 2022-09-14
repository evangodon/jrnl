package cmd

import (
	"fmt"

	"github.com/evangodon/jrnl/db"
	ui "github.com/evangodon/jrnl/ui"

	lg "github.com/charmbracelet/lipgloss"
	"github.com/urfave/cli/v2"
)

var ShowDBPathCmd = &cli.Command{
	Name:    "showdbpath",
	Aliases: []string{"sdp"},
	Usage:   "Show the path to the database",
	Action: func(_ *cli.Context) error {

		path := lg.NewStyle().
			Background(ui.Color.Primary).
			Padding(0, 2).
			Render(db.GetDBPath())

		fmt.Println("\nDb path: ", path)

		return nil
	},
}
