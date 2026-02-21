package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"media-assistant-os/internal/models"
	"media-assistant-os/internal/repos"
)

type ProjectService struct {
	projectRepo          *repos.ProjectRepo
	librarySourceRepo    *repos.LibrarySourceRepo
	projectSourceRepo    *repos.ProjectSourceRepo
	projectSourceJobRepo *repos.ProjectSourceBindJobRepo
	projectAssetRepo     *repos.ProjectAssetRepo
	assetRepo            *repos.AssetRepo
	activityRepo         *repos.ActivityRepo
	scanService          *ScanService
	licenseService       *LicenseService
	sourceBindWorkers    chan struct{}
	mu                   sync.Mutex
	HomeDir              string // For testing override
}

// NewProjectService 创建项目服务实例
func NewProjectService(
	projectRepo *repos.ProjectRepo,
	librarySourceRepo *repos.LibrarySourceRepo,
	projectSourceRepo *repos.ProjectSourceRepo,
	projectSourceJobRepo *repos.ProjectSourceBindJobRepo,
	projectAssetRepo *repos.ProjectAssetRepo,
	assetRepo *repos.AssetRepo,
	activityRepo *repos.ActivityRepo,
	scanService *ScanService,
	licenseService *LicenseService,
) *ProjectService {
	home, _ := os.UserHomeDir()
	return &ProjectService{
		projectRepo:          projectRepo,
		librarySourceRepo:    librarySourceRepo,
		projectSourceRepo:    projectSourceRepo,
		projectSourceJobRepo: projectSourceJobRepo,
		projectAssetRepo:     projectAssetRepo,
		assetRepo:            assetRepo,
		activityRepo:         activityRepo,
		scanService:          scanService,
		licenseService:       licenseService,
		sourceBindWorkers:    make(chan struct{}, projectSourceBindMaxWorkers),
		HomeDir:              home,
	}
}

// CreateProject 创建新项目
func (s *ProjectService) CreateProject(ctx context.Context, name string, projectType string, path string) (*models.Project, error) {
	// Check license limits before creating a new project
	license, err := s.licenseService.GetLicense(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get license info: %w", err)
	}

	if license.Type == LicenseTypeFree {
		projectCount, err := s.projectRepo.GetCount(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to count existing projects: %w", err)
		}
		if projectCount >= int64(license.MaxProjects) {
			return nil, fmt.Errorf("free tier limit reached: cannot create more than %d projects", license.MaxProjects)
		}
	}

	// 1. Create project in DB
	p, err := s.projectRepo.Create(ctx, name, projectType)
	if err != nil {
		return nil, err
	}

	// 2. Update path if provided
	if path != "" {
		if err := s.UpdateProjectPath(ctx, p.ID, path); err != nil {
			return nil, err
		}
		p.Path = path
	}

	return p, nil
}

func (s *ProjectService) UpdateProjectPath(ctx context.Context, projectID string, path string) error {
	projectID = strings.TrimSpace(projectID)
	path = strings.TrimSpace(path)
	if projectID == "" {
		return fmt.Errorf("project id is required")
	}

	if err := s.projectRepo.UpdatePath(ctx, projectID, path); err != nil {
		return err
	}

	if s.projectSourceRepo != nil && path != "" {
		_, err := s.projectSourceRepo.SetPrimary(ctx, projectID, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ProjectService) ListSources(ctx context.Context, projectID string) ([]models.ProjectSource, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []models.ProjectSource{}, nil
	}
	if s.projectSourceRepo == nil {
		return []models.ProjectSource{}, nil
	}
	return s.projectSourceRepo.ListByProject(ctx, projectID)
}

const (
	projectSourceBindSyncThreshold = 2000
	projectSourceBindBatchSize     = 500
	projectSourceBindMaxWorkers    = 2
)

type ProjectSourceBindResult struct {
	Source       *models.ProjectSource        `json:"source"`
	BindMode     string                       `json:"bind_mode"` // sync | async
	TotalAssets  int                          `json:"total_assets"`
	LinkedAssets int                          `json:"linked_assets"`
	Job          *models.ProjectSourceBindJob `json:"job,omitempty"`
	Message      string                       `json:"message,omitempty"`
}

func (s *ProjectService) ListLibrarySources(ctx context.Context) ([]models.LibrarySource, error) {
	if s.librarySourceRepo == nil {
		return []models.LibrarySource{}, nil
	}
	return s.librarySourceRepo.List(ctx)
}

func (s *ProjectService) UpsertLibrarySource(ctx context.Context, rootPath string, watchEnabled *bool) (*models.LibrarySource, error) {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil, fmt.Errorf("root_path is required")
	}
	if s.librarySourceRepo == nil {
		return nil, fmt.Errorf("library source repo is not available")
	}
	enabled := true
	if watchEnabled != nil {
		enabled = *watchEnabled
	}
	return s.librarySourceRepo.Upsert(ctx, rootPath, enabled)
}

