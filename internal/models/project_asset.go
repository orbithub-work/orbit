package models

import "github.com/uptrace/bun"

type ProjectAsset struct {
	bun.BaseModel `bun:"table:project_assets"`

	ProjectID           string  `bun:",pk" json:"project_id"`
	AssetID             string  `bun:",pk" json:"asset_id"`
	SourceID            string  `json:"source_id,omitempty"`
	Role                string  `json:"role"`
	BindMode            string  `json:"bind_mode"`
	Confidence          float64 `json:"confidence"`
	Alias               *string `json:"alias,omitempty"`
	TagsJSON            *string `json:"tags_json,omitempty"`
	Status              *string `json:"status,omitempty"`
	ProjectMetadataJSON *string `json:"project_metadata_json,omitempty"`
	CreatedAt           int64   `json:"created_at"`
	UpdatedAt           int64   `json:"updated_at"`
}
