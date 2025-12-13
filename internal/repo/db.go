package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nedpals/supabase-go"
)

// DB 封装 Supabase 客户端（中文说明：使用 HTTP API 替代直连数据库）
type DB struct {
	Client *supabase.Client
}

// Internal structs for mapping snake_case DB fields
type facilityDB struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	IsActive bool   `json:"is_active"`
}

func (f *facilityDB) toAPI() Facility {
	return Facility{
		ID:       f.ID,
		Name:     f.Name,
		Type:     f.Type,
		IsActive: f.IsActive,
	}
}

type resourceUnitDB struct {
	ID         int64  `json:"id"`
	FacilityID int64  `json:"facility_id"`
	Label      string `json:"label"`
	IsActive   bool   `json:"is_active"`
}

func (r *resourceUnitDB) toAPI() ResourceUnit {
	return ResourceUnit{
		ID:         r.ID,
		FacilityID: r.FacilityID,
		Label:      r.Label,
		IsActive:   r.IsActive,
	}
}

type bookingDB struct {
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
	ID             int64     `json:"id"`
	ResourceUnitID int64     `json:"resource_unit_id"`
	UserID         string    `json:"user_id"`
	Status         string    `json:"status"`
	Price          float64   `json:"price"`
	Notes          string    `json:"notes,omitempty"`
}

func (b *bookingDB) toAPI() Booking {
	return Booking{
		StartTime:      b.StartTime,
		EndTime:        b.EndTime,
		ID:             b.ID,
		ResourceUnitID: b.ResourceUnitID,
		UserID:         b.UserID,
		Status:         b.Status,
		Price:          b.Price,
		Notes:          b.Notes,
	}
}

// NewDB 创建 Supabase 客户端连接
func NewDB(url, key string) (*DB, error) {
	if url == "" || key == "" {
		return nil, errors.New("missing supabase url or key")
	}

	client := supabase.CreateClient(url, key)
	return &DB{Client: client}, nil
}

// Close 关闭连接（Supabase HTTP 客户端无需显式关闭，保留接口兼容性）
func (d *DB) Close() {
	// No-op
}

// Facility 设施实体
type Facility struct {
	ID       int64  `json:"ID"`
	Name     string `json:"Name"`
	Type     string `json:"Type"`
	IsActive bool   `json:"IsActive"`
}

// ResourceUnit 单元实体
type ResourceUnit struct {
	ID         int64  `json:"ID"`
	FacilityID int64  `json:"FacilityID"`
	Label      string `json:"Label"`
	IsActive   bool   `json:"IsActive"`
}

// Booking 预约实体
type Booking struct {
	StartTime      time.Time `json:"StartTime"`
	EndTime        time.Time `json:"EndTime"`
	ID             int64     `json:"ID"`
	ResourceUnitID int64     `json:"ResourceUnitID"`
	UserID         string    `json:"UserID"`
	Status         string    `json:"Status"`
	Price          float64   `json:"Price"`
	Notes          string    `json:"Notes,omitempty"`
}

// PricingRule 价格规则
type PricingRule struct {
	FacilityType string  `json:"FacilityType"`
	DayOfWeek    int     `json:"DayOfWeek"`
	StartHour    int     `json:"StartHour"`
	EndHour      int     `json:"EndHour"`
	PricePerHour float64 `json:"PricePerHour"`
}

// BlackoutRequest 封场请求
type BlackoutRequest struct {
	FacilityID     *int64 `json:"facility_id,omitempty"`
	ResourceUnitID *int64 `json:"resource_unit_id,omitempty"`
	StartTime      string `json:"start_time"` // ISO8601 for insert
	EndTime        string `json:"end_time"`   // ISO8601 for insert
	Reason         string `json:"reason"`
}

// Blackout 简化结构用于可用性计算
type Blackout struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ListFacilities 查询设施列表
func (d *DB) ListFacilities(ctx context.Context) ([]Facility, error) {
	var out []facilityDB
	err := d.Client.DB.From("facilities").
		Select("*").
		Execute(&out)
	if err != nil {
		return nil, err
	}

	res := make([]Facility, len(out))
	for i, v := range out {
		res[i] = v.toAPI()
	}
	return res, nil
}

// GetFacilityByID 查询单个设施
func (d *DB) GetFacilityByID(ctx context.Context, id int64) (*Facility, error) {
	var out []facilityDB
	err := d.Client.DB.From("facilities").
		Select("*").
		Eq("id", fmt.Sprintf("%d", id)).
		Execute(&out)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, errors.New("facility not found")
	}
	f := out[0].toAPI()
	return &f, nil
}

