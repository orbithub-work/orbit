//go:build windows

package utils

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

func enrichSystemInfo(info *SystemInfo) {
	info.TotalRAM = getTotalRAMWindows()
	info.CPU = getCPUName()
}

func getTotalRAMWindows() uint64 {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	globalMemoryStatusEx := kernel32.NewProc("GlobalMemoryStatusEx")

	type memoryStatusEx struct {
		Length               uint32
		MemoryLoad           uint32
		TotalPhys            uint64
		AvailPhys            uint64
		TotalPageFile        uint64
		AvailPageFile        uint64
		TotalVirtual         uint64
		AvailVirtual         uint64
		AvailExtendedVirtual uint64
	}

	var memStatus memoryStatusEx
	memStatus.Length = uint32(unsafe.Sizeof(memStatus))

	ret, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memStatus)))
	if ret == 0 {
		return 0
	}

	return memStatus.TotalPhys
}

func getCPUName() string {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `HARDWARE\DESCRIPTION\System\CentralProcessor\0`, registry.READ)
	if err != nil {
		return "Unknown CPU"
	}
	defer k.Close()

	s, _, err := k.GetStringValue("ProcessorNameString")
	if err != nil {
		return "Unknown CPU"
	}
	return s
}

// GetUserDefaultDirs returns the standard user directories for Windows
// Returns a map of Label -> Path
func GetUserDefaultDirs() map[string]string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	dirs := make(map[string]string)

	// Standard Windows User Directories
	dirs["Downloads"] = filepath.Join(home, "Downloads")
	dirs["Documents"] = filepath.Join(home, "Documents")
	dirs["Pictures"] = filepath.Join(home, "Pictures")
	dirs["Videos"] = filepath.Join(home, "Videos")
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
