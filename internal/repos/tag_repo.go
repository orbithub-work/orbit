package repos

import (
	"context"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type TagRepo struct {
	db *bun.DB
}

func NewTagRepo(db *bun.DB) *TagRepo {
	return &TagRepo{db: db}
}

// Create 创建标签
func (r *TagRepo) Create(ctx context.Context, name string, color, icon, parentID *string) (*models.Tag, error) {
	now := time.Now().Unix()
	tag := &models.Tag{
		ID:        utils.NewID(),
		Name:      name,
		Color:     color,
		Icon:      icon,
		ParentID:  parentID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	_, err := r.db.NewInsert().Model(tag).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// Update 更新标签
func (r *TagRepo) Update(ctx context.Context, id string, name *string, color, icon, parentID *string) (*models.Tag, error) {
	tag, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if name != nil {
		tag.Name = *name
	}
	if color != nil {
		tag.Color = color
	}
	if icon != nil {
		tag.Icon = icon
	}
	if parentID != nil {
		tag.ParentID = parentID
	}

	tag.UpdatedAt = time.Now().Unix()

	_, err = r.db.NewUpdate().Model(tag).WherePK().Exec(ctx)
	if err != nil {
		return nil, err
	}

	return tag, nil
}

// Delete 删除标签
func (r *TagRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.NewDelete().Model((*models.Tag)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

// GetByID 根据ID获取标签
func (r *TagRepo) GetByID(ctx context.Context, id string) (*models.Tag, error) {
	tag := new(models.Tag)
	err := r.db.NewSelect().Model(tag).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// GetByName 根据名称获取标签
func (r *TagRepo) GetByName(ctx context.Context, name string) (*models.Tag, error) {
	tag := new(models.Tag)
	err := r.db.NewSelect().Model(tag).Where("name = ?", name).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tag, nil
}

// List 获取所有标签
func (r *TagRepo) List(ctx context.Context) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.NewSelect().Model(&tags).Order("created_at DESC").Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// Search 搜索标签
func (r *TagRepo) Search(ctx context.Context, query string, limit int) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.NewSelect().
		Model(&tags).
		Where("name LIKE ?", "%"+query+"%").
		Order("created_at DESC").
		Limit(limit).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetFileCount 获取标签关联的文件数量
func (r *TagRepo) GetFileCount(ctx context.Context, tagID string) (int, error) {
	count, err := r.db.NewSelect().
		Model((*models.AssetTag)(nil)).
		Where("tag_id = ?", tagID).
		Count(ctx)
	return count, err
}

// AddTagToAsset 为资产添加标签
func (r *TagRepo) AddTagToAsset(ctx context.Context, assetID, tagID string) error {
	assetTag := &models.AssetTag{
		AssetID:   assetID,
		TagID:     tagID,
		CreatedAt: time.Now().Unix(),
	}

	_, err := r.db.NewInsert().
		Model(assetTag).
		On("CONFLICT DO NOTHING").
		Exec(ctx)
	return err
}

// RemoveTagFromAsset 从资产移除标签
func (r *TagRepo) RemoveTagFromAsset(ctx context.Context, assetID, tagID string) error {
	_, err := r.db.NewDelete().
		Model((*models.AssetTag)(nil)).
		Where("asset_id = ? AND tag_id = ?", assetID, tagID).
		Exec(ctx)
	return err
}

// GetAssetTags 获取资产的所有标签
func (r *TagRepo) GetAssetTags(ctx context.Context, assetID string) ([]models.Tag, error) {
	var tags []models.Tag
	err := r.db.NewSelect().
		Model(&tags).
		Join("INNER JOIN asset_tags ON asset_tags.tag_id = tag.id").
		Where("asset_tags.asset_id = ?", assetID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagAssets 获取标签下的所有资产ID
func (r *TagRepo) GetTagAssets(ctx context.Context, tagID string) ([]string, error) {
	var assetIDs []string
	err := r.db.NewSelect().
		Model((*models.AssetTag)(nil)).
		Column("asset_id").
		Where("tag_id = ?", tagID).
		Scan(ctx, &assetIDs)
	return assetIDs, err
}