func (s *ProjectService) RemoveLibrarySource(ctx context.Context, rootPath string) error {
	rootPath = strings.TrimSpace(rootPath)
	if rootPath == "" {
		return nil
	}
	if s.librarySourceRepo == nil {
		return nil
	}
	if s.projectSourceRepo != nil {
		count, err := s.projectSourceRepo.CountByRootPath(ctx, rootPath)
		if err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("root_path is still bound to %d project source(s)", count)
		}
	}
	return s.librarySourceRepo.Remove(ctx, rootPath)
}

func (s *ProjectService) GetSourceBindJob(ctx context.Context, jobID string) (*models.ProjectSourceBindJob, error) {
	jobID = strings.TrimSpace(jobID)
	if jobID == "" {
		return nil, fmt.Errorf("job_id is required")
	}
	if s.projectSourceJobRepo == nil {
		return nil, fmt.Errorf("project source bind job repo is not available")
	}
	return s.projectSourceJobRepo.Get(ctx, jobID)
}

func (s *ProjectService) ResumeSourceBindJobs(ctx context.Context) error {
	if s.projectSourceJobRepo == nil {
		return nil
	}
	jobs, err := s.projectSourceJobRepo.ListInflight(ctx, 200)
	if err != nil {
		return err
	}
	for _, job := range jobs {
		jobID := strings.TrimSpace(job.ID)
		if jobID == "" {
			continue
		}
		go s.runSourceBindJob(jobID)
	}
	return nil
}

