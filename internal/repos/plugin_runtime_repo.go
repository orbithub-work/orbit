package repos

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"media-assistant-os/internal/models"

	"github.com/uptrace/bun"
)

type PluginRuntimeRepo struct {
	db *bun.DB
}

func NewPluginRuntimeRepo(db *bun.DB) *PluginRuntimeRepo {
	return &PluginRuntimeRepo{db: db}
}

func (r *PluginRuntimeRepo) Upsert(ctx context.Context, item models.PluginRuntime) error {
	item.PluginID = strings.TrimSpace(item.PluginID)
	item.Token = strings.TrimSpace(item.Token)
	if item.PluginID == "" || item.Token == "" {
		return nil
	}
	if item.CreatedAt <= 0 {
		item.CreatedAt = time.Now().Unix()
	}
	item.UpdatedAt = time.Now().Unix()

	_, err := r.db.NewInsert().
		Model(&item).
		On("CONFLICT (plugin_id) DO UPDATE").
		Set("token = EXCLUDED.token").
		Set("info_json = EXCLUDED.info_json").
		Set("issued_at = EXCLUDED.issued_at").
		Set("last_used_at = EXCLUDED.last_used_at").
		Set("expires_at = EXCLUDED.expires_at").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func (r *PluginRuntimeRepo) List(ctx context.Context) ([]models.PluginRuntime, error) {
	var out []models.PluginRuntime
	err := r.db.NewSelect().
		Model(&out).
		OrderExpr("updated_at DESC").
		Scan(ctx)
	return out, err
}

func (r *PluginRuntimeRepo) Delete(ctx context.Context, pluginID string) error {
	pluginID = strings.TrimSpace(pluginID)
	if pluginID == "" {
		return nil
	}
	_, err := r.db.NewDelete().
		Model((*models.PluginRuntime)(nil)).
		Where("plugin_id = ?", pluginID).
		Exec(ctx)
	return err
}

func (r *PluginRuntimeRepo) DeleteExpired(ctx context.Context, now time.Time) (int, error) {
	res, err := r.db.NewDelete().
		Model((*models.PluginRuntime)(nil)).
		Where("expires_at > 0").
		Where("expires_at <= ?", now.Unix()).
		Exec(ctx)
	if err != nil {
		return 0, err
	}
	affected, _ := res.RowsAffected()
	return int(affected), nil
}

func (r *PluginRuntimeRepo) GetByPluginID(ctx context.Context, pluginID string) (*models.PluginRuntime, error) {
	pluginID = strings.TrimSpace(pluginID)
	if pluginID == "" {
		return nil, nil
	}
	var out models.PluginRuntime
	err := r.db.NewSelect().
		Model(&out).
		Where("plugin_id = ?", pluginID).
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
