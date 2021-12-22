package sqldb

import (
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

func GetDbPath() string {
	if isDev {
		return "./testdb/jrnl.db"
	} else {
		home := os.Getenv("HOME")
		path := filepath.Join(home, ".data/jrnl", "jrnl.db")

		return path
	}
}

func logQueries() bool {
	var logENV = os.Getenv("JRNL_LOG_QUERIES")
	var shouldLog = logENV == "true"
	return shouldLog && isDev
}

func Connect() *bun.DB {
	DB_PATH := GetDbPath()

	if _, err := os.Stat(DB_PATH); err != nil {
		CreateNewDB(DB_PATH)
		fmt.Println("Created new database")
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, DB_PATH)
	util.CheckError(err)
	db = bun.NewDB(sqlite, sqlitedialect.New())

	verbose := logQueries()
	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(verbose)))

	return db
}

func CreateId() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 16)
	util.CheckError(err)

	return id
}
