package model

import (
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

type Category struct {
	Id          int64 `bun:",pk,autoincrement"`
	Name        string
	Description string
	Status      string
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt   time.Time
}

type CategoryInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *Category) ToCategoryInfo() *CategoryInfo {
	return &CategoryInfo{
		Id:   c.Id,
		Name: c.Name,
	}
}

type CategoryItem struct {
	bun.BaseModel `bun:"table:categories_items"`
	Id            int64 `bun:",pk,autoincrement"`
	ItemId        int64
	CategoryId    int64
	CreatedAt     time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt     time.Time
	Category      *Category `bun:"rel:has-one,join:item_id=id"`
}

type CategoryRequest struct {
	Type []string `json:"type" query:"type"`
}

func (r *CategoryRequest) Validate() error {
	typeSet := map[string]bool{
		"剧情":   true,
		"惊悚":   true,
		"学园":   true,
		"流氓":   true,
		"玄幻":   true,
		"外语":   true,
		"言情":   true,
		"乱伦":   true,
		"偷情":   true,
		"真人":   true,
		"3D":   true,
		"同人志":  true,
		"猎奇":   true,
		"男男":   true,
		"多人":   true,
		"强奸":   true,
		"兽交":   true,
		"游戏CG": true,
		"角色扮演": true,
	}

	for i := 0; i < len(r.Type); i++ {
		if !typeSet[r.Type[i]] {
			return fmt.Errorf("invalid type: %s", r.Type[i])
		}
	}

	return nil
}
