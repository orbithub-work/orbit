//go:build !windows
package utils

import (
	"os"
	"path/filepath"
)

func enrichSystemInfo(info *SystemInfo) {
	// TODO: Implement specific Unix/macOS system info enrichment if needed
}

// GetUserDefaultDirs returns the standard user directories for macOS/Linux
// Returns a map of Label -> Path
func GetUserDefaultDirs() map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	dirs := make(map[string]string)
	
	// Standard macOS/Linux XDG directories
	dirs["Downloads"] = filepath.Join(home, "Downloads")
	dirs["Documents"] = filepath.Join(home, "Documents")
	dirs["Pictures"] = filepath.Join(home, "Pictures")
	dirs["Movies"] = filepath.Join(home, "Movies")
	dirs["Music"] = filepath.Join(home, "Music")
	
	// Filter out non-existent directories
	validDirs := make(map[string]string)
	for label, path := range dirs {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			validDirs[label] = path
		}
	}
	
	return validDirs
}
