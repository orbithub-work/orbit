package httpapi

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const pluginRuntimeRoutePrefix = "/api/plugin-runtime/"

func (h *Handler) handlePluginRuntimeProxy(w http.ResponseWriter, r *http.Request) {
	if h.deps.ResolvePluginRuntimeEndpoint == nil {
		writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "not implemented"})
		return
	}
	if strings.TrimSpace(r.Header.Get("Upgrade")) != "" {
		writeJSON(w, http.StatusNotImplemented, APIResponse{Success: false, Error: "upgrade protocol is not supported"})
		return
	}

	pluginID, subPath, err := parsePluginRuntimePath(r.URL.Path)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	endpoint, err := h.deps.ResolvePluginRuntimeEndpoint(r.Context(), pluginID)
	if err != nil {
		status := http.StatusBadRequest
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			status = http.StatusNotFound
		}
		writeJSON(w, status, APIResponse{Success: false, Error: err.Error()})
		return
	}

	baseURL, err := parseAndValidatePluginEndpoint(endpoint)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: err.Error()})
		return
	}

	targetURL := buildPluginRuntimeURL(baseURL, subPath, r.URL.RawQuery)
	proxyReq, err := http.NewRequestWithContext(r.Context(), r.Method, targetURL.String(), r.Body)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: "failed to build proxy request"})
		return
	}
	copyProxyRequestHeaders(proxyReq.Header, r.Header)
	proxyReq.Header.Set("X-Smart-Archive-Plugin-ID", pluginID)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(proxyReq)
	if err != nil {
		writeJSON(w, http.StatusBadGateway, APIResponse{Success: false, Error: "plugin endpoint is unavailable"})
		return
	}
	defer resp.Body.Close()

	copyProxyResponseHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func parsePluginRuntimePath(path string) (string, string, error) {
	if !strings.HasPrefix(path, pluginRuntimeRoutePrefix) {
		return "", "", errors.New("invalid plugin runtime path")
	}
	raw := strings.TrimPrefix(path, pluginRuntimeRoutePrefix)
	raw = strings.TrimPrefix(raw, "/")
	if raw == "" {
		return "", "", errors.New("plugin_id is required")
	}
	parts := strings.SplitN(raw, "/", 2)
	pluginID := strings.TrimSpace(parts[0])
	if !isValidPluginID(pluginID) {
		return "", "", errors.New("invalid plugin_id")
	}
	subPath := "/"
	if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
		subPath = "/" + parts[1]
	}
	return pluginID, subPath, nil
}

func isValidPluginID(pluginID string) bool {
	if pluginID == "" {
		return false
	}
	for _, ch := range pluginID {
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || (ch >= '0' && ch <= '9') {
			continue
		}
		if ch == '-' || ch == '_' || ch == '.' {
			continue
		}
		return false
	}
	return true
}

func parseAndValidatePluginEndpoint(raw string) (*url.URL, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, errors.New("plugin endpoint is empty")
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, errors.New("invalid plugin endpoint")
	}
	scheme := strings.ToLower(strings.TrimSpace(u.Scheme))
	if scheme != "http" && scheme != "https" {
		return nil, errors.New("plugin endpoint must use http or https")
	}
	host := strings.TrimSpace(u.Hostname())
	if host == "" {
		return nil, errors.New("plugin endpoint host is required")
	}
	if !isLoopbackHost(host) {
		return nil, errors.New("plugin endpoint must use loopback address")
	}
	return u, nil
}

func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func buildPluginRuntimeURL(base *url.URL, subPath string, rawQuery string) *url.URL {
	baseCopy := *base
	if !strings.HasSuffix(baseCopy.Path, "/") {
		baseCopy.Path += "/"
	}
	ref := &url.URL{Path: strings.TrimPrefix(subPath, "/")}
	target := baseCopy.ResolveReference(ref)

	values := target.Query()
	if rawQuery != "" {
		if reqValues, err := url.ParseQuery(rawQuery); err == nil {
			for key, items := range reqValues {
				for _, item := range items {
					values.Add(key, item)
				}
			}
		}
	}
	target.RawQuery = values.Encode()
	return target
}

var hopByHopHeaders = map[string]struct{}{
	"Connection":          {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"Te":                  {},
	"Trailer":             {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

func copyProxyRequestHeaders(dst http.Header, src http.Header) {
	for key, values := range src {
		if _, blocked := hopByHopHeaders[http.CanonicalHeaderKey(key)]; blocked {
			continue
		}
		// Avoid leaking core auth tokens into plugin services.
		if strings.EqualFold(key, "Authorization") {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}

func copyProxyResponseHeaders(dst http.Header, src http.Header) {
	for key, values := range src {
		if _, blocked := hopByHopHeaders[http.CanonicalHeaderKey(key)]; blocked {
			continue
		}
		// Keep CORS response under host control.
		if strings.HasPrefix(strings.ToLower(key), "access-control-") {
			continue
		}
		for _, value := range values {
			dst.Add(key, value)
		}
	}
}
