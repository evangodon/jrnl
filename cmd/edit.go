package cmd

import (
	"fmt"
	"jrnl/pkg/sqldb"
	"jrnl/pkg/util"
	"strconv"

	"github.com/urfave/cli/v2"
)

var EditCmd = &cli.Command{
	Name:    "edit",
	Aliases: []string{"e"},
	Usage:   "Edit an entry",
	Subcommands: []*cli.Command{
		// JOURNAL
		{
			Name:    "journal",
			Aliases: []string{"j"},
			Usage:   "Edit a journal entry",
			Action: func(c *cli.Context) error {

				if c.NArg() == 0 {
					return cli.Exit("Please provide an identifier", 1)
				}

				var (
					db         sqldb.DB      = sqldb.Connect()
					identifier               = c.Args().Get(0)
					item       sqldb.Journal = sqldb.Journal{}
				)

				// Using row number as identifier
				if rowNumber, err := strconv.Atoi(identifier); err == nil {

					i, err := db.SelectEntryByRowNumber(&item, rowNumber)

					util.CheckIfNoRowsFound(err, "No entry found at row "+identifier)

					item.Id = i.Id
					item.Content = i.Content
				}

				//  Using id as identifier
				if len(identifier) == sqldb.ID_LENGTH {
					i, err := db.SelectEntryById(&item, identifier)

					util.CheckIfNoRowsFound(err, "No entry found at row "+identifier)
					item.Id = i.Id
					item.Content = i.Content
				}

				if item.Id == "" {
					return cli.Exit("No entry found", 1)
				}

				editedContent := util.GetNewEntry(item.Content)

				if editedContent == item.Content {
					fmt.Println("No changes made")
					return nil
				}

				err := db.UpdateEntryContent(&item, sqldb.Item{
					Id:      item.Id,
					Content: editedContent,
				})

				util.CheckError(err)

				return nil
			},
		},
		// TIL
		{
			Name:    "til",
			Aliases: []string{"t"},
			Usage:   "Edit a TIL entry",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					return cli.Exit("Please provide an identifier", 1)
				}

				var (
					db         sqldb.DB  = sqldb.Connect()
					identifier           = c.Args().Get(0)
					item       sqldb.TIL = sqldb.TIL{}
				)

				// Using row number as identifier
				if rowNumber, err := strconv.Atoi(identifier); err == nil {

					i, err := db.SelectEntryByRowNumber(&item, rowNumber)

					util.CheckIfNoRowsFound(err, "No entry found at row "+identifier)

					item.Id = i.Id
					item.Content = i.Content
				}

				//  Using id as identifier
				if len(identifier) == sqldb.ID_LENGTH {
					i, err := db.SelectEntryById(&item, identifier)

					util.CheckIfNoRowsFound(err, "No entry found at row "+identifier)
					item.Id = i.Id
					item.Content = i.Content
				}

				if item.Id == "" {
					return cli.Exit("No entry found", 1)
				}

				editedContent := util.GetNewEntry(item.Content)

				if editedContent == item.Content {
					fmt.Println("No changes made")
					return nil
				}

				err := db.UpdateEntryContent(&item, sqldb.Item{
					Id:      item.Id,
					Content: editedContent,
				})

				util.CheckError(err)

				return nil
			},
		},
	},
}
