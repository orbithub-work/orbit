package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"media-assistant-os/internal/models"
	"media-assistant-os/internal/pkg/logger"
	"media-assistant-os/internal/repos"
	"strings"
	"time"

	"go.uber.org/zap"
)

type TaskService struct {
	taskRepo  *repos.MediaTaskRepo
	assetRepo *repos.AssetRepo
	eventHub  *EventHub
	// We might need to call assetService to update asset status
}

const (
	defaultTaskLease = 30 * time.Minute
)

var coreTaskPlan = []struct {
	taskType string
	priority int
}{
	{taskType: "fingerprint", priority: 200},
	{taskType: "metadata", priority: 100},
	{taskType: "thumbnail", priority: 50},
}

func coreTaskTypes() []string {
	taskTypes := make([]string, 0, len(coreTaskPlan))
	for i := range coreTaskPlan {
		taskTypes = append(taskTypes, coreTaskPlan[i].taskType)
	}
	return taskTypes
}

func NewTaskService(taskRepo *repos.MediaTaskRepo, assetRepo *repos.AssetRepo, eventHub *EventHub) *TaskService {
	return &TaskService{
		taskRepo:  taskRepo,
		assetRepo: assetRepo,
		eventHub:  eventHub,
	}
}

// ResetAllProcessingTasks resets all tasks in processing state to pending.
// This should be called on application startup to recover orphaned tasks.
func (s *TaskService) ResetAllProcessingTasks(ctx context.Context) (int64, error) {
	return s.taskRepo.ResetProcessingTasks(ctx)
}

// CreateInitialTasks generates the necessary tasks for a new asset
func (s *TaskService) CreateInitialTasks(ctx context.Context, assetID string) error {
	for i := range coreTaskPlan {
		_, err := s.taskRepo.Create(ctx, assetID, coreTaskPlan[i].taskType, coreTaskPlan[i].priority)
		if err != nil {
			return err
		}
	}

	logger.Info("Initial tasks created for asset", zap.String("asset_id", assetID))
	return nil
}

func (s *TaskService) EnqueueTask(ctx context.Context, assetID, taskType string, priority int, force bool) (*models.MediaTask, error) {
	taskType = strings.TrimSpace(taskType)
	if taskType == "" {
		return nil, errors.New("task_type is required")
	}
	if priority <= 0 {
		priority = 50
	}

	if !force {
		tasks, err := s.taskRepo.GetByAssetID(ctx, assetID)
		if err != nil {
			return nil, err
		}
		for i := range tasks {
			task := tasks[i]
			if task.TaskType == taskType && (task.Status == models.TaskStatusPending || task.Status == models.TaskStatusProcessing) {
				return &task, nil
			}
		}
	}
	return s.taskRepo.Create(ctx, assetID, taskType, priority)
}

// ListPendingTasks for external workers
func (s *TaskService) ListPendingTasks(ctx context.Context, taskTypes []string) (*models.MediaTask, error) {
	_, _ = s.taskRepo.RequeueExpiredLeases(ctx)
	return s.taskRepo.GetNextPending(ctx, taskTypes)
}

// ClaimTask for workers
func (s *TaskService) ClaimTask(ctx context.Context, taskID, workerID string) (bool, error) {
	ok, err := s.taskRepo.Claim(ctx, taskID, workerID, time.Now().Add(defaultTaskLease))
	if ok && s.eventHub != nil {
		s.broadcastTaskUpdate(ctx, taskID, "task_claimed")
	}
	return ok, err
}

func (s *TaskService) HeartbeatTask(ctx context.Context, taskID, workerID string) (bool, error) {
	ok, err := s.taskRepo.RenewLease(ctx, taskID, workerID, time.Now().Add(defaultTaskLease))
	if ok && s.eventHub != nil {
		s.broadcastTaskUpdate(ctx, taskID, "task_heartbeat")
	}
	return ok, err
}

func (s *TaskService) ReportProgress(ctx context.Context, taskID string, progress int) error {
	if err := s.taskRepo.UpdateProgress(ctx, taskID, progress); err != nil {
		return err
	}
	if s.eventHub != nil {
		s.broadcastTaskUpdate(ctx, taskID, "task_progress")
	}
	return nil
}

