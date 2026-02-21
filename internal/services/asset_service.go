package services

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/repos"
	"media-assistant-os/internal/utils"
)

type AssetService struct {
	assets            *repos.AssetRepo
	historyEvents     *repos.AssetHistoryEventRepo
	searchHistoryRepo *repos.SearchHistoryRepo
	projectAssets     *repos.ProjectAssetRepo
	projects          *repos.ProjectRepo
	lineage           *repos.AssetLineageRepo
	lineageCandidates *repos.LineageCandidateRepo
	activities        *ActivityService
	eventHub          *EventHub
	taskService       *TaskService
	cache             *AssetCache
	bloom             *utils.BloomFilter
}

// NewAssetService 创建资产服务实例
func NewAssetService(assets *repos.AssetRepo, historyEvents *repos.AssetHistoryEventRepo, searchHistoryRepo *repos.SearchHistoryRepo, projectAssets *repos.ProjectAssetRepo, projects *repos.ProjectRepo, lineage *repos.AssetLineageRepo, lineageCandidates *repos.LineageCandidateRepo, activities *ActivityService, eventHub *EventHub, taskService *TaskService) *AssetService {
	// 预估 100 万文件，误判率 1%
	// 内存占用: ~1.2MB
	bf := utils.NewBloomFilter(1000000, 0.01)

	return &AssetService{
		assets:            assets,
		historyEvents:     historyEvents,
		searchHistoryRepo: searchHistoryRepo,
		projectAssets:     projectAssets,
		projects:          projects,
		lineage:           lineage,
		lineageCandidates: lineageCandidates,
		activities:        activities,
		eventHub:          eventHub,
		taskService:       taskService,
		cache:             NewAssetCache(200000),
		bloom:             bf,
	}
}

// InitBloomFilter 初始化布隆过滤器（需在服务启动时调用）
func (s *AssetService) InitBloomFilter(ctx context.Context) error {
	paths, err := s.assets.GetAllPaths(ctx)
	if err != nil {
		return err
	}
	s.bloom.Reset()
	for _, p := range paths {
		s.bloom.AddString(p)
	}
	return nil
}

// FileEntry represents a file in a directory listing with asset status
type FileEntry struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	IsDirectory bool      `json:"is_directory"`
	Size        int64     `json:"size"`
	ModifiedAt  time.Time `json:"modified_at"`
	Status      string    `json:"status,omitempty"`   // PENDING, READY, INDEXED, etc.
	AssetID     string    `json:"asset_id,omitempty"` // ID of the asset if indexed
	Width       int       `json:"width,omitempty"`    // Image/video width from media_meta
	Height      int       `json:"height,omitempty"`   // Image/video height from media_meta
}

type ListAssetsRequest struct {
	ProjectID string
	Directory string
	Query     string
	TagIDs    []string
	Types     []string
	Shapes    []string
	// Quick filter hints for user-facing smart search (without reordering).
	QuickFilter string
	DatePreset  string

	SizeMin int64
	SizeMax int64

	RatingMin int
	RatingMax int

	MtimeFrom int64
	MtimeTo   int64

	WidthMin  int
	WidthMax  int
	HeightMin int
	HeightMax int

	SortBy    string
	SortOrder string

	Limit  int
	Cursor string
}

type AssetListItem struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Path            string `json:"path"`
	Size            int64  `json:"size"`
	Mtime           int64  `json:"mtime"`
	ModifiedAt      int64  `json:"modified_at"`
	FileType        string `json:"file_type,omitempty"`
	Status          string `json:"status,omitempty"`
	Shape           string `json:"shape"`
	SuggestedRating *int   `json:"suggested_rating,omitempty"`
	UserRating      *int   `json:"user_rating,omitempty"`
	ThumbnailPath   string `json:"thumbnail_path,omitempty"`
	Width           int    `json:"width,omitempty"`
	Height          int    `json:"height,omitempty"`
	CreatedAt       int64  `json:"created_at"`
}

type ListAssetsResult struct {
	Items      []AssetListItem `json:"items"`
	NextCursor *string         `json:"nextCursor"`
	HasMore    bool            `json:"hasMore"`
	Total      int             `json:"total"`
}

