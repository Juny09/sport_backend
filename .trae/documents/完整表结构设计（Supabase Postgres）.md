# 目标
- 设计体育场馆预约系统的完整数据库表（Supabase Postgres），覆盖设施/场地、预约、可用性、价格、封场、营业时间、权限与审计。
- 保持易扩展与数据一致性，采用约束、索引与范围排除防止冲突。

## 设计原则
- 使用 UTC 时间；客户端传入 RFC3339（ISO8601）字符串
- 防重叠：预约表用 `TSTZRANGE` + `EXCLUDE USING gist`
- 设施类型：使用 CHECK 约束代替枚举，便于扩展
- 外键采用 `ON DELETE CASCADE`，重要表含 `created_at/updated_at`

## 表总览
1. `facilities` — 设施（羽毛球、网球、健身房、多功能厅、其他）
2. `resource_units` — 具体可预约单元（球场/区域）
3. `bookings` — 预约（时间范围、状态、价格、防重叠）
4. `pricing_rules` — 价格规则（按设施类型、星期与小时段）
5. `blackouts` — 封场（设施或单元级范围禁用）
6. `opening_hours` — 营业时间（设施维度，星期日到星期六）
7. `profiles` — 用户资料（映射 Supabase `auth.users`）与角色
8. `facility_admins` — 设施管理员授权映射（按设施）
9. `audit_logs` — 审计日志（便于追踪关键变更）
10. `reservation_policies` — 预约规则（最短/最长时长、粒度、取消时间等）

## DDL（可直接作为迁移执行）
```sql
-- 基础扩展：支持排除约束所需 gist 能力
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- 1) 设施
CREATE TABLE IF NOT EXISTS facilities (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('badminton','tennis','gym','multipurpose','other')),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (name, type)
);

-- 2) 资源单元（球场/区域）
CREATE TABLE IF NOT EXISTS resource_units (
  id BIGSERIAL PRIMARY KEY,
  facility_id BIGINT NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  label TEXT NOT NULL,
  capacity INT NOT NULL DEFAULT 1 CHECK (capacity >= 1),
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (facility_id, label)
);
CREATE INDEX IF NOT EXISTS idx_units_facility ON resource_units(facility_id);

-- 3) 预约（使用范围避免重叠）
CREATE TABLE IF NOT EXISTS bookings (
  id BIGSERIAL PRIMARY KEY,
  resource_unit_id BIGINT NOT NULL REFERENCES resource_units(id) ON DELETE CASCADE,
  user_id UUID NOT NULL, -- Supabase auth.users.id
  time_range TSTZRANGE NOT NULL,
  start_time TIMESTAMPTZ GENERATED ALWAYS AS (lower(time_range)) STORED,
  end_time TIMESTAMPTZ GENERATED ALWAYS AS (upper(time_range)) STORED,
  status TEXT NOT NULL CHECK (status IN ('pending','confirmed','cancelled')) DEFAULT 'confirmed',
  price NUMERIC(10,2) NOT NULL DEFAULT 0,
  currency TEXT NOT NULL DEFAULT 'USD',
  notes TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (lower(time_range) < upper(time_range))
);
-- 同一单元时间不重叠
ALTER TABLE bookings ADD CONSTRAINT bookings_no_overlap EXCLUDE USING gist (
  resource_unit_id WITH =,
  time_range WITH &&
);
CREATE INDEX IF NOT EXISTS idx_bookings_user ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_unit ON bookings(resource_unit_id);

-- 4) 价格规则（简单时段计价）
CREATE TABLE IF NOT EXISTS pricing_rules (
  id BIGSERIAL PRIMARY KEY,
  facility_type TEXT NOT NULL CHECK (facility_type IN ('badminton','tennis','gym','multipurpose','other')),
  day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sunday
  start_hour INT NOT NULL CHECK (start_hour BETWEEN 0 AND 23),
  end_hour INT NOT NULL CHECK (end_hour BETWEEN 1 AND 24 AND end_hour > start_hour),
  price_per_hour NUMERIC(10,2) NOT NULL CHECK (price_per_hour >= 0),
  currency TEXT NOT NULL DEFAULT 'USD',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (facility_type, day_of_week, start_hour, end_hour)
);

-- 5) 封场（设施/单元级别的禁用时间）
CREATE TABLE IF NOT EXISTS blackouts (
  id BIGSERIAL PRIMARY KEY,
  facility_id BIGINT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  resource_unit_id BIGINT NULL REFERENCES resource_units(id) ON DELETE CASCADE,
  time_range TSTZRANGE NOT NULL,
  reason TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (facility_id IS NOT NULL OR resource_unit_id IS NOT NULL)
);
CREATE INDEX IF NOT EXISTS idx_blackouts_range ON blackouts USING gist (time_range);

-- 6) 营业时间（每个设施按星期配置）
CREATE TABLE IF NOT EXISTS opening_hours (
  id BIGSERIAL PRIMARY KEY,
  facility_id BIGINT NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
  open_time TIME NOT NULL,
  close_time TIME NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CHECK (open_time < close_time),
  UNIQUE (facility_id, day_of_week)
);

-- 7) 用户资料（镜像），包含角色
CREATE TABLE IF NOT EXISTS profiles (
  user_id UUID PRIMARY KEY,
  display_name TEXT,
  phone TEXT,
  role TEXT NOT NULL DEFAULT 'user' CHECK (role IN ('user','admin')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_profiles_role ON profiles(role);

-- 8) 设施管理员映射（细粒度到设施）
CREATE TABLE IF NOT EXISTS facility_admins (
  user_id UUID NOT NULL,
  facility_id BIGINT NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, facility_id)
);

-- 9) 审计日志（可选，但推荐）
CREATE TABLE IF NOT EXISTS audit_logs (
  id BIGSERIAL PRIMARY KEY,
  actor_user_id UUID NULL,
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id BIGINT NULL,
  payload JSONB NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_audit_entity ON audit_logs(entity_type, entity_id);

-- 10) 预约策略（时长限制、粒度、提前预订天数、取消截止等）
CREATE TABLE IF NOT EXISTS reservation_policies (
  id BIGSERIAL PRIMARY KEY,
  facility_type TEXT NOT NULL CHECK (facility_type IN ('badminton','tennis','gym','multipurpose','other')),
  min_duration_minutes INT NOT NULL CHECK (min_duration_minutes > 0),
  max_duration_minutes INT NOT NULL CHECK (max_duration_minutes >= min_duration_minutes),
  slot_granularity_minutes INT NOT NULL CHECK (slot_granularity_minutes >= 5),
  advance_booking_days INT NOT NULL CHECK (advance_booking_days >= 0),
  cancellation_cutoff_minutes INT NOT NULL CHECK (cancellation_cutoff_minutes >= 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (facility_type)
);
```

