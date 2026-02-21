package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	gRuntime "runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"media-assistant-os/internal/core"
	"media-assistant-os/internal/httpapi"
	"media-assistant-os/internal/infra"
	"media-assistant-os/internal/pkg/logger"
	"media-assistant-os/internal/repos"
	"media-assistant-os/internal/services"
	"media-assistant-os/internal/utils"

	"go.uber.org/zap"
)

func main() {
	debugMode := flag.Bool("debug", false, "Enable debug monitor mode in console")
	logsMode := flag.Bool("logs", false, "Enable real-time log viewer mode")
	logDir := flag.String("log-dir", "", "Directory to store log files")
	enablePro := flag.Bool("enable-pro", false, "Enable built-in Pro routes and services")
	flag.Parse()

	proEnabled := *enablePro
	if raw := strings.TrimSpace(os.Getenv("SMART_ARCHIVE_ENABLE_PRO")); raw != "" {
		if parsed, err := strconv.ParseBool(raw); err == nil {
			proEnabled = parsed
		}
	}

	// 初始化日志
	if *logDir == "" {
		// 默认日志目录: ./data/logs
		dataDir, err := infra.ResolveDataDir()
		if err != nil {
			// 如果无法解析数据目录，回退到临时目录或当前目录
			log.Printf("failed to resolve data dir for logs: %v", err)
			*logDir = "logs"
		} else {
			*logDir = filepath.Join(dataDir, "logs")
		}
	}
	if err := logger.Init(*logDir, *debugMode); err != nil {
		log.Fatalf("logger init failed: %v", err)
	}
	defer logger.Sync()

	// Single Instance Check
	if infra.CheckAndNotifyExistingInstance() {
		logger.Info("Another instance is running. Activating it and exiting...")
		os.Exit(0)
	}

	sysInfo := utils.GetSystemInfo()
	logger.Info("Media Assistant Core Starting...",
		zap.String("version", "v0.0.1"),
		zap.Bool("debug", *debugMode),
		zap.Bool("pro_enabled", proEnabled),
		zap.String("log_dir", *logDir),
		zap.String("os", sysInfo.OS),
		zap.String("cpu", sysInfo.CPU),
		zap.Int("cores", sysInfo.CPUCores),
		zap.String("ram", utils.FormatBytes(sysInfo.TotalRAM)),
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	system := core.NewSystem()
	if err := system.Startup(ctx); err != nil {
		logger.Fatal("core startup failed", zap.Error(err))
	}
	defer system.Shutdown()

	// 记录启动日志
	system.ActivityService.Log(ctx, "INFO", fmt.Sprintf("System started successfully (pro_enabled=%t)", proEnabled))

	// 创建日志查看器
	logViewer := services.NewLogViewer(1000)
	defer logViewer.Stop()

	// 启动实时日志查看模式
	if *logsMode {
		logViewer.LogInfo("system", "启动实时日志查看模式")
		logViewer.StartRealTimeDisplay(ctx)
		return
	}

	deps := httpapi.Deps{
		EnableProFeatures: proEnabled,
		// (keep existing deps := httpapi.Deps{
		ShowMainWindow: func() {
			logger.Info("ShowMainWindow called via IPC")
			system.EventHub.Broadcast(map[string]any{
				"type":      "SHOW_MAIN_WINDOW",
				"timestamp": time.Now().UnixMilli(),
			})
		},
		QuitCore: func() {
			cancel()
		},
		SetUIReady: func(ready bool) {
			_ = ready
		},
		IndexFile: func(ctx context.Context, path string, projectID string) (any, error) {
			return system.AssetService.IndexFile(ctx, services.IndexFileRequest{
				Path:      path,
				ProjectID: projectID,
				Trigger:   "api",
			})
		},
		ArchiveFiles: func(ctx context.Context, projectID string, paths []string) error {
			return system.AssetService.ArchiveFiles(ctx, projectID, paths)
		},

		GetDefaultDirs: func() map[string]string {
			return utils.GetUserDefaultDirs()
		},
		ListProjects: func(ctx context.Context) (any, error) {
			return system.ProjectRepo.List(ctx)
		},
		CreateProject: func(ctx context.Context, name string, projectType string, path string) (any, error) {
			return system.ProjectService.CreateProject(ctx, name, projectType, path)
		},
		GetProject: func(ctx context.Context, id string) (any, error) {
			return system.ProjectRepo.Get(ctx, id)
		},
		UpdateProject: func(ctx context.Context, id string, name string, projectType string, status string, description string, path string) (any, error) {
			project, err := system.ProjectRepo.Get(ctx, id)
			if err != nil {
				return nil, err
			}
			if project == nil {
				return nil, errors.New("project not found")
			}
			pathChanged := false
			if trimmed := strings.TrimSpace(name); trimmed != "" {
				project.Name = trimmed
			}
			if trimmed := strings.TrimSpace(projectType); trimmed != "" {
				project.ProjectType = trimmed
			}
			if trimmed := strings.TrimSpace(status); trimmed != "" {
				project.Status = trimmed
			}
			if description != "" {
				project.Description = description
			}
			if path != "" {
				project.Path = path
				pathChanged = true
			}
			if strings.TrimSpace(project.ProjectType) == "" {
				project.ProjectType = "custom"
			}
			if strings.TrimSpace(project.Status) == "" {
				project.Status = "active"
			}
			if err := system.ProjectRepo.Update(ctx, *project); err != nil {
				return nil, err
			}
			if pathChanged {
				_, err := system.ProjectSourceRepo.SetPrimary(ctx, id, project.Path)
				if err != nil {
					return nil, err
				}
			}
			return project, nil
		},
		DeleteProject: func(ctx context.Context, id string) error {
			if err := system.ProjectAssetRepo.UnlinkProject(ctx, id); err != nil {
				return err
			}
			if err := system.ProjectSourceRepo.RemoveByProject(ctx, id); err != nil {
				return err
			}
			return system.ProjectRepo.Delete(ctx, id)
		},
		UpdateProjectPath: func(ctx context.Context, id string, path string) error {
			return system.ProjectService.UpdateProjectPath(ctx, id, path)
		},
		ListProjectSources: func(ctx context.Context, projectID string) (any, error) {
			return system.ProjectService.ListSources(ctx, projectID)
		},
		AddProjectSource: func(ctx context.Context, projectID string, rootPath string, sourceType string, watchEnabled *bool) (any, error) {
			return system.ProjectService.AddSource(ctx, projectID, rootPath, sourceType, watchEnabled)
		},
		RemoveProjectSource: func(ctx context.Context, projectID string, rootPath string) error {
			return system.ProjectService.RemoveSource(ctx, projectID, rootPath)
		},
		GetProjectSourceBindJob: func(ctx context.Context, jobID string) (any, error) {
			return system.ProjectService.GetSourceBindJob(ctx, jobID)
		},
		ListLibrarySources: func(ctx context.Context) (any, error) {
			return system.ProjectService.ListLibrarySources(ctx)
		},
		AddLibrarySource: func(ctx context.Context, rootPath string, watchEnabled *bool) (any, error) {
			return system.ProjectService.UpsertLibrarySource(ctx, rootPath, watchEnabled)
		},
		RemoveLibrarySource: func(ctx context.Context, rootPath string) error {
			return system.ProjectService.RemoveLibrarySource(ctx, rootPath)
		},
		ListProjectBoundDirectories: func(ctx context.Context, projectID string) (any, error) {
			return system.ProjectService.ListBoundDirectories(ctx, projectID)
		},
		ListProjectDirectoryChildren: func(ctx context.Context, projectID string, path string) (any, error) {
			return system.ProjectService.ListDirectoryChildren(ctx, projectID, path)
		},
		ListProjectDirectoryWarnings: func(ctx context.Context, projectID string, path string) (any, error) {
			return system.ProjectService.ListDirectoryWarnings(ctx, projectID, path)
		},
		StartProjectDirectoryWatch: func(ctx context.Context, projectID string, path string, ttlSeconds int) (any, error) {
			if system.WatcherService == nil {
				return nil, errors.New("watcher service is not available")
			}
			return system.WatcherService.StartSessionWatch(ctx, projectID, path, ttlSeconds)
		},
		StopProjectDirectoryWatch: func(ctx context.Context, projectID string, path string) error {
			if system.WatcherService == nil {
				return errors.New("watcher service is not available")
			}
			return system.WatcherService.StopSessionWatch(ctx, projectID, path)
		},
		GetProjectHealth: func(ctx context.Context, projectID string) (any, error) {
			return system.ProjectService.GetProjectHealth(ctx, projectID)
		},
		ScanProject: func(ctx context.Context, id string, path string) error {
			return system.ScanService.ScanProject(ctx, id, path)
		},
		GetProjectStats: func(ctx context.Context, id string) (any, error) {
			assets, err := system.AssetService.GetProjectFiles(ctx, id)
			if err != nil {
				return nil, err
			}
			totalSize := int64(0)
			for _, a := range assets {
				totalSize += a.Size
			}
			return map[string]any{
				"file_count":       len(assets),
				"total_size":       totalSize,
				"archive_progress": 0,
			}, nil
		},
		AddFileToProject: func(ctx context.Context, projectID string, fileID string) error {
			return system.ProjectAssetRepo.LinkWithOptions(ctx, projectID, fileID, &repos.ProjectAssetLinkOptions{
				Role:       "source",
				BindMode:   "manual",
				Confidence: 1.0,
			})
		},
		RemoveFileFromProject: func(ctx context.Context, projectID string, fileID string) error {
			return system.ProjectAssetRepo.Unlink(ctx, projectID, fileID)
		},
		StartInitialImport: func(ctx context.Context, directories []string, quickScanLimit int) error {
			return system.ScanService.StartInitialImport(ctx, directories, quickScanLimit)
		},
		GetImportProgress: func(ctx context.Context) (any, error) {
			return system.ScanService.GetImportProgress(), nil
		},
		RunOnboarding: func(ctx context.Context, req services.OnboardingRequest) error {
			if err := system.ProjectService.RunOnboarding(ctx, req); err != nil {
				return err
			}
			return system.SettingsService.SetFirstLaunchCompleted(true)
		},
		GetCommonDirectories: func(ctx context.Context) (any, error) {
			return getCommonDirectories()
		},
		IsFirstLaunch: func(ctx context.Context) (any, error) {
			return system.SettingsService.IsFirstLaunch(), nil
		},
		CompleteFirstLaunch: func(ctx context.Context) (any, error) {
			return map[string]any{"completed": true}, system.SettingsService.SetFirstLaunchCompleted(true)
		},
		GetProjectFiles: func(ctx context.Context, projectID string) (any, error) {
			return system.AssetService.GetProjectFiles(ctx, projectID)
		},
		ListAssets: func(ctx context.Context, req services.ListAssetsRequest) (any, error) {
			return system.AssetService.ListAssets(ctx, req)
		},
		ListAssetHistory: func(ctx context.Context, assetID string, limit int) (any, error) {
			return system.AssetService.ListAssetHistory(ctx, assetID, limit)
		},
		GetSearchHistory: func(ctx context.Context, limit int) (any, error) {
			return system.AssetService.GetSearchHistory(ctx, limit)
		},
		ClearSearchHistory: func(ctx context.Context) error {
			return system.AssetService.ClearSearchHistory(ctx)
		},
		GetAsset: func(ctx context.Context, id string) (any, error) {
			return system.AssetService.GetAsset(ctx, id)
		},
		ImportPath: func(ctx context.Context, paths []string) error {
			// This is typically called by the file watcher or manual import
			return nil // TODO: implement if needed for manual API import
		},
		UpdateAssetMeta: func(ctx context.Context, id string, mediaMeta string) error {
			return system.AssetService.UpdateMediaMeta(ctx, id, mediaMeta)
		},
		SetAssetUserRating: func(ctx context.Context, id string, userRating *int) error {
			return system.AssetService.SetUserRating(ctx, id, userRating)
		},
		DeleteAsset: func(ctx context.Context, id string) error {
			return system.AssetService.DeleteAsset(ctx, id)
		},
		BatchDeleteAssets: func(ctx context.Context, ids []string) error {
			return system.AssetService.BatchDeleteAssets(ctx, ids)
		},
		ValidateToken: func(token string) bool {
			return system.PluginService.ValidateToken(token)
		},
		ListFiles: func(ctx context.Context, path string) (any, error) {
			return system.AssetService.ListDirectory(ctx, path)
		},
		OpenFile: func(ctx context.Context, path string, assetID string) error {
			resolved, err := resolvePath(path, assetID, system.AssetRepo)
			if err != nil {
				return err
			}
			return openFilePath(resolved)
		},
		OpenInFolder: func(ctx context.Context, path string, assetID string) error {
			resolved, err := resolvePath(path, assetID, system.AssetRepo)
			if err != nil {
				return err
			}
			return openInFolder(resolved)
		},
		CreateLineage: func(ctx context.Context, ancestorID string, descendantID string, relationType string) (any, error) {
			return system.AssetService.CreateLineage(ctx, ancestorID, descendantID, relationType)
		},
		UpdateLineage: func(ctx context.Context, id string, ancestorID string, descendantID string, relationType string) error {
			return system.AssetService.UpdateLineage(ctx, id, ancestorID, descendantID, relationType)
		},
		DeleteLineage: func(ctx context.Context, id string) error {
			return system.AssetService.DeleteLineage(ctx, id)
		},
		DeleteLineageByPair: func(ctx context.Context, ancestorID string, descendantID string, relationType string) error {
			return system.AssetService.DeleteLineageByPair(ctx, ancestorID, descendantID, relationType)
		},
		ListLineage: func(ctx context.Context, assetID string) (any, error) {
			return system.AssetService.ListLineage(ctx, assetID)
		},
		ListLineageCandidates: func(ctx context.Context, assetID string, projectID string, status string, limit int) (any, error) {
			return system.AssetService.ListLineageCandidates(ctx, assetID, projectID, status, limit)
		},
		ConfirmLineageCandidate: func(ctx context.Context, candidateID string) (any, error) {
			return system.AssetService.ConfirmLineageCandidate(ctx, candidateID)
		},
		RejectLineageCandidate: func(ctx context.Context, candidateID string, reason string) error {
			return system.AssetService.RejectLineageCandidate(ctx, candidateID, reason)
		},
		RegisterPlugin: func(ctx context.Context, req services.PluginRegistrationRequest) (any, error) {
			if system.CapabilityService != nil && len(req.Mounts) > 0 {
				if err := system.CapabilityService.ValidateMounts(ctx, req.Mounts); err != nil {
					return nil, err
				}
			}
			return system.PluginService.Register(ctx, req)
		},
		ListPlugins: func(ctx context.Context) (any, error) {
			return system.PluginService.List(ctx)
		},
		ListPluginMounts: func(ctx context.Context) (any, error) {
			return system.PluginService.ListMountedSlots(ctx)
		},
		ListPluginTaskTypes: func(ctx context.Context) (any, error) {
			return system.PluginService.ListTaskTypes(ctx)
		},
		HeartbeatPlugin: func(ctx context.Context, req services.PluginHeartbeatRequest) (any, error) {
			return system.PluginService.Heartbeat(ctx, req)
		},
		ResolvePluginRuntimeEndpoint: func(ctx context.Context, pluginID string) (string, error) {
			return system.PluginService.ResolveRuntimeEndpoint(ctx, pluginID)
		},
		GetCapabilities: func(ctx context.Context) (any, error) {
			if system.CapabilityService == nil {
				return nil, errors.New("capability service is not available")
			}
			return system.CapabilityService.Get(ctx)
		},
		ListExtensionSlots: func(ctx context.Context) (any, error) {
			if system.CapabilityService == nil {
				return nil, errors.New("capability service is not available")
			}
			return system.CapabilityService.ListExtensionSlots(ctx)
		},
		ListActivityLogs: func(ctx context.Context, limit int) (any, error) {
			return system.ActivityService.List(ctx, limit)
		},
		ListReplayEvents: func(ctx context.Context, sinceID int64, limit int) (any, error) {
			return system.EventHub.ReplaySince(ctx, sinceID, limit)
		},
		GetMetrics: func(ctx context.Context) (any, error) {
			pendingTasks, err := system.MediaTaskRepo.CountByStatus(ctx, "pending")
			if err != nil {
				return nil, err
			}
			activeTasks, err := system.MediaTaskRepo.CountByStatus(ctx, "processing")
			if err != nil {
				return nil, err
			}
			failedTasks, err := system.MediaTaskRepo.CountByStatus(ctx, "failed")
			if err != nil {
				return nil, err
			}
			dlqCount, err := system.MediaTaskRepo.CountDLQ(ctx)
			if err != nil {
				return nil, err
			}
			progress := system.ScanService.GetImportProgress()
			return map[string]any{
				"tasks": map[string]any{
					"pending": pendingTasks,
					"active":  activeTasks,
					"failed":  failedTasks,
					"dlq":     dlqCount,
				},
				"scan": map[string]any{
					"queue_size":    progress.QueueSize,
					"pending_count": progress.PendingCount,
					"status":        progress.Status,
				},
				"events": map[string]any{
					"ws_connections": system.EventHub.ConnectionCount(),
				},
				"plugins": system.PluginService.Stats(),
			}, nil
		},
		ListArtifacts: func(ctx context.Context, projectID string, kind string, limit int) (any, error) {
			return system.ArtifactService.List(ctx, projectID, kind, limit)
		},
		GetArtifact: func(ctx context.Context, id string) (any, error) {
			return system.ArtifactService.Get(ctx, id)
		},
		CreateArtifact: func(ctx context.Context, req services.CreateArtifactRequest) (any, error) {
			return system.ArtifactService.Create(ctx, req)
		},
		UpdateArtifact: func(ctx context.Context, req services.UpdateArtifactRequest) (any, error) {
			return system.ArtifactService.Update(ctx, req)
		},
		DeleteArtifact: func(ctx context.Context, id string) error {
			return system.ArtifactService.Delete(ctx, id)
		},
		ListPendingTasks: func(ctx context.Context, taskTypes []string) (any, error) {
			return system.TaskService.ListPendingTasks(ctx, taskTypes)
		},
		GetActiveTasks: func(ctx context.Context) (any, error) {
			return system.TaskService.GetActiveTasks(ctx)
		},
		ClaimTask: func(ctx context.Context, taskID, workerID string) (bool, error) {
			return system.TaskService.ClaimTask(ctx, taskID, workerID)
		},
		HeartbeatTask: func(ctx context.Context, taskID, workerID string) (bool, error) {
			return system.TaskService.HeartbeatTask(ctx, taskID, workerID)
		},
		ReportTaskProgress: func(ctx context.Context, taskID string, workerID string, success bool, errMsg string, resultData map[string]any) error {
			return system.TaskService.ReportTaskProgress(ctx, taskID, workerID, success, errMsg, resultData)
		},
		EnqueueTask: func(ctx context.Context, assetID string, taskType string, priority int, force bool) (any, error) {
			return system.TaskService.EnqueueTask(ctx, assetID, taskType, priority, force)
		},
		GenerateThumbnail: func(ctx context.Context, assetID string, force bool) (any, error) {
			return system.TaskService.EnqueueThumbnailTask(ctx, assetID, force)
		},
		EventsWS: func(w http.ResponseWriter, r *http.Request) {
			system.EventHub.ServeWS(w, r)
		},
		TagService:            system.TagService,
		WorkflowService:       system.WorkflowService,
		PublishMetricsService: system.PublishMetricsService,
	}

	srv, err := httpapi.Start(ctx, 32000, 5, deps)
	if err != nil {
		log.Fatalf("http server start failed: %v", err)
	}
	log.Printf("core server listening on %s", srv.BaseURL())

	// Write server port to data directory for discovery
	if system.DataDir != "" {
		portFile := filepath.Join(system.DataDir, "server.port")
		if err := os.WriteFile(portFile, []byte(strconv.Itoa(srv.Port())), 0644); err != nil {
			logger.Error("failed to write server port file", zap.Error(err))
		}
	}

	if *debugMode {
		go runDebugMonitor(system.ScanService)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	select {
	case <-stop:
		cancel()
	case <-ctx.Done():
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	_ = srv.Close(shutdownCtx)
}

func getCommonDirectories() ([]map[string]any, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	type CommonDirectory struct {
		Type   string `json:"type"`
		Name   string `json:"name"`
		Path   string `json:"path"`
		Exists bool   `json:"exists"`
		Color  string `json:"color"`
		Icon   string `json:"icon"`
	}
	dirs := []CommonDirectory{
		{Type: "documents", Name: "下载", Path: filepath.Join(homeDir, "Downloads"), Color: "#8b5cf6", Icon: "download"},
		{Type: "pictures", Name: "图片", Path: filepath.Join(homeDir, "Pictures"), Color: "#3b82f6", Icon: "image"},
		{Type: "videos", Name: "视频", Path: filepath.Join(homeDir, "Videos"), Color: "#ef4444", Icon: "video"},
		{Type: "music", Name: "音乐", Path: filepath.Join(homeDir, "Music"), Color: "#f59e0b", Icon: "music"},
		{Type: "documents", Name: "桌面", Path: filepath.Join(homeDir, "Desktop"), Color: "#6b7280", Icon: "desktop"},
	}
	for i := range dirs {
		if _, err := os.Stat(dirs[i].Path); err == nil {
			dirs[i].Exists = true
		}
	}
	out := make([]map[string]any, 0, len(dirs))
	for _, d := range dirs {
		out = append(out, map[string]any{
			"type":   d.Type,
			"name":   d.Name,
			"path":   d.Path,
			"exists": d.Exists,
			"color":  d.Color,
			"icon":   d.Icon,
		})
	}
	return out, nil
}

func resolvePath(path string, id string, assetRepo *repos.AssetRepo) (string, error) {
	if path != "" {
		return path, nil
	}
	if id == "" {
		return "", errors.New("path is required")
	}
	if exists(id) {
		return id, nil
	}
	if assetRepo == nil {
		return "", errors.New("asset repo not initialized")
	}
	asset, err := assetRepo.GetByID(context.Background(), id)
	if err != nil {
		return "", err
	}
	if asset == nil {
		return "", errors.New("file not found")
	}
	return asset.Path, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func openFilePath(path string) error {
	switch gRuntime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", "", path).Start()
	case "darwin":
		return exec.Command("open", path).Start()
	default:
		return exec.Command("xdg-open", path).Start()
	}
}

func openInFolder(path string) error {
	if path == "" {
		return errors.New("path is required")
	}
	switch gRuntime.GOOS {
	case "windows":
		return exec.Command("explorer", "/select,", path).Start()
	case "darwin":
		return exec.Command("open", "-R", path).Start()
	default:
		dir := path
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			dir = filepath.Dir(path)
		}
		return exec.Command("xdg-open", dir).Start()
	}
}

func runDebugMonitor(scanService *services.ScanService) {
	ticker := time.NewTicker(1 * time.Second)
	fmt.Println("\n>>> Debug Monitor Started (Press Ctrl+C to exit) <<<")
	for range ticker.C {
		progress := scanService.GetImportProgress()

		// Clear current line or use ANSI escape to move cursor
		fmt.Printf("\r\033[K") // Clear line

		statusLine := fmt.Sprintf("[Status: %s] Processed: %d | Pending: %d | Queue: %d",
			progress.Status,
			progress.Processed,
			progress.PendingCount,
			progress.QueueSize)

		if progress.Detail != "" {
			detail := progress.Detail
			if len(detail) > 40 {
				detail = "..." + detail[len(detail)-37:]
			}
			statusLine += fmt.Sprintf(" | Current: %s", detail)
		}

		fmt.Print(statusLine)
	}
}
