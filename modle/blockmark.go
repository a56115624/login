package model

import (
	"fmt"
	"time"
)

const (
	BookmarkActionAdd = "add"
	BookmarkStatusAdd = "added"

	BookmarkActionRemove = "remove"
	BookmarkStatusRemove = "removed"
)

type Bookmark struct {
	Id         int64 `bun:",pk,autoincrement"`
	CustomerId int64
	ItemId     int64
	Status     string    `bun:",default:added"`
	CreatedAt  time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt  time.Time
}

type BookmarkForm struct {
	CustomerId int64  `json:"-"`
	ItemId     int64  `json:"item_id"`
	Action     string `json:"action"`
}

func (f *BookmarkForm) Validate() error {
	if f.ItemId == 0 {
		return fmt.Errorf("item_id is required")
	}

	if f.Action != BookmarkActionAdd && f.Action != BookmarkActionRemove {
		return fmt.Errorf("action is required")
	}

	return nil
}

func (f *BookmarkForm) NewBookmark() *Bookmark {
	status := BookmarkStatusAdd
	if f.Action == BookmarkActionRemove {
		status = BookmarkStatusRemove
	}

	return &Bookmark{
		CustomerId: f.CustomerId,
		ItemId:     f.ItemId,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}
}
