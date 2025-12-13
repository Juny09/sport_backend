package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/user/sport-booking/internal/auth"
	"github.com/user/sport-booking/internal/handlers"
	"github.com/user/sport-booking/internal/repo"
)

// NewRouter 构建 HTTP 路由（中文说明：集中管理所有 API 路由）
func NewRouter(db *repo.DB, jwtSecret string, authClient *auth.Client) *gin.Engine {
	r := gin.Default()

	// 添加 CORS 中间件
	r.Use(func(c *gin.Context) {
		// 动态设置 Allow-Origin 为请求的 Origin，以支持 Allow-Credentials=true
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().UTC()})
	})

	// 根路径提示 (用于 Supabase 重定向测试)
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, `
			<h1>Sport Booking Backend</h1>
			<p>Server is running correctly.</p>
			<p>If you were redirected here from Supabase, your email verification might have been processed.</p>
		`)
	})

	// 注册 Auth 路由（登录/注册）
	handlers.RegisterAuthRoutes(r, authClient)

	// 用户信息（需要鉴权）
	authMW := auth.NewJWTMiddleware(jwtSecret)
	r.GET("/me", authMW, func(c *gin.Context) {
		userID, _ := auth.GetUserID(c)
		c.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	// 设施路由
	handlers.RegisterFacilityRoutes(r, db, jwtSecret)

	// 预留路由组（后续逐步实现）
	// /availability, /bookings, /admin
	handlers.RegisterAvailabilityRoutes(r, db)
	handlers.RegisterBookingRoutes(r, db, jwtSecret)
	handlers.RegisterAdminRoutes(r, db, jwtSecret)

	return r
}
