package config

import (
	"errors"
	"os"
)

// Config 用于保存服务运行所需的环境配置
// 中文说明：从系统环境变量读取 Supabase 连接、JWT 密钥、服务端口等
type Config struct {
	Port              string
	SupabaseDBURL     string
	SupabaseJWTSecret string
	SupabaseURL       string
	SupabaseAnonKey   string
}

// Load 读取并校验配置
func Load() (Config, error) {
	cfg := Config{
		Port:              getenvDefault("PORT", "8080"),
		SupabaseDBURL:     firstNonEmpty(os.Getenv("SUPABASE_DB_URL"), os.Getenv("DATABASE_URL")),
		SupabaseJWTSecret: os.Getenv("SUPABASE_JWT_SECRET"),
		SupabaseURL:       os.Getenv("SUPABASE_URL"),
		SupabaseAnonKey:   os.Getenv("SUPABASE_ANON_KEY"),
	}
	// 允许无 DB 情况启动（便于本地先跑起来），但提示缺失
	if cfg.SupabaseDBURL == "" {
		// 仅警告；部分接口会受限
	}
	if cfg.SupabaseJWTSecret == "" {
		// 仅警告；受保护接口不可用
	}
	if cfg.Port == "" {
		return cfg, errors.New("port not set")
	}
	return cfg, nil
}

func getenvDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