func (s *ProjectService) AddSource(ctx context.Context, projectID string, rootPath string, sourceType string, watchEnabled *bool) (*ProjectSourceBindResult, error) {
	projectID = strings.TrimSpace(projectID)
	rootPath = strings.TrimSpace(rootPath)
	if projectID == "" || rootPath == "" {
		return nil, fmt.Errorf("project_id and root_path are required")
	}
	if s.projectSourceRepo == nil {
		return nil, fmt.Errorf("project source repo is not available")
	}

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	enabled := true
	if watchEnabled != nil {
		enabled = *watchEnabled
	}

	var src *models.ProjectSource
	if strings.EqualFold(strings.TrimSpace(sourceType), "primary") {
		if err := s.UpdateProjectPath(ctx, projectID, rootPath); err != nil {
			return nil, err
		}
		src, err = s.projectSourceRepo.GetByProjectAndPath(ctx, projectID, rootPath)
		if err != nil {
			return nil, err
		}
		if src != nil && src.WatchEnabled != enabled {
			src.WatchEnabled = enabled
			src, err = s.projectSourceRepo.Upsert(ctx, src)
			if err != nil {
				return nil, err
			}
		}
	} else {
		src, err = s.projectSourceRepo.Upsert(ctx, &models.ProjectSource{
			ProjectID:    projectID,
			RootPath:     rootPath,
			SourceType:   sourceType,
			WatchEnabled: enabled,
		})
		if err != nil {
			return nil, err
		}
	}

	if src == nil {
		return nil, fmt.Errorf("failed to resolve created source")
	}

	// Keep physical directory registry independent from project containers.
	if s.librarySourceRepo != nil {
		if _, err := s.librarySourceRepo.Upsert(ctx, src.RootPath, src.WatchEnabled); err != nil {
			return nil, err
		}
	}

	res := &ProjectSourceBindResult{
		Source:   src,
		BindMode: "sync",
	}
	if s.assetRepo == nil || s.projectAssetRepo == nil {
		res.Message = "asset repositories are not available, skip source binding"
		return res, nil
	}

	total, err := s.assetRepo.CountByDirectory(ctx, src.RootPath)
	if err != nil {
		return nil, err
	}
	res.TotalAssets = total
	if total == 0 {
		res.Message = "no indexed assets found under source path"
		return res, nil
	}

	if total <= projectSourceBindSyncThreshold || s.projectSourceJobRepo == nil {
		linked, err := s.bindSourceAssetsSync(ctx, projectID, src.ID, src.RootPath, total)
		if err != nil {
			return nil, err
		}
		res.LinkedAssets = linked
		if total > projectSourceBindSyncThreshold {
			res.Message = "job repo unavailable, fallback to sync binding"
		}
		return res, nil
	}

	inflight, err := s.projectSourceJobRepo.GetInflightByProjectAndSource(ctx, projectID, src.ID)
	if err != nil {
		return nil, err
	}
	if inflight != nil {
		res.BindMode = "async"
		res.Job = inflight
		res.Message = "source bind job is already running"
		return res, nil
	}

	job, err := s.projectSourceJobRepo.Create(ctx, projectID, src.ID, src.RootPath, total)
	if err != nil {
		inflight, lookupErr := s.projectSourceJobRepo.GetInflightByProjectAndSource(ctx, projectID, src.ID)
		if lookupErr != nil {
			return nil, err
		}
		if inflight == nil {
			return nil, err
		}
		res.BindMode = "async"
		res.Job = inflight
		res.Message = "source bind job is already running"
		return res, nil
	}
	res.BindMode = "async"
	res.Job = job
	res.Message = "source bind scheduled"
	go s.runSourceBindJob(job.ID)
	return res, nil
}

func (s *ProjectService) RemoveSource(ctx context.Context, projectID string, rootPath string) error {
	projectID = strings.TrimSpace(projectID)
	rootPath = strings.TrimSpace(rootPath)
	if projectID == "" || rootPath == "" {
		return nil
	}
	if s.projectSourceRepo == nil {
		return nil
	}

	src, err := s.projectSourceRepo.GetByProjectAndPath(ctx, projectID, rootPath)
	if err != nil {
		return err
	}
	if src != nil && s.projectSourceJobRepo != nil {
		inflight, err := s.projectSourceJobRepo.GetInflightByProjectAndSource(ctx, projectID, src.ID)
		if err != nil {
			return err
		}
		if inflight != nil {
			_ = s.projectSourceJobRepo.MarkFailed(ctx, inflight.ID, inflight.ProcessedAssets, "source removed")
		}
	}
	if err := s.projectSourceRepo.Remove(ctx, projectID, rootPath); err != nil {
		return err
	}
	if src != nil && s.projectAssetRepo != nil {
		if err := s.projectAssetRepo.UnlinkBySource(ctx, projectID, src.ID); err != nil {
			return err
		}
	}

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return err
	}
	if project == nil {
		return nil
	}

	currentPath := normalizeProjectRootPath(project.Path)
	removedPath := normalizeProjectRootPath(rootPath)
	if currentPath != removedPath {
		return nil
	}

	sources, err := s.projectSourceRepo.ListByProject(ctx, projectID)
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return s.projectRepo.UpdatePath(ctx, projectID, "")
	}

	next := sources[0]
	for _, src := range sources {
		if src.SourceType == "primary" {
			next = src
			break
		}
	}
	return s.UpdateProjectPath(ctx, projectID, next.RootPath)
}

