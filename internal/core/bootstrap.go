package core

import (
	"context"
	"fmt"

	"media-assistant-os/internal/db"
	"media-assistant-os/internal/infra"
	"media-assistant-os/internal/processor"
	parsersfallback "media-assistant-os/internal/processor/parsers/fallback"
	parsersimage "media-assistant-os/internal/processor/parsers/image"
	parserspsd "media-assistant-os/internal/processor/parsers/psd"
	parsersraw "media-assistant-os/internal/processor/parsers/raw"
	parsersvideo "media-assistant-os/internal/processor/parsers/video"
	"media-assistant-os/internal/repos"
	"media-assistant-os/internal/services"
)

// System represents the Base System (Stone)
// It holds all the core logic, repositories, and services.
type System struct {
	DataDir string
	DB      *db.DB

	// Repos
	ProjectRepo              *repos.ProjectRepo
	LibrarySourceRepo        *repos.LibrarySourceRepo
	ProjectSourceRepo        *repos.ProjectSourceRepo
	ProjectSourceBindJobRepo *repos.ProjectSourceBindJobRepo
	AssetRepo                *repos.AssetRepo
	AssetHistoryEventRepo    *repos.AssetHistoryEventRepo
	SearchHistoryRepo        *repos.SearchHistoryRepo
	ProjectAssetRepo         *repos.ProjectAssetRepo
	PluginRuntimeRepo        *repos.PluginRuntimeRepo
	AssetLineageRepo         *repos.AssetLineageRepo
	LineageCandidateRepo     *repos.LineageCandidateRepo
	ActivityRepo             *repos.ActivityRepo
	MediaTaskRepo            *repos.MediaTaskRepo
	EventLogRepo             *repos.EventLogRepo
	ArtifactRepo             *repos.ProjectArtifactRepo
	TagRepo                  *repos.TagRepo
	WorkflowTemplateRepo     *repos.WorkflowTemplateRepo
	ProjectWorkflowRepo      *repos.ProjectWorkflowRepo
	ProjectWorkflowStepRepo  *repos.ProjectWorkflowStepRepo
	ProjectRoadmapRepo       *repos.ProjectRoadmapRepo
	ProjectNoteRepo          *repos.ProjectNoteRepo
	PublishRepo              *repos.PublishRepo
	MetricsRepo              *repos.MetricsRepo

	// Services
	AssetService          *services.AssetService
	ScanService           *services.ScanService
	PluginService         *services.PluginService
	CapabilityService     *services.CapabilityService
	ProjectService        *services.ProjectService
	ArtifactService       *services.ArtifactService
	EventHub              *services.EventHub
	MediaQueue            *services.MediaQueue
	ActivityService       *services.ActivityService
	TaskService           *services.TaskService
	SettingsService       *services.SettingsService
	WatcherService        *services.WatcherService
	LicenseService        *services.LicenseService
	TagService            *services.TagService
	WorkflowService       *services.WorkflowService
	PublishMetricsService *services.PublishMetricsService
}

func NewSystem() *System {
	return &System{}
}