// ListDirectory lists files in a directory and enriches them with asset status
func (s *AssetService) ListDirectory(ctx context.Context, path string) ([]FileEntry, error) {
	if path == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		path = homeDir
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(abs)
	if err != nil {
		return nil, err
	}

	var filePaths []string
	var fileEntries []FileEntry

	// First pass: collect files and basic info
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		fullPath := filepath.Join(abs, info.Name())
		item := FileEntry{
			ID:          fullPath, // Use full path as ID for non-indexed files to ensure uniqueness
			Name:        info.Name(),
			Path:        fullPath,
			IsDirectory: entry.IsDir(),
			Size:        info.Size(),
			ModifiedAt:  info.ModTime(),
		}
		if !entry.IsDir() {
			filePaths = append(filePaths, fullPath)
		}
		fileEntries = append(fileEntries, item)
	}

	// Batch query assets
	assets, err := s.assets.GetByPaths(ctx, filePaths)
	if err != nil {
		return nil, err
	}

	// Create a map for quick lookup
	assetMap := make(map[string]models.Asset)
	for _, a := range assets {
		assetMap[a.Path] = a
	}

	// Second pass: enrich with asset info
	for i := range fileEntries {
		if !fileEntries[i].IsDirectory {
			if asset, ok := assetMap[fileEntries[i].Path]; ok {
				fileEntries[i].Status = asset.Status
				fileEntries[i].AssetID = asset.ID

				// Parse media_meta for width/height
				if asset.MediaMeta != "" {
					var meta struct {
						Width  int `json:"width"`
						Height int `json:"height"`
					}
					if err := json.Unmarshal([]byte(asset.MediaMeta), &meta); err == nil {
						fileEntries[i].Width = meta.Width
						fileEntries[i].Height = meta.Height
					}
				}
			}
		}
	}

	return fileEntries, nil
}

// GetCache 返回资产缓存（用于日志查看器）
func (s *AssetService) GetCache() *AssetCache {
	return s.cache
}

// ArchiveFiles 将文件逻辑归档到项目中（非侵入，不复制/移动源文件）
func (s *AssetService) ArchiveFiles(ctx context.Context, projectID string, paths []string) error {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return errors.New("project id is required")
	}

	// 1. Ensure project exists
	project, err := s.projects.Get(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return errors.New("project not found")
	}

	// 2. Archive each file by indexing and binding to project (no file IO mutation)
	var failedFiles []string
	for _, p := range paths {
		path := strings.TrimSpace(p)
		if path == "" {
			continue
		}
		if _, err := s.IndexFile(ctx, IndexFileRequest{
			Path:      path,
			ProjectID: projectID,
			Trigger:   "api",
		}); err != nil {
			failedFiles = append(failedFiles, path+": "+err.Error())
		}
	}

	if len(failedFiles) > 0 {
		return errors.New("some files failed to archive: " + strings.Join(failedFiles, "; "))
	}
	return nil
}

