package model

import (
	"time"
)

const (
	ManagerStatusActive   = "active"
	ManagerStatusInactive = "inactive"
)

type Manager struct {
	Id           int64 `bun:",pk,autoincrement"`
	Username     string
	Password     string
	Name         string
	Scope        any `bun:"type:jsonb,nullzero"`
	Status       string
	LastLoggedAt time.Time `bun:",nullzero"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m *Manager) UpdateLastLoggedAt() {
	m.LastLoggedAt = time.Now().UTC()
}

// CheckPassword check user password
func (m *Manager) CheckPassword(password string) bool {
	return comparePasswords(m.Password, password)
}

// SetNewPassword change user password
func (m *Manager) SetNewPassword(password string) error {
	hash, err := hashAndSalt(password)
	if err != nil {
		return err
	}

	m.Password = hash
	return nil
}

type ManagerInfo struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Username     string `json:"username"`
	Status       string `json:"status"`
	Scope        any    `json:"scope"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	LastLoggedAt string `json:"last_logged_at"`
}

func (m *Manager) ToManagerInfo() *ManagerInfo {
	var LoggedAt string
	if !m.LastLoggedAt.IsZero() {
		LoggedAt = m.LastLoggedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05")
	}

	return &ManagerInfo{
		Id:           m.Id,
		Name:         m.Name,
		Username:     m.Username,
		Status:       m.Status,
		Scope:        m.Scope,
		LastLoggedAt: LoggedAt,
		CreatedAt:    m.CreatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
		UpdatedAt:    m.UpdatedAt.Add(8 * time.Hour).Format("2006-01-02 15:04:05"),
	}
}

type ManagerUpdateForm struct {
	Id     int64    `json:"-"`
	Name   string   `json:"name"`
	Scope  []string `json:"scope"`
	Status string   `json:"status"`
}

func (m *ManagerUpdateForm) Validate() error {
	return nil
}

func (m *ManagerUpdateForm) ToUpdateManager() *Manager {
	return &Manager{
		Id:        m.Id,
		Name:      m.Name,
		Scope:     m.Scope,
		Status:    m.Status,
		UpdatedAt: time.Now().UTC(),
	}
}
