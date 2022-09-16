package cmd

import (
	"fmt"

	"github.com/evangodon/jrnl/db"

	"github.com/urfave/cli/v2"
)

var ShowDBPathCmd = &cli.Command{
	Name:    "showdbpath",
	Aliases: []string{"sdp"},
	Usage:   "Show the path to the database",
	Action: func(_ *cli.Context) error {
		fmt.Println("\nDb path: ", db.GetDBPath())
		return nil
	},
}
