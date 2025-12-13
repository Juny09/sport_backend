package config

import (
    "log/slog"
    "os"
)

// NewLogger 创建结构化日志记录器（中文说明：使用 JSON 格式便于生产环境收集）
func NewLogger() *slog.Logger {
    h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
    return slog.New(h)
}