// ListUnitsByFacility 查询设施下的单元列表
func (d *DB) ListUnitsByFacility(ctx context.Context, facilityID int64) ([]ResourceUnit, error) {
	var out []resourceUnitDB
	err := d.Client.DB.From("resource_units").
		Select("*").
		Eq("facility_id", fmt.Sprintf("%d", facilityID)).
		Execute(&out)
	if err != nil {
		return nil, err
	}

	res := make([]ResourceUnit, len(out))
	for i, v := range out {
		res[i] = v.toAPI()
	}
	return res, nil
}

// ListUnitsByFacilityType 根据设施类型查询激活单元
func (d *DB) ListUnitsByFacilityType(ctx context.Context, facilityType string) ([]ResourceUnit, error) {
	var out []resourceUnitDB
	err := d.Client.DB.From("resource_units").
		Select("*,facilities!inner(type)").
		Eq("facilities.type", facilityType).
		Eq("is_active", "true").
		Eq("facilities.is_active", "true").
		Execute(&out)
	if err != nil {
		return nil, err
	}

	res := make([]ResourceUnit, len(out))
	for i, v := range out {
		res[i] = v.toAPI()
	}
	return res, nil
}

// ListBookingsForUnitOnDay 查询某单元在指定日期的所有预约
func (d *DB) ListBookingsForUnitOnDay(ctx context.Context, unitID int64, day time.Time) ([]Booking, error) {
	start := time.Date(day.Year(), day.Month(), day.Day(), 0, 0, 0, 0, time.UTC)
	end := start.Add(24 * time.Hour)

	var out []bookingDB
	err := d.Client.DB.From("bookings").
		Select("start_time,end_time,id,resource_unit_id,user_id,status,price").
		Eq("resource_unit_id", fmt.Sprintf("%d", unitID)).
		Neq("status", "cancelled").
		Lt("start_time", end.Format(time.RFC3339)).
		Gt("end_time", start.Format(time.RFC3339)).
		Execute(&out)
	if err != nil {
		return nil, err
	}

	res := make([]Booking, len(out))
	for i, v := range out {
		res[i] = v.toAPI()
	}
	return res, nil
}

// ListBlackoutsForUnitOrFacilityOnDay 查询封场（单元或设施级别）
func (d *DB) ListBlackoutsForUnitOrFacilityOnDay(ctx context.Context, facilityID int64, unitID int64, day time.Time) ([]Blackout, error) {
	// TODO: 实现 Blackout 查询
	return []Blackout{}, nil
}

// CreateBooking 创建预约
func (d *DB) CreateBooking(ctx context.Context, unitID int64, userID string, start, end time.Time, notes string) (*Booking, error) {
	if !start.Before(end) {
		return nil, errors.New("invalid time range: start must be before end")
	}

	day := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	existing, err := d.ListBookingsForUnitOnDay(ctx, unitID, day)
	if err != nil {
		return nil, err
	}
	for _, b := range existing {
		if start.Before(b.EndTime) && end.After(b.StartTime) {
			return nil, fmt.Errorf("slot unavailable: overlaps with booking ID %d [%s - %s]", b.ID, b.StartTime.Format(time.RFC3339), b.EndTime.Format(time.RFC3339))
		}
	}

	payload := map[string]interface{}{
		"resource_unit_id": unitID,
		"user_id":          userID,
		"start_time":       start.Format(time.RFC3339),
		"end_time":         end.Format(time.RFC3339),
		"notes":            notes,
		"status":           "confirmed", // Default status
	}

	var out []bookingDB
	err = d.Client.DB.From("bookings").
		Insert(payload).
		Execute(&out)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, errors.New("failed to create booking")
	}
	res := out[0].toAPI()
	return &res, nil
}

// GetBookingByID 获取单个预约
func (d *DB) GetBookingByID(ctx context.Context, id int64) (*Booking, error) {
	var out []bookingDB
	err := d.Client.DB.From("bookings").
		Select("*").
		Eq("id", fmt.Sprintf("%d", id)).
		Execute(&out)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		return nil, errors.New("booking not found")
	}
	res := out[0].toAPI()
	return &res, nil
}