// IndexFile 索引文件到系统，如果文件已存在则返回现有资产ID
func (s *AssetService) IndexFile(ctx context.Context, req IndexFileRequest) (*IndexFileResult, error) {
	if req.Path == "" {
		return nil, errors.New("path is required")
	}

	abs, err := filepath.Abs(req.Path)
	if err != nil {
		return nil, err
	}
	trigger := normalizeIndexTrigger(req.Trigger)

	// 1. Check cache first
	if cached, ok := s.cache.GetByPath(abs); ok {
		if req.ProjectID != "" {
			_ = s.projectAssets.Link(ctx, req.ProjectID, cached.ID)
		}
		return &IndexFileResult{AssetID: cached.ID}, nil
	}

	// 2. Check Bloom Filter
	// 如果布隆过滤器说 "不在"，那大概率是不在。
	// 但要注意：Bloom Filter 只覆盖 "Path"。如果是 "文件改名" (Relink)，路径是新的，BF 也会说不在。
	// 所以 BF 只能帮我们跳过 s.assets.GetByPath 的查询，不能跳过 Relink 的逻辑。
	var existing *models.Asset

	if s.bloom.ContainsString(abs) {
		// 只有 Bloom 说 "可能在" 时，才去查 DB
		existing, err = s.assets.GetByPath(ctx, abs)
		if err != nil {
			return nil, err
		}
	} else {
		// Bloom 说 "不在"，那就肯定不在，跳过 DB 查询，existing = nil
	}

	// 3. Get current metadata from disk
	info, err := os.Stat(abs)
	if err != nil {
		return nil, err
	}
	currentSize := info.Size()
	currentMtime := info.ModTime().Unix()

	if existing != nil {
		// 如果资产处于被忽略（黑名单）状态，直接跳过，不建立任何关联
		if existing.Status == "IGNORED" {
			return &IndexFileResult{AssetID: existing.ID}, nil
		}

		// 检查元数据是否变更（数据一致性校验）
		hasChanged := existing.Size != currentSize || existing.Mtime != currentMtime

		if hasChanged {
			// 文件内容变更，重置状态为 PENDING，触发重新处理
			_ = s.assets.UpdateMetadata(ctx, existing.ID, currentSize, currentMtime, "PENDING", utils.FormatNow()+": File change detected, re-indexing")

			// 记录变更日志
			if s.activities != nil {
				s.activities.LogEx(ctx, "INFO", "File modified externally", existing.ID, req.ProjectID)
			}
			s.recordHistoryEvent(ctx, repos.CreateAssetHistoryEventInput{
				AssetID:    existing.ID,
				ProjectID:  req.ProjectID,
				EventType:  models.AssetHistoryEventModified,
				SourcePath: abs,
				TargetPath: abs,
				Confidence: historyConfidence(trigger, false),
				IsInferred: trigger != "watcher",
				Detail:     "file content or metadata changed",
			}, 8)

			existing.Status = "PENDING"
			existing.Size = currentSize
			existing.Mtime = currentMtime
		} else if existing.Status == "MISSING" {
			// 如果资产之前是 MISSING 状态，现在找回来了，自动恢复
			_ = s.assets.UpdateStatusWithLog(ctx, existing.ID, "READY", utils.FormatNow()+": File rediscovered during scan")
			existing.Status = "READY"
		}

		s.cache.Put(existing)
		if req.ProjectID != "" {
			_ = s.projectAssets.Link(ctx, req.ProjectID, existing.ID)
		}

		// 如果发生了变更，需要重新创建任务
		if hasChanged && s.taskService != nil {
			_ = s.taskService.CreateInitialTasks(ctx, existing.ID)
		}

		return &IndexFileResult{AssetID: existing.ID}, nil
	}

	// 4. Try to Relink or Detect Duplicates
	// 我们需要指纹来辅助判断。为了性能，只有在可能有匹配（有同大小文件）时才计算
	candidates, _ := s.assets.FindMissingAssetsBySize(ctx, currentSize)

	// 无论有没有 MISSING，我们都可能需要检测副本（血缘）
	// 这里我们优化一下：如果文件比较大，且我们需要做重连或副本检测
	fp := ""
	fpComputed := false

	// 如果有“失踪人口”，尝试重连
	if len(candidates) > 0 {
		fp, _, _, _, err = ComputeAdaptiveFingerprint(abs)
		if err == nil {
			fpComputed = true
			for _, candidate := range candidates {
				if candidate.Fingerprint != nil && *candidate.Fingerprint == fp {
					oldPath := candidate.Path
					_ = s.assets.RelinkAsset(ctx, candidate.ID, abs, currentMtime)
					s.bloom.AddString(abs) // 更新 Bloom Filter

					if s.activities != nil {
						s.activities.LogEx(ctx, "INFO", "Asset moved/renamed: "+filepath.Base(abs), candidate.ID, req.ProjectID)
					}
					eventType := models.AssetHistoryEventMoved
					if filepath.Dir(oldPath) == filepath.Dir(abs) {
						eventType = models.AssetHistoryEventRenamed
					}
					s.recordHistoryEvent(ctx, repos.CreateAssetHistoryEventInput{
						AssetID:    candidate.ID,
						ProjectID:  req.ProjectID,
						EventType:  eventType,
						SourcePath: oldPath,
						TargetPath: abs,
						Confidence: historyConfidence(trigger, true),
						IsInferred: true,
						Detail:     "path relinked by fingerprint match",
					}, 8)

					candidate.Path = abs
					candidate.Status = "READY"
					candidate.Mtime = currentMtime
					s.cache.Put(&candidate)

					if req.ProjectID != "" {
						_ = s.projectAssets.Link(ctx, req.ProjectID, candidate.ID)
					}
					return &IndexFileResult{AssetID: candidate.ID}, nil
				}
			}
		}
	}

	// 如果到这里还没返回，说明不是简单的移动。
	// 检查是否是现有资产的副本（血缘推断）
	if !fpComputed {
		// 如果还没计算指纹，现在计算一下（仅针对新文件）
		fp, _, _, _, err = ComputeAdaptiveFingerprint(abs)
		if err == nil {
			fpComputed = true
		}
	}

	if fpComputed {
		activeCopies, _ := s.assets.FindActiveAssetsByFingerprint(ctx, fp)
		copySourcePath := ""
		if len(activeCopies) > 0 {
			parent := activeCopies[0]
			copySourcePath = parent.Path
			if s.activities != nil {
				s.activities.LogEx(ctx, "INFO", "Found suspected derivative: "+filepath.Base(abs)+" (Copy of "+filepath.Base(parent.Path)+")", "", req.ProjectID)
			}
		}

		// 5. Insert as PENDING
		asset, err := s.assets.Create(ctx, abs, currentSize, currentMtime)
		if err != nil {
			return nil, err
		}
		s.bloom.AddString(abs) // 添加新文件到 Bloom Filter

		s.cache.Put(asset)
		s.recordHistoryEvent(ctx, repos.CreateAssetHistoryEventInput{
			AssetID:    asset.ID,
			ProjectID:  req.ProjectID,
			EventType:  models.AssetHistoryEventCreated,
			TargetPath: abs,
			Confidence: historyConfidence(trigger, false),
			IsInferred: trigger != "watcher",
			Detail:     "new file indexed",
		}, 8)
		if copySourcePath != "" {
			s.recordHistoryEvent(ctx, repos.CreateAssetHistoryEventInput{
				AssetID:    asset.ID,
				ProjectID:  req.ProjectID,
				EventType:  models.AssetHistoryEventCopied,
				SourcePath: copySourcePath,
				TargetPath: abs,
				Confidence: historyConfidence(trigger, true),
				IsInferred: true,
				Detail:     "copy inferred by fingerprint",
			}, 8)
		}

		// Broadcast pending event for UI to show placeholder
		if s.eventHub != nil {
			s.eventHub.Broadcast(map[string]any{
				"type": "asset_pending",
				"data": map[string]any{
					"id":   asset.ID,
					"path": asset.Path,
					"name": filepath.Base(asset.Path),
				},
			})
		}

		// 5. Link project if needed
		if req.ProjectID != "" {
			_ = s.projectAssets.Link(ctx, req.ProjectID, asset.ID)
		}

		// 6. Create processing tasks
		if s.taskService != nil {
			_ = s.taskService.CreateInitialTasks(ctx, asset.ID)
		}

		return &IndexFileResult{AssetID: asset.ID}, nil
	}

	// No fingerprint computed. Treat as plain create.
	asset, err := s.assets.Create(ctx, abs, currentSize, currentMtime)
	if err != nil {
		return nil, err
	}
	s.bloom.AddString(abs)
	s.cache.Put(asset)
	s.recordHistoryEvent(ctx, repos.CreateAssetHistoryEventInput{
		AssetID:    asset.ID,
		ProjectID:  req.ProjectID,
		EventType:  models.AssetHistoryEventCreated,
		TargetPath: abs,
		Confidence: historyConfidence(trigger, false),
		IsInferred: trigger != "watcher",
		Detail:     "new file indexed",
	}, 8)
	if s.eventHub != nil {
		s.eventHub.Broadcast(map[string]any{
			"type": "asset_pending",
			"data": map[string]any{
				"id":   asset.ID,
				"path": asset.Path,
				"name": filepath.Base(asset.Path),
			},
		})
	}
	if req.ProjectID != "" {
		_ = s.projectAssets.Link(ctx, req.ProjectID, asset.ID)
	}
	if s.taskService != nil {
		_ = s.taskService.CreateInitialTasks(ctx, asset.ID)
	}
	return &IndexFileResult{AssetID: asset.ID}, nil
}

