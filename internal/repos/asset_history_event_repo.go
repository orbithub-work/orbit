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

type AssetHistoryEventRepo struct {
	db *bun.DB
}

type CreateAssetHistoryEventInput struct {
	AssetID    string
	ProjectID  string
	EventType  string
	SourcePath string
	TargetPath string
	Confidence string
	IsInferred bool
	Detail     string
	OccurredAt int64
}

func NewAssetHistoryEventRepo(db *bun.DB) *AssetHistoryEventRepo {
	return &AssetHistoryEventRepo{db: db}
}

func (r *AssetHistoryEventRepo) CreateDedup(ctx context.Context, in CreateAssetHistoryEventInput, dedupWindowSec int64) (*models.AssetHistoryEvent, error) {
	norm := normalizeAssetHistoryEventInput(in)
	if norm.AssetID == "" || norm.EventType == "" {
		return nil, nil
	}

	if dedupWindowSec > 0 {
		var latest models.AssetHistoryEvent
		err := r.db.NewSelect().
			Model(&latest).
			Where("asset_id = ?", norm.AssetID).
			Where("event_type = ?", norm.EventType).
			Where("source_path = ?", norm.SourcePath).
			Where("target_path = ?", norm.TargetPath).
			OrderExpr("occurred_at DESC").
			Limit(1).
			Scan(ctx)
		if err == nil {
			if norm.OccurredAt-latest.OccurredAt <= dedupWindowSec {
				return &latest, nil
			}
		} else if !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
	}

	item := &models.AssetHistoryEvent{
		ID:         utils.NewID(),
		AssetID:    norm.AssetID,
		ProjectID:  norm.ProjectID,
		EventType:  norm.EventType,
		SourcePath: norm.SourcePath,
		TargetPath: norm.TargetPath,
		Confidence: norm.Confidence,
		IsInferred: norm.IsInferred,
		Detail:     norm.Detail,
		OccurredAt: norm.OccurredAt,
		CreatedAt:  time.Now().Unix(),
	}
	if _, err := r.db.NewInsert().Model(item).Exec(ctx); err != nil {
		return nil, err
	}
	return item, nil
}

func (r *AssetHistoryEventRepo) ListByAsset(ctx context.Context, assetID string, limit int) ([]models.AssetHistoryEvent, error) {
	assetID = strings.TrimSpace(assetID)
	if assetID == "" {
		return []models.AssetHistoryEvent{}, nil
	}
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	var out []models.AssetHistoryEvent
	err := r.db.NewSelect().
		Model(&out).
		Where("asset_id = ?", assetID).
		OrderExpr("occurred_at DESC, created_at DESC").
		Limit(limit).
		Scan(ctx)
	return out, err
}

func normalizeAssetHistoryEventInput(in CreateAssetHistoryEventInput) CreateAssetHistoryEventInput {
	out := in
	out.AssetID = strings.TrimSpace(out.AssetID)
	out.ProjectID = strings.TrimSpace(out.ProjectID)
	out.SourcePath = normalizeHistoryPath(out.SourcePath)
	out.TargetPath = normalizeHistoryPath(out.TargetPath)
	out.EventType = normalizeHistoryEventType(out.EventType)
	out.Confidence = normalizeHistoryConfidence(out.Confidence)
	out.Detail = strings.TrimSpace(out.Detail)
	if out.OccurredAt <= 0 {
		out.OccurredAt = time.Now().Unix()
	}
	return out
}

func normalizeHistoryPath(path string) string {
	v := strings.TrimSpace(path)
	if v == "" {
		return ""
	}
	if abs, err := filepath.Abs(v); err == nil {
		v = abs
	}
	return filepath.Clean(v)
}

func normalizeHistoryEventType(eventType string) string {
	switch strings.ToLower(strings.TrimSpace(eventType)) {
	case models.AssetHistoryEventCreated:
		return models.AssetHistoryEventCreated
	case models.AssetHistoryEventCopied:
		return models.AssetHistoryEventCopied
	case models.AssetHistoryEventRenamed:
		return models.AssetHistoryEventRenamed
	case models.AssetHistoryEventMoved:
		return models.AssetHistoryEventMoved
	case models.AssetHistoryEventDeleted:
		return models.AssetHistoryEventDeleted
	case models.AssetHistoryEventModified:
		return models.AssetHistoryEventModified
	default:
		return ""
	}
}

func normalizeHistoryConfidence(confidence string) string {
	switch strings.ToLower(strings.TrimSpace(confidence)) {
	case "high":
		return "high"
	case "low":
		return "low"
	default:
		return "medium"
	}
}
