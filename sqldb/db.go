package sqldb

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"jrnl/util"

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
		return "./testdb/jrnl.db"
	} else {
		home := os.Getenv("HOME")
		path := filepath.Join(home, ".data/jrnl", "jrnl.db")

		return path
	}
}

func Connect() DB {
	DB_PATH := GetDBPath()

	if _, err := os.Stat(DB_PATH); err != nil {
		CreateNewDB(DB_PATH)
		fmt.Println("Created new database")
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, DB_PATH)
	util.CheckError(err)
	db.DB = bun.NewDB(sqlite, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(enableLogs)))

	return db
}

const ID_LENGTH = 16

func CreateId() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", ID_LENGTH)
	util.CheckError(err)

	return id
}
