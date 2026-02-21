package httpapi

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

type serverMetrics struct {
	mu             sync.Mutex
	startUnix      int64
	requestTotal   int64
	request4xx     int64
	request5xx     int64
	inflight       int64
	latencyTotalMs int64
	byPath         map[string]int64
}

func newServerMetrics() *serverMetrics {
	return &serverMetrics{
		startUnix: time.Now().Unix(),
		byPath:    map[string]int64{},
	}
}

func (m *serverMetrics) begin(path string) func(statusCode int, elapsed time.Duration) {
	m.mu.Lock()
	m.inflight++
	m.requestTotal++
	m.byPath[path]++
	m.mu.Unlock()
	return func(statusCode int, elapsed time.Duration) {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.inflight--
		if statusCode >= 500 {
			m.request5xx++
		} else if statusCode >= 400 {
			m.request4xx++
		}
		m.latencyTotalMs += elapsed.Milliseconds()
	}
}

func (m *serverMetrics) snapshot() map[string]any {
	m.mu.Lock()
	defer m.mu.Unlock()
	avgLatency := float64(0)
	if m.requestTotal > 0 {
		avgLatency = float64(m.latencyTotalMs) / float64(m.requestTotal)
	}
	pathCopy := make(map[string]int64, len(m.byPath))
	for k, v := range m.byPath {
		pathCopy[k] = v
	}
	return map[string]any{
		"uptime_sec":          time.Now().Unix() - m.startUnix,
		"http_requests_total": m.requestTotal,
		"http_4xx_total":      m.request4xx,
		"http_5xx_total":      m.request5xx,
		"http_inflight":       m.inflight,
		"http_avg_latency_ms": avgLatency,
		"http_requests_by":    pathCopy,
	}
}

type responseCapture struct {
	http.ResponseWriter
	statusCode int
}

func (w *responseCapture) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseCapture) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

func (w *responseCapture) Flush() {
	if f, ok := w.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

func (w *responseCapture) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijacker not supported")
	}
	return h.Hijack()
}

func (w *responseCapture) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.ResponseWriter.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

type recordWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func (w *recordWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *recordWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

type idempotencyEntry struct {
	requestHash string
	record      responseRecord
	expiresAt   time.Time
}

type responseRecord struct {
	Status int
	Body   []byte
}

type idempotencyStore struct {
	mu    sync.Mutex
	items map[string]idempotencyEntry
	ttl   time.Duration
}

func newIdempotencyStore(ttl time.Duration) *idempotencyStore {
	if ttl <= 0 {
		ttl = 24 * time.Hour
	}
	return &idempotencyStore{
		items: make(map[string]idempotencyEntry, 1024),
		ttl:   ttl,
	}
}

func (s *idempotencyStore) load(storeKey string) (idempotencyEntry, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry, ok := s.items[storeKey]
	if !ok {
		return idempotencyEntry{}, false
	}
	if time.Now().After(entry.expiresAt) {
		delete(s.items, storeKey)
		return idempotencyEntry{}, false
	}
	return entry, true
}

func (s *idempotencyStore) save(storeKey string, entry idempotencyEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()
	entry.expiresAt = time.Now().Add(s.ttl)
	s.items[storeKey] = entry
}

func (s *idempotencyStore) purgeExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, v := range s.items {
		if now.After(v.expiresAt) {
			delete(s.items, k)
		}
	}
}

func (h *Handler) withIdempotency(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			next(w, r)
			return
		}
		key := strings.TrimSpace(r.Header.Get("Idempotency-Key"))
		if key == "" {
			next(w, r)
			return
		}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid request body"})
			return
		}
		_ = r.Body.Close()
		r.Body = io.NopCloser(bytes.NewReader(body))

		sum := sha256.Sum256(append([]byte(r.Method+"|"+r.URL.Path+"|"), body...))
		reqHash := hex.EncodeToString(sum[:])
		storeKey := r.Method + "|" + r.URL.Path + "|" + key

		if entry, ok := h.idempo.load(storeKey); ok {
			if entry.requestHash != reqHash {
				writeJSON(w, http.StatusConflict, APIResponse{Success: false, Error: "idempotency key reused with different payload"})
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(entry.record.Status)
			_, _ = w.Write(entry.record.Body)
			return
		}

		rw := &recordWriter{ResponseWriter: w}
		next(rw, r)
		if rw.statusCode == 0 {
			rw.statusCode = http.StatusOK
		}
		h.idempo.save(storeKey, idempotencyEntry{
			requestHash: reqHash,
			record: responseRecord{
				Status: rw.statusCode,
				Body:   rw.body.Bytes(),
			},
		})
	}
}

func (h *Handler) withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if h.deps.ValidateToken == nil {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("Authorization")
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		if !h.deps.ValidateToken(token) {
			writeJSON(w, http.StatusUnauthorized, APIResponse{Success: false, Error: "unauthorized: invalid or missing token"})
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (h *Handler) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		begin := h.metrics.begin(r.URL.Path)
		start := time.Now()
		capture := &responseCapture{ResponseWriter: w}

		traceID := strings.TrimSpace(r.Header.Get("X-Trace-Id"))
		if traceID == "" {
			traceID = newTraceID()
		}
		capture.Header().Set("Access-Control-Allow-Origin", "*")
		capture.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		capture.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Idempotency-Key, X-Trace-Id")
		capture.Header().Set("X-Trace-Id", traceID)

		if r.Method == http.MethodOptions {
			capture.WriteHeader(http.StatusOK)
			begin(capture.statusCode, time.Since(start))
			return
		}

		ctxWithTrace := context.WithValue(r.Context(), traceIDKey{}, traceID)
		next.ServeHTTP(capture, r.WithContext(ctxWithTrace))
		if capture.statusCode == 0 {
			capture.statusCode = http.StatusOK
		}
		begin(capture.statusCode, time.Since(start))
	})
}
