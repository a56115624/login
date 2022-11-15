package model

import (
	"github.com/uptrace/bun"
)

type Customerdata struct {
	bun.BaseModel `bun:"table:customer"`

	Id       int64 `bun:",pk,autoincrement"`
	Username string
	Password string
}
