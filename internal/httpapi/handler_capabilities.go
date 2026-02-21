package httpapi

import "net/http"

func (h *Handler) handleGetCapabilities(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.GetCapabilities == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.GetCapabilities(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}

func (h *Handler) handleListExtensionSlots(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if h.deps.ListExtensionSlots == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	res, err := h.deps.ListExtensionSlots(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, APIResponse{Success: true, Data: res})
}
