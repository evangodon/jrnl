package sqldb

import (
	"time"

	"github.com/uptrace/bun"
)

type EntryModel struct {
	bun.BaseModel `bun:"table:entry"`

	Id        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Content   string    `bun:"content,notnull"`
}
