# Goal / 范围
- Build a production-ready backend for a sport booking system in Go with Supabase (Postgres + Auth).
- Phase 1 includes facilities and courts, availability search, booking CRUD, admin seeding for: badminton (8 courts), tennis (8 courts), gym, multipurpose hall, and “other sport” types.
- Deliver a runnable server, SQL migrations, seed data, and a Postman v2.1 collection.

## Tech Stack / 技术栈
- Language: Go (>= 1.21)
- Web: `gin` (fast, minimal), or `echo` if preferred — default to `gin`
- DB: Supabase Postgres via `pgx` driver (`github.com/jackc/pgx/v5`), migrations using SQL files
- Auth: Supabase JWT verification (`github.com/golang-jwt/jwt/v5`), header `Authorization: Bearer <token>`
- Config: `.env` with Supabase connection string and JWT secret
- Logging: `log/slog`

## Project Structure / 项目结构
- `cmd/server/main.go`: server entry
- `internal/config`: env load & validation
- `internal/http`: routers, middlewares, DTOs
- `internal/handlers`: facilities, availability, bookings, admin
- `internal/repo`: DB access (pgx queries)
- `internal/service`: domain logic (availability, pricing, overlap checks)
- `internal/auth`: Supabase JWT verify, user context
- `db/migrations`: SQL migrations
- `db/seed`: initial data for courts
- `postman/collection.json`: Postman v2.1 collection

## Data Model / 数据模型
- `facilities`
  - `id BIGSERIAL PK`, `name TEXT`, `type TEXT` (enum-like: badminton|tennis|gym|multipurpose|other), `is_active BOOL`
- `resource_units`
  - Each bookable unit (court, gym area, hall)
  - `id BIGSERIAL PK`, `facility_id FK`, `label TEXT`, `is_active BOOL`
- `bookings`
  - `id BIGSERIAL PK`, `resource_unit_id FK`, `user_id UUID` (from `auth.users.id`), `time_range TSTZRANGE`, `start_time TIMESTAMPTZ`, `end_time TIMESTAMPTZ`, `status TEXT` (pending|confirmed|cancelled), `price NUMERIC(10,2)`, `notes TEXT`
  - Exclusion constraint to prevent overlaps: `EXCLUDE USING gist (resource_unit_id WITH =, time_range WITH &&)`
- `pricing_rules` (phase 1 minimal)
  - `id PK`, `facility_type TEXT`, `day_of_week INT`, `start_hour INT`, `end_hour INT`, `price_per_hour NUMERIC(10,2)`
- `blackouts`
  - `id PK`, `resource_unit_id FK NULL`, `facility_id FK NULL`, `time_range TSTZRANGE`, `reason TEXT`
- `profiles`
  - Mirror of Supabase `auth.users` profile (optional): `user_id UUID PK`, `display_name TEXT`, `phone TEXT`

## Migrations & Seed / 迁移与初始化
- Provide SQL migration files to create tables, indexes, exclusion constraints, and useful CHECK constraints (e.g., `start_time < end_time`).
- Seed script to insert facilities:
  - Badminton (8 units: Court 1–8)
  - Tennis (8 units: Court 1–8)
  - Gym (1 unit or configurable units)
  - Multipurpose Hall (1 unit)
  - Other (example unit for extensibility)

## Auth & Security / 鉴权与安全
- Authentication via Supabase JWT
  - Verify using Supabase JWT secret; extract `sub` as `user_id`
- Authorization
  - Public: `GET /health`, `GET /facilities`, `GET /availability`
  - Auth required: booking operations (`/bookings`), `GET /me`
  - Admin: facility/unit management, pricing, blackouts (role check via JWT claims or a `profiles.role` table)
- Input validation for time ranges, duration, and unit existence
- Store all timestamps in UTC; accept client time with timezone and normalize

## Core Endpoints / 核心接口
- `GET /health` — service health
- `GET /me` — current user profile (requires auth)

- Facilities / 场馆
  - `GET /facilities` — list facilities
  - `GET /facilities/:id` — facility detail
  - `GET /facilities/:id/units` — list units
  - `POST /facilities` — admin create facility
  - `POST /facilities/:id/units` — admin create units
  - `PATCH /units/:id` — admin update unit status

- Availability / 可用性
  - `GET /availability?facility_type=badminton&date=YYYY-MM-DD&duration=60` — per-unit free slots
  - Optional: `unit_id` filter, `start_hour`, `end_hour`

- Bookings / 预约
  - `POST /bookings` — create booking
    - body: `resource_unit_id`, `start_time`, `end_time`, optional `notes`
  - `GET /bookings/:id` — booking detail (owner/admin)
  - `GET /bookings?mine=true` — list my bookings
  - `PATCH /bookings/:id/cancel` — cancel booking
  - `PATCH /bookings/:id/reschedule` — change time (re-check overlap)

- Admin / 管理
  - `GET /admin/bookings?facility_type=...&date=...` — search
  - `POST /pricing_rules` — add pricing rule
  - `POST /blackouts` — add blackout range

## Booking Logic / 预约逻辑
- Overlap prevention enforced by DB exclusion constraint; additionally checked at service level for clear errors
- Slot granularity: 30 minutes by default (可配置)
- Min/Max booking duration per facility type (e.g., badminton 60–120min)
- Pricing computed from rules by summing hourly segments; fallback default price if no rule
- Cancellation policy: allow cancel before start; status transitions tracked

## Availability Algorithm / 可用性计算
- For a facility type and date:
  - Load active units
  - Build candidate slots based on opening hours (configurable per facility) and requested duration
  - Subtract `bookings` overlaps and `blackouts`
  - Return free ranges per unit with price estimate

## Postman Collection / Postman 集合
- Provide `postman/collection.json` (v2.1) including folders:
  - Health, Auth, Facilities, Availability, Bookings, Admin
- Pre-request script to inject `{{access_token}}` into `Authorization` header for protected endpoints
- Example bodies and test scripts that assert status codes and basic response shapes

## Local Setup (macOS) / 本地环境
- `.env` keys:
  - `SUPABASE_DB_URL=postgres://USER:PASSWORD@HOST:PORT/postgres`
  - `SUPABASE_JWT_SECRET=...`
  - `PORT=8080`
- Steps:
  - Install Go and run `go mod init` + add deps
  - Apply migrations to Supabase via SQL editor or `psql`
  - Run server: `go run ./cmd/server`

## Testing & Verification / 测试
- Unit tests for services: availability computation, overlap detection
- Handler tests using `httptest` with a test DB (transaction rollback per test)
- Smoke test: `GET /health`, seed presence, booking create and cancel flow

## Delivery Artifacts / 交付物
- Go source code with clear package structure（所有关键代码含中文注释）
- `db/migrations/*.sql` and `db/seed/*.sql`
- `postman/collection.json`
- README with setup and Supabase config notes（含中文说明）

## Next Iteration / 下一步
- Payments integration (optional)
- Opening hours per facility, holidays calendar
- Rate limiting / quotas per user
- Webhooks for notifications

请确认以上方案；确认后我将开始实现代码、SQL迁移、种子数据，以及Postman集合，并在关键逻辑处加入通俗易懂的中文注释。