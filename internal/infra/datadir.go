package infra

import (
	"os"
	"path/filepath"
)

func ResolveDataDir() (string, error) {
	if v := os.Getenv("MEDIA_ASSISTANT_DATA_DIR"); v != "" {
		if err := os.MkdirAll(v, 0o755); err != nil {
			return "", err
		}
		return v, nil
	}

	dataDir := filepath.Join(".", "data")
	
	// Create subdirectories for a professional structure
	subDirs := []string{"db", "cache", "temp", "logs"}
	for _, sub := range subDirs {
		path := filepath.Join(dataDir, sub)
		if err := os.MkdirAll(path, 0o755); err != nil {
			return "", err
		}
	}

	return dataDir, nil
}
