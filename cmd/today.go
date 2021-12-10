package cmd

import (
	"context"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"

	"github.com/urfave/cli/v2"
)

var TodayCmd = &cli.Command{
	Name:    "today",
	Aliases: []string{"n"},
	Usage:   "Create a new journal entry for today",
	Action: func(c *cli.Context) error {

		content := util.OpenEditor()

		db := sqldb.Connect()
		ctx := context.Background()

		entry := sqldb.EntryModel{
			Id:      sqldb.CreateId(),
			Content: content,
		}

		_, err := db.NewInsert().Model(&entry).Returning("id").Exec(ctx)

		util.CheckError(err)

		return nil
	},
}
