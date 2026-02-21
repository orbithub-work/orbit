package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"media-assistant-os/internal/infra"
	"media-assistant-os/internal/models"
	"media-assistant-os/internal/pkg/logger"
	"media-assistant-os/internal/processor"

	"go.uber.org/zap"
)

type MediaQueue struct {
	assetService *AssetService
	taskService  *TaskService
	eventHub     *EventHub
	taskHandlers map[string]taskHandler
	taskTypes    []string
	taskTypesMu  sync.RWMutex
	stopChan     chan struct{}
	wg           sync.WaitGroup
	workerCount  int
	pollInterval time.Duration
}

type taskHandler func(ctx context.Context, asset *models.Asset, task *models.MediaTask) (map[string]any, error)

func NewMediaQueue(assetService *AssetService, taskService *TaskService, eventHub *EventHub, workerCount int) *MediaQueue {
	if workerCount <= 0 {
		workerCount = 2
	}
	q := &MediaQueue{
		assetService: assetService,
		taskService:  taskService,
		eventHub:     eventHub,
		taskHandlers: make(map[string]taskHandler),
		stopChan:     make(chan struct{}),
		workerCount:  workerCount,
		pollInterval: 2 * time.Second,
	}
	q.registerBuiltinHandlers()
	return q
}

func (q *MediaQueue) registerBuiltinHandlers() {
	q.RegisterTaskHandler("fingerprint", func(ctx context.Context, asset *models.Asset, task *models.MediaTask) (map[string]any, error) {
		return nil, q.processFingerprint(ctx, asset, task)
	})
	q.RegisterTaskHandler("metadata", q.extractMetadata)
	q.RegisterTaskHandler("thumbnail", func(ctx context.Context, asset *models.Asset, task *models.MediaTask) (map[string]any, error) {
		return nil, q.generateThumbnail(ctx, asset, task)
	})
}

func (q *MediaQueue) RegisterTaskHandler(taskType string, handler taskHandler) {
	taskType = strings.TrimSpace(taskType)
	if taskType == "" || handler == nil {
		return
	}
	q.taskTypesMu.Lock()
	defer q.taskTypesMu.Unlock()
	q.taskHandlers[taskType] = handler
	q.rebuildTaskTypesLocked()
}

func (q *MediaQueue) SupportedTaskTypes() []string {
	q.taskTypesMu.RLock()
	defer q.taskTypesMu.RUnlock()
	out := make([]string, len(q.taskTypes))
	copy(out, q.taskTypes)
	return out
}

func (q *MediaQueue) rebuildTaskTypesLocked() {
	types := make([]string, 0, len(q.taskHandlers))
	for taskType := range q.taskHandlers {
		types = append(types, taskType)
	}
	sort.Strings(types)
	q.taskTypes = types
}

func (q *MediaQueue) Start() {
	for i := 0; i < q.workerCount; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}
}

func (q *MediaQueue) Stop() {
	close(q.stopChan)
	q.wg.Wait()
}

func (q *MediaQueue) worker(id int) {
	defer q.wg.Done()
	logger.Info("Media worker started", zap.Int("worker_id", id))

	// Initial check
	if q.processNextTask(id) {
		// If we found a task immediately, keep going
	}

	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Try to process as many as possible
			for {
				if !q.processNextTask(id) {
					break
				}
				// Check for stop signal between tasks to be responsive
				select {
				case <-q.stopChan:
					return
				default:
				}
			}
		case <-q.stopChan:
			return
		}
	}
}

