package models

import "github.com/uptrace/bun"

type ProjectSource struct {
	bun.BaseModel `bun:"table:project_sources"`

	ID          string `bun:",pk" json:"id"`
	ProjectID   string `bun:"project_id" json:"project_id"`
	RootPath    string `bun:"root_path" json:"root_path"`
	SourceType  string `bun:"source_type" json:"source_type"` // primary | extra
	WatchEnabled bool  `bun:"watch_enabled" json:"watch_enabled"`
	CreatedAt   int64  `bun:"created_at" json:"created_at"`
	UpdatedAt   int64  `bun:"updated_at" json:"updated_at"`
}
