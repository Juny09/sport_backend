package handlers

import (
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/user/sport-booking/internal/auth"
    "github.com/user/sport-booking/internal/repo"
)

// RegisterBookingRoutes 注册预约相关路由
func RegisterBookingRoutes(r *gin.Engine, db *repo.DB, jwtSecret string) {
    authMW := auth.NewJWTMiddleware(jwtSecret)

    // 创建预约
    r.POST("/bookings", authMW, func(c *gin.Context) {
        userID, _ := auth.GetUserID(c)
        var body struct {
            ResourceUnitID int64  `json:"resource_unit_id"`
            StartTime      string `json:"start_time"` // ISO8601
            EndTime        string `json:"end_time"`   // ISO8601
            Notes          string `json:"notes"`
        }
        if err := c.BindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
        st, err := time.Parse(time.RFC3339, body.StartTime)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"}); return }
        et, err := time.Parse(time.RFC3339, body.EndTime)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"}); return }
        b, err := db.CreateBooking(c.Request.Context(), body.ResourceUnitID, userID, st, et, body.Notes)
        if err != nil { c.JSON(http.StatusConflict, gin.H{"error": err.Error()}); return }
        c.JSON(http.StatusCreated, b)
    })

    // 获取单个预约
    r.GET("/bookings/:id", authMW, func(c *gin.Context) {
        userID, _ := auth.GetUserID(c)
        id, err := parseIDParam(c.Param("id"))
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
        b, err := db.GetBookingByID(c.Request.Context(), id)
        if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }
        if b.UserID != userID && !auth.IsAdmin(c) { c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return }
        c.JSON(http.StatusOK, b)
    })

    // 我的预约列表
    r.GET("/bookings", authMW, func(c *gin.Context) {
        if c.Query("mine") != "true" { c.JSON(http.StatusBadRequest, gin.H{"error": "mine=true required"}); return }
        userID, _ := auth.GetUserID(c)
        list, err := db.ListBookingsByUser(c.Request.Context(), userID)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
        c.JSON(http.StatusOK, list)
    })

    // 取消预约
    r.PATCH("/bookings/:id/cancel", authMW, func(c *gin.Context) {
        userID, _ := auth.GetUserID(c)
        id, err := parseIDParam(c.Param("id"))
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
        b, err := db.GetBookingByID(c.Request.Context(), id)
        if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }
        if b.UserID != userID && !auth.IsAdmin(c) { c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return }
        if err := db.CancelBooking(c.Request.Context(), id); err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
        c.Status(http.StatusOK)
    })

    // 改签预约
    r.PATCH("/bookings/:id/reschedule", authMW, func(c *gin.Context) {
        userID, _ := auth.GetUserID(c)
        id, err := parseIDParam(c.Param("id"))
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"}); return }
        b, err := db.GetBookingByID(c.Request.Context(), id)
        if err != nil { c.JSON(http.StatusNotFound, gin.H{"error": "not found"}); return }
        if b.UserID != userID && !auth.IsAdmin(c) { c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"}); return }
        var body struct {
            StartTime string `json:"start_time"`
            EndTime   string `json:"end_time"`
        }
        if err := c.BindJSON(&body); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"}); return }
        st, err := time.Parse(time.RFC3339, body.StartTime)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"}); return }
        et, err := time.Parse(time.RFC3339, body.EndTime)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"}); return }
        if err := db.RescheduleBooking(c.Request.Context(), id, st, et); err != nil { c.JSON(http.StatusConflict, gin.H{"error": err.Error()}); return }
        c.Status(http.StatusOK)
    })
}

func parseIDParam(s string) (int64, error) {
    var id int64
    _, err := fmt.Sscan(s, &id)
    return id, err
}
