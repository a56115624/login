package model

import (
	"errors"
	"time"
)

type AuthUserLoginForm struct {
	Referrer int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *AuthUserLoginForm) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (u *AuthUserLoginForm) NewManager() (*Manager, error) {
	hash, err := hashAndSalt(u.Password)
	if err != nil {
		return nil, err
	}

	return &Manager{
		Username:  u.Username,
		Password:  hash,
		Status:    ManagerStatusActive,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}

func (u *AuthUserLoginForm) NewCustomer(uid int64) (*Customer, error) {
	hash, err := hashAndSalt(u.Password)
	if err != nil {
		return nil, err
	}

	return &Customer{
		Username:  u.Username,
		Referrer:  u.Referrer,
		Uid:       uid,
		Password:  hash,
		Status:    CustomerStatusRegister,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}, nil
}
