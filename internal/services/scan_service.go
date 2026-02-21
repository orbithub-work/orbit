package services

import (
	"context"
	"errors"
	"io/fs"
	"media-assistant-os/internal/models"
	"media-assistant-os/internal/repos"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	integrityCheckInterval = 10 * time.Minute
	periodicRescanInterval = 6 * time.Hour
)

type ScanTask struct {
	ProjectID string
	Path      string
	Name      string
	Strict    bool
}

type ImportDirectoryProgress struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	Progress int    `json:"progress"`
	Status   string `json:"status"`
}

type ImportProgress struct {
	Processed    int                       `json:"processed"`
	Total        int                       `json:"total"`
	Directories  []ImportDirectoryProgress `json:"directories"`
	Status       string                    `json:"status"`
	Detail       string                    `json:"detail"`
	Complete     bool                      `json:"complete"`
	QueueSize    int                       `json:"queue_size"`    // Number of pending scan batches
	PendingCount int                       `json:"pending_count"` // Number of assets waiting for background analysis
}

type ScanService struct {
	assetService      *AssetService
	projectRepo       *repos.ProjectRepo
	projectSourceRepo *repos.ProjectSourceRepo
	scanQueue         chan []ScanTask
	eventHub          *EventHub
	stopChan          chan struct{}
	running           bool
	mu                sync.Mutex
	progress          ImportProgress
	policy            *scanSystemPolicy
}

func NewScanService(assetService *AssetService, projectRepo *repos.ProjectRepo, projectSourceRepo *repos.ProjectSourceRepo, eventHub *EventHub) *ScanService {
	s := &ScanService{
		assetService:      assetService,
		projectRepo:       projectRepo,
		projectSourceRepo: projectSourceRepo,
		scanQueue:         make(chan []ScanTask, 100),
		eventHub:          eventHub,
		stopChan:          make(chan struct{}),
		policy:            newScanSystemPolicy(),
	}
	go s.processQueue()

	go s.runIntegrityCheck()
	go s.runPeriodicRescan()
	go s.runPerformancePolicy()
	return s
}

// StartStartupScan triggers a full scan on startup with a delay
func (s *ScanService) StartStartupScan(ctx context.Context) {
	go func() {
		// Wait for system to settle down
		time.Sleep(3 * time.Second)
		s.rescanAllProjects(ctx)
	}()
}

// Stop gracefully shuts down the scan service
func (s *ScanService) Stop() {
	close(s.stopChan)
	close(s.scanQueue)
}

func (s *ScanService) processQueue() {
	for {
		select {
		case tasks, ok := <-s.scanQueue:
			if !ok {
				return // Channel closed, exit
			}

			s.mu.Lock()
			s.running = true
			s.progress = ImportProgress{
				Processed:   0,
				Total:       0,
				Directories: make([]ImportDirectoryProgress, 0, len(tasks)),
				Status:      "正在批量扫描...",
				Detail:      "",
				Complete:    false,
			}
			for _, t := range tasks {
				s.progress.Directories = append(s.progress.Directories, ImportDirectoryProgress{
					Name:     t.Name,
					Path:     t.Path,
					Progress: 0,
					Status:   "等待中",
				})
			}
			s.mu.Unlock()

			// Use Background context because the original request context might be cancelled
			s.runBatchScan(context.Background(), tasks)
		case <-s.stopChan:
			return // Graceful shutdown
		}
	}
}

func (s *ScanService) ScanMultipleProjects(ctx context.Context, tasks []ScanTask) error {
	select {
	case s.scanQueue <- tasks:
		return nil
	default:
		return errors.New("scan queue is full")
	}
}

func (s *ScanService) runBatchScan(ctx context.Context, tasks []ScanTask) {
	defer func() {
		s.mu.Lock()
		s.running = false
		s.progress.Complete = true
		s.progress.Status = "扫描完成"
		s.progress.Detail = ""
		s.mu.Unlock()
	}()

	for idx, task := range tasks {
		task.Path = strings.TrimSpace(task.Path)
		if task.Path == "" {
			s.updateDir(idx, func(p *ImportDirectoryProgress) {
				p.Status = "跳过"
				p.Progress = 0
			})
			continue
		}
		if info, err := os.Stat(task.Path); err != nil || !info.IsDir() {
			s.updateDir(idx, func(p *ImportDirectoryProgress) {
				p.Status = "路径不可用"
				p.Progress = 0
			})
			continue
		}

		s.updateDir(idx, func(p *ImportDirectoryProgress) {
			p.Status = "扫描中"
			p.Progress = 0
		})

		count := 0
		_ = filepath.WalkDir(task.Path, func(filePath string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				name := d.Name()
				if s.shouldSkipDir(name, task.Strict) {
					return filepath.SkipDir
				}
				return nil
			}

			if s.shouldSkipFileWithPolicy(d.Name(), task.Strict) {
				return nil
			}

			// Use IndexFile instead of full IndexFile
			_, _ = s.assetService.IndexFile(ctx, IndexFileRequest{
				Path:      filePath,
				ProjectID: task.ProjectID,
				Trigger:   "scan",
			})
			count++

			if count%10 == 0 {
				s.updateDir(idx, func(p *ImportDirectoryProgress) {
					p.Progress = count
				})
			}
			return nil
		})

		s.updateDir(idx, func(p *ImportDirectoryProgress) {
			p.Status = "完成"
			p.Progress = count
		})
	}
}

