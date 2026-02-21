package models

import "github.com/uptrace/bun"

type LibrarySource struct {
	bun.BaseModel `bun:"table:library_sources"`

	ID           string `bun:",pk" json:"id"`
	RootPath     string `bun:"root_path" json:"root_path"`
	WatchEnabled bool   `bun:"watch_enabled" json:"watch_enabled"`
	CreatedAt    int64  `bun:"created_at" json:"created_at"`
	UpdatedAt    int64  `bun:"updated_at" json:"updated_at"`
}
