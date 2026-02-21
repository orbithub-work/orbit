package models

import "github.com/uptrace/bun"

const (
	AssetHistoryEventCreated  = "created"
	AssetHistoryEventCopied   = "copied"
	AssetHistoryEventRenamed  = "renamed"
	AssetHistoryEventMoved    = "moved"
	AssetHistoryEventDeleted  = "deleted"
	AssetHistoryEventModified = "modified"
)

type AssetHistoryEvent struct {
	bun.BaseModel `bun:"table:asset_history_events"`

	ID         string `bun:",pk" json:"id"`
	AssetID    string `bun:"asset_id" json:"asset_id"`
	ProjectID  string `bun:"project_id" json:"project_id,omitempty"`
	EventType  string `bun:"event_type" json:"event_type"`
	SourcePath string `bun:"source_path" json:"source_path,omitempty"`
	TargetPath string `bun:"target_path" json:"target_path,omitempty"`
	Confidence string `bun:"confidence" json:"confidence"` // high | medium | low
	IsInferred bool   `bun:"is_inferred" json:"is_inferred"`
	Detail     string `bun:"detail" json:"detail,omitempty"`
	OccurredAt int64  `bun:"occurred_at" json:"occurred_at"`
	CreatedAt  int64  `bun:"created_at" json:"created_at"`
}
