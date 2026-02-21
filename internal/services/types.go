package services

import (
	"net/http"
	"sync"

	"media-assistant-os/internal/repos"

	"github.com/gorilla/websocket"
)

type OnboardingRequest struct {
	ImportDownloads bool `json:"import_downloads"`
	ImportPictures  bool `json:"import_pictures"`
	ImportVideos    bool `json:"import_videos"`
	ImportMusic     bool `json:"import_music"`
	ImportDesktop   bool `json:"import_desktop"`
}

// AssetService Structs (defined in asset_service.go)

type IndexFileRequest struct {
	Path      string `json:"path"`
	ProjectID string `json:"project_id"`
	Trigger   string `json:"trigger,omitempty"` // watcher | scan | api
}

type IndexFileResult struct {
	AssetID     string  `json:"asset_id"`
	Fingerprint string  `json:"fingerprint"`
	ParentID    *string `json:"parent_asset_id,omitempty"`
}

// EventHub Structs
type EventHub struct {
	mu        sync.Mutex
	conns     map[*websocket.Conn]struct{}
	upgrader  websocket.Upgrader
	eventRepo *repos.EventLogRepo
}

// WebSocket upgrader config
var defaultUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}
