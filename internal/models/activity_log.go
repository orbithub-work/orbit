package models

import "github.com/uptrace/bun"

type ActivityLog struct {
	bun.BaseModel `bun:"table:activity_logs"`

	ID        string `bun:",pk" json:"id"`
	Level     string `json:"level"`   // INFO, SUCCESS, WARN, ERROR
	Message   string `json:"message"` // User-friendly message
	AssetID   string `bun:"asset_id" json:"asset_id,omitempty"`
	ProjectID string `bun:"project_id" json:"project_id,omitempty"`
	CreatedAt int64  `json:"created_at"`
}
