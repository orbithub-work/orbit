package httpapi

import (
	"context"
	"net/http"
	"time"
)

type Handler struct {
	deps      Deps
	metrics   *serverMetrics
	idempo    *idempotencyStore
	startTime int64
	port      int
}

func NewHandler(deps Deps, port int) *Handler {
	return &Handler{
		deps:      deps,
		metrics:   newServerMetrics(),
		idempo:    newIdempotencyStore(24 * time.Hour),
		startTime: time.Now().Unix(),
		port:      port,
	}
}

func (h *Handler) StartBackgroundTasks(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				h.idempo.purgeExpired()
			}
		}
	}()
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// System & Health
	mux.HandleFunc("/api/health", h.handleHealth)
	mux.HandleFunc("/api/metrics", h.handleMetrics)
	mux.HandleFunc("/api/system/info", h.handleSystemInfo)
	mux.HandleFunc("/api/capabilities", h.handleGetCapabilities)
	mux.HandleFunc("/api/extensions/slots", h.handleListExtensionSlots)
	mux.HandleFunc("/api/info", h.handleInfoRedirect) // Legacy

	// Internal
	mux.HandleFunc("/api/internal/show_main", h.handleInternalShowMain)
	mux.HandleFunc("/api/internal/quit", h.handleInternalQuit)
	mux.HandleFunc("/api/internal/ui/ready", h.handleInternalUIReady)

	// Projects
	mux.HandleFunc("/api/projects", h.withIdempotency(h.handleProjects))
	mux.HandleFunc("/api/projects/get", h.handleGetProject)
	mux.HandleFunc("/api/projects/update", h.withIdempotency(h.handleUpdateProject))
	mux.HandleFunc("/api/projects/delete", h.withIdempotency(h.handleDeleteProject))
	mux.HandleFunc("/api/projects/update-path", h.withIdempotency(h.handleUpdateProjectPath))
	mux.HandleFunc("/api/projects/sources", h.handleListProjectSources)
	mux.HandleFunc("/api/projects/sources/add", h.withIdempotency(h.handleAddProjectSource))
	mux.HandleFunc("/api/projects/sources/remove", h.withIdempotency(h.handleRemoveProjectSource))
	mux.HandleFunc("/api/projects/sources/bind-jobs/get", h.handleGetProjectSourceBindJob)
	mux.HandleFunc("/api/projects/directories/bound", h.handleListProjectBoundDirectories)
	mux.HandleFunc("/api/projects/directories/children", h.handleListProjectDirectoryChildren)
	mux.HandleFunc("/api/projects/directories/warnings", h.handleListProjectDirectoryWarnings)
	mux.HandleFunc("/api/projects/directories/watch/start", h.withIdempotency(h.handleStartProjectDirectoryWatch))
	mux.HandleFunc("/api/projects/directories/watch/stop", h.withIdempotency(h.handleStopProjectDirectoryWatch))
	mux.HandleFunc("/api/projects/health", h.handleGetProjectHealth)
	mux.HandleFunc("/api/projects/scan", h.withIdempotency(h.handleScanProject))
	mux.HandleFunc("/api/projects/stats", h.handleGetProjectStats)
	mux.HandleFunc("/api/projects/assets/add", h.withIdempotency(h.handleAddFileToProject))
	mux.HandleFunc("/api/projects/assets/remove", h.withIdempotency(h.handleRemoveFileFromProject))

	// Assets
	mux.HandleFunc("/api/assets/archive", h.withIdempotency(h.handleArchiveFiles))
	mux.HandleFunc("/api/assets/index", h.withIdempotency(h.handleIndexFile))
	mux.HandleFunc("/api/assets/update", h.withIdempotency(h.withAuth(h.handleUpdateAssetMeta)))
	mux.HandleFunc("/api/import", h.withIdempotency(h.handleImportPath))
	mux.HandleFunc("/api/assets/get", h.withAuth(h.handleGetAsset))
	mux.HandleFunc("/api/assets", h.handleListAssets)
	mux.HandleFunc("/api/assets/history", h.handleListAssetHistory)
	mux.HandleFunc("/api/files", h.handleListFiles)
	mux.HandleFunc("/api/assets/delete", h.withIdempotency(h.withAuth(h.handleDeleteAsset)))
	mux.HandleFunc("/api/assets/batch-delete", h.withIdempotency(h.withAuth(h.handleBatchDeleteAssets)))
	mux.HandleFunc("/api/open_file", h.handleOpenFile)
	mux.HandleFunc("/api/open_in_folder", h.handleOpenInFolder)
	mux.HandleFunc("/api/search/history", h.handleGetSearchHistory)
	mux.HandleFunc("/api/search/history/clear", h.withIdempotency(h.handleClearSearchHistory))

	// Lineage
	mux.HandleFunc("/api/lineage/create", h.withIdempotency(h.handleCreateLineage))
	mux.HandleFunc("/api/lineage/update", h.withIdempotency(h.handleUpdateLineage))
	mux.HandleFunc("/api/lineage/delete", h.withIdempotency(h.handleDeleteLineage))
	mux.HandleFunc("/api/lineage/list", h.handleListLineage)
	if h.deps.EnableProFeatures {
		mux.HandleFunc("/api/lineage/candidates/list", h.handleListLineageCandidates)
		mux.HandleFunc("/api/lineage/candidates/confirm", h.withIdempotency(h.handleConfirmLineageCandidate))
		mux.HandleFunc("/api/lineage/candidates/reject", h.withIdempotency(h.handleRejectLineageCandidate))
	}

	// Plugins
	mux.HandleFunc("/api/plugins", h.handleListPlugins) // Alias for frontend plugin marketplace pages.
	mux.HandleFunc("/api/plugins/register", h.withIdempotency(h.handleRegisterPlugin))
	mux.HandleFunc("/api/plugins/list", h.handleListPlugins)
	mux.HandleFunc("/api/plugins/mounts", h.handleListPluginMounts)
	mux.HandleFunc("/api/plugins/task-types", h.handleListPluginTaskTypes)
	mux.HandleFunc("/api/plugins/heartbeat", h.withIdempotency(h.handleHeartbeatPlugin))
	mux.HandleFunc("/api/plugin-runtime/", h.handlePluginRuntimeProxy)

	// Activity & Events
	mux.HandleFunc("/api/activity/list", h.handleListActivityLogs)
	mux.HandleFunc("/api/events/replay", h.handleListReplayEvents)
	mux.HandleFunc("/events", h.handleEventsWS)

	if h.deps.EnableProFeatures {
		// Artifacts
		mux.HandleFunc("/api/artifacts/list", h.handleListArtifacts)
		mux.HandleFunc("/api/artifacts/get", h.handleGetArtifact)
		mux.HandleFunc("/api/artifacts/create", h.withIdempotency(h.handleCreateArtifact))
		mux.HandleFunc("/api/artifacts/update", h.withIdempotency(h.handleUpdateArtifact))
		mux.HandleFunc("/api/artifacts/delete", h.withIdempotency(h.handleDeleteArtifact))

		// Workflow Templates & Project Workflow
		mux.HandleFunc("/api/workflow/templates/list", h.handleListWorkflowTemplates)
		mux.HandleFunc("/api/workflow/templates/create", h.withIdempotency(h.handleCreateWorkflowTemplate))
		mux.HandleFunc("/api/workflow/templates/steps/list", h.handleListWorkflowTemplateSteps)
		mux.HandleFunc("/api/workflow/templates/steps/create", h.withIdempotency(h.handleCreateWorkflowTemplateStep))
		mux.HandleFunc("/api/workflow/project/init", h.withIdempotency(h.handleInitProjectWorkflow))
		mux.HandleFunc("/api/workflow/project/snapshot", h.handleGetProjectWorkflowSnapshot)
		mux.HandleFunc("/api/workflow/project/steps/list", h.handleListProjectWorkflowSteps)
		mux.HandleFunc("/api/workflow/project/steps/evaluate-kpi", h.withIdempotency(h.handleEvaluateWorkflowStepKPI))
		mux.HandleFunc("/api/workflow/roadmap/list", h.handleListRoadmapItems)
		mux.HandleFunc("/api/workflow/roadmap/create", h.withIdempotency(h.handleCreateRoadmapItem))
		mux.HandleFunc("/api/workflow/roadmap/update", h.withIdempotency(h.handleUpdateRoadmapItem))
		mux.HandleFunc("/api/workflow/notes/list", h.handleListProjectNotes)
		mux.HandleFunc("/api/workflow/notes/create", h.withIdempotency(h.handleCreateProjectNote))
		mux.HandleFunc("/api/workflow/notes/update", h.withIdempotency(h.handleUpdateProjectNote))

		// Publish & Metrics
		mux.HandleFunc("/api/publish/channels/list", h.handleListPublishChannels)
		mux.HandleFunc("/api/publish/channels/upsert", h.withIdempotency(h.handleUpsertPublishChannel))
		mux.HandleFunc("/api/publish/jobs/list", h.handleListPublishJobs)
		mux.HandleFunc("/api/publish/jobs/create", h.withIdempotency(h.handleCreatePublishJob))
		mux.HandleFunc("/api/publish/jobs/update-status", h.withIdempotency(h.handleUpdatePublishJobStatus))
		mux.HandleFunc("/api/publish/records/list", h.handleListPublishRecords)
		mux.HandleFunc("/api/publish/records/create", h.withIdempotency(h.handleCreatePublishRecord))
		mux.HandleFunc("/api/metrics/events/ingest", h.withIdempotency(h.handleIngestMetricsEvent))
		mux.HandleFunc("/api/metrics/snapshots/daily/list", h.handleListMetricsDailySnapshots)
		mux.HandleFunc("/api/metrics/snapshots/daily/upsert", h.withIdempotency(h.handleUpsertMetricsDailySnapshot))
	}

	// Tasks
	mux.HandleFunc("/api/tasks/pending", h.handleListPendingTasks)
	mux.HandleFunc("/api/tasks/active", h.handleGetActiveTasks)
	mux.HandleFunc("/api/tasks/enqueue", h.withIdempotency(h.handleEnqueueTask))
	mux.HandleFunc("/api/tasks/claim", h.withIdempotency(h.handleClaimTask))
	mux.HandleFunc("/api/tasks/heartbeat", h.withIdempotency(h.handleHeartbeatTask))
	mux.HandleFunc("/api/tasks/report", h.withIdempotency(h.handleReportTaskProgress))

	// Thumbnails
	mux.HandleFunc("/api/thumbnails/", h.handleGetThumbnail)
	mux.HandleFunc("/api/thumbnails/generate", h.handleGenerateThumbnail)

	// Asset File Serving
	mux.HandleFunc("/api/assets/file", h.handleServeAssetFile)

	// Tags
	mux.HandleFunc("/api/tags", h.handleListTags)
	mux.HandleFunc("/api/tags/tree", h.handleListTagTree)
	mux.HandleFunc("/api/tags/create", h.withIdempotency(h.handleCreateTag))
	mux.HandleFunc("/api/tags/update", h.withIdempotency(h.handleUpdateTag))
	mux.HandleFunc("/api/tags/delete", h.withIdempotency(h.handleDeleteTag))
	mux.HandleFunc("/api/tags/search", h.handleSearchTags)
	mux.HandleFunc("/api/tags/file", h.handleGetFileTags)
	mux.HandleFunc("/api/tags/batch-add", h.withIdempotency(h.handleBatchAddTags))
	mux.HandleFunc("/api/tags/batch-remove", h.withIdempotency(h.handleBatchRemoveTags))

	// Onboarding & Initial Import
	mux.HandleFunc("/api/start_initial_import", h.withIdempotency(h.handleStartInitialImport))
	mux.HandleFunc("/api/get_import_progress", h.handleGetImportProgress)
	mux.HandleFunc("/api/system/first-launch", h.handleIsFirstLaunch)
	mux.HandleFunc("/api/system/first-launch/complete", h.withIdempotency(h.handleCompleteFirstLaunch))
	mux.HandleFunc("/api/get_common_directories", h.handleGetCommonDirectories)
	mux.HandleFunc("/api/system/default-dirs", h.handleGetDefaultDirs)
	mux.HandleFunc("/api/system/initial-import", h.withIdempotency(h.handleStartInitialImport))
	mux.HandleFunc("/api/system/import-progress", h.handleGetImportProgress)
	mux.HandleFunc("/api/run_onboarding", h.withIdempotency(h.handleRunOnboarding))

	// Library Sources (Physical Directories)
	mux.HandleFunc("/api/library/sources", h.handleListLibrarySources)
	mux.HandleFunc("/api/library/sources/add", h.withIdempotency(h.handleAddLibrarySource))
	mux.HandleFunc("/api/library/sources/remove", h.withIdempotency(h.handleRemoveLibrarySource))
	mux.HandleFunc("/api/library/directories/children", h.handleListLibraryDirectoryChildren)

	// UI Interaction APIs
	mux.HandleFunc("/api/ui/notification", h.withIdempotency(h.handleUINotification))
	mux.HandleFunc("/api/ui/dialog", h.handleUIDialog)
	mux.HandleFunc("/api/ui/context", h.withIdempotency(h.handleUIContext))

	// API Documentation (Scalar)
	mux.HandleFunc("/docs/api/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "docs/api/openapi.yaml")
	})
	mux.HandleFunc("/docs/api", func(w http.ResponseWriter, r *http.Request) {
		html := `
<!doctype html>
<html>
  <head>
    <title>Media Assistant API</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
    <style>
      body {
        margin: 0;
      }
    </style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/docs/api/openapi.yaml"></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>
`
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})
}

func (h *Handler) handleGetDefaultDirs(w http.ResponseWriter, r *http.Request) {
	if h.deps.GetDefaultDirs == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	dirs := h.deps.GetDefaultDirs()
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: dirs})
}