func (s *AssetService) ListAssetHistory(ctx context.Context, assetID string, limit int) ([]models.AssetHistoryEvent, error) {
	if s.historyEvents == nil {
		return []models.AssetHistoryEvent{}, nil
	}
	return s.historyEvents.ListByAsset(ctx, assetID, limit)
}

func (s *AssetService) recordHistoryEvent(ctx context.Context, in repos.CreateAssetHistoryEventInput, dedupWindowSec int64) {
	if s.historyEvents == nil {
		return
	}
	_, _ = s.historyEvents.CreateDedup(ctx, in, dedupWindowSec)
}

func historyConfidence(trigger string, inferred bool) string {
	if inferred {
		if trigger == "watcher" {
			return "medium"
		}
		return "low"
	}
	if trigger == "watcher" {
		return "high"
	}
	if trigger == "api" {
		return "medium"
	}
	return "low"
}

func normalizeIndexTrigger(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "watcher":
		return "watcher"
	case "api":
		return "api"
	default:
		return "scan"
	}
}

// UpdateMediaMeta 更新资产的媒体元数据
func (s *AssetService) UpdateMediaMeta(ctx context.Context, id string, mediaMeta string) error {
	return s.assets.UpdateMediaMeta(ctx, id, mediaMeta)
}

func (s *AssetService) ApplyDerivedMetadata(ctx context.Context, id string, mediaMeta string, shape string, suggestedRating *int) error {
	shape = normalizeShape(shape)
	return s.assets.UpdateDerivedMetadata(ctx, id, mediaMeta, shape, suggestedRating)
}

