package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/evangodon/jrnl/internal/util"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

type DB struct {
	*bun.DB
}

var db DB

var isDev = os.Getenv("DEV") == "true"
var enableLogs = os.Getenv("JRNL_ENABLE_LOGS") == "true"

func GetDBPath() string {
	if isDev {
		fmt.Println("Using dev database")
		return "./tmp/devjrnl.db"
	} else {
		home := os.Getenv("HOME")
		path := filepath.Join(home, ".data/jrnl", "jrnl.db")

		return path
	}
}

func Connect() DB {
	dbPath := GetDBPath()

	if _, err := os.Stat(dbPath); err != nil {
		CreateNewDB(dbPath)
		fmt.Println("Created new database")
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, dbPath)
	util.CheckError(err)
	db.DB = bun.NewDB(sqlite, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(enableLogs)))

	return db
}

const IDLength = 16

func CreateID() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", IDLength)
	util.CheckError(err)

	return id
}

type Item struct {
	ID      string
	Content string
}

// SelectEntryByRowNumber selects an entry by row number
func (db *DB) SelectEntryByRowNumber(
	ctx context.Context,
	model interface{},
	rowNumber int,
) (item Item, err error) {

	var selectedItem = Item{}

	err = db.NewSelect().
		Model(model).
		Column("id", "content").
		Limit(1).
		Offset(rowNumber-1).
		Scan(ctx, &selectedItem)

	if err != nil {
		return selectedItem, err
	}

	return selectedItem, nil
}

// SelectEntryByID selects an entry by its id
func (db *DB) SelectEntryByID(
	ctx context.Context,
	model interface{},
	id string,
) (item Item, err error) {

	var selectedItem = Item{}

	err = db.NewSelect().
		Model(model).
		Column("id", "content").
		Where("id = ?", id).
		Scan(ctx, &selectedItem)

	if err != nil {
		return selectedItem, err
	}

	return selectedItem, nil
}

func (DB) UpdateEntryContent(ctx context.Context, model interface{}, item Item) (err error) {
	_, err = db.NewUpdate().
		Model(model).
		Set("content = ?", item.Content).
		Where("id = ?", item.ID).
		Exec(ctx)

	if err != nil {
		return err
	}

	return nil
}
