package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

// TODO: write to file if entry already exists
var TodayCmd = &cli.Command{
	Name:    "today",
	Aliases: []string{"n"},
	Usage:   "Create a new journal entry for today",
	Action: func(c *cli.Context) error {

		ctx := context.Background()
		db := sqldb.Connect()
		existingEntryId := ""
		var existingContent string

		err := db.NewSelect().
			Model(&sqldb.JournalEntry{}).
			Column("id", "content").
			Where("DATE(created_at) = DATE('now')").
			Scan(ctx, &existingEntryId, &existingContent)

		if err != nil {
			if err == sql.ErrNoRows {
			} else {
				log.Fatal(err)
			}
		}

		content := util.GetNewEntry(existingContent)

		if content == "" {
			fmt.Println("No content found. Exiting.")
			os.Exit(0)
		}

		var id string
		if existingEntryId != "" {
			id = existingEntryId
		} else {
			id = sqldb.CreateId()
		}

		entry := sqldb.JournalEntry{
			Id:      id,
			Content: content,
		}

		_, err = db.NewInsert().
			Model(&entry).
			On("CONFLICT (id) DO UPDATE").
			Returning("id").
			Exec(ctx)

		util.CheckError(err)

		return nil
	},
}