// GetAsset 根据ID获取资产（委托给GetAssetCached）
func (s *AssetService) GetAsset(ctx context.Context, id string) (*models.Asset, error) {
	return s.GetAssetCached(ctx, id)
}

// GetAssetCached 根据ID获取资产，优先使用缓存
func (s *AssetService) GetAssetCached(ctx context.Context, id string) (*models.Asset, error) {
	if cached, ok := s.cache.GetByID(id); ok {
		return cached, nil
	}
	asset, err := s.assets.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if asset != nil {
		s.cache.Put(asset)
	}
	return asset, nil
}

// FindByFingerprintCached 根据指纹查询资产并缓存结果
func (s *AssetService) FindByFingerprintCached(ctx context.Context, fp string) ([]models.Asset, error) {
	results, err := s.assets.FindByFingerprint(ctx, fp)
	if err != nil {
		return nil, err
	}
	for i := range results {
		s.cache.Put(&results[i])
	}
	return results, nil
}

// LogicRemoveAsset 将资产标记为已移除（入黑名单），但不物理删除文件
func (s *AssetService) LogicRemoveAsset(ctx context.Context, id string, reason string) error {
	asset, err := s.assets.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	// 1. 更新状态为 IGNORED 并记录原因
	if reason == "" {
		reason = "User requested removal"
	}
	opLog := utils.FormatNow() + ": " + reason
	if err := s.assets.UpdateStatusWithLog(ctx, id, "IGNORED", opLog); err != nil {
		return err
	}

	// 2. 解除所有项目绑定
	if err := s.projectAssets.UnlinkAll(ctx, id); err != nil {
		return err
	}

	// 3. 清理缓存
	s.cache.Invalidate(id)

	// 4. 记录用户日志
	if s.activities != nil {
		msg := "资产已移除并加入黑名单: " + filepath.Base(asset.Path)
		s.activities.LogEx(ctx, "INFO", msg, id, "")
	}

	// 5. 广播事件
	if s.eventHub != nil {
		s.eventHub.Broadcast(map[string]any{
			"type": "asset_ignored",
			"data": map[string]any{
				"id":   id,
				"path": asset.Path,
			},
		})
	}

	return nil
}

