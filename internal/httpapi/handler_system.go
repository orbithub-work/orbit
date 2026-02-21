package httpapi

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"
)

func (h *Handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok"})
}

func (h *Handler) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	payload := map[string]any{
		"http": h.metrics.snapshot(),
	}
	if h.deps.GetMetrics != nil {
		extra, err := h.deps.GetMetrics(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
			return
		}
		payload["runtime"] = extra
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: payload})
}

func (h *Handler) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, SystemInfo{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Pid:       os.Getpid(),
		Port:      h.port,
		StartTime: h.startTime,
	})
}

func (h *Handler) handleInfoRedirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/api/system/info", http.StatusTemporaryRedirect)
}

func (h *Handler) handleInternalShowMain(w http.ResponseWriter, r *http.Request) {
	if h.deps.ShowMainWindow != nil {
		h.deps.ShowMainWindow()
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
	} else {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "handler not defined"})
	}
}

func (h *Handler) handleInternalQuit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.QuitCore != nil {
		go h.deps.QuitCore()
		writeJSON(w, http.StatusOK, map[string]any{"success": true})
		return
	}
	writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "handler not defined"})
}

func (h *Handler) handleInternalUIReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.SetUIReady == nil {
		writeJSON(w, http.StatusInternalServerError, map[string]any{"error": "handler not defined"})
		return
	}
	var req struct {
		Ready bool `json:"ready"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	h.deps.SetUIReady(req.Ready)
	writeJSON(w, http.StatusOK, map[string]any{"success": true})
}
