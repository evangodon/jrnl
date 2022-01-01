package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"
	"log"
	"regexp"

	"github.com/urfave/cli/v2"
)

var NewCmd = &cli.Command{
	Name:    "new",
	Aliases: []string{"n"},
	Usage:   "Create a new journal entry",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "date",
			Aliases:  []string{"d"},
			Value:    "date",
			Usage:    "Date of the entry",
			Required: true,
		},
	},
	Action: func(c *cli.Context) error {
		date := c.String("date")
		r, _ := regexp.Compile("[0-9]{4}-[0-9]{2}-[0-9]{2}")

		if !r.MatchString(date) {
			return cli.Exit("Invalid date format. Use YYYY-MM-DD", 1)
		}

		var (
			db              sqldb.DB        = sqldb.Connect()
			ctx             context.Context = context.Background()
			existingEntryId string          = ""
			existingContent string          = ""
		)

		err := db.NewSelect().
			Model(&sqldb.Journal{}).
			Column("id", "content").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", date)).
			Scan(ctx, &existingEntryId, &existingContent)

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
		if existingEntryId != "" {
			id = existingEntryId
		} else {
			id = sqldb.CreateId()
		}

		journalEntry := sqldb.Journal{
			Id:        id,
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

		if existingEntryId != "" {
			fmt.Println("Entry updated")
		} else {
			fmt.Println("Entry added")
		}

		return nil

	},
}
