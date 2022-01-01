package cmd

import (
	"context"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"

	"github.com/urfave/cli/v2"
)

var TILCmd = &cli.Command{
	Name:  "til",
	Usage: "Create a new entry entry for something you learnt.",
	Action: func(c *cli.Context) error {

		var (
			db  sqldb.DB        = sqldb.Connect()
			ctx context.Context = context.Background()
		)

		content := util.GetNewEntry("")

		id := sqldb.CreateId()

		entry := sqldb.TIL{
			Id:      id,
			Content: content,
		}
		_, err := db.NewInsert().
			Model(&entry).
			Returning("id").
			Exec(ctx)

		util.CheckError(err)

		fmt.Println("TIL entry added.")

		return nil
	},
}