## 关系与约束说明
- `resource_units.facility_id → facilities.id`（级联删除）
- `bookings.resource_unit_id → resource_units.id`（级联删除）；`bookings.user_id → auth.users.id`（由 Supabase 管理）
- `opening_hours.facility_id → facilities.id`；唯一约束每周仅一条
- `blackouts` 至少指定 `facility_id` 或 `resource_unit_id`；支持范围查询索引
- `reservation_policies` 与 `pricing_rules` 都按 `facility_type` 唯一

## 索引建议
- `bookings(user_id)`、`bookings(resource_unit_id)`
- `blackouts USING gist (time_range)`、`bookings` 由排除约束自带 gist 索引
- `resource_units(facility_id)`、`profiles(role)`、`audit_logs(entity_type,entity_id)`

## 数据填充建议（与现有 seed 对齐）
- `facilities`: Badminton, Tennis, Gym, Multipurpose Hall, Other
- `resource_units`: Badminton Court 1–8，Tennis Court 1–8；Gym Area 1；Multipurpose Area 1；Other Unit 1
- 可选：为 `opening_hours` 初始化 08:00–22:00（周一至周日）
- 可选：添加 `reservation_policies`（如羽毛球最短 60、最长 120 分钟，粒度 30 分钟）

## RLS（行级安全）建议（Supabase）
- `bookings`: 允许用户访问自己的记录；管理员可访问全部；创建/改签/取消需匹配 `user_id`
- `facility_admins` & `profiles`: 仅管理员或自身可读写
- `audit_logs`: 仅管理员可读

## 后续扩展（非必须）
- 支付表：`payments`、`refunds`、`invoices`
- 通知表：`notifications`（渠道、内容、状态）
- 资源日历缓存：`availability_cache`（提升大规模查询性能）

确认该表结构后，我将把这些 DDL 编入迁移文件并补充 `opening_hours` 与 `reservation_policies` 的初始 seed。