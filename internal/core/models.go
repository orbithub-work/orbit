package core

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	ProjectType string `json:"project_type"`
	Status      string `json:"status"`
	Description string `json:"description"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

type Asset struct {
	ID            string  `json:"id"`
	Path          string  `json:"path"`
	Size          int64   `json:"size"`
	Mtime         int64   `json:"mtime"`
	Fingerprint   *string `json:"fingerprint,omitempty"`
	ParentAssetID *string `json:"parent_asset_id,omitempty"`
	Scope         string  `json:"scope"`                  // "global" or "private"
	ProjectID     *string `json:"project_id,omitempty"`   // If scope is private
	CreatedAt     int64   `json:"created_at"`
	UpdatedAt     int64   `json:"updated_at"`
}

type Activity struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	ProjectID string `json:"project_id"`
	Details   string `json:"details"`
	CreatedAt int64  `json:"created_at"`
}
