package handlers

import (
	"net/http"
	"strconv"

	"github.com/Juny09/sport_backend/internal/auth"
	"github.com/Juny09/sport_backend/internal/repo"
	"github.com/gin-gonic/gin"
)

// RegisterFacilityRoutes 注册设施相关路由（中文说明：包含增查与单元管理）
func RegisterFacilityRoutes(r *gin.Engine, db *repo.DB, jwtSecret string) {
	// 公开路由：列表与详情
	r.GET("/facilities", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		list, err := db.ListFacilities(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, list)
	})
	r.GET("/facilities/:id", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		f, err := db.GetFacilityByID(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, f)
	})

	r.GET("/facilities/:id/units", func(c *gin.Context) {
		if db == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"})
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		units, err := db.ListUnitsByFacility(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, units)
	})

	// 管理路由：需要鉴权 + 管理员
	authMW := auth.NewJWTMiddleware(jwtSecret)
	r.POST("/facilities", authMW, func(c *gin.Context) {
		if !auth.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		var body struct {
			Name string `json:"name"`
			Type string `json:"type"`
		}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if body.Name == "" || body.Type == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name/type required"})
			return
		}
		// 简化：直接插入
		if err := db.CreateFacility(c.Request.Context(), body.Name, body.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.POST("/facilities/:id/units", authMW, func(c *gin.Context) {
		if !auth.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var body struct {
			Label string `json:"label"`
		}
		if err := c.BindJSON(&body); err != nil || body.Label == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid label"})
			return
		}

		if err := db.CreateResourceUnit(c.Request.Context(), id, body.Label); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusCreated)
	})

	r.PATCH("/units/:id", authMW, func(c *gin.Context) {
		if !auth.IsAdmin(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "admin only"})
			return
		}
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var body struct {
			IsActive *bool `json:"is_active"`
		}
		if err := c.BindJSON(&body); err != nil || body.IsActive == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		if err := db.UpdateResourceUnit(c.Request.Context(), id, *body.IsActive); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusOK)
	})
}
