package models

import (
	"time"

	"github.com/uptrace/bun"
)

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
)

type MediaTask struct {
	bun.BaseModel `bun:"table:media_tasks,alias:mt"`

	ID         string     `bun:"id,pk" json:"id"`
	AssetID    string     `bun:"asset_id,notnull" json:"asset_id"`
	TaskType   string     `bun:"task_type,notnull" json:"task_type"` // "metadata", "thumbnail", "ai_tags", etc.
	Status     TaskStatus `bun:"status,notnull,default:'pending'" json:"status"`
	Priority   int        `bun:"priority,notnull,default:0" json:"priority"`
	RetryCount int        `bun:"retry_count,notnull,default:0" json:"retry_count"`
	MaxRetries int        `bun:"max_retries,notnull,default:3" json:"max_retries"`

	// For external/distributed workers
	WorkerID   string    `bun:"worker_id" json:"worker_id,omitempty"`
	StartedAt  time.Time `bun:"started_at" json:"started_at,omitempty"`
	FinishedAt time.Time `bun:"finished_at" json:"finished_at,omitempty"`
	LeaseUntil time.Time `bun:"lease_until" json:"lease_until,omitempty"`

	ErrorMessage     string    `bun:"error_message" json:"error_message,omitempty"`
	Progress         int       `bun:"progress,notnull,default:0" json:"progress"`
	NextRetryAt      time.Time `bun:"next_retry_at" json:"next_retry_at,omitempty"`
	DeadLetterReason string    `bun:"dead_letter_reason" json:"dead_letter_reason,omitempty"`

	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
