package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *Handler) handleListPendingTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	taskTypes := r.URL.Query()["type"]
	if h.deps.ListPendingTasks == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListPendingTasks(r.Context(), taskTypes)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleGetActiveTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.GetActiveTasks == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetActiveTasks(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleEnqueueTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.EnqueueTask == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	var req struct {
		AssetID  string `json:"asset_id"`
		TaskType string `json:"task_type"`
		Priority int    `json:"priority"`
		Force    bool   `json:"force"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if strings.TrimSpace(req.AssetID) == "" || strings.TrimSpace(req.TaskType) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "asset_id and task_type are required"})
		return
	}
	res, err := h.deps.EnqueueTask(r.Context(), req.AssetID, req.TaskType, req.Priority, req.Force)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleClaimTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		TaskID   string `json:"task_id"`
		WorkerID string `json:"worker_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.ClaimTask == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if strings.TrimSpace(req.TaskID) == "" || strings.TrimSpace(req.WorkerID) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "task_id and worker_id are required"})
		return
	}
	success, err := h.deps.ClaimTask(r.Context(), req.TaskID, req.WorkerID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"success": success}})
}

func (h *Handler) handleReportTaskProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		TaskID     string         `json:"task_id"`
		WorkerID   string         `json:"worker_id"`
		Success    bool           `json:"success"`
		Error      string         `json:"error"`
		ResultData map[string]any `json:"result_data"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.ReportTaskProgress == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if strings.TrimSpace(req.TaskID) == "" || strings.TrimSpace(req.WorkerID) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "task_id and worker_id are required"})
		return
	}
	err := h.deps.ReportTaskProgress(r.Context(), req.TaskID, req.WorkerID, req.Success, req.Error, req.ResultData)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

func (h *Handler) handleHeartbeatTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		TaskID   string `json:"task_id"`
		WorkerID string `json:"worker_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.HeartbeatTask == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if strings.TrimSpace(req.TaskID) == "" || strings.TrimSpace(req.WorkerID) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "task_id and worker_id are required"})
		return
	}
	ok, err := h.deps.HeartbeatTask(r.Context(), req.TaskID, req.WorkerID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"success": ok}})
}
