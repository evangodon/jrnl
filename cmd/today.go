package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"
	"log"
	"os"
	"time"

	"github.com/uptrace/bun"
	"github.com/urfave/cli/v2"
)

var TodayCmd = &cli.Command{
	Name:    "today",
	Aliases: []string{"n"},
	Usage:   "Create a new journal entry for today",
	Action: func(c *cli.Context) error {

		var (
			db              *bun.DB         = sqldb.Connect()
			ctx             context.Context = context.Background()
			existingEntryId string          = ""
			existingContent string          = ""
		)

		err := db.NewSelect().
			Model(&sqldb.JournalEntry{}).
			Column("id", "content").
			Where("DATE(created_at, 'localtime') = DATE('now', 'localtime')").
			Scan(ctx, &existingEntryId, &existingContent)

		if err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
		}

		if existingContent == "" {
			todayDate := time.Now().Format("Monday, January 2 2006")
			existingContent = "# " + todayDate + "\n\n"
		}

		content := util.GetNewEntry(existingContent)

		if content == "" {
			os.Exit(0)
		}

		if content == existingContent {
			fmt.Println("No changes made.")
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
			Set("updated_at = EXCLUDED.updated_at").
			Set("content = EXCLUDED.content").
			Returning("id").
			Exec(ctx)

		util.CheckError(err)

		if existingEntryId != "" {
			fmt.Println("Today's entry updated.")
		} else {
			fmt.Println("Entry added for today.")
		}

		return nil
	},
}
