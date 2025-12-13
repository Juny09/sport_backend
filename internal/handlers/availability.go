package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/user/sport-booking/internal/repo"
    "github.com/user/sport-booking/internal/service"
)

// RegisterAvailabilityRoutes 注册可用性查询路由
func RegisterAvailabilityRoutes(r *gin.Engine, db *repo.DB) {
    r.GET("/availability", func(c *gin.Context) {
        if db == nil { c.JSON(http.StatusServiceUnavailable, gin.H{"error": "db not configured"}); return }
        facilityType := c.Query("facility_type")
        dateStr := c.Query("date")
        durStr := c.Query("duration")
        if facilityType == "" || dateStr == "" || durStr == "" { c.JSON(http.StatusBadRequest, gin.H{"error": "missing params"}); return }
        durationMin, err := strconv.Atoi(durStr)
        if err != nil || durationMin <= 0 { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid duration"}); return }
        day, err := time.Parse("2006-01-02", dateStr)
        if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date"}); return }

        // 查询该类型下的所有激活单元
        units, err := db.ListUnitsByFacilityType(c.Request.Context(), facilityType)
        if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }

        // 查询当天的预订与封场
        minDur := time.Duration(durationMin) * time.Minute
        oh := service.OpeningHoursForDay(day) // 暴露函数名便于理解
        resp := make([]gin.H, 0, len(units))
        for _, u := range units {
            bookings, err := db.ListBookingsForUnitOnDay(c.Request.Context(), u.ID, day)
            if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
            blackouts, err := db.ListBlackoutsForUnitOrFacilityOnDay(c.Request.Context(), u.FacilityID, u.ID, day)
            if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
            var blocks []service.TimeRange
            for _, b := range bookings { blocks = append(blocks, service.TimeRange{Start: b.StartTime, End: b.EndTime}) }
            for _, b := range blackouts { blocks = append(blocks, service.TimeRange{Start: b.StartTime, End: b.EndTime}) }
            free := service.SubtractRanges(service.TimeRange{Start: oh.Start, End: oh.End}, blocks, minDur)
            resp = append(resp, gin.H{
                "unit_id": u.ID,
                "label":   u.Label,
                "free":    free,
            })
        }
        c.JSON(http.StatusOK, resp)
    })
}

