package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/sqldb"
	"jrnl/util"
	"log"

	"github.com/urfave/cli/v2"
)

var NewCmd = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "Create a new journal entry",
	Flags: []cli.Flag{
		&cli.TimestampFlag{
			Name:     "date",
			Aliases:  []string{"d"},
			Usage:    "Date of the entry",
			Required: true,
			Layout:   "2006-01-02",
		},
	},
	Action: func(c *cli.Context) error {
		date := c.Timestamp("date").Format("2006-01-02")
		fmt.Println(date)

		var (
			db              = sqldb.Connect()
			ctx             = context.Background()
			existingEntryID = ""
			existingContent = ""
		)

		err := db.NewSelect().
			Model(&sqldb.Journal{}).
			Column("id", "content").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", date)).
			Scan(ctx, &existingEntryID, &existingContent)

		if err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
		}

		var entryDate = util.CreateTimeDate(date)
		if existingContent == "" {
			formattedDate := entryDate.Format("Monday, January 2 2006")
			existingContent = "# " + formattedDate + "\n\n"
		}

		content := util.GetNewEntry(existingContent)

		if content == existingContent {
			return cli.Exit("No changes were made", 0)
		}

		var id string
		if existingEntryID != "" {
			id = existingEntryID
		} else {
			id = sqldb.CreateId()
		}

		journalEntry := sqldb.Journal{
			ID:        id,
			CreatedAt: entryDate,
			Content:   content,
		}

		_, err = db.NewInsert().
			Model(&journalEntry).
			On("CONFLICT (id) DO UPDATE").
			Set("updated_at = EXCLUDED.updated_at").
			Set("content = EXCLUDED.content").
			Exec(ctx)

		util.CheckError(err)

		if existingEntryID != "" {
			fmt.Println("Entry updated")
		} else {
			fmt.Println("Entry added")
		}

		return nil

	},
}
