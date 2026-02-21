package services

import (
	"path/filepath"
	"strings"
)

func (s *ScanService) shouldSkipFile(filename string) bool {
	if len(filename) > 0 && filename[0] == '.' {
		return true
	}
	ext := strings.ToLower(filepath.Ext(filename))
	excludedExts := map[string]bool{
		".exe": true, ".msi": true, ".dll": true, ".sys": true,
		".zip": true, ".rar": true, ".7z": true, ".tar": true, ".gz": true,
		".tmp": true, ".log": true, ".part": true, ".bak": true,
		".iso": true, ".dmg": true, ".pkg": true,
		".ini": true, ".db": true, ".dat": true,
	}
	return excludedExts[ext]
}

func (s *ScanService) shouldSkipFileWithPolicy(filename string, strict bool) bool {
	if s.shouldSkipFile(filename) {
		return true
	}
	if !strict {
		return false
	}
	ext := strings.ToLower(filepath.Ext(filename))
	excludedExts := map[string]bool{
		".lnk": true, ".url": true, ".crdownload": true, ".download": true, ".partial": true,
	}
	if excludedExts[ext] {
		return true
	}
	name := strings.ToLower(strings.TrimSpace(filename))
	return name == "thumbs.db" || name == "desktop.ini"
}

func (s *ScanService) shouldSkipDir(name string, strict bool) bool {
	if len(name) > 0 && name[0] == '.' {
		return true
	}
	if name == "node_modules" || name == "target" || name == "dist" || name == "build" {
		return true
	}
	if !strict {
		return false
	}
	lower := strings.ToLower(name)
	if lower == "temp" || lower == "tmp" || lower == ".cache" || lower == "cache" {
		return true
	}
	return lower == "$recycle.bin" || lower == "system volume information"
}