func (s *ScanService) ScanProject(ctx context.Context, projectID string, path string) error {
	projectID = strings.TrimSpace(projectID)
	path = strings.TrimSpace(path)
	if projectID == "" {
		return errors.New("project id is required")
	}

	tasks := make([]ScanTask, 0, 4)
	if path != "" {
		tasks = append(tasks, ScanTask{
			ProjectID: projectID,
			Path:      path,
			Name:      filepath.Base(path),
			Strict:    false,
		})
	} else {
		project, err := s.projectRepo.Get(ctx, projectID)
		if err != nil {
			return err
		}
		if project == nil {
			return errors.New("project not found")
		}

		roots, err := s.resolveProjectRoots(ctx, *project)
		if err != nil {
			return err
		}
		if len(roots) == 0 {
			return errors.New("project path is empty")
		}
		for _, root := range roots {
			tasks = append(tasks, ScanTask{
				ProjectID: projectID,
				Path:      root,
				Name:      project.Name + ":" + filepath.Base(root),
				Strict:    false,
			})
		}
	}

	select {
	case s.scanQueue <- tasks:
		return nil
	default:
		return errors.New("scan queue is full")
	}
}

func (s *ScanService) resolveProjectRoots(ctx context.Context, project models.Project) ([]string, error) {
	roots := make([]string, 0, 4)
	seen := make(map[string]struct{}, 4)

	if s.projectSourceRepo != nil {
		sources, err := s.projectSourceRepo.ListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		for _, source := range sources {
			root := strings.TrimSpace(source.RootPath)
			if root == "" {
				continue
			}
			clean := filepath.Clean(root)
			if _, ok := seen[clean]; ok {
				continue
			}
			seen[clean] = struct{}{}
			roots = append(roots, clean)
		}
	}

	if len(roots) == 0 {
		root := strings.TrimSpace(project.Path)
		if root != "" {
			clean := filepath.Clean(root)
			if _, ok := seen[clean]; !ok {
				roots = append(roots, clean)
			}
		}
	}

	return roots, nil
}

func (s *ScanService) runInitialImport(ctx context.Context, directories []string, quickScanLimit int, strict bool) {
	projectID := ""

	defer func() {
		s.mu.Lock()
		s.running = false
		s.progress.Complete = true
		s.progress.Status = "导入完成"
		s.progress.Detail = ""
		s.mu.Unlock()
	}()

	for idx, dir := range directories {
		s.updateDir(idx, func(p *ImportDirectoryProgress) {
			p.Status = "扫描中"
			p.Progress = 0
		})

		var paths []string
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if d.IsDir() {
				name := d.Name()
				if s.shouldSkipDir(name, strict) {
					return filepath.SkipDir
				}
				return nil
			}

			if s.shouldSkipFileWithPolicy(d.Name(), strict) {
				return nil
			}

			paths = append(paths, path)
			if quickScanLimit > 0 && len(paths) >= quickScanLimit {
				return fs.SkipAll
			}
			return nil
		})

		for i, p := range paths {
			_, _ = s.assetService.IndexFile(ctx, IndexFileRequest{
				Path:      p,
				ProjectID: projectID,
				Trigger:   "scan",
			})
			s.mu.Lock()
			s.progress.Processed++
			s.progress.Status = "正在扫描文件..."
			s.progress.Detail = p

			// Update directory progress
			if len(paths) > 0 {
				dirIdx := idx // Capture loop variable
				progress := int(float64(i+1) / float64(len(paths)) * 100)
				// Access directories safely inside the lock
				if dirIdx >= 0 && dirIdx < len(s.progress.Directories) {
					s.progress.Directories[dirIdx].Progress = progress
				}
			}
			s.mu.Unlock()
		}

		s.updateDir(idx, func(p *ImportDirectoryProgress) {
			p.Status = "完成"
			p.Progress = 100
		})
	}
}

func (s *ScanService) StartInitialImport(ctx context.Context, directories []string, quickScanLimit int) error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return errors.New("scan already running")
	}
	s.running = true
	s.progress = ImportProgress{
		Processed:   0,
		Total:       0,
		Directories: make([]ImportDirectoryProgress, 0, len(directories)),
		Status:      "准备导入...",
		Detail:      "",
		Complete:    false,
	}
	for _, dir := range directories {
		s.progress.Directories = append(s.progress.Directories, ImportDirectoryProgress{
			Name:     filepath.Base(dir),
			Path:     dir,
			Progress: 0,
			Status:   "等待中",
		})
	}
	s.mu.Unlock()

	go s.runInitialImport(context.Background(), directories, quickScanLimit, true)
	return nil
}
