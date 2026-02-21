package repos

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"strings"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type ProjectSourceRepo struct {
	db *bun.DB
}

func NewProjectSourceRepo(db *bun.DB) *ProjectSourceRepo {
	return &ProjectSourceRepo{db: db}
}

func (r *ProjectSourceRepo) Upsert(ctx context.Context, source *models.ProjectSource) (*models.ProjectSource, error) {
	if source == nil {
		return nil, nil
	}
	projectID := strings.TrimSpace(source.ProjectID)
	rootPath := normalizeRootPath(source.RootPath)
	if projectID == "" || rootPath == "" {
		return nil, nil
	}

	now := time.Now().Unix()
	if source.ID == "" {
		source.ID = utils.NewID()
	}
	source.ProjectID = projectID
	source.RootPath = rootPath
	source.SourceType = normalizeSourceType(source.SourceType)
	source.UpdatedAt = now
	if source.CreatedAt == 0 {
		source.CreatedAt = now
	}

	_, err := r.db.NewInsert().
		Model(source).
		On("CONFLICT (project_id, root_path) DO UPDATE").
		Set("source_type = EXCLUDED.source_type").
		Set("watch_enabled = EXCLUDED.watch_enabled").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	out, err := r.GetByProjectAndPath(ctx, projectID, rootPath)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (r *ProjectSourceRepo) SetPrimary(ctx context.Context, projectID string, rootPath string) (*models.ProjectSource, error) {
	projectID = strings.TrimSpace(projectID)
	rootPath = normalizeRootPath(rootPath)
	if projectID == "" || rootPath == "" {
		return nil, nil
	}
	now := time.Now().Unix()

	_, err := r.db.NewUpdate().
		Model((*models.ProjectSource)(nil)).
		Set("source_type = ?", "extra").
		Set("updated_at = ?", now).
		Where("project_id = ?", projectID).
		Where("source_type = ?", "primary").
		Where("root_path != ?", rootPath).
		Exec(ctx)
	if err != nil {
		return nil, err
	}

	return r.Upsert(ctx, &models.ProjectSource{
		ID:           utils.NewID(),
		ProjectID:    projectID,
		RootPath:     rootPath,
		SourceType:   "primary",
		WatchEnabled: true,
		CreatedAt:    now,
		UpdatedAt:    now,
	})
}

func (r *ProjectSourceRepo) Remove(ctx context.Context, projectID string, rootPath string) error {
	projectID = strings.TrimSpace(projectID)
	rootPath = normalizeRootPath(rootPath)
	if projectID == "" || rootPath == "" {
		return nil
	}
	_, err := r.db.NewDelete().
		Model((*models.ProjectSource)(nil)).
		Where("project_id = ?", projectID).
		Where("root_path = ?", rootPath).
		Exec(ctx)
	return err
}

func (r *ProjectSourceRepo) RemoveByProject(ctx context.Context, projectID string) error {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil
	}
	_, err := r.db.NewDelete().
		Model((*models.ProjectSource)(nil)).
		Where("project_id = ?", projectID).
		Exec(ctx)
	return err
}

func (r *ProjectSourceRepo) ListByProject(ctx context.Context, projectID string) ([]models.ProjectSource, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []models.ProjectSource{}, nil
	}
	var out []models.ProjectSource
	err := r.db.NewSelect().
		Model(&out).
		Where("project_id = ?", projectID).
		OrderExpr("CASE WHEN source_type = 'primary' THEN 0 ELSE 1 END ASC").
		OrderExpr("created_at ASC").
		Scan(ctx)
	return out, err
}

func (r *ProjectSourceRepo) ListWatchEnabledByProject(ctx context.Context, projectID string) ([]models.ProjectSource, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []models.ProjectSource{}, nil
	}
	var out []models.ProjectSource
	err := r.db.NewSelect().
		Model(&out).
		Where("project_id = ?", projectID).
		Where("watch_enabled = ?", true).
		OrderExpr("CASE WHEN source_type = 'primary' THEN 0 ELSE 1 END ASC").
		OrderExpr("created_at ASC").
		Scan(ctx)
	return out, err
}

func (r *ProjectSourceRepo) GetByProjectAndPath(ctx context.Context, projectID string, rootPath string) (*models.ProjectSource, error) {
	projectID = strings.TrimSpace(projectID)
	rootPath = normalizeRootPath(rootPath)
	if projectID == "" || rootPath == "" {
		return nil, nil
	}
	var out models.ProjectSource
	err := r.db.NewSelect().
		Model(&out).
		Where("project_id = ?", projectID).
		Where("root_path = ?", rootPath).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *ProjectSourceRepo) GetByID(ctx context.Context, sourceID string) (*models.ProjectSource, error) {
	sourceID = strings.TrimSpace(sourceID)
	if sourceID == "" {
		return nil, nil
	}
	var out models.ProjectSource
	err := r.db.NewSelect().
		Model(&out).
		Where("id = ?", sourceID).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (r *ProjectSourceRepo) CountByRootPath(ctx context.Context, rootPath string) (int, error) {
	rootPath = normalizeRootPath(rootPath)
	if rootPath == "" {
		return 0, nil
	}
	return r.db.NewSelect().
		Model((*models.ProjectSource)(nil)).
		Where("root_path = ?", rootPath).
		Count(ctx)
}

func normalizeRootPath(path string) string {
	v := strings.TrimSpace(path)
	if v == "" {
		return ""
	}
	if abs, err := filepath.Abs(v); err == nil {
		v = abs
	}
	return filepath.Clean(v)
}

func normalizeSourceType(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "primary":
		return "primary"
	default:
		return "extra"
	}
}