func (s *ProjectService) bindSourceAssetsSync(ctx context.Context, projectID string, sourceID string, rootPath string, total int) (int, error) {
	if s.assetRepo == nil || s.projectAssetRepo == nil {
		return 0, nil
	}
	linked := 0
	offset := 0
	for {
		ids, err := s.assetRepo.ListIDsByDirectory(ctx, rootPath, projectSourceBindBatchSize, offset)
		if err != nil {
			return linked, err
		}
		if len(ids) == 0 {
			break
		}
		if err := s.projectAssetRepo.BatchLinkWithOptions(ctx, projectID, ids, &repos.ProjectAssetLinkOptions{
			Role:       "source",
			BindMode:   "source",
			Confidence: 1.0,
			SourceID:   sourceID,
		}); err != nil {
			return linked, err
		}
		linked += len(ids)
		offset += len(ids)
		// Avoid accidental endless loop if repository returns duplicate pages.
		if total > 0 && offset > total+projectSourceBindBatchSize {
			break
		}
	}
	return linked, nil
}

func (s *ProjectService) runSourceBindJob(jobID string) {
	ctx := context.Background()
	if s.projectSourceJobRepo == nil || s.assetRepo == nil || s.projectAssetRepo == nil {
		return
	}
	if s.sourceBindWorkers != nil {
		s.sourceBindWorkers <- struct{}{}
		defer func() { <-s.sourceBindWorkers }()
	}

	job, err := s.projectSourceJobRepo.Get(ctx, jobID)
	if err != nil || job == nil {
		return
	}
	if err := s.projectSourceJobRepo.MarkRunning(ctx, jobID); err != nil {
		_ = s.projectSourceJobRepo.MarkFailed(ctx, jobID, job.ProcessedAssets, err.Error())
		return
	}

	processed := job.ProcessedAssets
	offset := processed
	for {
		current, err := s.projectSourceJobRepo.Get(ctx, jobID)
		if err != nil {
			_ = s.projectSourceJobRepo.MarkFailed(ctx, jobID, processed, err.Error())
			return
		}
		if current == nil {
			return
		}
		if current.Status == models.ProjectSourceBindJobFailed || current.Status == models.ProjectSourceBindJobSucceeded {
			return
		}
		if s.projectSourceRepo != nil {
			src, err := s.projectSourceRepo.GetByID(ctx, job.SourceID)
			if err != nil {
				_ = s.projectSourceJobRepo.MarkFailed(ctx, jobID, processed, err.Error())
				return
			}
			if src == nil {
				_ = s.projectSourceJobRepo.MarkFailed(ctx, jobID, processed, "source no longer exists")
				return
			}
		}

		ids, err := s.assetRepo.ListIDsByDirectory(ctx, job.RootPath, projectSourceBindBatchSize, offset)
		if err != nil {
			_ = s.projectSourceJobRepo.MarkFailed(ctx, jobID, processed, err.Error())
			return
		}
		if len(ids) == 0 {
			break
		}
		if err := s.projectAssetRepo.BatchLinkWithOptions(ctx, job.ProjectID, ids, &repos.ProjectAssetLinkOptions{
			Role:       "source",
			BindMode:   "source",
			Confidence: 1.0,
			SourceID:   job.SourceID,
		}); err != nil {
			_ = s.projectSourceJobRepo.MarkFailed(ctx, jobID, processed, err.Error())
			return
		}
		processed += len(ids)
		offset += len(ids)
		_ = s.projectSourceJobRepo.UpdateProgress(ctx, jobID, processed)
	}
	_ = s.projectSourceJobRepo.MarkSucceeded(ctx, jobID, processed)
}

type ProjectDirectoryNode struct {
	Path        string `json:"path"`
	Name        string `json:"name"`
	SourceType  string `json:"source_type,omitempty"`
	IsRoot      bool   `json:"is_root"`
	Exists      bool   `json:"exists"`
	HasChildren bool   `json:"has_children"`
}

func (s *ProjectService) ListBoundDirectories(ctx context.Context, projectID string) ([]ProjectDirectoryNode, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return []ProjectDirectoryNode{}, nil
	}

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	roots, err := s.resolveProjectRootRefs(ctx, *project)
	if err != nil {
		return nil, err
	}
	out := make([]ProjectDirectoryNode, 0, len(roots))
	for _, root := range roots {
		exists := pathExists(root.Path)
		out = append(out, ProjectDirectoryNode{
			Path:        root.Path,
			Name:        directoryDisplayName(root.Path),
			SourceType:  root.SourceType,
			IsRoot:      true,
			Exists:      exists,
			HasChildren: exists && hasVisibleSubDirectory(root.Path),
		})
	}
	return out, nil
}

