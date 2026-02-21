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

type LibrarySourceRepo struct {
	db *bun.DB
}

func NewLibrarySourceRepo(db *bun.DB) *LibrarySourceRepo {
	return &LibrarySourceRepo{db: db}
}

func (r *LibrarySourceRepo) Upsert(ctx context.Context, rootPath string, watchEnabled bool) (*models.LibrarySource, error) {
	rootPath = normalizeLibraryRootPath(rootPath)
	if rootPath == "" {
		return nil, nil
	}

	now := time.Now().Unix()
	item := &models.LibrarySource{
		ID:           utils.NewID(),
		RootPath:     rootPath,
		WatchEnabled: watchEnabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	_, err := r.db.NewInsert().
		Model(item).
		On("CONFLICT (root_path) DO UPDATE").
		Set("watch_enabled = EXCLUDED.watch_enabled").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	return r.GetByPath(ctx, rootPath)
}

func (r *LibrarySourceRepo) GetByPath(ctx context.Context, rootPath string) (*models.LibrarySource, error) {
	rootPath = normalizeLibraryRootPath(rootPath)
	if rootPath == "" {
		return nil, nil
	}

	var out models.LibrarySource
	err := r.db.NewSelect().
		Model(&out).
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

func (r *LibrarySourceRepo) List(ctx context.Context) ([]models.LibrarySource, error) {
	var out []models.LibrarySource
	err := r.db.NewSelect().
		Model(&out).
		OrderExpr("created_at DESC").
		Scan(ctx)
	return out, err
}

func (r *LibrarySourceRepo) Remove(ctx context.Context, rootPath string) error {
	rootPath = normalizeLibraryRootPath(rootPath)
	if rootPath == "" {
		return nil
	}
	_, err := r.db.NewDelete().
		Model((*models.LibrarySource)(nil)).
		Where("root_path = ?", rootPath).
		Exec(ctx)
	return err
}

func normalizeLibraryRootPath(path string) string {
	v := strings.TrimSpace(path)
	if v == "" {
		return ""
	}
	if abs, err := filepath.Abs(v); err == nil {
		v = abs
	}
	return filepath.Clean(v)
}
