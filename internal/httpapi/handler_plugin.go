package httpapi

import (
	"encoding/json"
	"net/http"

	"media-assistant-os/internal/services"
)

func (h *Handler) handleRegisterPlugin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req services.PluginRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.RegisterPlugin == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.RegisterPlugin(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListPlugins(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListPlugins == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListPlugins(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListPluginMounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListPluginMounts == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListPluginMounts(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListPluginTaskTypes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListPluginTaskTypes == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListPluginTaskTypes(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleHeartbeatPlugin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.HeartbeatPlugin == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	var req services.PluginHeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	res, err := h.deps.HeartbeatPlugin(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}
