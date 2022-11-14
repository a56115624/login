package model

import "time"

type Product struct {
	Id          int64 `bun:",pk,autoincrement"`
	Name        string
	CashAmount  int64
	TokenAmount int64
	Options     map[string]any `bun:"type:jsonb"`
	IsTest      bool
	Ordering    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductInfo struct {
	Id          int64          `json:"id"`
	Name        string         `json:"name"`
	Options     map[string]any `json:"options"`
	CashAmount  int64          `json:"cash_amount"`
	TokenAmount int64          `json:"token_amount"`
	IsTest      bool           `json:"is_test"`
	Ordering    int64          `json:"ordering"`
}

func (p *Product) ToProductInfo() *ProductInfo {
	return &ProductInfo{
		Id:          p.Id,
		Name:        p.Name,
		Options:     p.Options,
		CashAmount:  p.CashAmount,
		TokenAmount: p.TokenAmount,
		IsTest:      p.IsTest,
		Ordering:    p.Ordering,
	}
}

type ProductForm struct {
	Name        string         `json:"name"`
	Options     map[string]any `json:"options"`
	CashAmount  int64          `json:"cash_amount"`
	TokenAmount int64          `json:"token_amount"`
	IsTest      bool           `json:"is_test"`
}

func (p *ProductForm) ToProduct() *Product {
	return &Product{
		Name:        p.Name,
		Options:     p.Options,
		CashAmount:  p.CashAmount,
		TokenAmount: p.TokenAmount,
		IsTest:      p.IsTest,
		UpdatedAt:   time.Now().UTC(),
		CreatedAt:   time.Now().UTC(),
	}
}
