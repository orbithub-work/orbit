package infra

import (
	"os"
	"path/filepath"
	"strconv"
)

func WritePortFile(dataDir string, port int) error {
	path := filepath.Join(dataDir, "server.port")
	return os.WriteFile(path, []byte(strconv.Itoa(port)), 0o644)
}