func (s *ProjectService) ListDirectoryChildren(ctx context.Context, projectID string, path string) ([]ProjectDirectoryNode, error) {
	projectID = strings.TrimSpace(projectID)
	path = normalizeProjectRootPath(path)
	if projectID == "" || path == "" {
		return []ProjectDirectoryNode{}, fmt.Errorf("project_id and path are required")
	}

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	roots, err := s.resolveProjectRootRefs(ctx, *project)
	if err != nil {
		return nil, err
	}

	allowed := false
	for _, root := range roots {
		if isPathWithinRoot(path, root.Path) {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, fmt.Errorf("path is not under project bound directories")
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []ProjectDirectoryNode{}, nil
		}
		return nil, err
	}
	if !info.IsDir() {
		return []ProjectDirectoryNode{}, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	children := make([]ProjectDirectoryNode, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := strings.TrimSpace(entry.Name())
		if name == "" {
			continue
		}
		// Hide dot directories by default to keep creator-facing tree clean.
		if strings.HasPrefix(name, ".") {
			continue
		}
		childPath := filepath.Clean(filepath.Join(path, name))
		children = append(children, ProjectDirectoryNode{
			Path:        childPath,
			Name:        name,
			IsRoot:      false,
			Exists:      true,
			HasChildren: hasVisibleSubDirectory(childPath),
		})
	}

	sort.Slice(children, func(i, j int) bool {
		return strings.ToLower(children[i].Name) < strings.ToLower(children[j].Name)
	})
	return children, nil
}

type projectRootRef struct {
	Path       string
	SourceType string
}

type DirectoryWarning struct {
	Code      string `json:"code"`
	Level     string `json:"level"`
	ProjectID string `json:"project_id"`
	Path      string `json:"path"`
	Message   string `json:"message"`
	CreatedAt int64  `json:"created_at"`
}

func (s *ProjectService) ListDirectoryWarnings(ctx context.Context, projectID string, path string) ([]DirectoryWarning, error) {
	projectID = strings.TrimSpace(projectID)
	path = normalizeProjectRootPath(path)
	if projectID == "" {
		return []DirectoryWarning{}, nil
	}
	if s.activityRepo == nil {
		return []DirectoryWarning{}, nil
	}

	logs, err := s.activityRepo.ListByProject(ctx, projectID, []string{"WARN", "ERROR"}, 200)
	if err != nil {
		return nil, err
	}

	const prefix = "monitor permission denied: "
	out := make([]DirectoryWarning, 0, len(logs))
	seen := make(map[string]struct{}, len(logs))
	for _, item := range logs {
		msg := strings.TrimSpace(item.Message)
		if !strings.HasPrefix(strings.ToLower(msg), prefix) {
			continue
		}
		warnPath := normalizeProjectRootPath(strings.TrimSpace(msg[len(prefix):]))
		if warnPath == "" {
			continue
		}
		if path != "" && !isPathWithinRoot(warnPath, path) && !isPathWithinRoot(path, warnPath) {
			continue
		}
		key := item.ProjectID + "::watch_permission_denied::" + warnPath
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}

		out = append(out, DirectoryWarning{
			Code:      "watch_permission_denied",
			Level:     strings.ToUpper(strings.TrimSpace(item.Level)),
			ProjectID: item.ProjectID,
			Path:      warnPath,
			Message:   msg,
			CreatedAt: item.CreatedAt,
		})
	}
	return out, nil
}

func (s *ProjectService) resolveProjectRootRefs(ctx context.Context, project models.Project) ([]projectRootRef, error) {
	roots := make([]projectRootRef, 0, 4)
	seen := make(map[string]struct{}, 4)

	if s.projectSourceRepo != nil {
		sources, err := s.projectSourceRepo.ListByProject(ctx, project.ID)
		if err != nil {
			return nil, err
		}
		for _, source := range sources {
			root := normalizeProjectRootPath(source.RootPath)
			if root == "" {
				continue
			}
			if _, ok := seen[root]; ok {
				continue
			}
			seen[root] = struct{}{}
			roots = append(roots, projectRootRef{
				Path:       root,
				SourceType: source.SourceType,
			})
		}
	}

	projectRoot := normalizeProjectRootPath(project.Path)
	if projectRoot != "" {
		if _, ok := seen[projectRoot]; !ok {
			roots = append(roots, projectRootRef{
				Path:       projectRoot,
				SourceType: "primary",
			})
		}
	}
	return roots, nil
}

