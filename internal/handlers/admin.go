package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/user/sport-booking/internal/auth"
    "github.com/user/sport-booking/internal/repo"
)

// RegisterAdminRoutes 注册管理端路由
func RegisterAdminRoutes(r *gin.Engine, db *repo.DB, jwtSecret string) {
    authMW := auth.NewJWTMiddleware(jwtSecret)

    // 添加价格规则
	r.POST("/pricing_rules", authMW, func(c *gin.Context) {
		if !auth.IsAdmin(c) { c.JSON(http.StatusForbidden, gin.H{"error": "admin only"}); return }
		var body repo.PricingRule
		if err := c.BindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
		
		if err := db.CreatePricingRule(c.Request.Context(), body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	// 添加封场时间
	r.POST("/blackouts", authMW, func(c *gin.Context) {
		if !auth.IsAdmin(c) { c.JSON(http.StatusForbidden, gin.H{"error": "admin only"}); return }
		var body repo.BlackoutRequest
		if err := c.BindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
		
		if err := db.CreateBlackout(c.Request.Context(), body); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	// 管理查询预约
	r.GET("/admin/bookings", authMW, func(c *gin.Context) {
		if !auth.IsAdmin(c) { c.JSON(http.StatusForbidden, gin.H{"error": "admin only"}); return }
		facilityType := c.Query("facility_type")
		dateStr := c.Query("date")
		day, err := time.Parse("2006-01-02", dateStr)
		if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"}); return }
		start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
		end := start.Add(24 * time.Hour)
		
		list, err := db.ListAdminBookings(c.Request.Context(), facilityType, start, end)
		if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
		c.JSON(http.StatusOK, list)
	})
}
