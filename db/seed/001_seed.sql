-- 初始化基础设施与球场（中文注释）：按需求创建 8 个羽毛球、8 个网球等

INSERT INTO facilities (name, type) VALUES
  ('Badminton', 'badminton'),
  ('Tennis', 'tennis'),
  ('Gym', 'gym'),
  ('Multipurpose Hall', 'multipurpose'),
  ('Other Sport', 'other')
ON CONFLICT DO NOTHING;

-- 选择插入的 id（假设顺序插入）；真实环境建议使用名称查询
WITH f AS (
  SELECT id, type FROM facilities
)
INSERT INTO resource_units (facility_id, label)
SELECT f.id, CONCAT('Court ', n)
FROM f CROSS JOIN generate_series(1, 8) AS n
WHERE f.type IN ('badminton','tennis')
ON CONFLICT DO NOTHING;

-- Gym 与多功能厅默认 1 单元
WITH f AS (
  SELECT id, type FROM facilities
)
INSERT INTO resource_units (facility_id, label)
SELECT f.id, 'Area 1'
FROM f WHERE f.type IN ('gym','multipurpose')
ON CONFLICT DO NOTHING;

-- Other sport 示例单元
WITH f AS (
  SELECT id FROM facilities WHERE type = 'other'
)
INSERT INTO resource_units (facility_id, label)
SELECT f.id, 'Unit 1' FROM f
ON CONFLICT DO NOTHING;

