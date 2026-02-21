package httpapi

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (h *Handler) handleServeAssetFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, APIResponse{Success: false, Error: "method not allowed"})
		return
	}

	filePath := r.URL.Query().Get("path")
	if filePath == "" {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "missing path"})
		return
	}

	filePath = filepath.Clean(filePath)
	if strings.Contains(filePath, "..") {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "invalid path"})
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			writeJSON(w, http.StatusNotFound, APIResponse{Success: false, Error: "file not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, APIResponse{Success: false, Error: err.Error()})
		return
	}

	if info.IsDir() {
		writeJSON(w, http.StatusBadRequest, APIResponse{Success: false, Error: "path is a directory"})
		return
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	contentType := getContentType(ext)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.Header().Set("Cache-Control", "max-age=3600")

	http.ServeFile(w, r, filePath)
}

func getContentType(ext string) string {
	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".bmp":  "image/bmp",
		".svg":  "image/svg+xml",
		".ico":  "image/x-icon",
		".tiff": "image/tiff",
		".tif":  "image/tiff",

		".mp4":  "video/mp4",
		".webm": "video/webm",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".mkv":  "video/x-matroska",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".m4v":  "video/mp4",
		".3gp":  "video/3gpp",

		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".flac": "audio/flac",
		".aac":  "audio/aac",
		".ogg":  "audio/ogg",
		".wma":  "audio/x-ms-wma",
		".m4a":  "audio/mp4",
		".aiff": "audio/aiff",

		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".rtf":  "application/rtf",

		".psd":    "image/vnd.adobe.photoshop",
		".ai":     "application/postscript",
		".eps":    "application/postscript",
		".sketch": "application/sketch",

		".zip": "application/zip",
		".rar": "application/vnd.rar",
		".7z":  "application/x-7z-compressed",
		".tar": "application/x-tar",
		".gz":  "application/gzip",

		".ttf":   "font/ttf",
		".otf":   "font/otf",
		".woff":  "font/woff",
		".woff2": "font/woff2",
		".eot":   "application/vnd.ms-fontobject",

		".json": "application/json",
		".xml":  "application/xml",
		".html": "text/html",
		".css":  "text/css",
		".js":   "application/javascript",
	}

	if ct, ok := contentTypes[ext]; ok {
		return ct
	}
	return "application/octet-stream"
}
