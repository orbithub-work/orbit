//go:build windows

package utils

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	modole32    = syscall.NewLazyDLL("ole32.dll")
	modshell32  = syscall.NewLazyDLL("shell32.dll")
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procOleInitialize = modole32.NewProc("OleInitialize")
	procDoDragDrop    = modole32.NewProc("DoDragDrop")
	procGlobalAlloc   = modkernel32.NewProc("GlobalAlloc")
	procGlobalLock    = modkernel32.NewProc("GlobalLock")
	procGlobalUnlock  = modkernel32.NewProc("GlobalUnlock")
)

const (
	GHND            = 0x0042
	CF_HDROP        = 15
	DROPEFFECT_COPY = 1
)

type DROPFILES struct {
	PFiles uint32
	Pt     struct{ X, Y int32 }
	FNC    int32
	FWide  int32
}

// StartDragFiles 发起系统原生拖拽 (CF_HDROP)
func StartDragFiles(paths []string) error {
	// 1. 初始化 OLE
	procOleInitialize.Call(0)

	// 2. 准备文件路径数据 (NULL 隔开，双 NULL 结尾)
	var data []uint16
	for _, p := range paths {
		abs, err := syscall.FullPath(p)
		if err != nil {
			abs = p
		}
		u, _ := syscall.UTF16FromString(abs)
		data = append(data, u...)
	}
	data = append(data, 0) // 双 NULL 结尾

	// 3. 分配全局内存
	hGlobal, _, _ := procGlobalAlloc.Call(GHND, uintptr(unsafe.Sizeof(DROPFILES{})+uintptr(len(data)*2)))
	if hGlobal == 0 {
		return fmt.Errorf("failed to allocate global memory")
	}

	ptr, _, _ := procGlobalLock.Call(hGlobal)
	df := (*DROPFILES)(unsafe.Pointer(ptr))
	df.PFiles = uint32(unsafe.Sizeof(DROPFILES{}))
	df.FWide = 1 // Unicode

	dest := unsafe.Pointer(ptr + uintptr(df.PFiles))
	for i, v := range data {
		*(*uint16)(unsafe.Pointer(uintptr(dest) + uintptr(i*2))) = v
	}
	procGlobalUnlock.Call(hGlobal)

	// 注意：完整的 DoDragDrop 需要实现 IDataObject 和 IDropSource 接口。
	// 在纯 Go 中手动构造 COM vtable 非常复杂。
	// 暂时返回错误，提示需要 CGO 或专门的 COM 桥接。
	// 但我们已经准备好了数据格式，后续可以通过集成一个小型 C++ 辅助库或使用 CGO 解决。

	fmt.Printf("System drag requested for %d files: %v\n", len(paths), paths)
	return fmt.Errorf("native DoDragDrop implementation pending (requires COM interface vtable construction in Go)")
}
