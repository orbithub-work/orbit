package repos

import (
	"context"
	"strings"
	"time"

	"media-assistant-os/internal/db"
	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type ActivityRepo struct {
	db *db.DB
}

func NewActivityRepo(db *db.DB) *ActivityRepo {
	return &ActivityRepo{db: db}
}

func (r *ActivityRepo) Create(ctx context.Context, level, message string) error {
	return r.CreateEx(ctx, level, message, "", "")
}

func (r *ActivityRepo) CreateEx(ctx context.Context, level, message, assetID, projectID string) error {
	log := &models.ActivityLog{
		ID:        utils.NewID(),
		Level:     level,
		Message:   message,
		AssetID:   assetID,
		ProjectID: projectID,
		CreatedAt: time.Now().UnixMilli(),
	}

	_, err := r.db.ORM().NewInsert().Model(log).Exec(ctx)
	return err
}

func (r *ActivityRepo) List(ctx context.Context, limit int) ([]models.ActivityLog, error) {
	var logs []models.ActivityLog
	err := r.db.ORM().NewSelect().
		Model(&logs).
		Order("created_at DESC").
		Limit(limit).
		Scan(ctx)
	return logs, err
}

func (r *ActivityRepo) ListByProject(ctx context.Context, projectID string, levels []string, limit int) ([]models.ActivityLog, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []models.ActivityLog{}, nil
	}
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	var logs []models.ActivityLog
	q := r.db.ORM().NewSelect().
		Model(&logs).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Limit(limit)

	if len(levels) > 0 {
		cleaned := make([]string, 0, len(levels))
		for _, lv := range levels {
			v := strings.TrimSpace(lv)
			if v == "" {
				continue
			}
			cleaned = append(cleaned, v)
		}
		if len(cleaned) > 0 {
			q = q.Where("level IN (?)", bun.In(cleaned))
		}
	}
	if err := q.Scan(ctx); err != nil {
		return nil, err
	}
	return logs, nil
}
