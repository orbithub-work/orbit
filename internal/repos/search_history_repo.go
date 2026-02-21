package repos

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type SearchHistoryRepo struct {
	db *bun.DB
}

func NewSearchHistoryRepo(db *bun.DB) *SearchHistoryRepo {
	return &SearchHistoryRepo{db: db}
}

func (r *SearchHistoryRepo) RecordSearch(ctx context.Context, query string, filters map[string]any) error {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}
	if filters == nil {
		filters = map[string]any{}
	}
	filtersJSON, err := json.Marshal(filters)
	if err != nil {
		return err
	}

	hash := computeSearchHash(query, string(filtersJSON))
	now := time.Now().Unix()
	item := &models.SearchHistory{
		ID:        utils.NewID(),
		QueryHash: hash,
		Query:     query,
		Filters:   string(filtersJSON),
		Count:     1,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err = r.db.NewInsert().
		Model(item).
		On("CONFLICT (query_hash) DO UPDATE").
		Set("query = EXCLUDED.query").
		Set("filters = EXCLUDED.filters").
		Set("count = search_history.count + 1").
		Set("updated_at = EXCLUDED.updated_at").
		Exec(ctx)
	return err
}

func (r *SearchHistoryRepo) GetRecent(ctx context.Context, limit int) ([]models.SearchHistory, error) {
	if limit <= 0 || limit > 200 {
		limit = 20
	}
	var out []models.SearchHistory
	err := r.db.NewSelect().
		Model(&out).
		OrderExpr("updated_at DESC").
		Limit(limit).
		Scan(ctx)
	return out, err
}

func (r *SearchHistoryRepo) Clear(ctx context.Context) error {
	_, err := r.db.NewDelete().
		Model((*models.SearchHistory)(nil)).
		Where("1=1").
		Exec(ctx)
	return err
}

func (r *SearchHistoryRepo) GetByHash(ctx context.Context, hash string) (*models.SearchHistory, error) {
	hash = strings.TrimSpace(hash)
	if hash == "" {
		return nil, nil
	}
	var out models.SearchHistory
	err := r.db.NewSelect().
		Model(&out).
		Where("query_hash = ?", hash).
		Limit(1).
		Scan(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func computeSearchHash(query string, filters string) string {
	h := sha256.New()
	h.Write([]byte(strings.TrimSpace(query)))
	h.Write([]byte{0})
	h.Write([]byte(strings.TrimSpace(filters)))
	return hex.EncodeToString(h.Sum(nil))[:16]
}
