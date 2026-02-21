package services

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"media-assistant-os/internal/repos"

	"github.com/gorilla/websocket"
)

func NewEventHub(eventRepo *repos.EventLogRepo) *EventHub {
	return &EventHub{
		conns:     map[*websocket.Conn]struct{}{},
		upgrader:  defaultUpgrader,
		eventRepo: eventRepo,
	}
}

func (h *EventHub) ServeWS(w http.ResponseWriter, r *http.Request) {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	h.mu.Lock()
	h.conns[c] = struct{}{}
	h.mu.Unlock()

	for {
		if _, _, err := c.ReadMessage(); err != nil {
			break
		}
	}

	h.mu.Lock()
	delete(h.conns, c)
	h.mu.Unlock()
	_ = c.Close()
}

func (h *EventHub) Broadcast(event any) {
	b, err := json.Marshal(event)
	if err != nil {
		return
	}

	eventType := "unknown"
	var decoded map[string]any
	if err := json.Unmarshal(b, &decoded); err == nil {
		if t, ok := decoded["type"].(string); ok && t != "" {
			eventType = t
		}
	}
	if h.eventRepo != nil {
		_, _ = h.eventRepo.Append(context.Background(), eventType, string(b))
	}

	h.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(h.conns))
	for c := range h.conns {
		conns = append(conns, c)
	}
	h.mu.Unlock()

	for _, c := range conns {
		// Use NextWriter for concurrent safety if needed, but WriteMessage is thread-safe for *different* connections.
		// However, gorilla/websocket requires one concurrent write per connection.
		// Since we iterate and write sequentially here, it is safe for *this* broadcast.
		// But if multiple Broadcasts happen concurrently, we might write to the same connection concurrently.
		// So we need a lock per connection or a central broadcast loop.
		// For simplicity, we wrap WriteMessage in a mutex or use a channel per connection.
		// Here, we'll just ignore errors for now, but in prod use a dedicated pump.
		h.mu.Lock()
		_ = c.WriteMessage(websocket.TextMessage, b)
		h.mu.Unlock()
	}
}

func (h *EventHub) ReplaySince(ctx context.Context, sinceID int64, limit int) ([]map[string]any, error) {
	if h.eventRepo == nil {
		return []map[string]any{}, nil
	}
	items, err := h.eventRepo.ListSince(ctx, sinceID, limit)
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		var payload map[string]any
		if err := json.Unmarshal([]byte(item.Payload), &payload); err != nil {
			continue
		}
		payload["_event_id"] = strconv.FormatInt(item.ID, 10)
		out = append(out, payload)
	}
	return out, nil
}

func (h *EventHub) ConnectionCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.conns)
}
