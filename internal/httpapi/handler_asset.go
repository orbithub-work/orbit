package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"media-assistant-os/internal/services"
)

func (h *Handler) handleArchiveFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ProjectID string   `json:"project_id"`
		Paths     []string `json:"paths"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.ArchiveFiles == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.ArchiveFiles(r.Context(), req.ProjectID, req.Paths); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

func (h *Handler) handleIndexFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path      string `json:"path"`
		ProjectID string `json:"project_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.IndexFile == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.IndexFile(r.Context(), req.Path, req.ProjectID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleUpdateAssetMeta(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID              string `json:"id"`
		MediaMeta       string `json:"media_meta"`
		UserRating      *int   `json:"user_rating,omitempty"`
		ClearUserRating bool   `json:"clear_user_rating,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if strings.TrimSpace(req.ID) == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "id is required"})
		return
	}

	hasMeta := strings.TrimSpace(req.MediaMeta) != ""
	hasRating := req.UserRating != nil || req.ClearUserRating
	if !hasMeta && !hasRating {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "nothing to update"})
		return
	}

	if hasMeta && h.deps.UpdateAssetMeta == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if hasRating && h.deps.SetAssetUserRating == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "rating update not implemented"})
		return
	}

	if hasMeta {
		if err := h.deps.UpdateAssetMeta(r.Context(), req.ID, req.MediaMeta); err != nil {
			writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
			return
		}
	}

	if hasRating {
		rating := req.UserRating
		if req.ClearUserRating {
			rating = nil
		}
		if err := h.deps.SetAssetUserRating(r.Context(), req.ID, rating); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

func (h *Handler) handleImportPath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Paths []string `json:"paths"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.ImportPath == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.ImportPath(r.Context(), req.Paths); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

func (h *Handler) handleGetAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "missing id"})
		return
	}
	if h.deps.GetAsset == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetAsset(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleDeleteAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.DeleteAsset == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.DeleteAsset(r.Context(), req.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

func (h *Handler) handleListFiles(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if h.deps.ListFiles == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListFiles(r.Context(), path)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListAssets == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	q := r.URL.Query()
	cursor := strings.TrimSpace(q.Get("cursor"))
	if cursor != "" {
		if n, err := strconv.Atoi(cursor); err != nil || n < 0 {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid cursor"})
			return
		}
	}
	req := services.ListAssetsRequest{
		ProjectID: firstNonEmpty(
			strings.TrimSpace(q.Get("projectId")),
			strings.TrimSpace(q.Get("project_id")),
		),
		Directory: firstNonEmpty(
			strings.TrimSpace(q.Get("directory")),
			strings.TrimSpace(q.Get("dir")),
			strings.TrimSpace(q.Get("path")),
		),
		Query: firstNonEmpty(
			strings.TrimSpace(q.Get("search")),
			strings.TrimSpace(q.Get("q")),
			strings.TrimSpace(q.Get("keyword")),
		),
		TagIDs:      splitCSVParams(q.Get("tagIds"), q.Get("tagId"), q.Get("tags")),
		Types:       splitCSVParams(q.Get("types"), q.Get("type"), q.Get("fileType")),
		Shapes:      splitCSVParams(q.Get("shapes"), q.Get("shape")),
		SortBy:      strings.TrimSpace(q.Get("sortBy")),
		SortOrder:   strings.TrimSpace(q.Get("sortOrder")),
		Cursor:      cursor,
		QuickFilter: strings.TrimSpace(q.Get("quickFilter")),
		DatePreset:  strings.TrimSpace(q.Get("datePreset")),
	}
	if req.SortBy == "" {
		req.SortBy = "name"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	req.Limit = parseIntWithDefault(q.Get("limit"), 50)
	req.SizeMin = parseInt64WithDefault(q.Get("sizeMin"), 0)
	req.SizeMax = parseInt64WithDefault(q.Get("sizeMax"), 0)
	req.RatingMin = parseIntWithDefault(q.Get("ratingMin"), 0)
	req.RatingMax = parseIntWithDefault(q.Get("ratingMax"), 0)
	if rv := parseIntWithDefault(q.Get("rating"), 0); rv > 0 {
		req.RatingMin = rv
		req.RatingMax = rv
	}
	req.MtimeFrom = parseInt64WithDefault(firstNonEmpty(q.Get("mtimeFrom"), q.Get("from")), 0)
	req.MtimeTo = parseInt64WithDefault(firstNonEmpty(q.Get("mtimeTo"), q.Get("to")), 0)
	req.WidthMin = parseIntWithDefault(q.Get("widthMin"), 0)
	req.WidthMax = parseIntWithDefault(q.Get("widthMax"), 0)
	req.HeightMin = parseIntWithDefault(q.Get("heightMin"), 0)
	req.HeightMax = parseIntWithDefault(q.Get("heightMax"), 0)

	res, err := h.deps.ListAssets(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListAssetHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListAssetHistory == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	assetID := strings.TrimSpace(firstNonEmpty(r.URL.Query().Get("assetId"), r.URL.Query().Get("id")))
	if assetID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "assetId is required"})
		return
	}
	limit := parseIntWithDefault(r.URL.Query().Get("limit"), 100)
	res, err := h.deps.ListAssetHistory(r.Context(), assetID, limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleOpenFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path string `json:"path"`
		ID   string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.OpenFile == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.OpenFile(r.Context(), req.Path, req.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

func (h *Handler) handleOpenInFolder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Path string `json:"path"`
		ID   string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.OpenInFolder == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.OpenInFolder(r.Context(), req.Path, req.ID); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

// handleBatchDeleteAssets 批量删除资产
func (h *Handler) handleBatchDeleteAssets(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req struct {
		IDs []string `json:"ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	if len(req.IDs) == 0 {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "ids is required"})
		return
	}

	if h.deps.BatchDeleteAssets == nil {
		writeJSON(w, http.StatusNotImplemented, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	if err := h.deps.BatchDeleteAssets(r.Context(), req.IDs); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"deleted": len(req.IDs)}})
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func splitCSVParams(values ...string) []string {
	out := make([]string, 0)
	seen := make(map[string]struct{})
	for _, raw := range values {
		if strings.TrimSpace(raw) == "" {
			continue
		}
		parts := strings.Split(raw, ",")
		for _, p := range parts {
			v := strings.TrimSpace(p)
			if v == "" {
				continue
			}
			if _, ok := seen[v]; ok {
				continue
			}
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

func parseIntWithDefault(raw string, def int) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return def
	}
	v, err := strconv.Atoi(raw)
	if err != nil {
		return def
	}
	return v
}

func parseInt64WithDefault(raw string, def int64) int64 {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return def
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return def
	}
	return v
}
