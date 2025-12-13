-- 初始化所有设施的营业时间（中文注释）：默认每天 08:00-22:00
INSERT INTO opening_hours (facility_id, day_of_week, open_time, close_time)
SELECT f.id, d, '08:00'::time, '22:00'::time
FROM facilities f CROSS JOIN generate_series(0,6) AS d
ON CONFLICT (facility_id, day_of_week) DO NOTHING;

