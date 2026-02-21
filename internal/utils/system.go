package utils

import (
	"fmt"
	"runtime"
)

type SystemInfo struct {
	OS        string `json:"os"`
	Arch      string `json:"arch"`
	CPU       string `json:"cpu"`
	CPUCores  int    `json:"cpu_cores"`
	TotalRAM  uint64 `json:"total_ram"`
	MachineID string `json:"machine_id"`
}

func GetSystemInfo() SystemInfo {
	info := SystemInfo{
		OS:       runtime.GOOS,
		Arch:     runtime.GOARCH,
		CPUCores: runtime.NumCPU(),
	}

	enrichSystemInfo(&info)

	return info
}

func FormatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
