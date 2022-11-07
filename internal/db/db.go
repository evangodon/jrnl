package db

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/evangodon/jrnl/internal/cfg"
	"github.com/evangodon/jrnl/internal/logger"

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
	env := cfg.GetEnv()
	switch env {
	case "DEV":
		return filepath.Join(cfg.GetProjectRoot(), "/tmp/devjrnl.db")
	case "TEST":
		return filepath.Join(cfg.GetProjectRoot(), "/tmp/testjrnl.db")
	default:
		dbFile := "/jrnl/jrnl.db"
		path := xdg.DataHome + dbFile
		return path
	}
}

func Connect() DB {
	dbPath := GetDBPath()
	logger := logger.NewLogger(os.Stdout)

	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		logger.PrintInfo("No database found")
		if err = CreateNewDB(dbPath); err != nil {
			logger.PrintFatal(err.Error())
		}

		logger.PrintSuccess("Created new database at " + dbPath)
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		logger.PrintFatal(err.Error())
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

func CreateNewDB(dbPath string) error {
	ctx := context.Background()
	dir := filepath.Dir(dbPath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		return err
	}
	db.DB = bun.NewDB(sqlite, sqlitedialect.New())

	db.NewCreateTable().Model(&Journal{}).Exec(ctx)
	return nil
}
