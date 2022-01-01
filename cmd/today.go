package cmd

import (
	"flag"
	"time"

	"github.com/urfave/cli/v2"
)

var TodayCmd = &cli.Command{
	Name:    "today",
	Aliases: []string{"t"},
	Usage:   "Create a new journal entry for today",
	Action: func(c *cli.Context) error {
		todayDate := time.Now().Format("2006-01-02")

		set := flag.NewFlagSet("date", 0)
		set.String("date", todayDate, "Date of today")

		nc := cli.NewContext(c.App, set, nil)

		nc.Command.Name = "new"

		return NewCmd.Action(nc)
	},
}
