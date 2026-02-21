package httpapi

import (
	"encoding/json"
	"net/http"
)

// UI Notification
type UINotificationRequest struct {
	Message  string `json:"message"`
	Type     string `json:"type"`     // success, error, warning, info
	Duration int    `json:"duration"` // milliseconds, 0 = auto
}

func (h *Handler) handleUINotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UINotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request"})
		return
	}

	if req.Message == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "message is required"})
		return
	}

	// TODO: Broadcast notification via WebSocket when EventHub is available
	// For now, just acknowledge the request
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}

// UI Dialog
type UIDialogRequest struct {
	Title   string                 `json:"title"`
	Message string                 `json:"message"`
	Type    string                 `json:"type"` // alert, confirm, prompt
	Buttons []string               `json:"buttons,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (h *Handler) handleUIDialog(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UIDialogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request"})
		return
	}

	// TODO: Broadcast dialog request via WebSocket when EventHub is available
	// For now, just acknowledge the request
	writeJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data: map[string]interface{}{
			"message": "Dialog request received",
		},
	})
}

// UI Context (current selection state)
type UIContextRequest struct {
	AssetIDs []int64  `json:"asset_ids,omitempty"`
	TagIDs   []int64  `json:"tag_ids,omitempty"`
	Path     string   `json:"path,omitempty"`
	View     string   `json:"view,omitempty"` // pool, workspace, artifact, rights
	Action   string   `json:"action,omitempty"`
	Data     map[string]interface{} `json:"data,omitempty"`
}

func (h *Handler) handleUIContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UIContextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request"})
		return
	}

	// TODO: Store context for plugins to query when EventHub is available
	// For now, just acknowledge the request
	writeJSON(w, http.StatusOK, APIResponse{Success: true})
}
