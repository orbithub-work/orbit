package httpapi

import (
	"encoding/json"
	"net/http"

	"media-assistant-os/internal/services"
)

func (h *Handler) handleStartInitialImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Directories    []string `json:"directories"`
		QuickScanLimit int      `json:"quickScanLimit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.StartInitialImport == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.StartInitialImport(r.Context(), req.Directories, req.QuickScanLimit); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"started": true}})
}

func (h *Handler) handleGetImportProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.GetImportProgress == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetImportProgress(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleIsFirstLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.IsFirstLaunch == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.IsFirstLaunch(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleCompleteFirstLaunch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.CompleteFirstLaunch == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.CompleteFirstLaunch(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleGetCommonDirectories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.GetCommonDirectories == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetCommonDirectories(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleRunOnboarding(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req services.OnboardingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.RunOnboarding == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.RunOnboarding(r.Context(), req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"started": true}})
}
