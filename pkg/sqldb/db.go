package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"jrnl/pkg/util"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

const DB_FILE = "./jrnl.db"

var db *bun.DB

func Connect() *bun.DB {
	if f, err := os.Stat(DB_FILE); err != nil || f.Size() == 0 {
		fmt.Println("Creating new database")
		CreateDB()
	}

	sqlite, err := sql.Open(sqliteshim.ShimName, DB_FILE)
	util.CheckError(err)
	db = bun.NewDB(sqlite, sqlitedialect.New())

	db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))

	return db

}

func CreateDB() {
	ctx := context.Background()

	sqlite, err := sql.Open(sqliteshim.ShimName, DB_FILE)
	util.CheckError(err)
	db = bun.NewDB(sqlite, sqlitedialect.New())

	db.NewCreateTable().Model(&JournalEntry{}).Exec(ctx)

}

func CreateId() string {
	id, err := gonanoid.Generate("abcdefghijklmnopqrstuvwxyz0123456789", 16)
	util.CheckError(err)

	return id
}
