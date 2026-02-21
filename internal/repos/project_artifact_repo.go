package repos

import (
	"context"
	"database/sql"

	"media-assistant-os/internal/models"

	"github.com/uptrace/bun"
)

type ProjectArtifactRepo struct {
	db *bun.DB
}

func NewProjectArtifactRepo(db *bun.DB) *ProjectArtifactRepo {
	return &ProjectArtifactRepo{db: db}
}

func (r *ProjectArtifactRepo) ListByProject(ctx context.Context, projectID string, kind string, limit int) ([]models.ProjectArtifact, error) {
	if limit <= 0 || limit > 500 {
		limit = 200
	}
	var out []models.ProjectArtifact
	q := r.db.NewSelect().
		Model(&out).
		Where("project_id = ?", projectID).
		OrderExpr("created_at DESC").
		Limit(limit)
	if kind != "" {
		q = q.Where("kind = ?", kind)
	}
	err := q.Scan(ctx)
	return out, err
}

func (r *ProjectArtifactRepo) Get(ctx context.Context, id string) (*models.ProjectArtifact, error) {
	var out models.ProjectArtifact
	err := r.db.NewSelect().
		Model(&out).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (r *ProjectArtifactRepo) Create(ctx context.Context, a *models.ProjectArtifact) error {
	_, err := r.db.NewInsert().Model(a).Exec(ctx)
	return err
}

func (r *ProjectArtifactRepo) Update(ctx context.Context, a *models.ProjectArtifact) error {
	_, err := r.db.NewUpdate().
		Model(a).
		Where("id = ?", a.ID).
		Exec(ctx)
	return err
}

func (r *ProjectArtifactRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.ProjectArtifact)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}
