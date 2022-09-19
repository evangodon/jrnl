package api

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/evangodon/jrnl/internal/db"
	"github.com/uptrace/bunrouter"
)

func (app Application) getDailyHandler() bunrouter.HandlerFunc {
	var (
		ctx             = context.Background()
		existingEntryID = ""
		existingContent = ""
	)

	return func(w http.ResponseWriter, req bunrouter.Request) error {
		// TODO: Get date or entry identifier, or use current date as default
		// Return entry if it exists
		var date time.Time
		if true {
			date = time.Now()
		}

		err := app.DbClient.NewSelect().
			Model(&db.Journal{}).
			Column("id", "content").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", date.Format("2006-01-02"))).
			Scan(ctx, &existingEntryID, &existingContent)

		if err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
				return err
			}
		}

		if existingContent == "" {
			app.writeJSON(w, http.StatusNotFound, envelope{"daily": nil}, nil)
			return nil
		}

		dailyEntry := db.Journal{
			ID:        existingEntryID,
			CreatedAt: date,
			Content:   existingContent,
		}

		app.writeJSON(w, http.StatusNotFound, envelope{"daily": dailyEntry}, nil)
		return nil
	}
}

func (Application) newDailyHandler() bunrouter.HandlerFunc {
	return func(_ http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}

func (Application) listDailyHandler() bunrouter.HandlerFunc {
	return func(_ http.ResponseWriter, req bunrouter.Request) error {
		return nil
	}
}
