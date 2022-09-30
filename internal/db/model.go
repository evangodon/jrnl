package db

import (
	"context"
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

var _ bun.BeforeAppendModelHook = (*Journal)(nil)

func (j *Journal) BeforeAppendModel(ctx context.Context, query bun.Query) error {
	switch query.(type) {
	case *bun.InsertQuery:
		j.CreatedAt = time.Now()
		j.UpdatedAt = time.Now()
	case *bun.UpdateQuery:
		j.UpdatedAt = time.Now()
	}
	return nil
}
