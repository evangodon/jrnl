package cmd

import (
	"fmt"
	"jrnl/sqldb"
	"jrnl/ui"

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
			Render(sqldb.GetDBPath())

		fmt.Println("\nDb path: ", path)

		return nil
	},
}
