package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/evangodon/jrnl/internal/db"
	"github.com/evangodon/jrnl/internal/util"
	"github.com/uptrace/bunrouter"
)

// GET daily entry
func (srv Server) handleGetDaily() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var ctx = context.Background()
		params := req.Params()
		dateParam := params.ByName("date")

		var date = time.Now()
		if dateParam != "" {
			var t time.Time
			var err error
			if t, err = util.CreateTimeDate(dateParam); err != nil {
				return srv.BadRequest(err)
			}
			date = t
		}

		daily := new(db.Journal)
		err := srv.dbClient.NewSelect().
			Model(daily).
			Column("id", "updated_at", "content").
			Where(fmt.Sprintf("DATE(date) = DATE('%s')", date.Format("2006-01-02"))).
			Scan(ctx)
		if err != nil {
			return err
		}

		srv.JSON(w, http.StatusOK, Envelope{"daily": daily}, nil)
		return nil
	}
}

// POST daily entry
func (srv Server) handleNewDaily() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		var ctx = context.Background()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}

		var dailyEntry db.Journal
		if err = json.Unmarshal(body, &dailyEntry); err != nil {
			return err
		}

		exists, err := srv.dbClient.NewSelect().
			Model(&db.Journal{}).
			Column("id").
			Where(fmt.Sprintf("DATE(date, 'localtime') = DATE('%s')", dailyEntry.CreatedAt.Format("2006-01-02"))).
			Exists(ctx)
		if err != nil {
			return err
		}

		if exists {
			return srv.BadRequest(errors.New("daily entry already exists for this date"))
		}

		// Create a new one
		dailyEntry.ID = db.CreateID()

		_, err = srv.dbClient.NewInsert().
			Model(&dailyEntry).
			Exec(ctx)

		if err != nil {
			return srv.UnexpectedError(err)
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
func (srv Server) handleUpdateDaily() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return srv.UnexpectedError(err)
		}

		var dailyEntry db.Journal
		err = json.Unmarshal(body, &dailyEntry)
		if err != nil {
			return err
		}

		var ctx = context.Background()

		_, err = srv.dbClient.NewUpdate().
			OmitZero().
			Model(&dailyEntry).
			WherePK().
			Exec(ctx)

		if err != nil {
			return srv.UnexpectedError(err)
		}
		srv.JSON(
			w,
			http.StatusNoContent,
			Envelope{"daily": dailyEntry},
			nil,
		)
		return nil
	}
}

// Formats a new entry
func (srv *Server) handleGetDailyTemplate() bunrouter.HandlerFunc {

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
func (srv *Server) handleListDaily() bunrouter.HandlerFunc {
	return func(w http.ResponseWriter, req bunrouter.Request) error {
		pageParam := req.URL.Query().Get("page")
		perPageParam := req.URL.Query().Get("perPage")

		page, err := strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}

		perPage, err := strconv.Atoi(perPageParam)
		if err != nil {
			perPage = 20
		}

		dailyEntries := new([]db.Journal)
		count, err := srv.dbClient.NewSelect().
			Model(dailyEntries).
			Order("created_at DESC").
			Offset((page - 1) * perPage).
			Limit(perPage).
			ScanAndCount(context.Background())

		if err != nil {
			return err
		}

		srv.JSON(
			w,
			http.StatusCreated,
			Envelope{"daily_entries": dailyEntries, "total": count},
			nil,
		)
		return nil
	}
}