// RemoveFromProject 将资产从特定项目中移除，但不进入全局黑名单
func (s *AssetService) RemoveFromProject(ctx context.Context, projectID string, assetID string) error {
	asset, err := s.assets.GetByID(ctx, assetID)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	// 1. 解除绑定
	if err := s.projectAssets.Unlink(ctx, projectID, assetID); err != nil {
		return err
	}

	// 2. 记录用户日志
	if s.activities != nil {
		msg := "资产已从项目中移除: " + filepath.Base(asset.Path)
		s.activities.LogEx(ctx, "INFO", msg, assetID, projectID)
	}

	return nil
}

// RestoreAsset 将资产从黑名单（IGNORED）中恢复
func (s *AssetService) RestoreAsset(ctx context.Context, id string) error {
	asset, err := s.assets.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if asset == nil {
		return errors.New("asset not found")
	}

	if asset.Status != "IGNORED" {
		return nil // 只有 IGNORED 状态需要恢复
	}

	// 1. 检查物理文件是否存在，决定恢复到什么状态
	newStatus := "READY"
	if _, err := os.Stat(asset.Path); os.IsNotExist(err) {
		newStatus = "MISSING"
	}

	// 2. 更新状态并清除操作日志
	opLog := utils.FormatNow() + ": Restored from blacklist"
	if err := s.assets.UpdateStatusWithLog(ctx, id, newStatus, opLog); err != nil {
		return err
	}

	// 3. 记录用户日志
	if s.activities != nil {
		msg := "资产已从黑名单恢复: " + filepath.Base(asset.Path)
		s.activities.LogEx(ctx, "SUCCESS", msg, id, "")
	}

	// 4. 广播事件
	if s.eventHub != nil {
		s.eventHub.Broadcast(map[string]any{
			"type": "asset_restored",
			"data": map[string]any{
				"id":     id,
				"status": newStatus,
			},
		})
	}

	return nil
}

// DeleteAsset 删除资产（实际使用逻辑删除）
func (s *AssetService) DeleteAsset(ctx context.Context, id string) error {
	// 默认使用逻辑删除，除非特殊情况
	return s.LogicRemoveAsset(ctx, id, "Physical delete replaced by logic remove")
}

// BatchDeleteAssets 批量删除资产
func (s *AssetService) BatchDeleteAssets(ctx context.Context, ids []string) error {
	for _, id := range ids {
		if err := s.LogicRemoveAsset(ctx, id, "Batch delete"); err != nil {
			// 继续处理其他文件，不中断
			continue
		}
	}
	return nil
}

// GetProjectFiles 获取项目的所有文件资产
func (s *AssetService) GetProjectFiles(ctx context.Context, projectID string) ([]models.Asset, error) {
	if projectID == "" {
		return nil, errors.New("projectID is required")
	}
	return s.projectAssets.GetAssetsByProject(ctx, projectID)
}

