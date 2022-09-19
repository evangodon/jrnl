package cmd

import (
	"time"

	"github.com/urfave/cli/v2"
)

var TodayCmd = &cli.Command{
	Name:    "today",
	Aliases: []string{"t"},
	Usage:   "Create a new journal entry for today",
	Action: func(c *cli.Context) error {
		todayDate := time.Now().Format("2006-01-02")

		return c.App.Run([]string{c.App.Name, "new", "--date", todayDate})
	},
}
