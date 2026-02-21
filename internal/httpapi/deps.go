package httpapi

import (
	"context"
	"net/http"

	"media-assistant-os/internal/services"
)

type Deps struct {
	EnableProFeatures            bool
	IndexFile                    func(ctx context.Context, path string, projectID string) (any, error)
	ArchiveFiles                 func(ctx context.Context, projectID string, paths []string) error
	ListProjects                 func(ctx context.Context) (any, error)
	CreateProject                func(ctx context.Context, name string, projectType string, path string) (any, error)
	GetProject                   func(ctx context.Context, id string) (any, error)
	UpdateProject                func(ctx context.Context, id string, name string, projectType string, status string, description string, path string) (any, error)
	DeleteProject                func(ctx context.Context, id string) error
	UpdateProjectPath            func(ctx context.Context, id string, path string) error
	ListProjectSources           func(ctx context.Context, projectID string) (any, error)
	AddProjectSource             func(ctx context.Context, projectID string, rootPath string, sourceType string, watchEnabled *bool) (any, error)
	RemoveProjectSource          func(ctx context.Context, projectID string, rootPath string) error
	GetProjectSourceBindJob      func(ctx context.Context, jobID string) (any, error)
	ListLibrarySources           func(ctx context.Context) (any, error)
	AddLibrarySource             func(ctx context.Context, rootPath string, watchEnabled *bool) (any, error)
	RemoveLibrarySource          func(ctx context.Context, rootPath string) error
	ListProjectBoundDirectories  func(ctx context.Context, projectID string) (any, error)
	ListProjectDirectoryChildren func(ctx context.Context, projectID string, path string) (any, error)
	ListProjectDirectoryWarnings func(ctx context.Context, projectID string, path string) (any, error)
	StartProjectDirectoryWatch   func(ctx context.Context, projectID string, path string, ttlSeconds int) (any, error)
	StopProjectDirectoryWatch    func(ctx context.Context, projectID string, path string) error
	GetProjectHealth             func(ctx context.Context, projectID string) (any, error)
	ScanProject                  func(ctx context.Context, id string, path string) error
	GetProjectStats              func(ctx context.Context, id string) (any, error)
	AddFileToProject             func(ctx context.Context, projectID string, fileID string) error
	RemoveFileFromProject        func(ctx context.Context, projectID string, fileID string) error
	StartInitialImport           func(ctx context.Context, directories []string, quickScanLimit int) error
	GetImportProgress            func(ctx context.Context) (any, error)
	IsFirstLaunch                func(ctx context.Context) (any, error)
	CompleteFirstLaunch          func(ctx context.Context) (any, error)
	GetCommonDirectories         func(ctx context.Context) (any, error)
	RunOnboarding                func(ctx context.Context, req services.OnboardingRequest) error
	GetDefaultDirs               func() map[string]string
	GetProjectFiles              func(ctx context.Context, projectID string) (any, error)
	ListAssets                   func(ctx context.Context, req services.ListAssetsRequest) (any, error)
	ListAssetHistory             func(ctx context.Context, assetID string, limit int) (any, error)
	GetSearchHistory             func(ctx context.Context, limit int) (any, error)
	ClearSearchHistory           func(ctx context.Context) error
	GetAsset                     func(ctx context.Context, id string) (any, error)
	ImportPath                   func(ctx context.Context, paths []string) error
	UpdateAssetMeta              func(ctx context.Context, id string, mediaMeta string) error
	SetAssetUserRating           func(ctx context.Context, id string, userRating *int) error
	DeleteAsset                  func(ctx context.Context, id string) error
	BatchDeleteAssets            func(ctx context.Context, ids []string) error
	ValidateToken                func(token string) bool
	ListFiles                    func(ctx context.Context, path string) (any, error)
	OpenFile                     func(ctx context.Context, path string, assetID string) error
	OpenInFolder                 func(ctx context.Context, path string, assetID string) error
	CreateLineage                func(ctx context.Context, ancestorID string, descendantID string, relationType string) (any, error)
	UpdateLineage                func(ctx context.Context, id string, ancestorID string, descendantID string, relationType string) error
	DeleteLineage                func(ctx context.Context, id string) error
	DeleteLineageByPair          func(ctx context.Context, ancestorID string, descendantID string, relationType string) error
	ListLineage                  func(ctx context.Context, assetID string) (any, error)
	ListLineageCandidates        func(ctx context.Context, assetID string, projectID string, status string, limit int) (any, error)
	ConfirmLineageCandidate      func(ctx context.Context, candidateID string) (any, error)
	RejectLineageCandidate       func(ctx context.Context, candidateID string, reason string) error
	RegisterPlugin               func(ctx context.Context, req services.PluginRegistrationRequest) (any, error)
	ListPlugins                  func(ctx context.Context) (any, error)
	ListPluginMounts             func(ctx context.Context) (any, error)
	ListPluginTaskTypes          func(ctx context.Context) (any, error)
	HeartbeatPlugin              func(ctx context.Context, req services.PluginHeartbeatRequest) (any, error)
	ResolvePluginRuntimeEndpoint func(ctx context.Context, pluginID string) (string, error)
	GetCapabilities              func(ctx context.Context) (any, error)
	ListExtensionSlots           func(ctx context.Context) (any, error)
	ListActivityLogs             func(ctx context.Context, limit int) (any, error)
	ListReplayEvents             func(ctx context.Context, sinceID int64, limit int) (any, error)
	GetMetrics                   func(ctx context.Context) (any, error)
	ListArtifacts                func(ctx context.Context, projectID string, kind string, limit int) (any, error)
	GetArtifact                  func(ctx context.Context, id string) (any, error)
	CreateArtifact               func(ctx context.Context, req services.CreateArtifactRequest) (any, error)
	UpdateArtifact               func(ctx context.Context, req services.UpdateArtifactRequest) (any, error)
	DeleteArtifact               func(ctx context.Context, id string) error
	ListPendingTasks             func(ctx context.Context, taskTypes []string) (any, error)
	GetActiveTasks               func(ctx context.Context) (any, error)
	ClaimTask                    func(ctx context.Context, taskID, workerID string) (bool, error)
	HeartbeatTask                func(ctx context.Context, taskID, workerID string) (bool, error)
	ReportTaskProgress           func(ctx context.Context, taskID string, workerID string, success bool, errMsg string, resultData map[string]any) error
	EnqueueTask                  func(ctx context.Context, assetID string, taskType string, priority int, force bool) (any, error)
	GenerateThumbnail            func(ctx context.Context, assetID string, force bool) (any, error)
	EventsWS                     func(w http.ResponseWriter, r *http.Request)
	ShowMainWindow               func() // 用于 Dock 唤醒主窗口
	QuitCore                     func() // 请求 Core 退出
	SetUIReady                   func(ready bool)
	TagService                   *services.TagService // 标签服务
	WorkflowService              *services.WorkflowService
	PublishMetricsService        *services.PublishMetricsService
}
