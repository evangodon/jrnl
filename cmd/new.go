package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/db"
	"github.com/evangodon/jrnl/internal/util"

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
		date := c.Timestamp("date")

		payload := struct {
			Daily db.Journal `json:"daily"`
		}{}

		res, err := api.MakeRequest(http.MethodGet, "/daily", nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(res.Body, &payload)
		if err != nil {
			return err
		}

		existingContent := func() string {
			if payload.Daily.ID != "" && strings.TrimSpace(payload.Daily.Content) != "" {
				println(len(payload.Daily.Content))
				return payload.Daily.Content
			}
			formattedDate := date.Format("Monday, January 2 2006")
			return "# " + formattedDate + "\n\n"
		}()

		newContent := util.GetNewEntry(existingContent)

		if newContent == existingContent {
			fmt.Println("No changes made")
			return nil
		}

		// Daily already exists for this day
		if payload.Daily.ID != "" {
			updatedEntry, err := json.Marshal(db.Journal{
				ID:        payload.Daily.ID,
				CreatedAt: payload.Daily.CreatedAt,
				Content:   newContent,
			})
			if err != nil {
				return err
			}
			_, err = api.MakeRequest(
				http.MethodPatch,
				"/daily/"+payload.Daily.ID,
				bytes.NewBuffer(updatedEntry),
			)
			if err != nil {
				return err
			}
			fmt.Println("Entry updated")
			return nil
		}

		// Create new daily
		newEntry, err := json.Marshal(db.Journal{
			Content: newContent,
		})
		if err != nil {
			return err
		}

		_, err = api.MakeRequest(http.MethodPost, "/daily/new", bytes.NewBuffer(newEntry))
		if err != nil {
			return err
		}

		fmt.Printf("New entry created for %s", date.Format("Monday, January 2 2006"))
		return nil
	},
}
