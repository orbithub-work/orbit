package services

import (
	"context"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/pkg/logger"
	"media-assistant-os/internal/repos"

	"go.uber.org/zap"
)

type ActivityService struct {
	repo *repos.ActivityRepo
}

// NewActivityService 创建活动日志服务实例
func NewActivityService(repo *repos.ActivityRepo) *ActivityService {
	return &ActivityService{repo: repo}
}

// Log 记录一条用户可见的操作日志
// level: INFO, SUCCESS, WARN, ERROR
func (s *ActivityService) Log(ctx context.Context, level, message string) {
	s.LogEx(ctx, level, message, "", "")
}

// LogEx 记录扩展的操作日志，包含资产ID和项目ID
func (s *ActivityService) LogEx(ctx context.Context, level, message, assetID, projectID string) {
	// 1. 写入数据库 (用户可见)
	if err := s.repo.CreateEx(ctx, level, message, assetID, projectID); err != nil {
		// 如果数据库写入失败，记录到系统日志
		logger.Error("failed to create activity log",
			zap.String("level", level),
			zap.String("message", message),
			zap.String("asset_id", assetID),
			zap.String("project_id", projectID),
			zap.Error(err),
		)
	}

	// 2. 同时写入系统日志 (开发可见)
	fields := []zap.Field{
		zap.String("source", "activity"),
		zap.String("level", level),
	}
	if assetID != "" {
		fields = append(fields, zap.String("asset_id", assetID))
	}
	if projectID != "" {
		fields = append(fields, zap.String("project_id", projectID))
	}

	switch level {
	case "ERROR":
		logger.Error(message, fields...)
	case "WARN":
		logger.Warn(message, fields...)
	default:
		logger.Info(message, fields...)
	}
}

// List 获取最近的操作日志
func (s *ActivityService) List(ctx context.Context, limit int) ([]models.ActivityLog, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.repo.List(ctx, limit)
}
