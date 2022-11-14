package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

// Block 分類列表
type Block struct {
	Id          int64 `bun:",pk,autoincrement"`
	Description string
	Status      string
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time
}

// BlockItem 分類與書本的對應列表
type BlockItem struct {
	bun.BaseModel `bun:"table:blocks_items"`
	Id            int64 `bun:",pk,autoincrement"`
	BlockId       int64
	ItemId        int64
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time

	Block *Block `bun:"rel:has-one,join:item_id=id"`
}

// BlockRequest 前端的區塊請求
type BlockRequest struct {
	Type []string `json:"type" query:"type"`
	//Target   string   `json:"target" query:"target"`
}

// Validate 區塊請求的驗證
func (r *BlockRequest) Validate() error {

	typeSet := map[string]bool{
		"banner":    true,
		"3D主打":      true,
		"私人收藏":      true,
		"吸睛首选":      true,
		"热门Cosplay": true,
		"新书强推":      true,
		"本周更新":      true,
	}

	for i := 0; i < len(r.Type); i++ {
		if !typeSet[r.Type[i]] {
			return fmt.Errorf("invalid type: %s", r.Type[i])
		}
	}

	return nil
}

type BlockNameOrdering struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

func NewBlockNameOrdering(id int64, name string, ordering int64) *BlockNameOrdering {
	return &BlockNameOrdering{
		Id:    id,
		Name:  name,
		Count: ordering,
	}
}
