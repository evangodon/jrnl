package cmd

import (
	"github.com/urfave/cli/v2"
)

var ListCmd = &cli.Command{
	Name:    "list",
	Aliases: []string{"l"},
	Usage:   "List entries",
	Subcommands: []*cli.Command{
		ListJournalsCmd,
		ListTilsCmd,
	},
}
