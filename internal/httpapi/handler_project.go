package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) handleProjects(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if h.deps.ListProjects == nil {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
			return
		}
		res, err := h.deps.ListProjects(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
	case http.MethodPost:
		if h.deps.CreateProject == nil {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
			return
		}
		var req struct {
			Name        string `json:"name"`
			ProjectType string `json:"project_type"`
			Path        string `json:"path,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
			return
		}
		res, err := h.deps.CreateProject(r.Context(), req.Name, req.ProjectType, req.Path)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) handleGetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("id")
	if projectID == "" {
		projectID = r.URL.Query().Get("projectId")
	}
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id is required"})
		return
	}
	if h.deps.GetProject == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetProject(r.Context(), projectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleUpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		ProjectType string `json:"project_type"`
		Status      string `json:"status"`
		Description string `json:"description"`
		Path        string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id is required"})
		return
	}
	if h.deps.UpdateProject == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.UpdateProject(r.Context(), req.ID, req.Name, req.ProjectType, req.Status, req.Description, req.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleDeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("id")
	if projectID == "" {
		var req struct {
			ID string `json:"id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&req)
		projectID = req.ID
	}
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id is required"})
		return
	}
	if h.deps.DeleteProject == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.DeleteProject(r.Context(), projectID); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"id": projectID}})
}

func (h *Handler) handleUpdateProjectPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID   string `json:"id"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ID == "" || req.Path == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id and path are required"})
		return
	}
	if h.deps.UpdateProjectPath == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.UpdateProjectPath(r.Context(), req.ID, req.Path); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"id": req.ID, "path": req.Path}})
}

func (h *Handler) handleListProjectSources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.URL.Query().Get("id")
	}
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId is required"})
		return
	}
	if h.deps.ListProjectSources == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListProjectSources(r.Context(), projectID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleAddProjectSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProjectID    string `json:"projectId"`
		RootPath     string `json:"rootPath"`
		SourceType   string `json:"sourceType"`
		WatchEnabled *bool  `json:"watchEnabled,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ProjectID == "" || req.RootPath == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and rootPath are required"})
		return
	}
	if h.deps.AddProjectSource == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.AddProjectSource(r.Context(), req.ProjectID, req.RootPath, req.SourceType, req.WatchEnabled)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	// Async bind jobs should be explicitly visible to callers.
	status := http.StatusOK
	if isAsyncSourceBindResponse(res) {
		status = http.StatusAccepted
	}
	writeJSON(w, status, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleRemoveProjectSource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProjectID string `json:"projectId"`
		RootPath  string `json:"rootPath"`
	}
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
			return
		}
	}
	if req.ProjectID == "" {
		req.ProjectID = r.URL.Query().Get("projectId")
	}
	if req.RootPath == "" {
		req.RootPath = r.URL.Query().Get("rootPath")
	}
	if req.ProjectID == "" || req.RootPath == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and rootPath are required"})
		return
	}
	if h.deps.RemoveProjectSource == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.RemoveProjectSource(r.Context(), req.ProjectID, req.RootPath); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{
		"projectId": req.ProjectID,
		"rootPath":  req.RootPath,
	}})
}

func (h *Handler) handleGetProjectSourceBindJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	jobID := strings.TrimSpace(r.URL.Query().Get("jobId"))
	if jobID == "" {
		jobID = strings.TrimSpace(r.URL.Query().Get("id"))
	}
	if jobID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "jobId is required"})
		return
	}
	if h.deps.GetProjectSourceBindJob == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetProjectSourceBindJob(r.Context(), jobID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func isAsyncSourceBindResponse(res any) bool {
	if res == nil {
		return false
	}
	var payload struct {
		BindMode string `json:"bind_mode"`
	}
	raw, err := json.Marshal(res)
	if err != nil {
		return false
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return false
	}
	return strings.EqualFold(strings.TrimSpace(payload.BindMode), "async")
}

func (h *Handler) handleListProjectBoundDirectories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.URL.Query().Get("id")
	}
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId is required"})
		return
	}
	if h.deps.ListProjectBoundDirectories == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListProjectBoundDirectories(r.Context(), projectID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListProjectDirectoryChildren(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.URL.Query().Get("id")
	}
	path := r.URL.Query().Get("path")
	if projectID == "" || path == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and path are required"})
		return
	}
	if h.deps.ListProjectDirectoryChildren == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListProjectDirectoryChildren(r.Context(), projectID, path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListProjectDirectoryWarnings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.URL.Query().Get("id")
	}
	path := r.URL.Query().Get("path")
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId is required"})
		return
	}
	if h.deps.ListProjectDirectoryWarnings == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListProjectDirectoryWarnings(r.Context(), projectID, path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleStartProjectDirectoryWatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.StartProjectDirectoryWatch == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	var req struct {
		ProjectID  string `json:"projectId"`
		Path       string `json:"path"`
		TTLSeconds int    `json:"ttlSeconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if strings.TrimSpace(req.ProjectID) == "" || strings.TrimSpace(req.Path) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and path are required"})
		return
	}
	res, err := h.deps.StartProjectDirectoryWatch(r.Context(), req.ProjectID, req.Path, req.TTLSeconds)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleStopProjectDirectoryWatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.StopProjectDirectoryWatch == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	var req struct {
		ProjectID string `json:"projectId"`
		Path      string `json:"path"`
	}
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
			return
		}
	}
	if strings.TrimSpace(req.ProjectID) == "" {
		req.ProjectID = strings.TrimSpace(r.URL.Query().Get("projectId"))
	}
	if strings.TrimSpace(req.Path) == "" {
		req.Path = strings.TrimSpace(r.URL.Query().Get("path"))
	}
	if strings.TrimSpace(req.ProjectID) == "" || strings.TrimSpace(req.Path) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and path are required"})
		return
	}
	if err := h.deps.StopProjectDirectoryWatch(r.Context(), req.ProjectID, req.Path); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{
		"projectId": req.ProjectID,
		"path":      req.Path,
		"stopped":   true,
	}})
}

func (h *Handler) handleGetProjectHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		projectID = r.URL.Query().Get("id")
	}
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId is required"})
		return
	}
	if h.deps.GetProjectHealth == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetProjectHealth(r.Context(), projectID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleScanProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID   string `json:"id"`
		Path string `json:"path"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id is required"})
		return
	}
	if h.deps.ScanProject == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.ScanProject(r.Context(), req.ID, req.Path); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"started": true, "id": req.ID}})
}

func (h *Handler) handleGetProjectStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	projectID := r.URL.Query().Get("id")
	if projectID == "" {
		projectID = r.URL.Query().Get("projectId")
	}
	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id is required"})
		return
	}
	if h.deps.GetProjectStats == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetProjectStats(r.Context(), projectID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleAddFileToProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProjectID string `json:"projectId"`
		FileID    string `json:"fileId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ProjectID == "" || req.FileID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and fileId are required"})
		return
	}
	if h.deps.AddFileToProject == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.AddFileToProject(r.Context(), req.ProjectID, req.FileID); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"projectId": req.ProjectID, "fileId": req.FileID}})
}

func (h *Handler) handleRemoveFileFromProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProjectID string `json:"projectId"`
		FileID    string `json:"fileId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ProjectID == "" || req.FileID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId and fileId are required"})
		return
	}
	if h.deps.RemoveFileFromProject == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.RemoveFileFromProject(r.Context(), req.ProjectID, req.FileID); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"projectId": req.ProjectID, "fileId": req.FileID}})
}

func (h *Handler) handleGetProjectFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	projectID := r.URL.Query().Get("projectId")
	if projectID == "" {
		// Try POST body if not in query
		var body struct {
			ProjectID string `json:"projectId"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		projectID = body.ProjectID
	}

	if projectID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "projectId is required"})
		return
	}

	if h.deps.GetProjectFiles == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	res, err := h.deps.GetProjectFiles(r.Context(), projectID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}
