package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"jrnl/pkg/util"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

var db *bun.DB

var devENV = os.Getenv("DEV")
var isDev = devENV == "true"

func getDBPATH() string {
	if isDev {
		return "./.dev/jrnl.db"
	} else {
		home := os.Getenv("HOME")
		path := filepath.Join(home, ".data/jrnl", "jrnl.db")

		return path
	}
}

func Connect() *bun.DB {
	DB_PATH := getDBPATH()

	if f, err := os.Stat(DB_PATH); f.Size() == 0 {
		util.CheckError(err)
		CreateDB(DB_PATH)
		fmt.Println("Created new database")
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, DB_PATH)
	util.CheckError(err)
	db = bun.NewDB(sqlite, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(isDev)))

	return db

}

func CreateDB(dbfile string) {
	ctx := context.Background()
	dir := filepath.Dir(dbfile)

	os.MkdirAll(dir, 0755)

	sqlite, err := sql.Open(sqliteshim.ShimName, dbfile)
	util.CheckError(err)
	db = bun.NewDB(sqlite, sqlitedialect.New())

	db.NewCreateTable().Model(&JournalEntry{}).Exec(ctx)

}

func CreateId() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 16)
	util.CheckError(err)

	return id
}
