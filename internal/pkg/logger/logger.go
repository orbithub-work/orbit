package logger

import (
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *zap.Logger
	once         sync.Once
)

// Init 初始化全局 Logger
// logDir: 日志文件存储目录
// debug: 是否开启调试模式（输出到控制台）
func Init(logDir string, debug bool) error {
	var err error
	once.Do(func() {
		// 确保日志目录存在
		if err = os.MkdirAll(logDir, 0755); err != nil {
			return
		}

		// 日志文件切割配置
		logFile := filepath.Join(logDir, "media-assistant.log")
		rotator := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    10,   // 每个日志文件最大 10MB
			MaxBackups: 5,    // 保留最近 5 个文件
			MaxAge:     30,   // 保留最近 30 天
			Compress:   true, // 压缩旧日志
		}

		// 编码器配置
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 时间格式: 2023-10-01T12:00:00.000Z
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder // 级别格式: INFO, WARN
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 调用者格式: package/file.go:line

		// 核心配置
		var core zapcore.Core

		// 文件输出核心
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		fileWriter := zapcore.AddSync(rotator)
		fileCore := zapcore.NewCore(fileEncoder, fileWriter, zap.InfoLevel)

		if debug {
			// 控制台输出核心 (开发模式用 Console 编码器，更易读)
			consoleEncoderConfig := zap.NewDevelopmentEncoderConfig()
			consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)
			consoleWriter := zapcore.Lock(os.Stdout)
			consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zap.DebugLevel)

			// 双路输出
			core = zapcore.NewTee(fileCore, consoleCore)
		} else {
			// 仅文件输出
			core = fileCore
		}

		// 创建 Logger
		// AddCaller: 添加调用行号
		// AddStacktrace: Error 级别以上添加堆栈
		globalLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	})

	return err
}

// Get 获取全局 Logger
func Get() *zap.Logger {
	if globalLogger == nil {
		// 如果未初始化，返回一个不输出的 logger 防止 panic，或者简单的控制台 logger
		return zap.NewExample()
	}
	return globalLogger
}

// Sync 刷新缓冲
func Sync() {
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}
}

// 快捷方法
func Debug(msg string, fields ...zap.Field) { Get().Debug(msg, fields...) }
func Info(msg string, fields ...zap.Field) { Get().Info(msg, fields...) }
func Warn(msg string, fields ...zap.Field) { Get().Warn(msg, fields...) }
func Error(msg string, fields ...zap.Field) { Get().Error(msg, fields...) }
func Fatal(msg string, fields ...zap.Field) { Get().Fatal(msg, fields...) }
