package repos

import (
	"context"
	"database/sql"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/utils"

	"github.com/uptrace/bun"
)

type MediaTaskRepo struct {
	db *bun.DB
}

type TaskFailOutcome struct {
	Retried      bool
	DeadLettered bool
	RetryCount   int
}

func NewMediaTaskRepo(db *bun.DB) *MediaTaskRepo {
	return &MediaTaskRepo{db: db}
}

func (r *MediaTaskRepo) Create(ctx context.Context, assetID, taskType string, priority int) (*models.MediaTask, error) {
	existing, err := r.GetInflightByAssetAndType(ctx, assetID, taskType)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	task := &models.MediaTask{
		ID:         utils.NewID(),
		AssetID:    assetID,
		TaskType:   taskType,
		Status:     models.TaskStatusPending,
		Priority:   priority,
		MaxRetries: 3,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	_, err = r.db.NewInsert().Model(task).Exec(ctx)
	if err != nil {
		existing, findErr := r.GetInflightByAssetAndType(ctx, assetID, taskType)
		if findErr == nil && existing != nil {
			return existing, nil
		}
		return nil, err
	}
	return task, nil
}

func (r *MediaTaskRepo) GetNextPending(ctx context.Context, taskTypes []string) (*models.MediaTask, error) {
	var task models.MediaTask
	q := r.db.NewSelect().
		Model(&task).
		Where("status = ?", models.TaskStatusPending).
		Where("(next_retry_at IS NULL OR next_retry_at <= CURRENT_TIMESTAMP)")

	if len(taskTypes) > 0 {
		q = q.Where("task_type IN (?)", bun.In(taskTypes))
	}

	err := q.Order("priority DESC", "created_at ASC").
		Limit(1).
		Scan(ctx)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (r *MediaTaskRepo) Claim(ctx context.Context, taskID, workerID string, leaseUntil time.Time) (bool, error) {
	now := time.Now()
	res, err := r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("status = ?", models.TaskStatusProcessing).
		Set("worker_id = ?", workerID).
		Set("progress = ?", 0).
		Set("next_retry_at = NULL").
		Set("started_at = ?", now).
		Set("updated_at = ?", now).
		Set("lease_until = ?", leaseUntil).
		Where("id = ?", taskID).
		Where("status = ?", models.TaskStatusPending).
		Where("(next_retry_at IS NULL OR next_retry_at <= CURRENT_TIMESTAMP)").
		Exec(ctx)

	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *MediaTaskRepo) RenewLease(ctx context.Context, taskID, workerID string, leaseUntil time.Time) (bool, error) {
	res, err := r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("lease_until = ?", leaseUntil).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", taskID).
		Where("status = ?", models.TaskStatusProcessing).
		Where("worker_id = ?", workerID).
		Exec(ctx)
	if err != nil {
		return false, err
	}
	rows, _ := res.RowsAffected()
	return rows > 0, nil
}

func (r *MediaTaskRepo) UpdateProgress(ctx context.Context, taskID string, progress int) error {
	_, err := r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("progress = ?", progress).
		Set("updated_at = ?", time.Now()).
		Where("id = ?", taskID).
		Exec(ctx)
	return err
}

func (r *MediaTaskRepo) Complete(ctx context.Context, taskID string, workerID string) error {
	_, err := r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("status = ?", models.TaskStatusCompleted).
		Set("progress = ?", 100).
		Set("next_retry_at = NULL").
		Set("dead_letter_reason = ?", "").
		Set("finished_at = ?", time.Now()).
		Set("updated_at = ?", time.Now()).
		Set("lease_until = NULL").
		Where("id = ?", taskID).
		Where("status = ?", models.TaskStatusProcessing).
		Where("worker_id = ?", workerID).
		Exec(ctx)
	return err
}

func (r *MediaTaskRepo) Fail(ctx context.Context, taskID, workerID, errMsg string) (TaskFailOutcome, error) {
	task, err := r.GetByID(ctx, taskID)
	if err != nil {
		return TaskFailOutcome{}, err
	}
	if task == nil {
		return TaskFailOutcome{}, sql.ErrNoRows
	}

	if task.Status != models.TaskStatusProcessing {
		return TaskFailOutcome{}, sql.ErrNoRows
	}
	if task.WorkerID != workerID {
		return TaskFailOutcome{}, sql.ErrNoRows
	}

	now := time.Now()
	nextRetryCount := task.RetryCount + 1
	if nextRetryCount <= task.MaxRetries {
		nextRetryAt := now.Add(backoffDuration(nextRetryCount))
		_, err = r.db.NewUpdate().
			Model((*models.MediaTask)(nil)).
			Set("status = ?", models.TaskStatusPending).
			Set("retry_count = ?", nextRetryCount).
			Set("error_message = ?", errMsg).
			Set("next_retry_at = ?", nextRetryAt).
			Set("worker_id = ?", "").
			Set("started_at = NULL").
			Set("finished_at = NULL").
			Set("lease_until = NULL").
			Set("updated_at = ?", now).
			Where("id = ?", taskID).
			Exec(ctx)
		if err != nil {
			return TaskFailOutcome{}, err
		}
		return TaskFailOutcome{
			Retried:    true,
			RetryCount: nextRetryCount,
		}, nil
	}

	_, err = r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("status = ?", models.TaskStatusFailed).
		Set("retry_count = ?", nextRetryCount).
		Set("error_message = ?", errMsg).
		Set("dead_letter_reason = ?", "max retries exceeded").
		Set("next_retry_at = NULL").
		Set("finished_at = ?", now).
		Set("updated_at = ?", now).
		Set("lease_until = NULL").
		Where("id = ?", taskID).
		Exec(ctx)
	if err != nil {
		return TaskFailOutcome{}, err
	}

	_, err = r.db.NewInsert().
		Model(map[string]any{
			"id":            utils.NewID(),
			"task_id":       taskID,
			"asset_id":      task.AssetID,
			"task_type":     task.TaskType,
			"error_message": errMsg,
			"retry_count":   nextRetryCount,
			"max_retries":   task.MaxRetries,
			"created_at":    now,
		}).
		Table("media_task_dlq").
		Exec(ctx)
	if err != nil {
		return TaskFailOutcome{}, err
	}

	return TaskFailOutcome{
		DeadLettered: true,
		RetryCount:   nextRetryCount,
	}, nil
}

func (r *MediaTaskRepo) RequeueExpiredLeases(ctx context.Context) (int64, error) {
	res, err := r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("status = ?", models.TaskStatusPending).
		Set("worker_id = ?", "").
		Set("started_at = NULL").
		Set("lease_until = NULL").
		Set("updated_at = ?", time.Now()).
		Where("status = ?", models.TaskStatusProcessing).
		Where("lease_until IS NOT NULL AND lease_until <= CURRENT_TIMESTAMP").
		Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *MediaTaskRepo) GetActiveTasks(ctx context.Context) ([]models.MediaTask, error) {
	var tasks []models.MediaTask
	err := r.db.NewSelect().
		Model(&tasks).
		Where("status = ?", models.TaskStatusProcessing).
		Order("updated_at DESC").
		Scan(ctx)
	return tasks, err
}

func backoffDuration(retryCount int) time.Duration {
	if retryCount < 1 {
		return 5 * time.Second
	}
	seconds := 5 << (retryCount - 1)
	if seconds > 300 {
		seconds = 300
	}
	return time.Duration(seconds) * time.Second
}

func (r *MediaTaskRepo) GetByAssetID(ctx context.Context, assetID string) ([]models.MediaTask, error) {
	var tasks []models.MediaTask
	err := r.db.NewSelect().
		Model(&tasks).
		Where("asset_id = ?", assetID).
		Scan(ctx)
	return tasks, err
}

func (r *MediaTaskRepo) GetByID(ctx context.Context, id string) (*models.MediaTask, error) {
	var task models.MediaTask
	err := r.db.NewSelect().Model(&task).Where("id = ?", id).Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}

func (r *MediaTaskRepo) AreAllTasksCompleted(ctx context.Context, assetID string) (bool, error) {
	count, err := r.db.NewSelect().
		Model((*models.MediaTask)(nil)).
		Where("asset_id = ?", assetID).
		Where("status != ?", models.TaskStatusCompleted).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *MediaTaskRepo) AreTaskTypesCompleted(ctx context.Context, assetID string, taskTypes []string) (bool, error) {
	if len(taskTypes) == 0 {
		return true, nil
	}
	count, err := r.db.NewSelect().
		Model((*models.MediaTask)(nil)).
		Where("asset_id = ?", assetID).
		Where("task_type IN (?)", bun.In(taskTypes)).
		Where("status != ?", models.TaskStatusCompleted).
		Count(ctx)
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (r *MediaTaskRepo) CountByStatus(ctx context.Context, status models.TaskStatus) (int, error) {
	return r.db.NewSelect().
		Model((*models.MediaTask)(nil)).
		Where("status = ?", status).
		Count(ctx)
}

func (r *MediaTaskRepo) CountDLQ(ctx context.Context) (int, error) {
	return r.db.NewSelect().
		Table("media_task_dlq").
		Count(ctx)
}

func (r *MediaTaskRepo) ResetProcessingTasks(ctx context.Context) (int64, error) {
	res, err := r.db.NewUpdate().
		Model((*models.MediaTask)(nil)).
		Set("status = ?", models.TaskStatusPending).
		Set("worker_id = ?", "").
		Set("started_at = NULL").
		Set("finished_at = NULL").
		Set("lease_until = NULL").
		Set("updated_at = ?", time.Now()).
		Where("status = ?", models.TaskStatusProcessing).
		Exec(ctx)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (r *MediaTaskRepo) GetInflightByAssetAndType(ctx context.Context, assetID, taskType string) (*models.MediaTask, error) {
	var task models.MediaTask
	err := r.db.NewSelect().
		Model(&task).
		Where("asset_id = ?", assetID).
		Where("task_type = ?", taskType).
		Where("status IN (?)", bun.In([]models.TaskStatus{models.TaskStatusPending, models.TaskStatusProcessing})).
		Order("created_at ASC").
		Limit(1).
		Scan(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &task, err
}
