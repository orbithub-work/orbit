package services

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/repos"
	"media-assistant-os/internal/utils"
)

func (s *ScanService) runIntegrityCheck() {
	ticker := time.NewTicker(integrityCheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !s.shouldRunHeavyTasks() {
				continue
			}
			s.mu.Lock()
			isRunning := s.running
			s.mu.Unlock()
			if isRunning {
				continue
			}
			_ = s.reconcileAllProjects(context.Background())
		case <-s.stopChan:
			return
		}
	}
}

func (s *ScanService) runPeriodicRescan() {
	ticker := time.NewTicker(periodicRescanInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !s.shouldRunHeavyTasks() {
				continue
			}
			s.mu.Lock()
			isRunning := s.running
			s.mu.Unlock()
			if isRunning {
				continue
			}
			s.rescanAllProjects(context.Background())
		case <-s.stopChan:
			return
		}
	}
}

func (s *ScanService) reconcileAllProjects(ctx context.Context) error {
	projects, err := s.projectRepo.List(ctx)
	if err != nil {
		return err
	}
	for _, project := range projects {
		_ = s.reconcileProject(ctx, project)
	}
	return nil
}

func (s *ScanService) reconcileProject(ctx context.Context, project models.Project) error {
	assets, err := s.assetService.GetProjectFiles(ctx, project.ID)
	if err != nil {
		return err
	}

	for _, asset := range assets {
		if asset.Status == "IGNORED" {
			continue
		}
		_, statErr := os.Stat(asset.Path)
		if statErr != nil {
			if os.IsNotExist(statErr) && asset.Status != "MISSING" {
				_ = s.assetService.assets.UpdateStatusWithLog(ctx, asset.ID, "MISSING", utils.FormatNow()+": File not found on disk")
				s.assetService.recordHistoryEvent(ctx, repos.CreateAssetHistoryEventInput{
					AssetID:    asset.ID,
					ProjectID:  project.ID,
					EventType:  models.AssetHistoryEventDeleted,
					SourcePath: asset.Path,
					Confidence: "low",
					IsInferred: true,
					Detail:     "file missing inferred by reconcile scan",
				}, 30)
				if s.assetService.activities != nil {
					s.assetService.activities.LogEx(ctx, "WARN", "资产文件丢失: "+filepath.Base(asset.Path), asset.ID, project.ID)
				}
			}
			continue
		}
		if asset.Status == "MISSING" {
			_ = s.assetService.assets.UpdateStatusWithLog(ctx, asset.ID, "READY", utils.FormatNow()+": File found on disk")
		}
	}
	return nil
}

func (s *ScanService) rescanAllProjects(ctx context.Context) {
	projects, err := s.projectRepo.ListByActivity(ctx)
	if err != nil {
		return
	}
	tasks := make([]ScanTask, 0, len(projects))
	for _, project := range projects {
		roots, _ := s.resolveProjectRoots(ctx, project)
		if len(roots) == 0 {
			continue
		}
		for _, root := range roots {
			tasks = append(tasks, ScanTask{
				ProjectID: project.ID,
				Path:      root,
				Name:      project.Name + ":" + filepath.Base(root),
				Strict:    false,
			})
		}
	}
	if len(tasks) == 0 {
		return
	}
	_ = s.ScanMultipleProjects(ctx, tasks)
}
