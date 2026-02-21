package repos

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"media-assistant-os/internal/models"

	"github.com/uptrace/bun"
)

type ProjectAssetRepo struct {
	db *bun.DB
}

type ProjectAssetLinkOptions struct {
	Role       string
	BindMode   string
	Confidence float64
	SourceID   string
}

type ProjectAssetBindingDetail struct {
	ProjectID  string  `json:"project_id"`
	AssetID    string  `json:"asset_id"`
	SourceID   string  `json:"source_id,omitempty"`
	Role       string  `json:"role"`
	BindMode   string  `json:"bind_mode"`
	Confidence float64 `json:"confidence"`
	LinkStatus string  `json:"link_status"`

	Path        string `json:"path"`
	AssetStatus string `json:"asset_status"`
	Mtime       int64  `json:"mtime"`
	Size        int64  `json:"size"`
}

func NewProjectAssetRepo(db *bun.DB) *ProjectAssetRepo {
	return &ProjectAssetRepo{db: db}
}

func (r *ProjectAssetRepo) Link(ctx context.Context, projectID string, assetID string) error {
	return r.LinkWithOptions(ctx, projectID, assetID, nil)
}

func (r *ProjectAssetRepo) LinkWithOptions(ctx context.Context, projectID string, assetID string, opts *ProjectAssetLinkOptions) error {
	projectID = strings.TrimSpace(projectID)
	assetID = strings.TrimSpace(assetID)
	if projectID == "" || assetID == "" {
		return nil
	}
	norm := normalizeProjectAssetLinkOptions(opts)
	now := time.Now().Unix()
	link := models.ProjectAsset{
		ProjectID:  projectID,
		AssetID:    assetID,
		SourceID:   strings.TrimSpace(norm.SourceID),
		Role:       norm.Role,
		BindMode:   norm.BindMode,
		Confidence: norm.Confidence,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err := r.db.NewInsert().
		Model(&link).
		On("CONFLICT (project_id, asset_id) DO UPDATE").
		Set("source_id = EXCLUDED.source_id").
		Set("role = EXCLUDED.role").
		Set("bind_mode = EXCLUDED.bind_mode").
		Set("confidence = EXCLUDED.confidence").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func (r *ProjectAssetRepo) GetAssetsByProject(ctx context.Context, projectID string) ([]models.Asset, error) {
	var assets []models.Asset
	err := r.db.NewSelect().
		Model(&assets).
		Join("JOIN project_assets AS pa ON pa.asset_id = asset.id").
		Where("pa.project_id = ?", projectID).
		Scan(ctx)
	return assets, err
}

func (r *ProjectAssetRepo) UpdateStatus(ctx context.Context, projectID string, assetID string, status *string) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.ProjectAsset)(nil)).
		Set("status = ?", status).
		Set("updated_at = ?", now).
		Where("project_id = ?", projectID).
		Where("asset_id = ?", assetID).
		Exec(ctx)
	return err
}

func (r *ProjectAssetRepo) UpdateBinding(ctx context.Context, projectID string, assetID string, role *string, status *string, bindMode *string, confidence *float64) error {
	projectID = strings.TrimSpace(projectID)
	assetID = strings.TrimSpace(assetID)
	if projectID == "" || assetID == "" {
		return nil
	}

	now := time.Now().Unix()
	q := r.db.NewUpdate().
		Model((*models.ProjectAsset)(nil)).
		Set("updated_at = ?", now).
		Where("project_id = ?", projectID).
		Where("asset_id = ?", assetID)

	if role != nil {
		v := normalizeProjectAssetRole(*role)
		q = q.Set("role = ?", v)
	}
	if status != nil {
		q = q.Set("status = ?", status)
	}
	if bindMode != nil {
		v := normalizeProjectAssetBindMode(*bindMode)
		q = q.Set("bind_mode = ?", v)
	}
	if confidence != nil {
		v := *confidence
		if v < 0 {
			v = 0
		}
		if v > 1 {
			v = 1
		}
		q = q.Set("confidence = ?", v)
	}

	_, err := q.Exec(ctx)
	return err
}

func (r *ProjectAssetRepo) Unlink(ctx context.Context, projectID string, assetID string) error {
	_, err := r.db.NewDelete().
		Model((*models.ProjectAsset)(nil)).
		Where("project_id = ?", projectID).
		Where("asset_id = ?", assetID).
		Exec(ctx)
	return err
}

func (r *ProjectAssetRepo) UnlinkProject(ctx context.Context, projectID string) error {
	_, err := r.db.NewDelete().
		Model((*models.ProjectAsset)(nil)).
		Where("project_id = ?", projectID).
		Exec(ctx)
	return err
}

func (r *ProjectAssetRepo) UnlinkAll(ctx context.Context, assetID string) error {
	_, err := r.db.NewDelete().
		Model((*models.ProjectAsset)(nil)).
		Where("asset_id = ?", assetID).
		Exec(ctx)
	return err
}

// BatchLink links multiple assets to a project in one insert.
func (r *ProjectAssetRepo) BatchLink(ctx context.Context, projectID string, assetIDs []string) error {
	return r.BatchLinkWithOptions(ctx, projectID, assetIDs, nil)
}