type ProjectHealthSource struct {
	RootPath     string `json:"root_path"`
	SourceType   string `json:"source_type"`
	WatchEnabled bool   `json:"watch_enabled"`
	Exists       bool   `json:"exists"`
}

type ProjectHealthIssue struct {
	Type      string `json:"type"`
	Role      string `json:"role,omitempty"`
	AssetID   string `json:"asset_id,omitempty"`
	Path      string `json:"path,omitempty"`
	AssetName string `json:"asset_name,omitempty"`
	Message   string `json:"message"`
}

type ProjectHealthReport struct {
	ProjectID string `json:"project_id"`
	CheckedAt int64  `json:"checked_at"`

	TotalBindings    int `json:"total_bindings"`
	SourceCount      int `json:"source_count"`
	EngineCount      int `json:"engine_count"`
	DeliverableCount int `json:"deliverable_count"`

	MissingBindings     int `json:"missing_bindings"`
	MissingSources      int `json:"missing_sources"`
	MissingEngines      int `json:"missing_engines"`
	MissingDeliverables int `json:"missing_deliverables"`
	MissingRoots        int `json:"missing_roots"`

	UnarchivedDeliverables bool `json:"unarchived_deliverables"`

	Sources []ProjectHealthSource `json:"sources"`
	Issues  []ProjectHealthIssue  `json:"issues"`
}

func (s *ProjectService) GetProjectHealth(ctx context.Context, projectID string) (*ProjectHealthReport, error) {
	projectID = strings.TrimSpace(projectID)
	if projectID == "" {
		return nil, fmt.Errorf("project id is required")
	}

	project, err := s.projectRepo.Get(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	report := &ProjectHealthReport{
		ProjectID: projectID,
		CheckedAt: time.Now().Unix(),
		Sources:   []ProjectHealthSource{},
		Issues:    []ProjectHealthIssue{},
	}

	if s.projectSourceRepo != nil {
		sources, err := s.projectSourceRepo.ListByProject(ctx, projectID)
		if err != nil {
			return nil, err
		}
		for _, src := range sources {
			exists := pathExists(src.RootPath)
			report.Sources = append(report.Sources, ProjectHealthSource{
				RootPath:     src.RootPath,
				SourceType:   src.SourceType,
				WatchEnabled: src.WatchEnabled,
				Exists:       exists,
			})
			if !exists {
				report.MissingRoots++
				report.Issues = append(report.Issues, ProjectHealthIssue{
					Type:    "source_root_missing",
					Path:    src.RootPath,
					Message: "source directory is missing or inaccessible",
				})
			}
		}
	}

	if s.projectAssetRepo == nil {
		return report, nil
	}

	details, err := s.projectAssetRepo.ListBindingDetailsByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	report.TotalBindings = len(details)
	for _, detail := range details {
		role := normalizeBindingRole(detail.Role)
		switch role {
		case "engine":
			report.EngineCount++
		case "deliverable":
			report.DeliverableCount++
		default:
			report.SourceCount++
		}

		missing, reason := bindingMissingReason(detail)
		if !missing {
			continue
		}

		report.MissingBindings++
		switch role {
		case "engine":
			report.MissingEngines++
		case "deliverable":
			report.MissingDeliverables++
		default:
			report.MissingSources++
		}

		report.Issues = append(report.Issues, ProjectHealthIssue{
			Type:      "binding_missing",
			Role:      role,
			AssetID:   detail.AssetID,
			Path:      detail.Path,
			AssetName: filepath.Base(detail.Path),
			Message:   reason,
		})
	}

	if report.DeliverableCount == 0 {
		report.UnarchivedDeliverables = true
		report.Issues = append(report.Issues, ProjectHealthIssue{
			Type:    "deliverable_not_bound",
			Role:    "deliverable",
			Message: "no deliverable assets bound to this project",
		})
	}

	return report, nil
}

func bindingMissingReason(detail repos.ProjectAssetBindingDetail) (bool, string) {
	if strings.EqualFold(strings.TrimSpace(detail.AssetStatus), "MISSING") {
		return true, "asset status is MISSING"
	}
	cleanPath := strings.TrimSpace(detail.Path)
	if cleanPath == "" {
		return true, "asset path is empty"
	}
	if _, err := os.Stat(cleanPath); err != nil {
		if os.IsNotExist(err) {
			return true, "asset file not found on disk"
		}
		return true, "asset path inaccessible: " + err.Error()
	}
	return false, ""
}

func normalizeBindingRole(role string) string {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "engine":
		return "engine"
	case "deliverable":
		return "deliverable"
	default:
		return "source"
	}
}

