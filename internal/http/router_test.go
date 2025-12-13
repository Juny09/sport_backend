package httpserver

import (
    "net/http"
    "net/http/httptest"
    "testing"
)

// 测试基础路由（中文说明：确保健康检查与鉴权行为正常）
func TestHealthAndAuth(t *testing.T) {
    r := NewRouter(nil, "", nil)

    // /health 应返回 200
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    w := httptest.NewRecorder()
    r.ServeHTTP(w, req)
    if w.Code != http.StatusOK {
        t.Fatalf("expected 200, got %d", w.Code)
    }

    // /me 未配置 JWT 密钥时，应返回 401
    req2 := httptest.NewRequest(http.MethodGet, "/me", nil)
    w2 := httptest.NewRecorder()
    r.ServeHTTP(w2, req2)
    if w2.Code != http.StatusUnauthorized {
        t.Fatalf("expected 401, got %d", w2.Code)
    }
}

