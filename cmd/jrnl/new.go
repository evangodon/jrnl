package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/evangodon/jrnl/internal/api"
	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/db"
	"github.com/evangodon/jrnl/internal/logger"
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
		config := cfg.GetConfig()
		logger := logger.NewLogger(os.Stdout)

		client := api.Client{
			Config: config,
		}
		entryDate := date.Format("2006-01-02")

		res, err := client.MakeRequest(http.MethodGet, "/daily/"+entryDate, nil)
		if err != nil {
			return err
		}

		payload := struct {
			Daily db.Journal `json:"daily"`
		}{
			Daily: db.Journal{
				CreatedAt: *date,
			},
		}

		err = json.Unmarshal(res.Body, &payload)
		if err != nil {
			return err
		}

		existingContent := util.FormatContent(payload.Daily, time.Now())
		newContent := util.OpenEditorWithContent(existingContent)

		if newContent == existingContent {
			logger.PrintInfo("No changes made")
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
			_, err = client.MakeRequest(
				http.MethodPatch,
				"/daily/"+payload.Daily.ID,
				bytes.NewBuffer(updatedEntry),
			)
			if err != nil {
				return err
			}
			logger.PrintSuccess("Entry updated")
			return nil
		}

		// Create new daily
		newEntry, err := json.Marshal(db.Journal{
			Content: newContent,
		})
		if err != nil {
			return err
		}

		_, err = client.MakeRequest(http.MethodPost, "/daily/new", bytes.NewBuffer(newEntry))
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("New entry created for %s", date.Format("Monday, January 2 2006"))
		logger.PrintSuccess(msg)
		return nil
	},
}