// BatchLinkWithOptions links multiple assets to a project in one insert.
func (r *ProjectAssetRepo) BatchLinkWithOptions(ctx context.Context, projectID string, assetIDs []string, opts *ProjectAssetLinkOptions) error {
	if len(assetIDs) == 0 {
		return nil
	}
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil
	}
	norm := normalizeProjectAssetLinkOptions(opts)
	now := time.Now().Unix()
	links := make([]models.ProjectAsset, 0, len(assetIDs))
	for _, id := range assetIDs {
		assetID := strings.TrimSpace(id)
		if assetID == "" {
			continue
		}
		links = append(links, models.ProjectAsset{
			ProjectID:  projectID,
			AssetID:    assetID,
			SourceID:   strings.TrimSpace(norm.SourceID),
			Role:       norm.Role,
			BindMode:   norm.BindMode,
			Confidence: norm.Confidence,
			CreatedAt:  now,
			UpdatedAt:  now,
		})
	}
	if len(links) == 0 {
		return nil
	}
	const batchSize = 500
	for i := 0; i < len(links); i += batchSize {
		end := i + batchSize
		if end > len(links) {
			end = len(links)
		}
		batch := links[i:end]
		_, err := r.db.NewInsert().
			Model(&batch).
			On("CONFLICT (project_id, asset_id) DO UPDATE").
			Set("source_id = EXCLUDED.source_id").
			Set("role = EXCLUDED.role").
			Set("bind_mode = EXCLUDED.bind_mode").
			Set("confidence = EXCLUDED.confidence").
			Set("updated_at = EXCLUDED.updated_at").
			Exec(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *ProjectAssetRepo) UnlinkBySource(ctx context.Context, projectID string, sourceID string) error {
	projectID = strings.TrimSpace(projectID)
	sourceID = strings.TrimSpace(sourceID)
	if projectID == "" || sourceID == "" {
		return nil
	}
	_, err := r.db.NewDelete().
		Model((*models.ProjectAsset)(nil)).
		Where("project_id = ?", projectID).
		Where("source_id = ?", sourceID).
		Exec(ctx)
	return err
}

func (r *ProjectAssetRepo) ListProjectIDsByAsset(ctx context.Context, assetID string) ([]string, error) {
	var projectIDs []string
	err := r.db.NewSelect().
		Model((*models.ProjectAsset)(nil)).
		Column("project_id").
		Where("asset_id = ?", assetID).
		Scan(ctx, &projectIDs)
	return projectIDs, err
}

func (r *ProjectAssetRepo) ListRecentAssetsByProject(ctx context.Context, projectID string, excludeAssetID string, limit int) ([]models.Asset, error) {
	if limit <= 0 || limit > 300 {
		limit = 100
	}
	var assets []models.Asset
	q := r.db.NewSelect().
		Model(&assets).
		Join("JOIN project_assets AS pa ON pa.asset_id = asset.id").
		Where("pa.project_id = ?", projectID).
		OrderExpr("asset.mtime DESC, asset.created_at DESC").
		Limit(limit)
	if excludeAssetID != "" {
		q = q.Where("asset.id != ?", excludeAssetID)
	}
	err := q.Scan(ctx)
	return assets, err
}

func (r *ProjectAssetRepo) ListBindingDetailsByProject(ctx context.Context, projectID string) ([]ProjectAssetBindingDetail, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []ProjectAssetBindingDetail{}, nil
	}

	var details []ProjectAssetBindingDetail
	err := r.db.NewSelect().
		TableExpr("project_assets AS pa").
		Join("JOIN assets AS a ON a.id = pa.asset_id").
		ColumnExpr("pa.project_id AS project_id").
		ColumnExpr("pa.asset_id AS asset_id").
		ColumnExpr("COALESCE(pa.source_id, '') AS source_id").
		ColumnExpr("COALESCE(pa.role, 'source') AS role").
		ColumnExpr("COALESCE(pa.bind_mode, 'auto') AS bind_mode").
		ColumnExpr("COALESCE(pa.confidence, 1.0) AS confidence").
		ColumnExpr("COALESCE(pa.status, '') AS link_status").
		ColumnExpr("a.path AS path").
		ColumnExpr("a.status AS asset_status").
		ColumnExpr("a.mtime AS mtime").
		ColumnExpr("a.size AS size").
		Where("pa.project_id = ?", projectID).
		OrderExpr("a.mtime DESC, a.created_at DESC").
		Scan(ctx, &details)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].Path = filepath.Clean(details[i].Path)
		details[i].Role = normalizeProjectAssetRole(details[i].Role)
		details[i].BindMode = normalizeProjectAssetBindMode(details[i].BindMode)
		if details[i].Confidence < 0 {
			details[i].Confidence = 0
		}
		if details[i].Confidence > 1 {
			details[i].Confidence = 1
		}
	}
	return details, nil
}

func normalizeProjectAssetLinkOptions(opts *ProjectAssetLinkOptions) ProjectAssetLinkOptions {
	out := ProjectAssetLinkOptions{
		Role:       "source",
		BindMode:   "auto",
		Confidence: 1.0,
	}
	if opts == nil {
		return out
	}
	out.Role = normalizeProjectAssetRole(opts.Role)
	out.BindMode = normalizeProjectAssetBindMode(opts.BindMode)
	out.SourceID = strings.TrimSpace(opts.SourceID)
	if opts.Confidence >= 0 && opts.Confidence <= 1 {
		out.Confidence = opts.Confidence
	}
	return out
}

func normalizeProjectAssetRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "engine":
		return "engine"
	case "deliverable":
		return "deliverable"
	default:
		return "source"
	}
}

func normalizeProjectAssetBindMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case "manual":
		return "manual"
	case "source":
		return "source"
	default:
		return "auto"
	}
}
