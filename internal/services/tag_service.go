package services

import (
	"context"
	"errors"
	"strings"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/repos"
)

type TagService struct {
	tagRepo *repos.TagRepo
}

func NewTagService(tagRepo *repos.TagRepo) *TagService {
	return &TagService{
		tagRepo: tagRepo,
	}
}

// TagWithCount 带文件计数的标签
type TagWithCount struct {
	models.Tag
	FileCount int `json:"file_count"`
}

type TagTreeNode struct {
	TagWithCount
	Children []*TagTreeNode `json:"children"`
}

// CreateTag 创建标签
func (s *TagService) CreateTag(ctx context.Context, name string, color, icon, parentID *string) (*models.Tag, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("tag name cannot be empty")
	}

	// 检查是否已存在
	existing, _ := s.tagRepo.GetByName(ctx, name)
	if existing != nil {
		return nil, errors.New("tag with this name already exists")
	}

	normalizedParentID, err := s.normalizeParentID(parentID)
	if err != nil {
		return nil, err
	}

	if normalizedParentID != nil {
		if _, err := s.tagRepo.GetByID(ctx, *normalizedParentID); err != nil {
			return nil, errors.New("parent tag does not exist")
		}
	}

	return s.tagRepo.Create(ctx, name, color, icon, normalizedParentID)
}

// UpdateTag 更新标签
func (s *TagService) UpdateTag(ctx context.Context, id string, name *string, color, icon, parentID *string) (*models.Tag, error) {
	if id == "" {
		return nil, errors.New("tag id is required")
	}

	var normalizedName *string
	if name != nil {
		n := strings.TrimSpace(*name)
		if n == "" {
			return nil, errors.New("tag name cannot be empty")
		}
		normalizedName = &n
	}

	normalizedParentID, err := s.normalizeParentID(parentID)
	if err != nil {
		return nil, err
	}

	if normalizedParentID != nil {
		if *normalizedParentID == id {
			return nil, errors.New("tag cannot be parent of itself")
		}
		if _, err := s.tagRepo.GetByID(ctx, *normalizedParentID); err != nil {
			return nil, errors.New("parent tag does not exist")
		}
		if err := s.ensureNoCycle(ctx, id, *normalizedParentID); err != nil {
			return nil, err
		}
	}

	return s.tagRepo.Update(ctx, id, normalizedName, color, icon, normalizedParentID)
}

// DeleteTag 删除标签
func (s *TagService) DeleteTag(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("tag id is required")
	}

	return s.tagRepo.Delete(ctx, id)
}

// GetTag 获取标签
func (s *TagService) GetTag(ctx context.Context, id string) (*models.Tag, error) {
	return s.tagRepo.GetByID(ctx, id)
}

// ListTags 获取所有标签（带文件计数）
func (s *TagService) ListTags(ctx context.Context) ([]TagWithCount, error) {
	tags, err := s.tagRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]TagWithCount, 0, len(tags))
	for _, tag := range tags {
		count, _ := s.tagRepo.GetFileCount(ctx, tag.ID)
		result = append(result, TagWithCount{
			Tag:       tag,
			FileCount: count,
		})
	}

	return result, nil
}

// ListTagTree 获取树形标签（支持多级）
func (s *TagService) ListTagTree(ctx context.Context) ([]*TagTreeNode, error) {
	flat, err := s.ListTags(ctx)
	if err != nil {
		return nil, err
	}

	nodes := make(map[string]*TagTreeNode, len(flat))
	for _, t := range flat {
		copyTag := t
		nodes[t.ID] = &TagTreeNode{
			TagWithCount: copyTag,
			Children:     make([]*TagTreeNode, 0),
		}
	}

	roots := make([]*TagTreeNode, 0)
	for _, t := range flat {
		n := nodes[t.ID]
		if t.ParentID == nil || *t.ParentID == "" {
			roots = append(roots, n)
			continue
		}
		parent, ok := nodes[*t.ParentID]
		if !ok {
			// 父节点缺失时回退为根节点，避免整棵树中断
			roots = append(roots, n)
			continue
		}
		parent.Children = append(parent.Children, n)
	}

	return roots, nil
}

// SearchTags 搜索标签
func (s *TagService) SearchTags(ctx context.Context, query string, limit int) ([]TagWithCount, error) {
	if limit <= 0 {
		limit = 20
	}

	tags, err := s.tagRepo.Search(ctx, query, limit)
	if err != nil {
		return nil, err
	}

	result := make([]TagWithCount, 0, len(tags))
	for _, tag := range tags {
		count, _ := s.tagRepo.GetFileCount(ctx, tag.ID)
		result = append(result, TagWithCount{
			Tag:       tag,
			FileCount: count,
		})
	}

	return result, nil
}

// AddTagsToFiles 批量为文件添加标签
func (s *TagService) AddTagsToFiles(ctx context.Context, fileIDs []string, tagIDs []string) error {
	if len(fileIDs) == 0 || len(tagIDs) == 0 {
		return errors.New("file_ids and tag_ids are required")
	}

	for _, fileID := range fileIDs {
		for _, tagID := range tagIDs {
			if err := s.tagRepo.AddTagToAsset(ctx, fileID, tagID); err != nil {
				return err
			}
		}
	}

	return nil
}

// RemoveTagsFromFiles 批量从文件移除标签
func (s *TagService) RemoveTagsFromFiles(ctx context.Context, fileIDs []string, tagIDs []string) error {
	if len(fileIDs) == 0 || len(tagIDs) == 0 {
		return errors.New("file_ids and tag_ids are required")
	}

	for _, fileID := range fileIDs {
		for _, tagID := range tagIDs {
			if err := s.tagRepo.RemoveTagFromAsset(ctx, fileID, tagID); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetFileTags 获取文件的所有标签
func (s *TagService) GetFileTags(ctx context.Context, fileID string) ([]models.Tag, error) {
	if fileID == "" {
		return nil, errors.New("file_id is required")
	}

	return s.tagRepo.GetAssetTags(ctx, fileID)
}

// GetTagFiles 获取标签下的所有文件
func (s *TagService) GetTagFiles(ctx context.Context, tagID string) ([]string, error) {
	if tagID == "" {
		return nil, errors.New("tag_id is required")
	}

	return s.tagRepo.GetTagAssets(ctx, tagID)
}

func (s *TagService) normalizeParentID(parentID *string) (*string, error) {
	if parentID == nil {
		return nil, nil
	}
	v := strings.TrimSpace(*parentID)
	if v == "" {
		return nil, nil
	}
	return &v, nil
}

func (s *TagService) ensureNoCycle(ctx context.Context, tagID, parentID string) error {
	seen := map[string]struct{}{tagID: {}}
	cur := parentID
	for cur != "" {
		if _, exists := seen[cur]; exists {
			return errors.New("tag parent relationship would create a cycle")
		}
		seen[cur] = struct{}{}

		parent, err := s.tagRepo.GetByID(ctx, cur)
		if err != nil {
			return errors.New("invalid parent tag chain")
		}
		if parent.ParentID == nil || *parent.ParentID == "" {
			break
		}
		cur = *parent.ParentID
	}
	return nil
}
