package cmd

import (
	"context"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"

	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"
)

var TILCmd = &cli.Command{
	Name:    "til",
	Aliases: []string{"t"},
	Usage:   "Create a new entry entry for something you learnt.",
	Action: func(c *cli.Context) error {

		var (
			db  *bun.DB         = sqldb.Connect()
			ctx context.Context = context.Background()
		)

		content := util.GetNewEntry("")

		id := sqldb.CreateId()

		entry := sqldb.Entry{
			Id:      id,
			Type:    sqldb.EntryType.TIL,
			Content: content,
		}
		_, err := db.NewInsert().
			Model(&entry).
			Returning("id").
			Exec(ctx)

		util.CheckError(err)

		return nil
	},
}
