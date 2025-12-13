package main

import (
	"context"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/user/sport-booking/internal/auth"
	"github.com/user/sport-booking/internal/config"
	httpserver "github.com/user/sport-booking/internal/http"
	"github.com/user/sport-booking/internal/repo"
)

// main 函数负责启动 HTTP 服务（后续会接入配置、路由、数据库等）
func main() {
	// 加载 .env 文件（如果存在）
	_ = godotenv.Load()

	// 使用默认日志；后续可替换为结构化日志
	logger := config.NewLogger()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load error", "err", err)
		os.Exit(1)
	}

	// 初始化 Auth Client
	authClient := auth.NewClient(cfg.SupabaseURL, cfg.SupabaseAnonKey)

	// 初始化数据库连接（使用 Supabase HTTP 客户端）
	// 注意：使用 URL 和 Anon Key，而不是 DB Connection String
	db, err := repo.NewDB(cfg.SupabaseURL, cfg.SupabaseAnonKey)
	if err != nil {
		logger.Error("Failed to init supabase client", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	// 初始化路由
	r := httpserver.NewRouter(db, cfg.SupabaseJWTSecret, authClient)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	logger.Info("starting server", "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}

	// 优雅关闭（示例，后续补充）
	_ = srv.Shutdown(context.Background())
	if db != nil {
		db.Close()
	}
}
