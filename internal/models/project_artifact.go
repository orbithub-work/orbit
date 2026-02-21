package models

import (
	"encoding/json"

	"github.com/uptrace/bun"
)

// ProjectArtifact represents a "产物" produced inside a project workflow.
// It can be a file artifact (path) or a text artifact (content), with extensible meta.
type ProjectArtifact struct {
	bun.BaseModel `bun:"table:project_artifacts"`

	ID             string          `bun:",pk" json:"id"`
	ProjectID      string          `bun:"project_id" json:"project_id"`
	Kind           string          `bun:"kind" json:"kind"` // e.g. "script", "cover", "export_bundle"
	Name           string          `bun:"name" json:"name"`
	Path           string          `bun:"path" json:"path"`
	Content        string          `bun:"content" json:"content"`
	MetaJSON       string          `bun:"meta_json" json:"-"`
	Meta           json.RawMessage `bun:"-" json:"meta,omitempty"`
	SourcePluginID string          `bun:"source_plugin_id" json:"source_plugin_id"`
	CreatedAt      int64           `bun:"created_at" json:"created_at"`
	UpdatedAt      int64           `bun:"updated_at" json:"updated_at"`
}
