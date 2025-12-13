package auth

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
)

// ContextKeyUserID 上下文键（中文说明：用于在请求上下文中存储用户ID）
const ContextKeyUserID = "user_id"
const ContextKeyUserRole = "user_role"

// NewJWTMiddleware 返回一个 Gin 中间件用于校验 Supabase JWT
// 中文说明：解析 Authorization Bearer Token，校验签名，提取 sub 作为用户ID
func NewJWTMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        if jwtSecret == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "auth disabled"})
            return
        }
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing bearer token"})
            return
        }
        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
            // Supabase 默认 HS256
            if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrTokenUnverifiable
            }
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
            return
        }
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
            return
        }
        sub, _ := claims["sub"].(string)
        if sub == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing sub"})
            return
        }
        c.Set(ContextKeyUserID, sub)
        if role, _ := claims["role"].(string); role != "" {
            c.Set(ContextKeyUserRole, role)
        }
        c.Next()
    }
}

// GetUserID 从上下文获取用户ID（中文说明：快捷读取中间件设置的用户ID）
func GetUserID(c *gin.Context) (string, bool) {
    v, ok := c.Get(ContextKeyUserID)
    if !ok {
        return "", false
    }
    id, _ := v.(string)
    return id, id != ""
}

// IsAdmin 判断是否管理员（中文说明：通过 JWT 的 role 字段简易判定）
func IsAdmin(c *gin.Context) bool {
    v, ok := c.Get(ContextKeyUserRole)
    if !ok { return false }
    role, _ := v.(string)
    return role == "admin"
}
