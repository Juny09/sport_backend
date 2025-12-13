-- 初始化预约策略（中文注释）：按设施类型设置时长限制与粒度
INSERT INTO reservation_policies (facility_type, min_duration_minutes, max_duration_minutes, slot_granularity_minutes, advance_booking_days, cancellation_cutoff_minutes)
VALUES
  ('badminton',    60, 120, 30, 30, 120),
  ('tennis',       60, 120, 30, 30, 120),
  ('gym',          30, 240, 30, 14,  60),
  ('multipurpose', 60, 360, 60, 30, 240),
  ('other',        60, 180, 30, 14, 120)
ON CONFLICT (facility_type) DO NOTHING;