// ReportTaskProgress updates the task status and result
func (s *TaskService) ReportTaskProgress(ctx context.Context, taskID string, workerID string, success bool, errMsg string, resultData map[string]any) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return fmt.Errorf("task not found: %s", taskID)
	}
	if task.Status != models.TaskStatusProcessing {
		return fmt.Errorf("task %s is not processing", taskID)
	}
	if task.WorkerID != workerID {
		return errors.New("worker mismatch for task report")
	}

	if !success {
		outcome, err := s.taskRepo.Fail(ctx, taskID, workerID, errMsg)
		if err != nil {
			return err
		}
		if outcome.Retried {
			s.broadcastTaskUpdate(ctx, taskID, "task_requeued")
		} else if outcome.DeadLettered {
			s.broadcastTaskUpdate(ctx, taskID, "task_dead_letter")
			_ = s.assetRepo.UpdateStatus(ctx, task.AssetID, "ERROR")
			if s.eventHub != nil {
				s.eventHub.Broadcast(map[string]any{
					"type": "asset_error",
					"data": map[string]any{
						"asset_id": task.AssetID,
						"reason":   "task_dead_letter",
					},
				})
			}
		} else {
			s.broadcastTaskUpdate(ctx, taskID, "task_failed")
		}
		return nil
	}

	// Handle successful completion
	if err := s.taskRepo.Complete(ctx, taskID, workerID); err != nil {
		return err
	}
	s.broadcastTaskUpdate(ctx, taskID, "task_completed")

	// Process resultData (e.g., update asset metadata in DB)
	if resultData != nil {
		switch task.TaskType {
		case "metadata":
			// If resultData contains metadata string
			if meta, ok := resultData["metadata"].(string); ok {
				_ = s.assetRepo.UpdateMediaMeta(ctx, task.AssetID, meta)
			}
		}
	}

	// Check if all tasks for this asset are done
	return s.HandleTaskCompletion(ctx, task.AssetID)
}

func (s *TaskService) EnqueueThumbnailTask(ctx context.Context, assetID string, force bool) (*models.MediaTask, error) {
	return s.EnqueueTask(ctx, assetID, "thumbnail", 50, force)
}

func (s *TaskService) GetActiveTasks(ctx context.Context) ([]models.MediaTask, error) {
	return s.taskRepo.GetActiveTasks(ctx)
}

func (s *TaskService) broadcastTaskUpdate(ctx context.Context, taskID string, eventType string) {
	if s.eventHub == nil {
		return
	}
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil || task == nil {
		return
	}

	s.eventHub.Broadcast(map[string]any{
		"type": eventType,
		"data": task,
	})
}

func (s *TaskService) HandleTaskCompletion(ctx context.Context, assetID string) error {
	allDone, err := s.taskRepo.AreTaskTypesCompleted(ctx, assetID, coreTaskTypes())
	if err != nil {
		return err
	}

	if allDone {
		// Determine final status
		status := "READY"

		// Check if it was processed by dummy parser
		asset, err := s.assetRepo.GetByID(ctx, assetID)
		if err == nil && asset != nil && asset.MediaMeta != "" {
			var metaMap map[string]any
			if json.Unmarshal([]byte(asset.MediaMeta), &metaMap) == nil {
				if extra, ok := metaMap["extra"].(map[string]any); ok {
					if parser, ok := extra["parser"].(string); ok && parser == "dummy" {
						status = "INDEXED"
					}
				}
			}
		}

		logger.Info("All tasks completed for asset",
			zap.String("asset_id", assetID),
			zap.String("status", status))

		err = s.assetRepo.UpdateStatus(ctx, assetID, status)
		if err == nil && s.eventHub != nil {
			s.eventHub.Broadcast(map[string]any{
				"type": "asset_ready",
				"data": map[string]any{
					"asset_id": assetID,
					"status":   status,
				},
			})
		}
		return err
	}

	return nil
}
