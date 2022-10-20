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
func (srv Server) getDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var ctx = context.Background()
		params := req.Params()
		dateParam := params.ByName("date")

		var date time.Time
		if dateParam != "" {
			t, err := util.CreateTimeDate(dateParam)
			if err != nil {
				srv.JSON(
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
		err := srv.dbClient.NewSelect().
			Model(daily).
			Column("id", "updated_at", "content").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", date.Format("2006-01-02"))).
			Scan(ctx)

		if err != nil {
			if err == sql.ErrNoRows {
				srv.JSON(w, http.StatusNotFound, Envelope{"msg": "not found"}, nil)
				return nil
			}

			srv.UnexpectedError(w, err)
			return err
		}

		srv.JSON(w, http.StatusOK, Envelope{"daily": daily}, nil)
		return nil
	}
}

// POST daily entry
func (srv Server) newDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var ctx = context.Background()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			srv.JSON(w, http.StatusBadRequest, Envelope{"msg": "can't read body"}, nil)
			return err
		}

		var dailyEntry db.Journal

		err = json.Unmarshal(body, &dailyEntry)
		if err != nil {
			srv.JSON(w, http.StatusBadRequest, Envelope{"msg": "error parsing body"}, nil)
			return err
		}

		exists, err := srv.dbClient.NewSelect().
			Model(&db.Journal{}).
			Column("id").
			Where(fmt.Sprintf("DATE(created_at, 'localtime') = DATE('%s')", time.Now().Format("2006-01-02"))).
			Exists(ctx)

		if err != nil {
			return err
		}

		if exists {
			srv.JSON(
				w,
				http.StatusBadRequest,
				Envelope{"error": "daily entry already exists for this date"},
				nil,
			)
			return err
		}

		// Create a new one
		dailyEntry.ID = db.CreateID()

		_, err = srv.dbClient.NewInsert().
			Model(&dailyEntry).
			Exec(ctx)

		if err != nil {
			srv.UnexpectedError(w, err)
			return err
		}
		srv.JSON(
			w,
			http.StatusCreated,
			Envelope{"daily": dailyEntry},
			nil,
		)
		return nil
	}
}

// PATCH daily
func (srv Server) updateDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			srv.JSON(w, http.StatusBadRequest, Envelope{"msg": "can't read body"}, nil)
			return err
		}

		var dailyEntry db.Journal
		err = json.Unmarshal(body, &dailyEntry)
		if err != nil {
			srv.UnexpectedError(w, err)
			return err
		}

		var ctx = context.Background()

		_, err = srv.dbClient.NewInsert().
			Model(&dailyEntry).
			On("CONFLICT (id) DO UPDATE").
			Set("updated_at = EXCLUDED.updated_at").
			Set("content = EXCLUDED.content").
			Exec(ctx)

		if err != nil {
			srv.UnexpectedError(w, err)
			return err
		}
		srv.JSON(
			w,
			http.StatusCreated,
			Envelope{"msg": "daily updated"},
			nil,
		)
		return nil
	}
}

// Formats a new entry
func (srv *Server) getDailyTemplate() bunrouter.HandlerFunc {

	return func(w http.ResponseWriter, _ bunrouter.Request) error {

		template := util.FormatContent(db.Journal{
			CreatedAt: time.Now(),
		}, time.Now())

		srv.JSON(
			w,
			http.StatusCreated,
			Envelope{"template": template},
			nil,
		)
		return nil
	}
}

// GET list daily entries
func (srv *Server) listDailyHandler() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, _ bunrouter.Request) error {
		dailyEntries := new([]db.Journal)

		err := srv.dbClient.NewSelect().
			Model(dailyEntries).
			Order("created_at DESC").
			Scan(context.Background())

		if err != nil {
			if err != sql.ErrNoRows {
				log.Fatal(err)
			}
		}

		srv.JSON(
			w,
			http.StatusCreated,
			Envelope{"daily_entries": dailyEntries},
			nil,
		)
		return nil
	}
}
