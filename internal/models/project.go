package models

import "github.com/uptrace/bun"

type Project struct {
	bun.BaseModel `bun:"table:projects"`

	ID          string `bun:",pk" json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	ProjectType string `json:"project_type"`
	Status      string `json:"status"`
	Description    string `json:"description"`
	LastActivityAt int64  `json:"last_activity_at"` // Unix timestamp for heuristic priority
	CreatedAt      int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}
