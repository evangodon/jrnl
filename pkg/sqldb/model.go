package sqldb

import (
	"time"

	"github.com/uptrace/bun"
)

type Journal struct {
	bun.BaseModel `bun:"table:journal"`

	Id        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Content   string    `bun:"content,notnull"`
}

type TIL struct {
	bun.BaseModel `bun:"table:today_i_learnt"`

	Id        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Content   string    `bun:"content,notnull"`
}

var EntryType = struct {
	Journal string
	TIL     string
}{
	Journal: "journal",
	TIL:     "til",
}
