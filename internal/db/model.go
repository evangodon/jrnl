package db

import (
	"time"

	"github.com/uptrace/bun"
)

type Journal struct {
	bun.BaseModel `bun:"table:journal"`

	ID        string    `bun:"id,pk"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	Content   string    `bun:"content,notnull"`
}

var EntryType = struct {
	Journal string
}{
	Journal: "journal",
}
