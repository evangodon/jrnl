package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func CreateNewDB(dbPath string) {
	ctx := context.Background()
	dir := filepath.Dir(dbPath)

	os.MkdirAll(dir, 0755)

	sqlite, err := sql.Open(sqliteshim.ShimName, dbPath)
	if err != nil {
		log.Fatal(err)
	}
	db.DB = bun.NewDB(sqlite, sqlitedialect.New())

	db.NewCreateTable().Model(&Journal{}).Exec(ctx)
}
