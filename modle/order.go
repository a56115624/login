package model

import (
	"fmt"
	"time"
)

const (
	OrderStatusUnpaid    = "unpaid"
	OrderStatusPaid      = "paid"
	OrderPlatformWeb     = "web"
	OrderPlatformIos     = "ios"
	OrderPlatformAndroid = "android"
)

type Order struct {
	Id              int64 `bun:",pk,autoincrement"`
	CustomerId      int64
	ProductId       int64
	PaymentMethodId int64
	Status          string
	OrderNum        string
	ProviderNum     string
	CashAmount      int64
	TokenAmount     int64
	OrderAmount     int64
	PayUrl          string
	Platform        string
	IsFirstOrder    bool
	PaidAt          time.Time `bun:",nullzero"`
	CreatedAt       time.Time
	UpdatedAt       time.Time

	// Relations
	PaymentMethod *PaymentMethod `bun:"rel:belongs-to,join:payment_method_id=id"`
	Customer      *Customer      `bun:"rel:belongs-to,join:customer_id=id"`
}

func NewOrder(customerId int64, orderNum string, platform string) *Order {
	return &Order{
		CustomerId: customerId,
		OrderNum:   orderNum,
		Status:     OrderStatusUnpaid,
		Platform:   platform,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (order *Order) AddMethodId(paymentMethod *PaymentMethod) {
	order.PaymentMethodId = paymentMethod.Id
}

func (order *Order) AddAmount(product *Product) {
	order.ProductId = product.Id
	order.CashAmount = product.CashAmount
	order.TokenAmount = product.TokenAmount
}

func (order *Order) AddOrderAmount(amount int64) {
	order.OrderAmount = amount
}

func (order *Order) AddFirstOrderStatus(firstOrder bool) {
	order.IsFirstOrder = firstOrder
}

func (order *Order) UpdatePayStatus() {
	order.Status = OrderStatusPaid
}

type OrderForm struct {
	CustomerId             int64  `json:"-"`
	ProductPaymentMethodId int64  `json:"product_payment_method_id"`
	Amount                 int64  `json:"amount"`
	Platform               string `json:"platform"`
	Host                   string `json:"-"`
	ClientIp               string `json:"-"`
}

func (order *OrderForm) Validate() error {
	if order.ProductPaymentMethodId == 0 {
		return fmt.Errorf("product_payment_method_id is required")
	}

	if order.Amount == 0 {
		return fmt.Errorf("amount is required")
	}

	if order.Platform == "" || (order.Platform != OrderPlatformWeb && order.Platform != OrderPlatformAndroid && order.Platform != OrderPlatformIos) {
		return fmt.Errorf("platform is required, and must be web, android or ios")
	}

	return nil
}

type OrderAdminInfo struct {
	Id          int64 `json:"id"`
	CustomerId  int64 `json:"customer_id"`
	OrderAmount int64 `json:"order_amount"`
	OrderInfo
}

func (order *Order) ToOrderAdminInfo() *OrderAdminInfo {
	return &OrderAdminInfo{
		Id:          order.Id,
		CustomerId:  order.CustomerId,
		OrderAmount: order.OrderAmount,
		OrderInfo:   *order.ToOrderInfo(),
	}
}

type OrderInfo struct {
	CashAmount  int64  `json:"cash_amount"`
	TokenAmount int64  `json:"token_amount"`
	OrderNum    string `json:"order_num"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
}

func (order *Order) ToOrderInfo() *OrderInfo {
	return &OrderInfo{
		CashAmount:  order.CashAmount,
		TokenAmount: order.TokenAmount,
		OrderNum:    order.OrderNum,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
	}
}

type OrderPaymentInfo struct {
	OrderNum string `json:"order_num"`
	//PayUrl   string `json:"pay_url"`
}

func NewOrderPaymentInfo(order *Order) *OrderPaymentInfo {
	return &OrderPaymentInfo{
		OrderNum: order.OrderNum,
	}
}
