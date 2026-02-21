package httpapi

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type SystemInfo struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	Pid       int    `json:"pid"`
	Port      int    `json:"port"`
	StartTime int64  `json:"start_time"`
}

type traceIDKey struct{}

func writeJSON(w http.ResponseWriter, status int, v any) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, "json marshal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write(b)
}

func newTraceID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(b[:])
}

func listenInRange(basePort int, tryCount int) (net.Listener, int, error) {
	if basePort == 0 {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return nil, 0, err
		}
		if addr, ok := ln.Addr().(*net.TCPAddr); ok {
			return ln, addr.Port, nil
		}
		_ = ln.Close()
		return nil, 0, errors.New("unexpected listener address")
	}

	var lastErr error
	for i := 0; i < tryCount; i++ {
		port := basePort + i
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			return ln, port, nil
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = errors.New("no available port")
	}
	return nil, 0, lastErr
}
