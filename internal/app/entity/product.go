package entity

import (
	"time"

	"github.com/gofrs/uuid"
	"github.com/uptrace/bun"
)

type Product struct {
	bun.BaseModel `bun:"table:product"`

	ID           int64     `bun:"id,autoincrement"`
	GUID         uuid.UUID `bun:"guid,pk"`
	Name         string    `bun:"name,notnull,unique"`
	Description  *string   `bun:"description,nullzero"`
	Price        int64     `bun:"price,notnull"`
	CategoryGUID uuid.UUID `bun:"category_guid,notnull"`
	CreatedAt    time.Time `bun:"created_at,notnull,default:now()"`
	UpdatedAt    time.Time `bun:"updated_at,notnull,default:now()"`
}
