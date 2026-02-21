package services

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
)

// PortPoolConfig is the configuration for satellite plugins
type PortPoolConfig struct {
	CorePort       int   `json:"core_port"`        // Go内核端口 (8848)
	SatelliteBase  int   `json:"satellite_base"`   // Satellite端口段起始 (8850)
	SatelliteSize  int   `json:"satellite_size"`   // 端口段大小 (50)
	AvailablePorts []int `json:"available_ports"`  // 可用端口列表
}

// InitPortPoolConfig initializes and writes port pool configuration
func InitPortPoolConfig(dataDir string, corePort int) error {
	config := PortPoolConfig{
		CorePort:      corePort,
		SatelliteBase: 32050, // Satellites use 32050-32099
		SatelliteSize: 50,
	}

	// Scan for available ports
	fmt.Printf("🔍 Scanning satellite port range %d-%d...\n", 
		config.SatelliteBase, config.SatelliteBase+config.SatelliteSize-1)

	for i := 0; i < config.SatelliteSize; i++ {
		port := config.SatelliteBase + i
		if isPortAvailable(port) {
			config.AvailablePorts = append(config.AvailablePorts, port)
		} else {
			fmt.Printf("⚠️  Port %d is in use, skipping\n", port)
		}
	}

	fmt.Printf("✅ Found %d available ports\n", len(config.AvailablePorts))

	// Write to config file
	configPath := filepath.Join(dataDir, "config.json")
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Port pool config written to: %s\n", configPath)
	return nil
}

// LoadPortPoolConfig loads port pool configuration
func LoadPortPoolConfig(dataDir string) (*PortPoolConfig, error) {
	configPath := filepath.Join(dataDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config PortPoolConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// isPortAvailable checks if a port is available
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// FindAvailablePort finds an available port from the pool
func FindAvailablePort(config *PortPoolConfig) (int, error) {
	for _, port := range config.AvailablePorts {
		if isPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available ports in pool")
}
