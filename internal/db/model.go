package db

import (
	"context"
	"time"

	"github.com/uptrace/bun"
)

type Journal struct {
	bun.BaseModel `bun:"table:journal"`

	ID        string    `bun:"id,pk"                                       json:"id"`
	Date      time.Time `bun:",notnull,unique"                             json:"date"`
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp" json:"updatedAt"`
	Content   string    `bun:"content,notnull"                             json:"content"`
}

var _ bun.BeforeAppendModelHook = (*Journal)(nil)

func (j *Journal) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		if j.CreatedAt.IsZero() {
			j.CreatedAt = time.Now()
		}
		j.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		j.UpdatedAt = time.Now()
	}
	return nil
}
