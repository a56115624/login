package model

import (
	"fmt"
	"strings"
	"time"
)

const (
	ProductPaymentMethodStatusActive   = "active"
	ProductPaymentMethodStatusInactive = "inactive"

	MethodStatusActive   = "active"
	MethodStatusInactive = "inactive"

	MethodTypeWechat = "WECHAT"
	MethodTypeAlipay = "ALIPAY"
	MethodTypeUSDT   = "USDT"

	MethodCurrencyCNY  = "CNY"
	MethodCurrencyHKD  = "HKD"
	MethodCurrencyUSDT = "USDT"
	MethodCurrencyTWD  = "TWD"
)

type PaymentMethod struct {
	Id        int64 `bun:",pk,autoincrement"`
	Name      string
	Options   map[string]any `bun:"type:jsonb"`
	Type      string
	Currency  string
	Ordering  int64
	Status    string
	Memo      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PaymentMethodForm struct {
	Name     string         `json:"name"`
	Options  map[string]any `json:"options"`
	Type     string         `json:"type"`
	Currency string         `json:"currency"`
	Ordering int64          `json:"ordering"`
	Status   string         `json:"status"`
	Memo     string         `json:"memo"`
}

func (p *PaymentMethodForm) UpdateValidate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	p.Type = strings.ToUpper(p.Type)
	if p.Type != "" && p.Type != MethodTypeAlipay && p.Type != MethodTypeWechat && p.Type != MethodTypeUSDT {
		return fmt.Errorf("invalid type, must be ALIPAY or WECHAT or USDT")
	}

	p.Status = strings.ToLower(p.Status)
	if p.Status != "" && p.Status != MethodStatusActive && p.Status != MethodStatusInactive {
		return fmt.Errorf("invalid status must be active or inactive")
	}

	p.Currency = strings.ToUpper(p.Currency)
	if p.Currency != "" && p.Currency != MethodCurrencyCNY && p.Currency != MethodCurrencyHKD && p.Currency != MethodCurrencyTWD && p.Currency != MethodCurrencyUSDT {
		return fmt.Errorf("invalid currency must be CNY or HKD or TWD or USDT")
	}

	if p.Ordering < 0 {
		return fmt.Errorf("invalid ordering")
	}

	return nil
}

func (p *PaymentMethodForm) CreateValidate() error {
	if p.Name == "" {
		return fmt.Errorf("name is required")
	}

	if p.Options == nil {
		return fmt.Errorf("options is required")
	}

	p.Type = strings.ToUpper(p.Type)
	if p.Type == "" || (p.Type != MethodTypeAlipay && p.Type != MethodTypeWechat && p.Type != MethodTypeUSDT) {
		return fmt.Errorf("invalid type, must be ALIPAY or WECHAT or USDT")
	}

	p.Currency = strings.ToUpper(p.Currency)
	if p.Currency == "" || (p.Currency != MethodCurrencyCNY && p.Currency != MethodCurrencyHKD && p.Currency != MethodCurrencyTWD && p.Currency != MethodCurrencyUSDT) {
		return fmt.Errorf("invalid currency must be CNY or HKD or TWD or USDT")
	}

	p.Status = strings.ToLower(p.Status)
	if p.Status == "" || (p.Status != MethodStatusActive && p.Status != MethodStatusInactive) {
		return fmt.Errorf("invalid status must be active or inactive")
	}

	if p.Ordering < 0 {
		return fmt.Errorf("invalid ordering")
	}

	return nil
}

func (p *PaymentMethodForm) ToNewPaymentMethod() *PaymentMethod {
	return &PaymentMethod{
		Name:      p.Name,
		Options:   p.Options,
		Type:      p.Type,
		Currency:  p.Currency,
		Status:    p.Status,
		UpdatedAt: time.Now().UTC(),
		CreatedAt: time.Now().UTC(),
	}
}

func (p *PaymentMethodForm) ToUpdatePaymentMethod() *PaymentMethod {
	return &PaymentMethod{
		Name:      p.Name,
		Options:   p.Options,
		Type:      p.Type,
		Currency:  p.Currency,
		Status:    p.Status,
		Memo:      p.Memo,
		UpdatedAt: time.Now().UTC(),
	}
}

type PaymentMethodInfoAdmin struct {
	Id       int64          `json:"id"`
	Name     string         `json:"name"`
	Options  map[string]any `json:"options"`
	Type     string         `json:"type"`
	Memo     string         `json:"memo"`
	Currency string         `json:"currency"`
	Status   string         `json:"status"`
}

func (p *PaymentMethod) ToPaymentMethodInfoAdmin() *PaymentMethodInfoAdmin {
	return &PaymentMethodInfoAdmin{
		Id:       p.Id,
		Name:     p.Name,
		Options:  p.Options,
		Type:     p.Type,
		Currency: p.Currency,
		Memo:     p.Memo,
		Status:   p.Status,
	}
}

type PaymentMethodInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
}

func (p *PaymentMethod) ToPaymentMethodInfo() *PaymentMethodInfo {
	return &PaymentMethodInfo{
		Id:       p.Id,
		Name:     p.Name,
		Currency: p.Currency,
	}
}

type ProductPaymentMethod struct {
	Id              int64 `bun:",pk,autoincrement"`
	ProductId       int64
	PaymentMethodId int64
	//Status          string
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relations
	Product       *Product       `bun:"rel:belongs-to,join:product_id=id"`
	PaymentMethod *PaymentMethod `bun:"rel:belongs-to,join:payment_method_id=id"`
}

func NewProductPaymentMethod(productId int64, paymentMethodId int64) *ProductPaymentMethod {
	return &ProductPaymentMethod{
		ProductId:       productId,
		PaymentMethodId: paymentMethodId,
		//Status:          ProductPaymentMethodStatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
}

func UpdateProductMethod(productId int64, paymentMethodId int64, status string) *ProductPaymentMethod {
	return &ProductPaymentMethod{
		ProductId:       productId,
		PaymentMethodId: paymentMethodId,
		//Status:          status,
		UpdatedAt: time.Now().UTC(),
	}
}

type ProductPaymentMethodForm struct {
	MethodId int64 `json:"method_id"`
}

type ProductPaymentMethodInfo struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Ordering int64  `json:"ordering"`
}

func (p *ProductPaymentMethod) ToProductPaymentMethodInfo() *ProductPaymentMethodInfo {
	return &ProductPaymentMethodInfo{
		Id:       p.Id,
		Name:     p.PaymentMethod.Name,
		Type:     p.PaymentMethod.Type,
		Ordering: p.PaymentMethod.Ordering,
	}
}