func (s *AssetService) ListAssets(ctx context.Context, req ListAssetsRequest) (*ListAssetsResult, error) {
	req = preprocessListAssetFilters(req)

	if strings.TrimSpace(req.Query) != "" && s.searchHistoryRepo != nil {
		query := strings.TrimSpace(req.Query)
		filters := map[string]any{
			"quick_filter": req.QuickFilter,
			"date_preset":  req.DatePreset,
			"project_id":   req.ProjectID,
			"directory":    req.Directory,
		}
		go func() {
			_ = s.searchHistoryRepo.RecordSearch(context.Background(), query, filters)
		}()
	}

	offset := 0
	if strings.TrimSpace(req.Cursor) != "" {
		v, err := strconv.Atoi(strings.TrimSpace(req.Cursor))
		if err != nil || v < 0 {
			return nil, errors.New("invalid cursor")
		}
		offset = v
	}

	assets, total, err := s.assets.ListByQuery(ctx, repos.AssetListQuery{
		ProjectID: req.ProjectID,
		Directory: req.Directory,
		Query:     req.Query,
		TagIDs:    req.TagIDs,
		Types:     req.Types,
		Shapes:    req.Shapes,
		SizeMin:   req.SizeMin,
		SizeMax:   req.SizeMax,
		RatingMin: req.RatingMin,
		RatingMax: req.RatingMax,
		MtimeFrom: req.MtimeFrom,
		MtimeTo:   req.MtimeTo,
		WidthMin:  req.WidthMin,
		WidthMax:  req.WidthMax,
		HeightMin: req.HeightMin,
		HeightMax: req.HeightMax,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
		Limit:     req.Limit,
		Offset:    offset,
	})
	if err != nil {
		return nil, err
	}

	items := make([]AssetListItem, 0, len(assets))
	for _, a := range assets {
		name := filepath.Base(a.Path)
		if name == "." || name == "/" || name == "" {
			name = a.Path
		}

		meta := struct {
			ThumbnailPath string `json:"thumbnail_path"`
			Width         int    `json:"width"`
			Height        int    `json:"height"`
		}{}
		if strings.TrimSpace(a.MediaMeta) != "" {
			_ = json.Unmarshal([]byte(a.MediaMeta), &meta)
		}

		items = append(items, AssetListItem{
			ID:              a.ID,
			Name:            name,
			Path:            a.Path,
			Size:            a.Size,
			Mtime:           a.Mtime,
			ModifiedAt:      a.Mtime,
			FileType:        detectAssetFileType(name),
			Status:          a.Status,
			Shape:           normalizeShape(a.Shape),
			SuggestedRating: a.SuggestedRating,
			UserRating:      a.UserRating,
			ThumbnailPath:   meta.ThumbnailPath,
			Width:           meta.Width,
			Height:          meta.Height,
			CreatedAt:       a.CreatedAt,
		})
	}

	hasMore := offset+len(items) < total
	var nextCursor *string
	if hasMore {
		next := strconv.Itoa(offset + len(items))
		nextCursor = &next
	}

	return &ListAssetsResult{
		Items:      items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Total:      total,
	}, nil
}

func (s *AssetService) GetSearchHistory(ctx context.Context, limit int) ([]models.SearchHistory, error) {
	if s.searchHistoryRepo == nil {
		return []models.SearchHistory{}, nil
	}
	return s.searchHistoryRepo.GetRecent(ctx, limit)
}

func (s *AssetService) ClearSearchHistory(ctx context.Context) error {
	if s.searchHistoryRepo == nil {
		return nil
	}
	return s.searchHistoryRepo.Clear(ctx)
}

func preprocessListAssetFilters(req ListAssetsRequest) ListAssetsRequest {
	req.QuickFilter = normalizeQuickFilter(req.QuickFilter)
	req.DatePreset = normalizeDatePreset(req.DatePreset)

	now := time.Now()
	if req.MtimeFrom == 0 && req.MtimeTo == 0 {
		switch req.DatePreset {
		case "today":
			start := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			req.MtimeFrom = start.Unix()
			req.MtimeTo = now.Unix()
		case "thisWeek":
			weekday := int(now.Weekday())
			if weekday == 0 {
				weekday = 7
			}
			start := now.AddDate(0, 0, -weekday+1)
			start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
			req.MtimeFrom = start.Unix()
			req.MtimeTo = now.Unix()
		case "thisMonth":
			start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			req.MtimeFrom = start.Unix()
			req.MtimeTo = now.Unix()
		case "lastMonth":
			lastMonth := now.AddDate(0, -1, 0)
			start := time.Date(lastMonth.Year(), lastMonth.Month(), 1, 0, 0, 0, 0, lastMonth.Location())
			end := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Add(-time.Second)
			req.MtimeFrom = start.Unix()
			req.MtimeTo = end.Unix()
		}
	}

	switch req.QuickFilter {
	case "recent":
		if req.MtimeFrom == 0 {
			req.MtimeFrom = now.AddDate(0, 0, -7).Unix()
		}
	case "unrated":
		// Preserve null-based semantic: user has not rated yet.
		req.RatingMin = 0
		req.RatingMax = -1
	case "large":
		if req.SizeMin < 100*1024*1024 {
			req.SizeMin = 100 * 1024 * 1024
		}
	case "vertical":
		req.Shapes = appendIfMissingFold(req.Shapes, "portrait")
	case "horizontal":
		req.Shapes = appendIfMissingFold(req.Shapes, "landscape")
	}
	return req
}

func normalizeQuickFilter(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "recent":
		return "recent"
	case "unrated":
		return "unrated"
	case "large":
		return "large"
	case "vertical":
		return "vertical"
	case "horizontal":
		return "horizontal"
	default:
		return ""
	}
}

