package model

import (
	"fmt"
	"math/rand"
	"time"
)

// Item 書本列表
type Item struct {
	Id          int64 `bun:",pk"`
	Source      int64
	Identifier  string
	Title       string
	Creator     string
	Status      string `bun:"default:連載"`
	Description string
	Covers      map[string]any `bun:"type:jsonb"`
	CreatedAt   time.Time      `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time

	BlockItem    *BlockItem    `bun:"rel:has-one,join:id=item_id"`
	CategoryItem *CategoryItem `bun:"rel:has-one,join:id=item_id"`
	Bookmark     *Bookmark     `bun:"rel:has-one,join:id=item_id"`
	Chapter      []*Chapter    `bun:"rel:has-many,join:id=item_id"`
}

// ItemInfo 書本簡介
type ItemInfo struct {
	Id          int64          `json:"id"`
	Title       string         `json:"title"`
	Creator     string         `json:"creator"`
	Description string         `json:"description"`
	Status      string         `json:"status"`
	Covers      map[string]any `json:"covers"`
	Hot         string         `json:"hot"`
	Views       string         `json:"views"`
}

// ToItemInfo 給前端的書本簡介
func (i *Item) ToItemInfo() *ItemInfo {
	rand.Seed(time.Now().UnixNano())

	hotList := []string{"8.4", "8.5", "8.6", "8.7", "8.8", "8.9", "9.0", "9.1", "9.2", "9.3", "9.4", "9.5", "9.6"}

	index := len(i.Title) * int(i.Id) % len(hotList)
	randomHot := hotList[index]

	min := 55
	views := min*index + int(i.Id)%min

	if views < 100 {
		views = views * 9
	}

	return &ItemInfo{
		Id:          i.Id,
		Title:       i.Title,
		Creator:     i.Creator,
		Description: i.Description,
		Covers:      i.Covers,
		Status:      i.Status,
		Hot:         randomHot,
		Views:       fmt.Sprintf("%d", views),
	}
}

type ItemDetail struct {
	ItemInfo
	BookmarkStatus bool           `json:"bookmark_status"`
	Chapter        []*ChapterInfo `json:"chapter"`
}

func (i *Item) ToItemDetail(chapterStatus, bookmarkStatus map[int64]bool) *ItemDetail {
	chapterInfoList := make([]*ChapterInfo, len(i.Chapter))
	for j := 0; j < len(i.Chapter); j++ {
		chapterInfoList[j] = i.Chapter[j].ToChapterInfo(chapterStatus, len(i.Chapter))
	}

	return &ItemDetail{
		ItemInfo:       *i.ToItemInfo(),
		BookmarkStatus: bookmarkStatus[i.Id],
		Chapter:        chapterInfoList,
	}
}

// Chapter 章節列表
type Chapter struct {
	Id        int64 `bun:",pk,autoincrement"`
	ItemId    int64
	Price     int64
	Covers    map[string]any `bun:"type:jsonb"`
	Contents  map[string]any `bun:"type:jsonb"`
	Position  int64
	CreatedAt time.Time `bun:",nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `bun:",nullzero"`

	Item *Item `bun:"rel:belongs-to,join:item_id=id"`
}

// ChapterBasic 章節基本的內容
type ChapterBasic struct {
	Id       int64          `json:"id"`
	Covers   map[string]any `json:"covers"`
	Position int64          `json:"position"`
	Price    int64          `json:"price"`
}

func (c *Chapter) ToChapterBasic() *ChapterBasic {
	return &ChapterBasic{
		Id:       c.Id,
		Covers:   c.Covers,
		Position: c.Position,
		Price:    c.Price,
	}
}

// ChapterInfo 章節簡介
type ChapterInfo struct {
	ChapterBasic
	Status bool `json:"status"`
}

// ToChapterInfo 給前端的章節簡介
func (c *Chapter) ToChapterInfo(chapterStatus map[int64]bool, totalChapter int) *ChapterInfo {
	var status bool
	if chapterStatus[c.Id] || (totalChapter > 2 && (c.Position == 1 || c.Position == 2)) {
		status = true
	}

	return &ChapterInfo{
		ChapterBasic: *c.ToChapterBasic(),
		Status:       status,
	}
}

type ChapterDetails struct {
	ChapterBasic
	Contents map[string]any `json:"contents"`
}

func (c *Chapter) ToChapterDetails() *ChapterDetails {

	return &ChapterDetails{
		ChapterBasic: *c.ToChapterBasic(),
		Contents:     c.Contents,
	}
}
