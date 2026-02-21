package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// CreateTagRequest 创建标签请求
type CreateTagRequest struct {
	Name     string  `json:"name"`
	Color    *string `json:"color"`
	Icon     *string `json:"icon"`
	ParentID *string `json:"parent_id"`
}

// UpdateTagRequest 更新标签请求
type UpdateTagRequest struct {
	ID       string  `json:"id"`
	Name     *string `json:"name"`
	Color    *string `json:"color"`
	Icon     *string `json:"icon"`
	ParentID *string `json:"parent_id"`
}

// BatchTagRequest 批量标签请求
type BatchTagRequest struct {
	FileIDs []string `json:"file_ids"`
	TagIDs  []string `json:"tag_ids"`
}

// handleListTags 获取所有标签
func (h *Handler) handleListTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	tags, err := h.deps.TagService.ListTags(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tags})
}

// handleListTagTree 获取树形标签
func (h *Handler) handleListTagTree(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	tags, err := h.deps.TagService.ListTagTree(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tags})
}

// handleCreateTag 创建标签
func (h *Handler) handleCreateTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req CreateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	tag, err := h.deps.TagService.CreateTag(r.Context(), req.Name, req.Color, req.Icon, req.ParentID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tag})
}

// handleUpdateTag 更新标签
func (h *Handler) handleUpdateTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req UpdateTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	tag, err := h.deps.TagService.UpdateTag(r.Context(), req.ID, req.Name, req.Color, req.Icon, req.ParentID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tag})
}

// handleDeleteTag 删除标签
func (h *Handler) handleDeleteTag(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	tagID := r.URL.Query().Get("id")
	if tagID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "tag id is required"})
		return
	}

	if err := h.deps.TagService.DeleteTag(r.Context(), tagID); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: nil})
}

// handleSearchTags 搜索标签
func (h *Handler) handleSearchTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	query := r.URL.Query().Get("query")
	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	tags, err := h.deps.TagService.SearchTags(r.Context(), query, limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tags})
}

// handleGetFileTags 获取文件的标签
func (h *Handler) handleGetFileTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "file_id is required"})
		return
	}

	tags, err := h.deps.TagService.GetFileTags(r.Context(), fileID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: tags})
}

// handleBatchAddTags 批量添加标签
func (h *Handler) handleBatchAddTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req BatchTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	if err := h.deps.TagService.AddTagsToFiles(r.Context(), req.FileIDs, req.TagIDs); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: nil})
}

// handleBatchRemoveTags 批量移除标签
func (h *Handler) handleBatchRemoveTags(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	var req BatchTagRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
		return
	}

	if err := h.deps.TagService.RemoveTagsFromFiles(r.Context(), req.FileIDs, req.TagIDs); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: nil})
}