func normalizeDatePreset(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "today":
		return "today"
	case "thisweek":
		return "thisWeek"
	case "thismonth":
		return "thisMonth"
	case "lastmonth":
		return "lastMonth"
	default:
		return ""
	}
}

func appendIfMissingFold(items []string, value string) []string {
	value = strings.TrimSpace(value)
	if value == "" {
		return items
	}
	for _, item := range items {
		if strings.EqualFold(strings.TrimSpace(item), value) {
			return items
		}
	}
	return append(items, value)
}

func (s *AssetService) SetUserRating(ctx context.Context, id string, userRating *int) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id is required")
	}
	if userRating != nil && (*userRating < 1 || *userRating > 5) {
		return errors.New("user_rating must be between 1 and 5")
	}
	return s.assets.UpdateUserRating(ctx, id, userRating)
}

// SetProjectAssetStatus 设置项目资产的状态
func (s *AssetService) SetProjectAssetStatus(ctx context.Context, projectID string, assetID string, status *string) error {
	if projectID == "" || assetID == "" {
		return errors.New("projectID and assetID are required")
	}
	return s.projectAssets.UpdateStatus(ctx, projectID, assetID, status)
}

func detectAssetFileType(name string) string {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	switch ext {
	case "jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "tif", "tiff", "psd", "psb", "exr":
		return "image"
	case "mp4", "avi", "mov", "mkv", "webm", "m4v", "flv", "wmv", "r3d", "braw":
		return "video"
	case "mp3", "wav", "flac", "aac", "ogg", "m4a", "aif", "aiff", "bwf":
		return "audio"
	case "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "md", "csv":
		return "document"
	default:
		return "file"
	}
}

func normalizeShape(shape string) string {
	v := strings.ToLower(strings.TrimSpace(shape))
	switch v {
	case "landscape", "portrait", "square", "panorama", "unknown":
		return v
	default:
		return "unknown"
	}
}

// CreateLineage 创建资产血缘关系
func (s *AssetService) CreateLineage(ctx context.Context, ancestorID string, descendantID string, relationType string) (*models.AssetLineage, error) {
	if ancestorID == "" || descendantID == "" {
		return nil, errors.New("ancestorID and descendantID are required")
	}
	if ancestorID == descendantID {
		return nil, errors.New("ancestorID and descendantID cannot be the same")
	}
	if relationType == "" {
		return nil, errors.New("relationType is required")
	}
	existing, err := s.lineage.GetByPair(ctx, ancestorID, descendantID, relationType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	return s.lineage.Create(ctx, ancestorID, descendantID, relationType)
}

// UpdateLineage 更新资产血缘关系
func (s *AssetService) UpdateLineage(ctx context.Context, id string, ancestorID string, descendantID string, relationType string) error {
	if id == "" {
		return errors.New("id is required")
	}
	if ancestorID == "" || descendantID == "" {
		return errors.New("ancestorID and descendantID are required")
	}
	if ancestorID == descendantID {
		return errors.New("ancestorID and descendantID cannot be the same")
	}
	if relationType == "" {
		return errors.New("relationType is required")
	}
	return s.lineage.Update(ctx, id, ancestorID, descendantID, relationType)
}

// DeleteLineage 删除资产血缘关系
func (s *AssetService) DeleteLineage(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}
	return s.lineage.Delete(ctx, id)
}

// DeleteLineageByPair 根据祖先和后代ID删除血缘关系
func (s *AssetService) DeleteLineageByPair(ctx context.Context, ancestorID string, descendantID string, relationType string) error {
	if ancestorID == "" || descendantID == "" {
		return errors.New("ancestorID and descendantID are required")
	}
	if relationType == "" {
		return errors.New("relationType is required")
	}
	return s.lineage.DeleteByPair(ctx, ancestorID, descendantID, relationType)
}

// ListLineage 列出资产的血缘关系
func (s *AssetService) ListLineage(ctx context.Context, assetID string) ([]models.AssetLineage, error) {
	if assetID == "" {
		return nil, errors.New("assetID is required")
	}
	return s.lineage.ListByAsset(ctx, assetID)
}
