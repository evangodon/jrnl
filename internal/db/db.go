package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/evangodon/jrnl/internal/cfg"

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

func GetDBPath() string {
	if cfg.IsDev {
		return "./tmp/devjrnl.db"
	} else {
		dbFile := "/jrnl/jrnl.db"
		path := xdg.DataHome + dbFile

		return path
	}
}

func Connect() DB {
	dbPath := GetDBPath()

	println("dbPath", dbPath)

	if _, err := os.Stat(dbPath); err != nil {
		CreateNewDB(dbPath)
		fmt.Println("Created new database")
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	db.DB = bun.NewDB(sqlite, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(cfg.EnableLogs)))

	return db
}

const IDLength = 16

func CreateID() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", IDLength)
	if err != nil {
		log.Fatal(err)
	}

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
