package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func (h *Handler) handleCreateLineage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		AncestorID   string `json:"ancestor_id"`
		DescendantID string `json:"descendant_id"`
		RelationType string `json:"relation_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.CreateLineage == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.CreateLineage(r.Context(), req.AncestorID, req.DescendantID, req.RelationType)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleUpdateLineage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID           string `json:"id"`
		AncestorID   string `json:"ancestor_id"`
		DescendantID string `json:"descendant_id"`
		RelationType string `json:"relation_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if h.deps.UpdateLineage == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.UpdateLineage(r.Context(), req.ID, req.AncestorID, req.DescendantID, req.RelationType); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"updated": true}})
}

func (h *Handler) handleDeleteLineage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		ID           string `json:"id"`
		AncestorID   string `json:"ancestor_id"`
		DescendantID string `json:"descendant_id"`
		RelationType string `json:"relation_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if req.ID != "" {
		if h.deps.DeleteLineage == nil {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
			return
		}
		if err := h.deps.DeleteLineage(r.Context(), req.ID); err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"deleted": true}})
		return
	}
	if h.deps.DeleteLineageByPair == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.DeleteLineageByPair(r.Context(), req.AncestorID, req.DescendantID, req.RelationType); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"deleted": true}})
}

func (h *Handler) handleListLineage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	assetID := r.URL.Query().Get("assetId")
	if assetID == "" {
		var body struct {
			AssetID string `json:"asset_id"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		assetID = body.AssetID
	}
	if assetID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "assetId is required"})
		return
	}
	if h.deps.ListLineage == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListLineage(r.Context(), assetID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListLineageCandidates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListLineageCandidates == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	assetID := r.URL.Query().Get("assetId")
	projectID := r.URL.Query().Get("projectId")
	status := r.URL.Query().Get("status")
	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if assetID == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "assetId is required"})
		return
	}
	res, err := h.deps.ListLineageCandidates(r.Context(), assetID, projectID, status, limit)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleConfirmLineageCandidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ConfirmLineageCandidate == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	var req struct {
		CandidateID string `json:"candidate_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	res, err := h.deps.ConfirmLineageCandidate(r.Context(), req.CandidateID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleRejectLineageCandidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.RejectLineageCandidate == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	var req struct {
		CandidateID string `json:"candidate_id"`
		Reason      string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid json"})
		return
	}
	if err := h.deps.RejectLineageCandidate(r.Context(), req.CandidateID, req.Reason); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"updated": true}})
}
