package sqldb

import (
	"time"

	"github.com/uptrace/bun"
)

type JournalEntry struct {
	bun.BaseModel `bun:"table:journal_entrys"`

	Id        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Content   string    `bun:"content,notnull"`
}
