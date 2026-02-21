package httpapi

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) handleListLibrarySources(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListLibrarySources == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListLibrarySources(r.Context())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleAddLibrarySource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		RootPath     string `json:"rootPath"`
		WatchEnabled *bool  `json:"watchEnabled,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if strings.TrimSpace(req.RootPath) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "rootPath is required"})
		return
	}
	if h.deps.AddLibrarySource == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.AddLibrarySource(r.Context(), req.RootPath, req.WatchEnabled)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleRemoveLibrarySource(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		RootPath string `json:"rootPath"`
	}
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
			return
		}
	}
	if strings.TrimSpace(req.RootPath) == "" {
		req.RootPath = strings.TrimSpace(r.URL.Query().Get("rootPath"))
	}
	if strings.TrimSpace(req.RootPath) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "rootPath is required"})
		return
	}
	if h.deps.RemoveLibrarySource == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.RemoveLibrarySource(r.Context(), req.RootPath); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"rootPath": req.RootPath}})
}

func (h *Handler) handleListLibraryDirectoryChildren(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimSpace(r.URL.Query().Get("path"))
	if path == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "path is required"})
		return
	}

	// 直接读取物理目录
	entries, err := os.ReadDir(path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	var children []map[string]any
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		// 跳过隐藏目录
		if strings.HasPrefix(name, ".") {
			continue
		}

		childPath := filepath.Join(path, name)
		children = append(children, map[string]any{
			"path":         childPath,
			"name":         name,
			"has_children": true,
		})
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: children})
}
