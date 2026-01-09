package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var log *zap.Logger

// Config 日志配置
type Config struct {
	Level  string `mapstructure:"level"`  // debug, info, warn, error
	Format string `mapstructure:"format"` // json, console
	Output string `mapstructure:"output"` // stdout, file, both
	Path   string `mapstructure:"path"`   // 日志文件路径
}

// Init 初始化日志模块
func Init(cfg *Config) error {
	// 设置日志级别
	var level zapcore.Level
	switch cfg.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn", "warning":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 设置编码器
	var encoderConfig zapcore.EncoderConfig
	if cfg.Format == "console" {
		// 控制台格式（人类可读）
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	} else {
		// JSON 格式（结构化）
		encoderConfig = zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stack",
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
	}

	var cores []zapcore.Core

	// 控制台输出
	if cfg.Output == "stdout" || cfg.Output == "both" {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
		cores = append(cores, core)
	}

	// 文件输出
	if cfg.Output == "file" || cfg.Output == "both" {
		// 确保日志目录存在
		if cfg.Path == "" {
			cfg.Path = "logs"
		}
		if err := os.MkdirAll(cfg.Path, 0755); err != nil {
			return fmt.Errorf("failed to create log directory: %w", err)
		}

		// 错误日志文件
		errorPath := filepath.Join(cfg.Path, "error.log")
		errorFile, err := os.OpenFile(errorPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open error log file: %w", err)
		}

		// 普通日志文件
		infoPath := filepath.Join(cfg.Path, "info.log")
		infoFile, err := os.OpenFile(infoPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open info log file: %w", err)
		}

		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

		// Info 级别日志
		infoCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(infoFile), level)
		cores = append(cores, infoCore)

		// Error 级别及以上日志
		errorCore := zapcore.NewCore(jsonEncoder, zapcore.AddSync(errorFile), zapcore.ErrorLevel)
		cores = append(cores, errorCore)
	}

	// 创建 logger
	log = zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)

	return nil
}

// Get 获取日志实例
func Get() *zap.Logger {
	if log == nil {
		// 默认配置
		Init(&Config{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		})
	}
	return log
}

// Sync 刷新日志
func Sync() {
	if log != nil {
		log.Sync()
	}
}

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Get().Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Get().Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Get().Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Get().Error(msg, fields...)
}

// Fatal 致命日志
func Fatal(msg string, fields ...zap.Field) {
	Get().Fatal(msg, fields...)
}

// ===== Zap Field 辅助函数 =====

// String 创建 string 类型字段
func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

// Int 创建 int 类型字段
func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

// Any 创建任意类型字段
func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

// Err 创建 error 类型字段
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Duration 创建 duration 类型字段
func Duration(key string, val interface{}) zap.Field {
	return zap.Duration(key, val.(time.Duration))
}