func (q *MediaQueue) processNextTask(workerID int) bool {
	ctx := context.Background()

	// 1. Try to claim a task
	taskTypes := q.SupportedTaskTypes()
	if len(taskTypes) == 0 {
		return false
	}
	task, err := q.taskService.ListPendingTasks(ctx, taskTypes)
	if err != nil || task == nil {
		return false
	}

	claimID := fmt.Sprintf("internal_worker_%d", workerID)
	success, err := q.taskService.ClaimTask(ctx, task.ID, claimID)
	if err != nil {
		// DB error, maybe backoff
		return false
	}
	if !success {
		// Task was claimed by someone else, but there are pending tasks.
		// Return true to try fetching the next one immediately.
		return true
	}

	logger.Debug("Task claimed", zap.String("task_id", task.ID), zap.String("worker", claimID))

	// 2. Fetch asset info
	asset, err := q.assetService.GetAssetCached(ctx, task.AssetID)
	if err != nil || asset == nil {
		_ = q.taskService.ReportTaskProgress(ctx, task.ID, claimID, false, "asset not found", nil)
		return true
	}

	// 3. Process based on type
	q.taskTypesMu.RLock()
	handler, ok := q.taskHandlers[task.TaskType]
	q.taskTypesMu.RUnlock()
	if !ok {
		_ = q.taskService.ReportTaskProgress(ctx, task.ID, claimID, false, "unsupported task type for internal worker: "+task.TaskType, nil)
		return true
	}
	resultData, processErr := handler(ctx, asset, task)

	// 4. Report results
	if processErr != nil {
		logger.Error("Task processing failed",
			zap.String("task_id", task.ID),
			zap.String("task_type", task.TaskType),
			zap.Error(processErr))
		_ = q.taskService.ReportTaskProgress(ctx, task.ID, claimID, false, processErr.Error(), nil)
	} else {
		logger.Info("Task completed successfully",
			zap.String("task_id", task.ID),
			zap.String("task_type", task.TaskType))
		_ = q.taskService.ReportTaskProgress(ctx, task.ID, claimID, true, "", resultData)
	}

	return true
}

func (q *MediaQueue) processFingerprint(ctx context.Context, asset *models.Asset, task *models.MediaTask) error {
	_ = q.taskService.ReportProgress(ctx, task.ID, 10)

	// 1. Compute fingerprint
	fp, _, _, _, err := ComputeAdaptiveFingerprint(asset.Path)
	if err != nil {
		return err
	}
	_ = q.taskService.ReportProgress(ctx, task.ID, 50)

	// 2. Update Asset in DB and Cache
	// Check for duplicates using cached lookup
	dupes, err := q.assetService.FindByFingerprintCached(ctx, fp)
	if err != nil {
		// Log error but continue
		logger.Error("Failed to check duplicates", zap.Error(err))
	}

	selected := chooseDuplicateParent(asset, dupes)
	var parentID *string
	if selected.Parent != nil {
		parentID = &selected.Parent.ID
		logger.Info("Duplicate detected",
			zap.String("asset_id", asset.ID),
			zap.String("parent_id", *parentID),
			zap.String("level", string(selected.Level)))
	}

	// Update Asset in DB
	// Note: We update directly via repo because we need to set parent_id
	if err := q.assetService.assets.UpdateFingerprint(ctx, asset.ID, fp, parentID); err != nil {
		return err
	}

	// Update cache
	q.assetService.cache.UpdateFingerprint(asset.ID, fp)

	// Create lineage if duplicate
	if parentID != nil {
		_, _ = q.assetService.CreateLineage(ctx, *parentID, asset.ID, lineageType(selected.Level))
	}

	// Build probabilistic lineage candidates for user confirmation.
	_ = q.assetService.BuildLineageCandidates(ctx, asset.ID)

	_ = q.taskService.ReportProgress(ctx, task.ID, 100)
	return nil
}

func (q *MediaQueue) extractMetadata(ctx context.Context, asset *models.Asset, task *models.MediaTask) (map[string]any, error) {
	_ = q.taskService.ReportProgress(ctx, task.ID, 10)
	ext := filepath.Ext(asset.Path)
	res, err := processor.GetManager().Process(ctx, asset.Path, ext)
	if err != nil {
		return nil, err
	}
	_ = q.taskService.ReportProgress(ctx, task.ID, 50)

	if res.Error != nil {
		return nil, res.Error
	}

	if res.Metadata != nil {
		logger.Info("Metadata extracted",
			zap.String("asset_id", asset.ID),
			zap.String("format", res.Metadata.Format))
		_ = q.taskService.ReportProgress(ctx, task.ID, 90)

		b, err := json.Marshal(res.Metadata)
		if err != nil {
			return nil, err
		}
		shape := deriveShape(res.Metadata.Width, res.Metadata.Height)
		suggestedRating := deriveSuggestedRating(res.Metadata)

		// Persist metadata early so lineage candidate scoring can use phash immediately.
		_ = q.assetService.ApplyDerivedMetadata(ctx, asset.ID, string(b), shape, &suggestedRating)
		asset.MediaMeta = string(b)
		_ = q.assetService.BuildLineageCandidates(ctx, asset.ID)
		return map[string]any{
			"metadata": string(b),
		}, nil
	}
	return nil, nil
}

