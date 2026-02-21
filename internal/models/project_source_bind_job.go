package models

import "github.com/uptrace/bun"

const (
	ProjectSourceBindJobPending   = "pending"
	ProjectSourceBindJobRunning   = "running"
	ProjectSourceBindJobSucceeded = "succeeded"
	ProjectSourceBindJobFailed    = "failed"
)

type ProjectSourceBindJob struct {
	bun.BaseModel `bun:"table:project_source_bind_jobs"`

	ID              string `bun:",pk" json:"id"`
	ProjectID       string `bun:"project_id" json:"project_id"`
	SourceID        string `bun:"source_id" json:"source_id"`
	RootPath        string `bun:"root_path" json:"root_path"`
	Status          string `bun:"status" json:"status"`
	TotalAssets     int    `bun:"total_assets" json:"total_assets"`
	ProcessedAssets int    `bun:"processed_assets" json:"processed_assets"`
	ErrorMessage    string `bun:"error_message" json:"error_message,omitempty"`
	CreatedAt       int64  `bun:"created_at" json:"created_at"`
	UpdatedAt       int64  `bun:"updated_at" json:"updated_at"`
	FinishedAt      int64  `bun:"finished_at" json:"finished_at,omitempty"`
}
