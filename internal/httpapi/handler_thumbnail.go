package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"media-assistant-os/internal/infra"
	"media-assistant-os/internal/models"
)

func (h *Handler) handleGetThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/thumbnails/")
	id = strings.TrimSpace(id)
	if id == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "missing id"})
		return
	}

	if h.deps.GetAsset == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	assetAny, err := h.deps.GetAsset(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: err.Error()})
		return
	}

	asset, ok := assetAny.(*models.Asset)
	if !ok || asset == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "asset not found"})
		return
	}

	dataDir, err := infra.ResolveDataDir()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "failed to resolve data dir"})
		return
	}

	thumbPath := extractThumbnailPath(asset.MediaMeta)
	if thumbPath != "" && !filepath.IsAbs(thumbPath) {
		thumbPath = filepath.Join(dataDir, "..", thumbPath)
	}
	if thumbPath == "" {
		thumbPath = filepath.Join(dataDir, "cache", "thumbnails", asset.ID+".jpg")
	}

	data, err := os.ReadFile(thumbPath)
	if err != nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "thumbnail not found"})
		return
	}

	encoded := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(data)
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: encoded})
}

func (h *Handler) handleGenerateThumbnail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		ID    string `json:"id"`
		Force bool   `json:"force"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if strings.TrimSpace(req.ID) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "missing id"})
		return
	}

	if h.deps.GenerateThumbnail == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	if h.deps.GetAsset != nil {
		if assetAny, err := h.deps.GetAsset(r.Context(), req.ID); err != nil || assetAny == nil {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "asset not found"})
			return
		}
	}

	res, err := h.deps.GenerateThumbnail(r.Context(), req.ID, req.Force)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func extractThumbnailPath(mediaMeta string) string {
	if strings.TrimSpace(mediaMeta) == "" {
		return ""
	}
	var meta map[string]any
	if err := json.Unmarshal([]byte(mediaMeta), &meta); err != nil {
		return ""
	}
	if v, ok := meta["thumbnail_path"].(string); ok {
		return v
	}
	return ""
}
