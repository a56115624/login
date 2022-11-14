package model

import (
	"github.com/uptrace/bun"
	"time"
)

type Purchased struct {
	bun.BaseModel `bun:"purchased"`
	Id            int64 `bun:",pk,autoincrement"`
	CustomerId    int64
	ChapterId     int64
	TxnId         int64
	Status        string
	CreatedAt     time.Time `bun:",notnull,default:current_timestamp"`
	UpdatedAt     time.Time `bun:",notnull,default:current_timestamp"`

	Chapter *Chapter `bun:"rel:belongs-to,join:chapter_id=id"`
	Txn     *Txn     `bun:"rel:belongs-to,join:txn_id=id"`
}

func NewPurchased(userId, chapterId, txnId int64) *Purchased {
	return &Purchased{
		CustomerId: userId,
		ChapterId:  chapterId,
		TxnId:      txnId,
		Status:     "paid",
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
