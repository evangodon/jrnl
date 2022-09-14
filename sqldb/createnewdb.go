package sqldb

import (
	"context"
	"database/sql"
	"fmt"
	"jrnl/util"
	"os"
	"path/filepath"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func CreateNewDB(dbfile string) {
	fmt.Println(dbfile)
	ctx := context.Background()
	dir := filepath.Dir(dbfile)

	os.MkdirAll(dir, 0755)

	sqlite, err := sql.Open(sqliteshim.ShimName, dbfile)
	util.CheckError(err)
	db.DB = bun.NewDB(sqlite, sqlitedialect.New())

	db.NewCreateTable().Model(&Journal{}).Exec(ctx)
}