// ListBookingsByUser 获取用户预约
func (d *DB) ListBookingsByUser(ctx context.Context, userID string) ([]Booking, error) {
	var out []bookingDB
	err := d.Client.DB.From("bookings").
		Select("*").
		Eq("user_id", userID).
		Execute(&out)
	if err != nil {
		return nil, err
	}

	res := make([]Booking, len(out))
	for i, v := range out {
		res[i] = v.toAPI()
	}
	return res, nil
}

// CancelBooking 取消预约
func (d *DB) CancelBooking(ctx context.Context, id int64) error {
	var out []bookingDB
	payload := map[string]interface{}{"status": "cancelled"}
	err := d.Client.DB.From("bookings").
		Update(payload).
		Eq("id", fmt.Sprintf("%d", id)).
		Execute(&out)
	return err
}

// RescheduleBooking 改签预约
func (d *DB) RescheduleBooking(ctx context.Context, id int64, start, end time.Time) error {
	payload := map[string]interface{}{
		"start_time": start.Format(time.RFC3339),
		"end_time":   end.Format(time.RFC3339),
	}
	var out []bookingDB
	err := d.Client.DB.From("bookings").
		Update(payload).
		Eq("id", fmt.Sprintf("%d", id)).
		Execute(&out)
	return err
}

// CreatePricingRule 创建价格规则
func (d *DB) CreatePricingRule(ctx context.Context, rule PricingRule) error {
	// Need to map PricingRule (PascalCase) to snake_case payload manually or use a struct with snake_case tags
	// Or just use a map
	payload := map[string]interface{}{
		"facility_type":  rule.FacilityType,
		"day_of_week":    rule.DayOfWeek,
		"start_hour":     rule.StartHour,
		"end_hour":       rule.EndHour,
		"price_per_hour": rule.PricePerHour,
	}
	var out []interface{}
	err := d.Client.DB.From("pricing_rules").Insert(payload).Execute(&out)
	return err
}

// CreateBlackout 创建封场
func (d *DB) CreateBlackout(ctx context.Context, req BlackoutRequest) error {
	// ... existing logic uses map, so it's fine ...
	payload := map[string]interface{}{
		"reason": req.Reason,
	}
	if req.FacilityID != nil {
		payload["facility_id"] = *req.FacilityID
	}
	if req.ResourceUnitID != nil {
		payload["resource_unit_id"] = *req.ResourceUnitID
	}

	// Construct PostgREST range literal
	rangeStr := fmt.Sprintf("[%s,%s)", req.StartTime, req.EndTime)
	payload["time_range"] = rangeStr

	var out []interface{}
	err := d.Client.DB.From("blackouts").Insert(payload).Execute(&out)
	return err
}

// ListAdminBookings 管理员查询预约
func (d *DB) ListAdminBookings(ctx context.Context, facilityType string, start, end time.Time) ([]Booking, error) {
	// Complex join + filter
	var out []bookingDB
	err := d.Client.DB.From("bookings").
		Select("*,resource_units!inner(facilities!inner(type))").
		Eq("resource_units.facilities.type", facilityType).
		// Overlap check: start < queryEnd AND end > queryStart
		Lt("start_time", end.Format(time.RFC3339)).
		Gt("end_time", start.Format(time.RFC3339)).
		Execute(&out)
	if err != nil {
		return nil, err
	}

	res := make([]Booking, len(out))
	for i, v := range out {
		res[i] = v.toAPI()
	}
	return res, nil
}

// CreateFacility 创建设施
func (d *DB) CreateFacility(ctx context.Context, name, type_ string) error {
	payload := map[string]interface{}{
		"name":      name,
		"type":      type_,
		"is_active": true,
	}
	var out []interface{}
	return d.Client.DB.From("facilities").Insert(payload).Execute(&out)
}

// CreateResourceUnit 创建单元
func (d *DB) CreateResourceUnit(ctx context.Context, facilityID int64, label string) error {
	payload := map[string]interface{}{
		"facility_id": facilityID,
		"label":       label,
		"is_active":   true,
	}
	var out []interface{}
	return d.Client.DB.From("resource_units").Insert(payload).Execute(&out)
}

// UpdateResourceUnit 更新单元状态
func (d *DB) UpdateResourceUnit(ctx context.Context, id int64, isActive bool) error {
	payload := map[string]interface{}{
		"is_active": isActive,
	}
	var out []interface{}
	return d.Client.DB.From("resource_units").Update(payload).Eq("id", fmt.Sprintf("%d", id)).Execute(&out)
}
