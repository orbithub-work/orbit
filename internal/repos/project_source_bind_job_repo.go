package repos

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type ProjectSourceBindJobRepo struct {
	db *bun.DB
}

func NewProjectSourceBindJobRepo(db *bun.DB) *ProjectSourceBindJobRepo {
	return &ProjectSourceBindJobRepo{db: db}
}

func (r *ProjectSourceBindJobRepo) Create(ctx context.Context, projectID, sourceID, rootPath string, totalAssets int) (*models.ProjectSourceBindJob, error) {
	projectID = strings.TrimSpace(projectID)
	sourceID = strings.TrimSpace(sourceID)
	rootPath = strings.TrimSpace(rootPath)
	if projectID == "" || sourceID == "" || rootPath == "" {
		return nil, nil
	}
	now := time.Now().Unix()
	job := &models.ProjectSourceBindJob{
		ID:              utils.NewID(),
		ProjectID:       projectID,
		SourceID:        sourceID,
		RootPath:        rootPath,
		Status:          models.ProjectSourceBindJobPending,
		TotalAssets:     totalAssets,
		ProcessedAssets: 0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if _, err := r.db.NewInsert().Model(job).Exec(ctx); err != nil {
		return nil, err
	}
	return job, nil
}

func (r *ProjectSourceBindJobRepo) Get(ctx context.Context, id string) (*models.ProjectSourceBindJob, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, nil
	}
	var out models.ProjectSourceBindJob
	err := r.db.NewSelect().
		Model(&out).
		Where("id = ?", id).
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

func (r *ProjectSourceBindJobRepo) ListByProject(ctx context.Context, projectID string, limit int) ([]models.ProjectSourceBindJob, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []models.ProjectSourceBindJob{}, nil
	}
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var out []models.ProjectSourceBindJob
	err := r.db.NewSelect().
		Model(&out).
		Where("project_id = ?", projectID).
		OrderExpr("created_at DESC").
		Limit(limit).
		Scan(ctx)
	return out, err
}

func (r *ProjectSourceBindJobRepo) ListInflight(ctx context.Context, limit int) ([]models.ProjectSourceBindJob, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var out []models.ProjectSourceBindJob
	err := r.db.NewSelect().
		Model(&out).
		Where("status IN (?)", bun.In([]string{
			models.ProjectSourceBindJobPending,
			models.ProjectSourceBindJobRunning,
		})).
		OrderExpr("created_at ASC").
		Limit(limit).
		Scan(ctx)
	return out, err
}

func (r *ProjectSourceBindJobRepo) GetInflightByProjectAndSource(ctx context.Context, projectID string, sourceID string) (*models.ProjectSourceBindJob, error) {
	projectID = strings.TrimSpace(projectID)
	sourceID = strings.TrimSpace(sourceID)
	if projectID == "" || sourceID == "" {
		return nil, nil
	}

	var out models.ProjectSourceBindJob
	err := r.db.NewSelect().
		Model(&out).
		Where("project_id = ?", projectID).
		Where("source_id = ?", sourceID).
		Where("status IN (?)", bun.In([]string{
			models.ProjectSourceBindJobPending,
			models.ProjectSourceBindJobRunning,
		})).
		OrderExpr("created_at DESC").
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

func (r *ProjectSourceBindJobRepo) MarkRunning(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.ProjectSourceBindJob)(nil)).
		Set("status = ?", models.ProjectSourceBindJobRunning).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Where("status IN (?)", bun.In([]string{
			models.ProjectSourceBindJobPending,
			models.ProjectSourceBindJobRunning,
		})).
		Exec(ctx)
	return err
}

func (r *ProjectSourceBindJobRepo) UpdateProgress(ctx context.Context, id string, processed int) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	if processed < 0 {
		processed = 0
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.ProjectSourceBindJob)(nil)).
		Set("processed_assets = ?", processed).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *ProjectSourceBindJobRepo) MarkSucceeded(ctx context.Context, id string, processed int) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	if processed < 0 {
		processed = 0
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.ProjectSourceBindJob)(nil)).
		Set("status = ?", models.ProjectSourceBindJobSucceeded).
		Set("processed_assets = ?", processed).
		Set("updated_at = ?", now).
		Set("finished_at = ?", now).
		Set("error_message = ?", "").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *ProjectSourceBindJobRepo) MarkFailed(ctx context.Context, id string, processed int, errMsg string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil
	}
	if processed < 0 {
		processed = 0
	}
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.ProjectSourceBindJob)(nil)).
		Set("status = ?", models.ProjectSourceBindJobFailed).
		Set("processed_assets = ?", processed).
		Set("updated_at = ?", now).
		Set("finished_at = ?", now).
		Set("error_message = ?", strings.TrimSpace(errMsg)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
