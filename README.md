# Sport Booking Backend (Go + Supabase)

中文说明：这是一个使用 Go 构建、数据库采用 Supabase Postgres 的体育场馆预约后端。支持羽毛球（8）、网球（8）、健身房、多功能厅以及其它运动类型。

## Quick Start
- Requirements: Go >= 1.22, Supabase project
- Env:
  - `SUPABASE_DB_URL=postgres://USER:PASSWORD@HOST:PORT/postgres` (或使用 `DATABASE_URL` 替代)
  - `SUPABASE_JWT_SECRET=...` (Supabase 项目设置中的 JWT 密钥)
  - `PORT=8080`

### Setup Steps (macOS)
1. 安装 Go 并拉取依赖：
   - `go mod tidy`
2. 在 Supabase SQL Editor 中执行 `db/migrations/001_init.sql`
3. 执行 `db/seed/001_seed.sql` 生成默认设施和球场
4. 运行：`go run ./cmd/server`
5. Postman 集合：导入 `postman/collection.json`，将 `{{base_url}}` 设置为 `http://localhost:8080`，`{{access_token}}` 填入 Supabase 登录获得的 JWT

### 快速启动（脚本方式）
- 复制环境模板：
  - `cp .env.example .env`
- 编辑 `.env`，填入数据库与 JWT 密钥
- 运行后端：
  - `./scripts/run_server.sh`

### Flutter 前端配置
- 复制配置模板：
  - `cp frontend_sport/lib/config.example.dart frontend_sport/lib/config.dart`
- 编辑 `frontend_sport/lib/config.dart`，填入 Supabase 项目 URL 与 Anon Key，以及后端地址（iOS 用 `localhost`，Android 用 `10.0.2.2`）
- 运行前端：
  - `cd frontend_sport && flutter run`

## Endpoints
- `GET /health` 健康检查
- `GET /me` 当前用户（需授权）
- `GET /facilities` 列出设施
- `GET /facilities/:id` 设施详情
- `GET /facilities/:id/units` 列出指定设施的场地单元
- `POST /facilities` 创建设施（管理员）
- `POST /facilities/:id/units` 创建单元（管理员）
- `PATCH /units/:id` 更新单元状态（管理员）
- `GET /availability?facility_type=badminton&date=YYYY-MM-DD&duration=60` 查询可用时段
- `POST /bookings` 创建预约（需授权）
- `GET /bookings/:id` 预约详情（本人或管理员）
- `GET /bookings?mine=true` 我的预约列表（需授权）
- `PATCH /bookings/:id/cancel` 取消预约（本人或管理员）
- `PATCH /bookings/:id/reschedule` 改签预约（本人或管理员）
- `GET /admin/bookings?facility_type=...&date=...` 管理员查询预约
- `POST /pricing_rules` 添加价格规则（管理员）
- `POST /blackouts` 添加封场时间（管理员）

## Design Notes
- 防重叠：`bookings` 使用 `TSTZRANGE` + `EXCLUDE USING gist` 防止同一场地时间冲突
- 时间：后端统一使用 UTC，客户端传入 ISO8601 字符串（RFC3339）
- 鉴权：使用 Supabase JWT，`Authorization: Bearer <token>`；`/me`、预订相关接口需要登录；管理接口要求 `role=admin`
- 可用性：支持 `opening_hours` 配置营业时间（默认 08:00-22:00），基于当天预订与封场计算空闲时段
- 角色与权限：`profiles.role` 以及 `facility_admins` 支持设施级管理员
- 预约策略：`reservation_policies` 统一配置每种设施类型的最短/最长时长与最小粒度

## Database Tables（数据库表）
- facilities：设施基础信息（类型、启用）
- resource_units：具体场地或区域（唯一 label、容量、启用）
- bookings：预约记录（时间范围、价格、状态）
- pricing_rules：价格规则（按设施类型、星期和小时段）
- blackouts：封场记录（设施或单元级）
- opening_hours：营业时间（每设施每日开闭）
- profiles：用户资料与角色（映射 Supabase 用户）
- facility_admins：设施管理员映射
- audit_logs：审计日志（关键操作记录）
- reservation_policies：预约策略（时长限制、粒度、提前预订与取消截止）

## Next
- 开放时段配置、节假日
- 价格计算（当前仅规则存储，未集成计算）
- 支付集成、配额和限流
### 环境变量示例（macOS zsh）
- 使用 `SUPABASE_DB_URL`：
  - `export SUPABASE_DB_URL='postgres://USER:PASSWORD@HOST:5432/postgres'`
- 或使用 `DATABASE_URL`（后端已支持回退读取）：
  - `export DATABASE_URL='postgres://USER:PASSWORD@HOST:5432/postgres'`
- 设置 JWT：
  - `export SUPABASE_JWT_SECRET='你的JWT密钥'`
