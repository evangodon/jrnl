package sqldb

import (
	"time"

	"github.com/uptrace/bun"
)

type Entry struct {
	bun.BaseModel `bun:"table:entries"`

	Id        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Type      string    `bun:",nullzero,notnull,default:'journal'"`
	Content   string    `bun:"content,notnull"`
}

var EntryType = struct {
	Journal string
	TIL     string
}{
	Journal: "journal",
	TIL:     "til",
}