func normalizeProjectRootPath(path string) string {
	v := strings.TrimSpace(path)
	if v == "" {
		return ""
	}
	abs, err := filepath.Abs(v)
	if err == nil {
		v = abs
	}
	return filepath.Clean(v)
}

func pathExists(path string) bool {
	path = strings.TrimSpace(path)
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func directoryDisplayName(path string) string {
	base := filepath.Base(path)
	if base == "." || base == string(filepath.Separator) || base == "" {
		return path
	}
	return base
}

func hasVisibleSubDirectory(path string) bool {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			return true
		}
	}
	return false
}

func isPathWithinRoot(path string, root string) bool {
	path = normalizeProjectRootPath(path)
	root = normalizeProjectRootPath(root)
	if path == "" || root == "" {
		return false
	}
	if path == root {
		return true
	}
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, ".."+string(filepath.Separator)) && rel != ".."
}

// RunOnboarding 执行用户引导流程，创建系统虚拟项目并扫描指定目录
func (s *ProjectService) RunOnboarding(ctx context.Context, req OnboardingRequest) error {
	homeDir := s.HomeDir
	if homeDir == "" {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home dir: %w", err)
		}
	}

	tasks := []struct {
		Enabled bool
		Name    string
		Path    string
	}{
		{req.ImportDownloads, "下载", filepath.Join(homeDir, "Downloads")},
		{req.ImportPictures, "图片", filepath.Join(homeDir, "Pictures")},
		{req.ImportVideos, "视频", filepath.Join(homeDir, "Videos")},
		{req.ImportMusic, "音乐", filepath.Join(homeDir, "Music")},
		{req.ImportDesktop, "桌面", filepath.Join(homeDir, "Desktop")},
	}

	var scanTasks []ScanTask

	for _, task := range tasks {
		if !task.Enabled {
			continue
		}

		// Check if directory exists
		if _, err := os.Stat(task.Path); os.IsNotExist(err) {
			continue
		}

		// Create Virtual Project
		// Note: Onboarding projects are system-generated and should bypass license limits
		// For now, we'll call the internal projectRepo.Create directly.
		// In a more robust system, ProjectService.CreateProject might have an internal flag
		// or a separate method for system projects.
		p, err := s.projectRepo.Create(ctx, task.Name, "system_virtual")
		if err != nil {
			// Log error but continue
			fmt.Printf("Failed to create onboarding project %s: %v\n", task.Name, err)
			continue
		}

		// Update path for system project
		if err := s.UpdateProjectPath(ctx, p.ID, task.Path); err != nil {
			fmt.Printf("Failed to update path for onboarding project %s: %v\n", task.Name, err)
			continue
		}
		p.Path = task.Path

		// Add to scan tasks
		scanTasks = append(scanTasks, ScanTask{
			ProjectID: p.ID,
			Path:      task.Path,
			Name:      task.Name,
			Strict:    true,
		})
	}

	if len(scanTasks) > 0 {
		if err := s.scanService.ScanMultipleProjects(ctx, scanTasks); err != nil {
			return fmt.Errorf("failed to start batch scan: %w", err)
		}
	}

	return nil
}
