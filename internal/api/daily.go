package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/evangodon/jrnl/internal/db"
	"github.com/evangodon/jrnl/internal/util"
	"github.com/uptrace/bunrouter"
)

// GET daily entry
func (app Application) getDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var ctx = context.Background()
		params := req.Params()
		dateParam := params.ByName("date")

		var date time.Time
		if dateParam != "" {
			t, err := util.CreateTimeDate(dateParam)
			if err != nil {
				app.writeJSON(
					w,
					http.StatusBadRequest,
					Envelope{"error": err.Error()},
					nil,
				)
				return err
			}
			date = t
		} else {
			date = time.Now()
		}

		daily := new(db.Journal)
		err := app.DBClient.NewSelect().
			Model(daily).
			Column("id", "updated_at", "content").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", date.Format("2006-01-02"))).
			Scan(ctx)

		if err != nil {
			if err == sql.ErrNoRows {
				app.writeJSON(w, http.StatusNotFound, Envelope{"msg": "not found"}, nil)
				return nil
			}

			app.UnexpectedError(w, err)
			return err
		}

		app.writeJSON(w, http.StatusOK, Envelope{"daily": daily}, nil)
		return nil
	}
}

// POST daily entry
func (app Application) newDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var ctx = context.Background()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			app.writeJSON(w, http.StatusBadRequest, Envelope{"msg": "can't read body"}, nil)
			return err
		}

		var dailyEntry db.Journal

		err = json.Unmarshal(body, &dailyEntry)
		if err != nil {
			app.writeJSON(w, http.StatusBadRequest, Envelope{"msg": "error parsing body"}, nil)
			return err
		}

		exists, err := app.DBClient.NewSelect().
			Model(&db.Journal{}).
			Column("id").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", time.Now().Format("2006-01-02"))).
			Exists(ctx)

		if err != nil {
			return err
		}

		if exists {
			app.writeJSON(
				w,
				http.StatusBadRequest,
				Envelope{"error": "daily entry already exists for this date"},
				nil,
			)
			return err
		}

		// Create a new one
		dailyEntry.ID = db.CreateID()

		_, err = app.DBClient.NewInsert().
			Model(&dailyEntry).
			Exec(ctx)

		if err != nil {
			app.UnexpectedError(w, err)
			return err
		}
		app.writeJSON(
			w,
			http.StatusCreated,
			Envelope{"msg": "daily created"},
			nil,
		)
		return nil
	}
}

// PATCH daily
func (app Application) updateDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			app.writeJSON(w, http.StatusBadRequest, Envelope{"msg": "can't read body"}, nil)
			return err
		}

		var dailyEntry db.Journal
		err = json.Unmarshal(body, &dailyEntry)
		if err != nil {
			app.UnexpectedError(w, err)
			return err
		}

		var ctx = context.Background()

		_, err = app.DBClient.NewInsert().
			Model(&dailyEntry).
			On("CONFLICT (id) DO UPDATE").
			Set("updated_at = EXCLUDED.updated_at").
			Set("content = EXCLUDED.content").
			Exec(ctx)

		if err != nil {
			app.UnexpectedError(w, err)
			return err
		}
		app.writeJSON(
			w,
			http.StatusCreated,
			Envelope{"msg": "daily updated"},
			nil,
		)
		return nil
	}
}

// GET list daily entries
func (app *Application) listDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, _ bunrouter.Request) error {
		dailyEntries := new([]db.Journal)

		err := app.DBClient.NewSelect().
			Model(dailyEntries).
			Order("created_at DESC").
			Scan(context.Background())

		if err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
		}

		app.writeJSON(
			w,
			http.StatusCreated,
			Envelope{"daily_entries": dailyEntries},
			nil,
		)
		return nil
	}
}
