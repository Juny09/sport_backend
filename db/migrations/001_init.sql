-- 数据库初始化（中文注释）：创建核心表与约束，禁止同一资源时间段重叠

-- 需要 btree_gist 扩展以支持排除约束（EXCLUDE）
CREATE EXTENSION IF NOT EXISTS btree_gist;

-- 设施表
CREATE TABLE IF NOT EXISTS facilities (
  id BIGSERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  type TEXT NOT NULL CHECK (type IN ('badminton','tennis','gym','multipurpose','other')),
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- 可预约单元（球场/区域）
CREATE TABLE IF NOT EXISTS resource_units (
  id BIGSERIAL PRIMARY KEY,
  facility_id BIGINT NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  label TEXT NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

-- 预约表，使用时间范围避免重叠
CREATE TABLE IF NOT EXISTS bookings (
  id BIGSERIAL PRIMARY KEY,
  resource_unit_id BIGINT NOT NULL REFERENCES resource_units(id) ON DELETE CASCADE,
  user_id UUID NOT NULL, -- Supabase auth.users.id
  time_range TSTZRANGE NOT NULL,
  start_time TIMESTAMPTZ GENERATED ALWAYS AS (lower(time_range)) STORED,
  end_time TIMESTAMPTZ GENERATED ALWAYS AS (upper(time_range)) STORED,
  status TEXT NOT NULL CHECK (status IN ('pending','confirmed','cancelled')) DEFAULT 'confirmed',
  price NUMERIC(10,2) DEFAULT 0,
  notes TEXT,
  CHECK (lower(time_range) < upper(time_range))
);

-- 排除约束：同一资源单元的时间范围不能重叠（&&）
ALTER TABLE bookings ADD CONSTRAINT bookings_no_overlap EXCLUDE USING gist (
  resource_unit_id WITH =,
  time_range WITH &&
);

-- 价格规则（简单版）
CREATE TABLE IF NOT EXISTS pricing_rules (
  id BIGSERIAL PRIMARY KEY,
  facility_type TEXT NOT NULL CHECK (facility_type IN ('badminton','tennis','gym','multipurpose','other')),
  day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sunday
  start_hour INT NOT NULL CHECK (start_hour BETWEEN 0 AND 23),
  end_hour INT NOT NULL CHECK (end_hour BETWEEN 1 AND 24 AND end_hour > start_hour),
  price_per_hour NUMERIC(10,2) NOT NULL CHECK (price_per_hour >= 0)
);

-- 封场黑名单
CREATE TABLE IF NOT EXISTS blackouts (
  id BIGSERIAL PRIMARY KEY,
  resource_unit_id BIGINT NULL REFERENCES resource_units(id) ON DELETE CASCADE,
  facility_id BIGINT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  time_range TSTZRANGE NOT NULL,
  reason TEXT
);

-- 用户资料（可选，镜像 Supabase 用户）
CREATE TABLE IF NOT EXISTS profiles (
  user_id UUID PRIMARY KEY,
  display_name TEXT,
  phone TEXT
);

