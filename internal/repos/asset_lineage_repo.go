package repos

import (
	"context"
	"database/sql"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type AssetLineageRepo struct {
	db *bun.DB
}

func NewAssetLineageRepo(db *bun.DB) *AssetLineageRepo {
	return &AssetLineageRepo{db: db}
}

func (r *AssetLineageRepo) GetByPair(ctx context.Context, ancestorID string, descendantID string, relationType string) (*models.AssetLineage, error) {
	var out models.AssetLineage
	err := r.db.NewSelect().
		Model(&out).
		Where("ancestor_id = ?", ancestorID).
		Where("descendant_id = ?", descendantID).
		Where("relation_type = ?", relationType).
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &out, err
}

func (r *AssetLineageRepo) Create(ctx context.Context, ancestorID string, descendantID string, relationType string) (*models.AssetLineage, error) {
	now := time.Now().Unix()
	item := models.AssetLineage{
		ID:           utils.NewID(),
		AncestorID:   ancestorID,
		DescendantID: descendantID,
		RelationType: relationType,
		CreatedAt:    now,
	}
	_, err := r.db.NewInsert().
		Model(&item).
		On("CONFLICT (ancestor_id, descendant_id, relation_type) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	existing, err := r.GetByPair(ctx, ancestorID, descendantID, relationType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return &item, nil
}

func (r *AssetLineageRepo) Update(ctx context.Context, id string, ancestorID string, descendantID string, relationType string) error {
	_, err := r.db.NewUpdate().
		Model((*models.AssetLineage)(nil)).
		Set("ancestor_id = ?", ancestorID).
		Set("descendant_id = ?", descendantID).
		Set("relation_type = ?", relationType).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetLineageRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().
		Model((*models.AssetLineage)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *AssetLineageRepo) DeleteByPair(ctx context.Context, ancestorID string, descendantID string, relationType string) error {
	_, err := r.db.NewDelete().
		Model((*models.AssetLineage)(nil)).
		Where("ancestor_id = ?", ancestorID).
		Where("descendant_id = ?", descendantID).
		Where("relation_type = ?", relationType).
		Exec(ctx)
	return err
}

func (r *AssetLineageRepo) ListByAncestor(ctx context.Context, ancestorID string) ([]models.AssetLineage, error) {
	var out []models.AssetLineage
	err := r.db.NewSelect().
		Model(&out).
		Where("ancestor_id = ?", ancestorID).
		OrderExpr("created_at DESC").
		Scan(ctx)
	return out, err
}

func (r *AssetLineageRepo) ListByDescendant(ctx context.Context, descendantID string) ([]models.AssetLineage, error) {
	var out []models.AssetLineage
	err := r.db.NewSelect().
		Model(&out).
		Where("descendant_id = ?", descendantID).
		OrderExpr("created_at DESC").
		Scan(ctx)
	return out, err
}

func (r *AssetLineageRepo) ListByAsset(ctx context.Context, assetID string) ([]models.AssetLineage, error) {
	var out []models.AssetLineage
	err := r.db.NewSelect().
		Model(&out).
		Where("ancestor_id = ? OR descendant_id = ?", assetID, assetID).
		OrderExpr("created_at DESC").
		Scan(ctx)
	return out, err
}

func (r *AssetLineageRepo) ExistsBetween(ctx context.Context, assetA string, assetB string) (bool, error) {
	count, err := r.db.NewSelect().
		Model((*models.AssetLineage)(nil)).
		Where("(ancestor_id = ? AND descendant_id = ?) OR (ancestor_id = ? AND descendant_id = ?)", assetA, assetB, assetB, assetA).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
