package service

import (
    "time"
)

// TimeRange 表示一个时间段
type TimeRange struct {
    Start time.Time
    End   time.Time
}

// openingHours 返回默认营业时间（中文说明：简化默认 08:00-22:00，可后续配置化）
func OpeningHoursForDay(day time.Time) TimeRange {
    y, m, d := day.Date()
    loc := time.UTC
    start := time.Date(y, m, d, 8, 0, 0, 0, loc)
    end := time.Date(y, m, d, 22, 0, 0, 0, loc)
    return TimeRange{Start: start, End: end}
}

// subtractRanges 从营业时间中扣除已占用与封场时间，得到可用的空闲段
// 中文说明：输入预订与封场的时间段列表，返回可预约的空闲时间段
func SubtractRanges(base TimeRange, blocks []TimeRange, minDuration time.Duration) []TimeRange {
    // 简化实现：将 blocks 按开始时间排序并线性扫描
    // 后续可优化为区间树
    sortByStart(blocks)
    curStart := base.Start
    var free []TimeRange
    for _, b := range blocks {
        // 如果 block 起始不晚于当前起点
        if !b.Start.After(curStart) {
            if b.End.After(curStart) { curStart = b.End }
            continue
        }
        gap := TimeRange{Start: curStart, End: b.Start}
        if gap.End.Sub(gap.Start) >= minDuration { free = append(free, gap) }
        if b.End.After(curStart) { curStart = b.End }
    }
    if base.End.Sub(curStart) >= minDuration {
        free = append(free, TimeRange{Start: curStart, End: base.End})
    }
    return free
}

func sortByStart(ranges []TimeRange) {
    if len(ranges) <= 1 { return }
    // 简单插入排序，避免引入额外依赖
    for i := 1; i < len(ranges); i++ {
        j := i
        for j > 0 && ranges[j-1].Start.After(ranges[j].Start) {
            ranges[j-1], ranges[j] = ranges[j], ranges[j-1]
            j--
        }
    }
}
