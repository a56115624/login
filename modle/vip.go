package model

import (
	"time"
)

const (
	SubscriptionStatusActive = "active"
	SubscriptionStatusCancel = "cancel"
)

type VipPlan struct {
	Id        int64 `bun:",pk,autoincrement"`
	Name      string
	Price     int64
	Days      int64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VipPlanInfo struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Days  int64  `json:"days"`
}

func (v *VipPlan) ToVipPlanInfo() *VipPlanInfo {
	return &VipPlanInfo{
		Id:    v.Id,
		Name:  v.Name,
		Price: v.Price,
		Days:  v.Days,
	}
}

type Subscription struct {
	Id         int64 `bun:",pk,autoincrement"`
	CustomerId int64
	TxnId      int64
	Period     []string
	//Period    []time.Time
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time

	Txn *Txn `bun:"rel:belongs-to,join:txn_id=id"`
}

func (s *Subscription) IsExpired() bool {
	return s.Period[1] < time.Now().UTC().Format("2006-01-02 15:04:05")
}

func (s *Subscription) GetExpiredAtUTC() string {
	return s.Period[1]
}

func (s *Subscription) GetExpiredAtTW() string {
	parse, err := time.Parse("2006-01-02 15:04:05+00", s.Period[1])
	if err != nil {
		return s.Period[1]
	}

	//轉換到台灣時間+8
	return parse.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
}

type SubscriptionForm struct {
	CustomerId int64 `json:"-"`
	VipPlanId  int64 `json:"vip_plan_id"`
}

func (s *SubscriptionForm) Validate() error {
	return nil
}

func NewSubscription(txnId, customerId int64, startTime string, days int64) *Subscription {
	timeArray := make([]string, 0)
	//timeArray := []string{}
	//var timeArray []time.Time

	var start time.Time
	// 2022-11-27 07:05:21+00
	start, err := time.Parse("2006-01-02 15:04:05+00", startTime)
	if err != nil {
		start = time.Now().UTC()
	}

	if startTime == "" || start.Before(time.Now().UTC()) {
		start = time.Now().UTC()
	}

	timeArray = append(timeArray, start.Format("2006-01-02 15:04:05"))
	timeArray = append(timeArray, start.AddDate(0, 0, int(days)).Format("2006-01-02 15:04:05"))

	//timeArray[0] = startTime
	//timeArray[1] = startTime.AddDate(0, 0, int(days))

	//timeArray[0] = startTime.Format("2006-01-02 15:04:05")
	//timeArray[1] = startTime.AddDate(0, 0, int(days)).Format("2006-01-02 15:04:05")

	return &Subscription{
		CustomerId: customerId,
		TxnId:      txnId,
		Status:     SubscriptionStatusActive,
		Period:     timeArray,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}

func (s *SubscriptionForm) ToNewTxn(price, balance int64) *Txn {
	return &Txn{
		CustomerId:  s.CustomerId,
		TokenAmount: -price,
		Balance:     balance,
		Description: "buy vip",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}
}
