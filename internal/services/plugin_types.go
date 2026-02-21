package services

import "time"

type PluginMode string

const (
	PluginModeLocalProcess  PluginMode = "local_process"
	PluginModeNetworkWorker PluginMode = "network_service"
	PluginModeFrontend      PluginMode = "frontend"
)

type PluginCapability struct {
	Name       string   `json:"name"`
	TaskTypes  []string `json:"task_types,omitempty"`
	Extensions []string `json:"extensions,omitempty"`
}

// PluginMount describes where a plugin UI entry should be mounted in host UI.
type PluginMount struct {
	Slot     string `json:"slot"`
	Entry    string `json:"entry"`
	Title    string `json:"title,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Order    int    `json:"order,omitempty"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Location string `json:"location,omitempty"` // Legacy alias for some frontend stores.
}

// PluginUIConfig keeps backward compatibility with current frontend micro-app format.
type PluginUIConfig struct {
	Entry    string `json:"entry"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
	Location string `json:"location,omitempty"`
}

type PluginRegistrationRequest struct {
	PluginID        string             `json:"plugin_id"`
	Name            string             `json:"name"`
	Version         string             `json:"version,omitempty"`
	Description     string             `json:"description,omitempty"`
	Mode            PluginMode         `json:"mode,omitempty"`
	Executable      string             `json:"executable,omitempty"` // local_process only
	Endpoint        string             `json:"endpoint,omitempty"`   // network_service only
	Extensions      []string           `json:"extensions,omitempty"`
	TaskTypes       []string           `json:"task_types,omitempty"`
	Capabilities    []PluginCapability `json:"capabilities,omitempty"`
	ProtocolVersion string             `json:"protocol_version,omitempty"`
	Permissions     []string           `json:"permissions,omitempty"`
	UI              *PluginUIConfig    `json:"ui,omitempty"`
	Mounts          []PluginMount      `json:"mounts,omitempty"`
}

type PluginRegistrationResponse struct {
	PluginID string `json:"plugin_id"`
	Token    string `json:"token"`
	Mode     string `json:"mode"`
}

type PluginHeartbeatRequest struct {
	PluginID string `json:"plugin_id"`
	Token    string `json:"token"`
}

type PluginInfo struct {
	PluginID        string             `json:"plugin_id"`
	ID              string             `json:"id"` // Frontend compatibility alias
	Name            string             `json:"name"`
	Version         string             `json:"version,omitempty"`
	Description     string             `json:"description,omitempty"`
	Mode            PluginMode         `json:"mode"`
	Endpoint        string             `json:"endpoint,omitempty"`
	Executable      string             `json:"executable,omitempty"`
	Extensions      []string           `json:"extensions,omitempty"`
	TaskTypes       []string           `json:"task_types,omitempty"`
	Capabilities    []PluginCapability `json:"capabilities,omitempty"`
	ProtocolVersion string             `json:"protocol_version,omitempty"`
	Permissions     []string           `json:"permissions,omitempty"`
	UI              *PluginUIConfig    `json:"ui,omitempty"`
	Mounts          []PluginMount      `json:"mounts,omitempty"`
	IssuedAt        time.Time          `json:"issued_at"`
	LastUsedAt      time.Time          `json:"last_used_at"`
	ExpiresAt       time.Time          `json:"expires_at"`
	Online          bool               `json:"online"`
	RegisteredAt    int64              `json:"registered_at"`  // Frontend compatibility.
	LastHeartbeat   int64              `json:"last_heartbeat"` // Frontend compatibility.
}

type PluginMountedSlot struct {
	PluginID   string      `json:"plugin_id"`
	PluginName string      `json:"plugin_name"`
	Online     bool        `json:"online"`
	Mount      PluginMount `json:"mount"`
}
