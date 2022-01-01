package cmd

import (
	"flag"
	"time"

	"github.com/urfave/cli/v2"
)

var YesterdayCmd = &cli.Command{
	Name:    "yesterday",
	Aliases: []string{"y"},
	Usage:   "Create a new journal entry for yesterday",
	Action: func(c *cli.Context) error {

		yesterdayDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		set := flag.NewFlagSet("date", 0)
		set.String("date", yesterdayDate, "Date of yesterday")

		nc := cli.NewContext(c.App, set, nil)

		nc.Command.Name = "new"

		return NewCmd.Action(nc)

	},
}
