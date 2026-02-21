package services

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// LogEntry 日志条目
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Source    string
}

// LogViewer 实时日志查看器
type LogViewer struct {
	mu         sync.RWMutex
	entries    []LogEntry
	maxEntries int
	stopChan   chan struct{}
}

// NewLogViewer 创建日志查看器
func NewLogViewer(maxEntries int) *LogViewer {
	if maxEntries <= 0 {
		maxEntries = 1000
	}
	return &LogViewer{
		entries:    make([]LogEntry, 0, maxEntries),
		maxEntries: maxEntries,
		stopChan:   make(chan struct{}),
	}
}

// AddLog 添加日志
func (lv *LogViewer) AddLog(level LogLevel, source, format string, args ...interface{}) {
	lv.mu.Lock()
	defer lv.mu.Unlock()

	message := fmt.Sprintf(format, args...)
	entry := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
		Source:    source,
	}

	lv.entries = append(lv.entries, entry)
	if len(lv.entries) > lv.maxEntries {
		lv.entries = lv.entries[len(lv.entries)-lv.maxEntries:]
	}
}

// GetLogs 获取日志
func (lv *LogViewer) GetLogs(count int) []LogEntry {
	lv.mu.RLock()
	defer lv.mu.RUnlock()

	if count <= 0 || count > len(lv.entries) {
		count = len(lv.entries)
	}
	return lv.entries[len(lv.entries)-count:]
}

// Clear 清空日志
func (lv *LogViewer) Clear() {
	lv.mu.Lock()
	defer lv.mu.Unlock()
	lv.entries = make([]LogEntry, 0, lv.maxEntries)
}

// Stop 停止日志收集
func (lv *LogViewer) Stop() {
	close(lv.stopChan)
}

// StartRealTimeDisplay 启动实时日志显示
func (lv *LogViewer) StartRealTimeDisplay(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	// 清屏并显示标题
	lv.clearScreen()
	lv.displayHeader()

	for {
		select {
		case <-ticker.C:
			lv.displayRealTimeData(ctx)
		case <-lv.stopChan:
			return
		case <-ctx.Done():
			return
		}
	}
}

// clearScreen 清屏
func (lv *LogViewer) clearScreen() {
	fmt.Print("\033[2J\033[H") // ANSI 清屏
}

// displayHeader 显示标题
func (lv *LogViewer) displayHeader() {
	fmt.Println("╔════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    Media Assistant - 实时日志监控                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

// displayRealTimeData 显示实时数据
func (lv *LogViewer) displayRealTimeData(ctx context.Context) {
	lv.clearScreen()
	lv.displayHeader()

	// 显示最近日志
	lv.displayRecentLogs()

	// 显示底部提示
	fmt.Println()
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")
	fmt.Println("按 Ctrl+C 退出 | 按 R 刷新 | 按 C 清空日志")
}

/*
// displaySystemStatus 显示系统状态
// 注意：由于循环依赖问题，移除了对 system 的直接依赖。
// 如果需要显示状态，应该通过回调接口或者单独的状态对象传入。
func (lv *LogViewer) displaySystemStatus(ctx context.Context) {
    // ... removed ...
}
*/

// displayRecentLogs 显示最近日志
func (lv *LogViewer) displayRecentLogs() {
	logs := lv.GetLogs(10)
	if len(logs) == 0 {
		fmt.Println("最近日志: 暂无日志")
		return
	}

	fmt.Println("最近日志:")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")

	for _, entry := range logs {
		levelStr := lv.getLevelString(entry.Level)
		timestamp := entry.Timestamp.Format("15:04:05")
		fmt.Printf("  [%s] [%s] %s: %s\n", timestamp, levelStr, entry.Source, entry.Message)
	}
}

// getLevelString 获取日志级别字符串
func (lv *LogViewer) getLevelString(level LogLevel) string {
	switch level {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO "
	case LogLevelWarn:
		return "WARN "
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogDebug 添加调试日志
func (lv *LogViewer) LogDebug(source, format string, args ...interface{}) {
	lv.AddLog(LogLevelDebug, source, format, args...)
}

// LogInfo 添加信息日志
func (lv *LogViewer) LogInfo(source, format string, args ...interface{}) {
	lv.AddLog(LogLevelInfo, source, format, args...)
}

// LogWarn 添加警告日志
func (lv *LogViewer) LogWarn(source, format string, args ...interface{}) {
	lv.AddLog(LogLevelWarn, source, format, args...)
}

// LogError 添加错误日志
func (lv *LogViewer) LogError(source, format string, args ...interface{}) {
	lv.AddLog(LogLevelError, source, format, args...)
}
