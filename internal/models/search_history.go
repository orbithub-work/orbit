package models

import "github.com/uptrace/bun"

type SearchHistory struct {
	bun.BaseModel `bun:"table:search_history"`

	ID        string `bun:",pk" json:"id"`
	QueryHash string `bun:"query_hash" json:"query_hash"`
	Query     string `bun:"query" json:"query"`
	Filters   string `bun:"filters" json:"filters,omitempty"`
	Count     int    `bun:"count" json:"count"`
	CreatedAt int64  `bun:"created_at" json:"created_at"`
	UpdatedAt int64  `bun:"updated_at" json:"updated_at"`
}