func mergeThumbnailMeta(mediaMeta string, thumbPath string) (string, error) {
	if strings.TrimSpace(mediaMeta) == "" {
		meta := map[string]any{"thumbnail_path": thumbPath}
		b, err := json.Marshal(meta)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}

	var meta map[string]any
	if err := json.Unmarshal([]byte(mediaMeta), &meta); err != nil {
		return "", err
	}
	meta["thumbnail_path"] = thumbPath
	b, err := json.Marshal(meta)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func deriveShape(width, height int) string {
	if width <= 0 || height <= 0 {
		return "unknown"
	}
	ratio := float64(width) / float64(height)
	switch {
	case ratio >= 1.8:
		return "panorama"
	case ratio >= 0.88 && ratio <= 1.12:
		return "square"
	case width > height:
		return "landscape"
	default:
		return "portrait"
	}
}

func deriveSuggestedRating(meta *processor.Metadata) int {
	if meta == nil {
		return 3
	}
	score := 3

	if meta.Width > 0 && meta.Height > 0 {
		shortSide := meta.Width
		if meta.Height < shortSide {
			shortSide = meta.Height
		}
		if shortSide >= 1080 {
			score++
		} else if shortSide < 480 {
			score--
		}
	} else {
		score--
	}

	if meta.Duration > 0 && meta.Duration < 0.3 {
		score--
	}

	if meta.Extra != nil {
		if _, ok := meta.Extra["probe_error"]; ok {
			score--
		}
		if _, ok := meta.Extra["thumbnail_error"]; ok {
			score--
		}
	}

	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}
	return score
}

func (q *MediaQueue) generateThumbnail(ctx context.Context, asset *models.Asset, task *models.MediaTask) error {
	_ = q.taskService.ReportProgress(ctx, task.ID, 10)

	ext := filepath.Ext(asset.Path)
	res, err := processor.GetManager().Process(ctx, asset.Path, ext)
	if err != nil {
		logger.Warn("Thumbnail parser failed, skip thumbnail generation",
			zap.String("asset_id", asset.ID),
			zap.String("path", asset.Path),
			zap.Error(err))
		return nil
	}
	if res != nil && res.Error != nil {
		logger.Warn("Thumbnail parse returned parser error, skip thumbnail generation",
			zap.String("asset_id", asset.ID),
			zap.String("path", asset.Path),
			zap.Error(res.Error))
		return nil
	}
	if res == nil || len(res.Thumbnail) == 0 {
		// No real thumbnail generated for this format/toolchain.
		// Frontend should fall back to extension/category icon rules.
		return nil
	}

	dataDir, err := infra.ResolveDataDir()
	if err != nil {
		return err
	}

	thumbDir := filepath.Join(dataDir, "cache", "thumbnails")
	if err := os.MkdirAll(thumbDir, 0o755); err != nil {
		return err
	}

	thumbPath := filepath.Join(thumbDir, asset.ID+".jpg")
	if err := os.WriteFile(thumbPath, res.Thumbnail, 0o644); err != nil {
		return err
	}

	updatedMeta, err := mergeThumbnailMeta(asset.MediaMeta, thumbPath)
	if err != nil {
		return err
	}
	if updatedMeta != "" {
		_ = q.assetService.UpdateMediaMeta(ctx, asset.ID, updatedMeta)
	}

	_ = q.taskService.ReportProgress(ctx, task.ID, 100)
	return nil
}
