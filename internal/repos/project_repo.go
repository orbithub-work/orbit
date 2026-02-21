package repos

import (
	"context"
	"database/sql"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type ProjectRepo struct {
	db *bun.DB
}

func NewProjectRepo(db *bun.DB) *ProjectRepo {
	return &ProjectRepo{db: db}
}

func (r *ProjectRepo) List(ctx context.Context) ([]models.Project, error) {
	var out []models.Project
	err := r.db.NewSelect().
		Model(&out).
		OrderExpr("created_at DESC").
		Scan(ctx)
	return out, err
}

func (r *ProjectRepo) ListByActivity(ctx context.Context) ([]models.Project, error) {
	var out []models.Project
	err := r.db.NewSelect().
		Model(&out).
		OrderExpr("last_activity_at DESC, created_at DESC").
		Scan(ctx)
	return out, err
}

func (r *ProjectRepo) UpdateActivity(ctx context.Context, id string) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Project)(nil)).
		Set("last_activity_at = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *ProjectRepo) Get(ctx context.Context, id string) (*models.Project, error) {
	var p models.Project
	err := r.db.NewSelect().
		Model(&p).
		Where("id = ?", id).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &p, err
}

func (r *ProjectRepo) GetCount(ctx context.Context) (int64, error) {
	count, err := r.db.NewSelect().Model((*models.Project)(nil)).Count(ctx)
	return int64(count), err
}

func (r *ProjectRepo) Create(ctx context.Context, name string, projectType string) (*models.Project, error) {
	now := time.Now().Unix()
	p := models.Project{
		ID:          utils.NewID(),
		Name:        name,
		Path:        "",
		ProjectType: projectType,
		Status:      "active",
		Description: "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := r.db.NewInsert().Model(&p).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *ProjectRepo) UpdatePath(ctx context.Context, id string, path string) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Project)(nil)).
		Set("path = ?", path).
		Set("updated_at = ?", now).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *ProjectRepo) Update(ctx context.Context, p models.Project) error {
	now := time.Now().Unix()
	_, err := r.db.NewUpdate().
		Model((*models.Project)(nil)).
		Set("name = ?", p.Name).
		Set("project_type = ?", p.ProjectType).
		Set("status = ?", p.Status).
		Set("description = ?", p.Description).
		Set("path = ?", p.Path).
		Set("updated_at = ?", now).
		Where("id = ?", p.ID).
		Exec(ctx)
	return err
}

func (r *ProjectRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.Project)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *ProjectRepo) EnsureDefaultProject(ctx context.Context) (*models.Project, error) {
	var p models.Project
	err := r.db.NewSelect().
		Model(&p).
		Where("name = ?", "Default Library").
		Limit(1).
		Scan(ctx)
	if err == nil {
		return &p, nil
	}
	if err != sql.ErrNoRows {
		return nil, err
	}
	return r.Create(ctx, "Default Library", "default")
}
