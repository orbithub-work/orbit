package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// PluginType defines the three types of components
type PluginType string

const (
	PluginTypeFrontend  PluginType = "frontend"  // UI components (static assets)
	PluginTypeBackend   PluginType = "backend"   // Process extensions (executables)
	PluginTypeSatellite PluginType = "satellite" // Independent apps (network services)
)

// PluginTier defines the commercial tier
type PluginTier string

const (
	PluginTierFree       PluginTier = "free"
	PluginTierPro        PluginTier = "pro"
	PluginTierEnterprise PluginTier = "enterprise"
)

// PluginManifest is the "constitution" - all components must have this
type PluginManifest struct {
	// Universal fields (required for all types)
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Version     string     `json:"version"`
	Type        PluginType `json:"type"`
	Description string     `json:"description,omitempty"`
	Author      string     `json:"author,omitempty"`
	License     string     `json:"license,omitempty"`
	Tier        PluginTier `json:"tier"`

	// Frontend-specific
	Entry       string        `json:"entry,omitempty"`        // Path to JS bundle
	Mounts      []PluginMount `json:"mounts,omitempty"`       // UI mount points
	Permissions []string      `json:"permissions,omitempty"`  // Required permissions

	// Backend-specific
	Executable   string               `json:"executable,omitempty"`   // Path to executable
	Capabilities []PluginCapability   `json:"capabilities,omitempty"` // What it can do

	// Satellite-specific
	Endpoint  string                 `json:"endpoint,omitempty"`  // HTTP endpoint (e.g., http://127.0.0.1:9090)
	Heartbeat *HeartbeatConfig       `json:"heartbeat,omitempty"` // Heartbeat configuration
}

// HeartbeatConfig for satellite apps
type HeartbeatConfig struct {
	Interval int    `json:"interval"` // Seconds
	Endpoint string `json:"endpoint"` // Health check endpoint
}

// PluginScanner discovers and parses manifest.json files
type PluginScanner struct {
	pluginsDir string
}

// NewPluginScanner creates a new scanner
func NewPluginScanner(pluginsDir string) *PluginScanner {
	return &PluginScanner{
		pluginsDir: pluginsDir,
	}
}

// ScanAll discovers all plugins in the plugins directory
func (s *PluginScanner) ScanAll() ([]PluginManifest, error) {
	var manifests []PluginManifest

	// Scan each plugin type directory
	for _, pluginType := range []PluginType{PluginTypeFrontend, PluginTypeBackend, PluginTypeSatellite} {
		typeDir := filepath.Join(s.pluginsDir, string(pluginType))
		
		// Skip if directory doesn't exist
		if _, err := os.Stat(typeDir); os.IsNotExist(err) {
			continue
		}

		// Read all subdirectories
		entries, err := os.ReadDir(typeDir)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s directory: %w", pluginType, err)
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			// Try to load manifest.json
			manifestPath := filepath.Join(typeDir, entry.Name(), "manifest.json")
			manifest, err := s.LoadManifest(manifestPath)
			if err != nil {
				fmt.Printf("⚠️  Failed to load manifest for %s/%s: %v\n", pluginType, entry.Name(), err)
				continue
			}

			// Validate type matches directory
			if manifest.Type != pluginType {
				fmt.Printf("⚠️  Plugin %s has type '%s' but is in '%s' directory\n", 
					manifest.ID, manifest.Type, pluginType)
				continue
			}

			manifests = append(manifests, *manifest)
			fmt.Printf("✅ Discovered %s plugin: %s (%s)\n", pluginType, manifest.Name, manifest.ID)
		}
	}

	return manifests, nil
}

// LoadManifest loads and parses a single manifest.json file
func (s *PluginScanner) LoadManifest(path string) (*PluginManifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest PluginManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Validate required fields
	if manifest.ID == "" {
		return nil, fmt.Errorf("manifest missing required field: id")
	}
	if manifest.Name == "" {
		return nil, fmt.Errorf("manifest missing required field: name")
	}
	if manifest.Version == "" {
		return nil, fmt.Errorf("manifest missing required field: version")
	}
	if manifest.Type == "" {
		return nil, fmt.Errorf("manifest missing required field: type")
	}
	if manifest.Tier == "" {
		manifest.Tier = PluginTierFree // Default to free
	}

	return &manifest, nil
}

// ValidateManifest checks if a manifest is valid
func (s *PluginScanner) ValidateManifest(manifest *PluginManifest) error {
	// Type-specific validation
	switch manifest.Type {
	case PluginTypeFrontend:
		if manifest.Entry == "" {
			return fmt.Errorf("frontend plugin must have 'entry' field")
		}
		if len(manifest.Mounts) == 0 {
			return fmt.Errorf("frontend plugin must have at least one mount point")
		}

	case PluginTypeBackend:
		if manifest.Executable == "" {
			return fmt.Errorf("backend plugin must have 'executable' field")
		}
		if len(manifest.Capabilities) == 0 {
			return fmt.Errorf("backend plugin must declare capabilities")
		}

	case PluginTypeSatellite:
		if manifest.Endpoint == "" {
			return fmt.Errorf("satellite plugin must have 'endpoint' field")
		}
		if manifest.Heartbeat == nil {
			return fmt.Errorf("satellite plugin must have 'heartbeat' configuration")
		}

	default:
		return fmt.Errorf("unknown plugin type: %s", manifest.Type)
	}

	return nil
}

// PrintDiscoveryReport prints a summary of discovered plugins
func (s *PluginScanner) PrintDiscoveryReport(manifests []PluginManifest) {
	fmt.Println("\n=== 智归档OS 插件发现报告 ===")
	fmt.Printf("插件目录: %s\n", s.pluginsDir)
	fmt.Printf("发现插件总数: %d\n\n", len(manifests))

	// Group by type
	byType := make(map[PluginType][]PluginManifest)
	for _, m := range manifests {
		byType[m.Type] = append(byType[m.Type], m)
	}

	// Print by type
	for _, pluginType := range []PluginType{PluginTypeFrontend, PluginTypeBackend, PluginTypeSatellite} {
		plugins := byType[pluginType]
		if len(plugins) == 0 {
			continue
		}

		fmt.Printf("【%s 组件】(%d个)\n", pluginType, len(plugins))
		for _, p := range plugins {
			tierIcon := "🆓"
			if p.Tier == PluginTierPro {
				tierIcon = "💎"
			} else if p.Tier == PluginTierEnterprise {
				tierIcon = "🏢"
			}
			fmt.Printf("  %s %s (%s) v%s\n", tierIcon, p.Name, p.ID, p.Version)
			if p.Description != "" {
				fmt.Printf("     └─ %s\n", p.Description)
			}
		}
		fmt.Println()
	}

	fmt.Println("=== 报告结束 ===\n")
}
