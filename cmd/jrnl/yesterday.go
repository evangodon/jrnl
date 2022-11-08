package main

import (
	"time"

	"github.com/urfave/cli/v2"
)

var YesterdayCmd = &cli.Command{
	Name:    "yesterday",
	Aliases: []string{"y"},
	Usage:   "Create a new journal entry for yesterday",
	Action: func(c *cli.Context) error {
		yesterdayDate := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

		return c.App.Run([]string{c.App.Name, "new", "--date", yesterdayDate})
	},
}