// Startup initializes the system components
func (s *System) Startup(ctx context.Context) error {
	dataDir, err := infra.ResolveDataDir()
	if err != nil {
		return fmt.Errorf("failed to resolve data dir: %w", err)
	}
	s.DataDir = dataDir

	d, err := db.Open(s.DataDir)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	s.DB = d

	if err := db.Migrate(ctx, d); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Init port pool configuration for satellite plugins
	// Core uses dynamic port (from portfile), satellites use 32050-32099
	corePort := 32000 // Default, will be overridden by actual port
	if err := services.InitPortPoolConfig(s.DataDir, corePort); err != nil {
		fmt.Printf("⚠️  Failed to init port pool config: %v\n", err)
	}

	// Init Settings (DB Based)
	s.SettingsService = services.NewSettingsService(d.ORM())

	// Init Repos
	s.ProjectRepo = repos.NewProjectRepo(d.ORM())
	s.LibrarySourceRepo = repos.NewLibrarySourceRepo(d.ORM())
	s.ProjectSourceRepo = repos.NewProjectSourceRepo(d.ORM())
	s.ProjectSourceBindJobRepo = repos.NewProjectSourceBindJobRepo(d.ORM())
	s.AssetRepo = repos.NewAssetRepo(d.ORM())
	s.AssetHistoryEventRepo = repos.NewAssetHistoryEventRepo(d.ORM())
	s.SearchHistoryRepo = repos.NewSearchHistoryRepo(d.ORM())
	s.ProjectAssetRepo = repos.NewProjectAssetRepo(d.ORM())
	s.PluginRuntimeRepo = repos.NewPluginRuntimeRepo(d.ORM())
	s.AssetLineageRepo = repos.NewAssetLineageRepo(d.ORM())
	s.LineageCandidateRepo = repos.NewLineageCandidateRepo(d.ORM())
	s.ActivityRepo = repos.NewActivityRepo(d)
	s.MediaTaskRepo = repos.NewMediaTaskRepo(d.ORM())
	s.EventLogRepo = repos.NewEventLogRepo(d.ORM())
	s.ArtifactRepo = repos.NewProjectArtifactRepo(d.ORM())
	s.TagRepo = repos.NewTagRepo(d.ORM())
	s.WorkflowTemplateRepo = repos.NewWorkflowTemplateRepo(d.ORM())
	s.ProjectWorkflowRepo = repos.NewProjectWorkflowRepo(d.ORM())
	s.ProjectWorkflowStepRepo = repos.NewProjectWorkflowStepRepo(d.ORM())
	s.ProjectRoadmapRepo = repos.NewProjectRoadmapRepo(d.ORM())
	s.ProjectNoteRepo = repos.NewProjectNoteRepo(d.ORM())
	s.PublishRepo = repos.NewPublishRepo(d.ORM())
	s.MetricsRepo = repos.NewMetricsRepo(d.ORM())

	// Init Processors
	procMgr := processor.GetManager()
	procMgr.Register(parsersimage.New(), 20)
	procMgr.Register(parsersvideo.New(), 20)
	procMgr.Register(parserspsd.New(), 15)
	procMgr.Register(parsersraw.New(), 15)
	procMgr.Register(parsersfallback.New(), 0) // Low priority fallback

	// Init Services
	s.EventHub = services.NewEventHub(s.EventLogRepo)
	s.ActivityService = services.NewActivityService(s.ActivityRepo)
	s.TaskService = services.NewTaskService(s.MediaTaskRepo, s.AssetRepo, s.EventHub)
	s.AssetService = services.NewAssetService(s.AssetRepo, s.AssetHistoryEventRepo, s.SearchHistoryRepo, s.ProjectAssetRepo, s.ProjectRepo, s.AssetLineageRepo, s.LineageCandidateRepo, s.ActivityService, s.EventHub, s.TaskService)

	// Init License Service
	s.LicenseService = services.NewLicenseService()

	// 初始化 Bloom Filter
	if err := s.AssetService.InitBloomFilter(ctx); err != nil {
		return fmt.Errorf("failed to init bloom filter: %w", err)
	}

	// Reset any stuck tasks from previous run
	if _, err := s.TaskService.ResetAllProcessingTasks(ctx); err != nil {
		return fmt.Errorf("failed to reset processing tasks: %w", err)
	}

	s.MediaQueue = services.NewMediaQueue(s.AssetService, s.TaskService, s.EventHub, 16)
	s.MediaQueue.Start()
	s.ScanService = services.NewScanService(s.AssetService, s.ProjectRepo, s.ProjectSourceRepo, s.EventHub)
	s.ScanService.StartStartupScan(ctx) // 启动自动对账扫描

	// 初始化实时文件监控
	watcher, err := services.NewWatcherService(s.AssetService, s.ProjectRepo, s.ProjectSourceRepo)
	if err == nil {
		s.WatcherService = watcher
		s.WatcherService.Start(ctx)
	} else {
		// Log warning but don't fail startup
		fmt.Printf("Warning: Failed to start file watcher: %v\n", err)
	}

	s.ProjectService = services.NewProjectService(
		s.ProjectRepo,
		s.LibrarySourceRepo,
		s.ProjectSourceRepo,
		s.ProjectSourceBindJobRepo,
		s.ProjectAssetRepo,
		s.AssetRepo,
		s.ActivityRepo,
		s.ScanService,
		s.LicenseService,
	)
	if err := s.ProjectService.ResumeSourceBindJobs(ctx); err != nil {
		return fmt.Errorf("failed to resume source bind jobs: %w", err)
	}
	s.ArtifactService = services.NewArtifactService(s.ProjectRepo, s.ArtifactRepo)
	s.PluginService = services.NewPluginService(s.PluginRuntimeRepo)
	if err := s.PluginService.Restore(ctx); err != nil {
		return fmt.Errorf("failed to restore plugin runtime: %w", err)
	}
	s.CapabilityService = services.NewCapabilityService(s.LicenseService, s.PluginService)
	s.TagService = services.NewTagService(s.TagRepo)
	s.WorkflowService = services.NewWorkflowService(
		s.ProjectRepo,
		s.WorkflowTemplateRepo,
		s.ProjectWorkflowRepo,
		s.ProjectWorkflowStepRepo,
		s.ProjectRoadmapRepo,
		s.ProjectNoteRepo,
	)
	s.PublishMetricsService = services.NewPublishMetricsService(
		s.ProjectRepo,
		s.ProjectWorkflowRepo,
		s.ProjectRoadmapRepo,
		s.PublishRepo,
		s.MetricsRepo,
	)
	if err := s.WorkflowService.EnsureSystemTemplates(ctx); err != nil {
		return fmt.Errorf("failed to ensure workflow templates: %w", err)
	}

	return nil
}

// Shutdown closes resources
func (s *System) Shutdown() {
	if s.MediaQueue != nil {
		s.MediaQueue.Stop()
	}
	if s.ScanService != nil {
		s.ScanService.Stop()
	}
	if s.WatcherService != nil {
		s.WatcherService.Stop()
	}
	if s.DB != nil {
		s.DB.Close()
	}
}
