package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/evangodon/jrnl/db"
	"github.com/evangodon/jrnl/util"

	"github.com/urfave/cli/v2"
)

var EditCmd = &cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit an entry",
	Before: func(c *cli.Context) error {
		if c.NArg() == 0 {
			return cli.Exit("Please provide an identifier", 1)
		}
		return nil
	},
	Action: func(c *cli.Context) error {
		var (
			dbClient   = db.Connect()
			identifier = c.Args().Get(0)
			item       = db.Journal{}
			ctx        = context.Background()
		)

		// Using row number as identifier
		if rowNumber, err := strconv.Atoi(identifier); err == nil {

			i, err := dbClient.SelectEntryByRowNumber(ctx, &item, rowNumber)
			util.CheckIfNoRowsFound(err, "No entry found at row "+identifier)

			item.ID = i.ID
			item.Content = i.Content
		}

		//  Using id as identifier
		if len(identifier) == db.IDLength {
			i, err := dbClient.SelectEntryByID(ctx, &item, identifier)

			util.CheckIfNoRowsFound(err, "No entry found at row "+identifier)
			item.ID = i.ID
			item.Content = i.Content
		}

		if item.ID == "" {
			return cli.Exit("No entry found", 1)
		}

		editedContent := util.GetNewEntry(item.Content)

		if editedContent == item.Content {
			fmt.Println("No changes made")
			return nil
		}

		err := dbClient.UpdateEntryContent(ctx, &item, db.Item{
			ID:      item.ID,
			Content: editedContent,
		})

		util.CheckError(err)

		return nil
	},
}
