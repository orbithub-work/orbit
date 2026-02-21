package httpapi

import (
	"net/http"
	"strconv"
)

func (h *Handler) handleGetSearchHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.GetSearchHistory == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}

	limit := 20
	if v := r.URL.Query().Get("limit"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 0 {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid limit"})
			return
		}
		if n > 0 {
			limit = n
		}
	}

	res, err := h.deps.GetSearchHistory(r.Context(), limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleClearSearchHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete && r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ClearSearchHistory == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if err := h.deps.ClearSearchHistory(r.Context()); err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: map[string]any{"cleared": true}})
}
