package models

import "github.com/uptrace/bun"

// PluginRuntime stores plugin runtime registration info so the host can recover state after restart.
type PluginRuntime struct {
	bun.BaseModel `bun:"table:plugin_runtimes"`

	PluginID   string `bun:"plugin_id,pk" json:"plugin_id"`
	Token      string `bun:"token" json:"token"`
	InfoJSON   string `bun:"info_json" json:"info_json"`
	IssuedAt   int64  `bun:"issued_at" json:"issued_at"`
	LastUsedAt int64  `bun:"last_used_at" json:"last_used_at"`
	ExpiresAt  int64  `bun:"expires_at" json:"expires_at"`
	CreatedAt  int64  `bun:"created_at" json:"created_at"`
	UpdatedAt  int64  `bun:"updated_at" json:"updated_at"`
}
