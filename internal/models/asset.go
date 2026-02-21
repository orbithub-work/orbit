package models

import "github.com/uptrace/bun"

type Asset struct {
	bun.BaseModel `bun:"table:assets"`

	ID            string  `bun:",pk" json:"id"`
	Path          string  `json:"path"`
	Size          int64   `json:"size"`
	Mtime         int64   `json:"mtime"`
	Fingerprint   *string `json:"fingerprint,omitempty"`
	ParentAssetID *string `json:"parent_asset_id,omitempty"`
	Scope         string  `json:"scope"`
	ProjectID     *string `json:"project_id,omitempty"`

	// Two-stage scan support
	Status    string `bun:"status,notnull,default:'PENDING'" json:"status"` // PENDING, READY, INDEXED, MISSING, IGNORED, ERROR
	MediaMeta string `bun:"media_meta" json:"media_meta"`                   // JSON string for extra metadata (width, height, duration)
	Shape     string `bun:"shape,notnull,default:'unknown'" json:"shape"`

	// Rating model:
	// - suggested_rating is generated once during first successful metadata extraction.
	// - user_rating is controlled only by user actions and always has higher priority.
	SuggestedRating *int `bun:"suggested_rating" json:"suggested_rating,omitempty"`
	UserRating      *int `bun:"user_rating" json:"user_rating,omitempty"`

	// Operation log
	LastOpLog string `bun:"last_op_log" json:"last_op_log,omitempty"` // Record last user operation (e.g. "Removed by user")

	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

type AssetLineage struct {
	bun.BaseModel `bun:"table:asset_lineage"`

	ID           string `bun:",pk" json:"id"`
	AncestorID   string `json:"ancestor_id"`
	DescendantID string `json:"descendant_id"`
	RelationType string `json:"relation_type"`
	CreatedAt    int64  `json:"created_at"`
}
