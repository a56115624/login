package model

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	CustomerStatusRegister = "entered"
	CustomerStatusActive   = "active"
	CustomerStatusBlocked  = "blocked"
)

type UserBindingForm struct {
	Id     int64  `json:"-"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

func (u *UserBindingForm) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}

	if u.Email == "" {
		return errors.New("email is required")
	}

	if u.Mobile == "" {
		return errors.New("mobile is required")
	}

	return nil
}

func (u *UserBindingForm) UpdateCustomer() *Customer {
	return &Customer{
		Id:        u.Id,
		Status:    CustomerStatusActive,
		Name:      u.Name,
		Email:     u.Email,
		Mobile:    u.Mobile,
		UpdatedAt: time.Now().UTC(),
	}
}

type ChangePasswordForm struct {
	Username           string `json:"username"`
	OldPassword        string `json:"old_password"`
	NewPassword        string `json:"new_password"`
	NewPasswordConfirm string `json:"new_password_confirm"`
}

type Customer struct {
	Id              int64 `bun:",pk,autoincrement"`
	Uid             int64
	Username        string
	Password        string
	Name            string
	Email           string `bun:",nullzero"`
	Mobile          string `bun:",nullzero"`
	Status          string `bun:",notnull,default:sign_up"`
	ConcurrentLimit int64  `bun:",notnull,default:3"`
	Referrer        int64
	CreatedAt       time.Time
	UpdatedAt       time.Time

	Txn          []*Txn          `bun:"rel:has-many,join:id=customer_id"`
	Order        []*Order        `bun:"rel:has-many,join:id=customer_id"`
	Subscription []*Subscription `bun:"rel:has-many,join:id=customer_id"`
}

// CheckPassword check user password
func (u *Customer) CheckPassword(password string) bool {
	return comparePasswords(u.Password, password)
}

// SetNewPassword change user password
func (u *Customer) SetNewPassword(password string) error {
	hash, err := hashAndSalt(password)
	if err != nil {
		return err
	}

	u.Password = hash
	return nil
}

// hashAndSalt 加密密碼
func hashAndSalt(pwdStr string) (pwdHash string, err error) {
	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	pwdHash = string(hash)
	return pwdHash, nil
}

// comparePasswords 驗證密碼
func comparePasswords(hashedPwd string, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		return false
	}
	return true
}

type CustomerBasicInfo struct {
	Id         int64  `json:"id"`
	Uid        int64  `json:"uid"`
	Username   string `json:"username"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Mobile     string `json:"mobile"`
	Referrer   int64  `json:"referrer"`
	Status     string `json:"status"`
	Concurrent int64  `json:"concurrent"`
	CreatedAt  string `json:"created_at"`
}

func (u *Customer) ToCustomerBasicInfo() *CustomerBasicInfo {
	return &CustomerBasicInfo{
		Id:        u.Id,
		Uid:       u.Uid,
		Username:  u.Username,
		Name:      u.Name,
		Email:     u.Email,
		Mobile:    u.Mobile,
		Referrer:  u.Referrer,
		Status:    u.Status,
		CreatedAt: u.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
	}
}

type CustomerInfo struct {
	CustomerBasicInfo
	Coins      int64  `json:"coins"`
	VipExpired string `json:"vip_expired"`
}

func (u *Customer) ToCustomerInfoWithBalanceAndVip(balance int64, vipExpired string) *CustomerInfo {
	return &CustomerInfo{
		CustomerBasicInfo: *u.ToCustomerBasicInfo(),
		Coins:             balance,
		VipExpired:        vipExpired,
	}
}

type CustomerDetailed struct {
	CustomerBasicInfo
	Txn          []*Txn          `json:"txn"`
	Order        []*Order        `json:"order"`
	Subscription []*Subscription `json:"subscription"`
}

func (u *Customer) ToCustomerDetailed() *CustomerDetailed {
	return &CustomerDetailed{
		CustomerBasicInfo: *u.ToCustomerBasicInfo(),
		Txn:               u.Txn,
		Order:             u.Order,
		Subscription:      u.Subscription,
	}
}
