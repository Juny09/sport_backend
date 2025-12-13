-- 扩展 schema：营业时间、管理员映射、审计日志、预约策略，以及 profiles 角色

-- 6) 营业时间：每个设施可配置每日开闭时间
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

-- 7) profiles 扩展：增加角色与创建时间
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'user';
ALTER TABLE profiles ADD COLUMN IF NOT EXISTS created_at TIMESTAMPTZ NOT NULL DEFAULT now();
DO $$ BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_constraint WHERE conname = 'profiles_role_check'
  ) THEN
    ALTER TABLE profiles ADD CONSTRAINT profiles_role_check CHECK (role IN ('user','admin'));
  END IF;
END $$;
CREATE INDEX IF NOT EXISTS idx_profiles_role ON profiles(role);

-- 8) 设施管理员映射：细粒度授权到设施
CREATE TABLE IF NOT EXISTS facility_admins (
  user_id UUID NOT NULL,
  facility_id BIGINT NOT NULL REFERENCES facilities(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, facility_id)
);

-- 9) 审计日志：记录关键操作
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

-- 10) 预约策略：时长限制、粒度、提前预订、取消截止
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

