package models

// Tag 标签模型
type Tag struct {
	ID        string  `json:"id" bun:"id,pk"`
	Name      string  `json:"name" bun:"name,notnull"`
	Color     *string `json:"color" bun:"color"`
	Icon      *string `json:"icon" bun:"icon"`
	ParentID  *string `json:"parent_id" bun:"parent_id"`
	CreatedAt int64   `json:"created_at" bun:"created_at,notnull"`
	UpdatedAt int64   `json:"updated_at" bun:"updated_at,notnull"`
}

// AssetTag 资产标签关联模型
type AssetTag struct {
	AssetID   string `json:"asset_id" bun:"asset_id,pk"`
	TagID     string `json:"tag_id" bun:"tag_id,pk"`
	CreatedAt int64  `json:"created_at" bun:"created_at,notnull"`
}
