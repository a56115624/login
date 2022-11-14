package model

import (
	"errors"
	"fmt"
	"time"
)

type Txn struct {
	Id          int64 `bun:",pk,autoincrement"`
	CustomerId  int64
	TokenAmount int64
	Balance     int64
	Description string
	CreatedAt   time.Time `bun:"created_at,notnull,default:current_timestamp"`
	UpdatedAt   time.Time

	Purchased *Purchased `bun:"rel:has-one,join:id=txn_id"`
}

type TxnForm struct {
	CustomerId int64 `json:"-"`
	ChapterId  int64 `json:"chapter_id"`
}

func (t *TxnForm) Validate() error {
	if t.ChapterId == 0 {
		return errors.New("chapter_id is required")
	}

	return nil
}

func NewChapterTxn(customerId int64, chapter *Chapter, balance int64) *Txn {
	return &Txn{
		CustomerId:  customerId,
		TokenAmount: -chapter.Price,
		Description: fmt.Sprintf("%s - 第%d章", chapter.Item.Title, chapter.Position),
		Balance:     balance,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func NewChargeTxn(order *Order, balance int64) *Txn {
	return &Txn{
		CustomerId:  order.CustomerId,
		TokenAmount: order.TokenAmount,
		Description: "充值",
		Balance:     balance,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func NewRegisterTxn(customerId int64) *Txn {
	return &Txn{
		CustomerId:  customerId,
		TokenAmount: 50,
		Description: "商户注册初始金额",
		Balance:     50,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func NewUpdateProfileTxn(customerId, balance int64) *Txn {
	return &Txn{
		CustomerId:  customerId,
		TokenAmount: 150,
		Description: "商户完善资料奖金",
		Balance:     balance + 150,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

func NewAdminTxn(customerId, tokenAmount, balance int64) *Txn {
	return &Txn{
		CustomerId:  customerId,
		TokenAmount: tokenAmount,
		Description: "管理员操作",
		Balance:     balance,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}

type TxnInfo struct {
	Amount  int64 `json:"amount"`
	Balance int64 `json:"balance"`
	//ItemTitle       string `json:"item_title"`
	Description string `json:"description"`
	//ChapterPosition int64  `json:"chapter_position"`
	CreatedAt string `json:"created_at"`
}

func (t *Txn) ToTxnInfo() *TxnInfo {

	return &TxnInfo{
		Amount:      t.TokenAmount,
		Balance:     t.Balance,
		Description: t.Description,
		//ItemTitle:       t.Purchased.Chapter.Item.Title,
		//ChapterPosition: t.Purchased.Chapter.Position,
		CreatedAt: t.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
	}
}
