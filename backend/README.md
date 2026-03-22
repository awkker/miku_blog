# NanaMiku Blog Backend

Go + Hertz + PostgreSQL + Redis

## Tech Stack

- **Framework**: [Hertz](https://github.com/cloudwego/hertz) (CloudWeGo)
- **Database**: PostgreSQL (full-text search via GIN index)
- **Cache/Rate Limit**: Redis (Lua script rate limiting)
- **ORM**: [sqlc](https://sqlc.dev/) (type-safe SQL -> Go codegen)
- **Auth**: JWT (access + refresh token rotation) + bcrypt

## Project Structure

```
backend/
├── main.go                     # Entry point
├── cmd/
│   ├── migrate/main.go         # DB migration runner
│   └── seed/main.go            # Admin user seeder
├── biz/
│   ├── bootstrap/              # Config, DB, Redis, Router
│   ├── dto/                    # API response structs
│   ├── errcode/                # Error code constants
│   ├── handler/
│   │   ├── admin/              # Admin-only handlers (auth, dashboard, moderation, posts, friends)
│   │   └── public/             # Public handlers (health, guestbook, moments, friends, posts, comments)
│   ├── jobs/                   # Background jobs (health check, scheduled publish)
│   ├── middleware/             # Recovery, RequestID, Logger, CORS, Auth, Visitor, RateLimit
│   └── service/                # Business logic layer
├── query/                      # sqlc generated Go code (DO NOT EDIT)
├── sql/
│   ├── migrations/             # SQL migration files
│   └── queries/                # SQL query files for sqlc
├── docker-compose.yml
├── sqlc.yaml
└── .env.example
```

## Quick Start

### 1. Start dependencies

```bash
docker-compose up -d
```

### 2. Configure environment

```bash
cp .env.example .env
# Edit .env as needed (especially JWT_SECRET for production)
# Optional: set GEOIP_DB_PATH to GeoLite2-City.mmdb absolute path
```

### 3. Run migrations

```bash
go run cmd/migrate/main.go
```

Migration runner now tracks versions in `schema_migrations`.
- Re-running `go run cmd/migrate/main.go` is safe.
- Already applied versions are skipped automatically.
- Legacy databases without migration records are backfilled on duplicate-object detection.

### GeoIP fallback for analytics (optional)

When request headers do not provide geo fields (for example local development), analytics can fallback to local MaxMind DB:

1. Download `GeoLite2-City.mmdb` from MaxMind.
2. Set `GEOIP_DB_PATH` in `.env` to an absolute file path.
3. Restart backend service.

Notes:
- Geo headers from reverse proxies (Cloudflare/Vercel) still take precedence.
- Private/loopback IPs (e.g. `127.0.0.1`, `192.168.x.x`) remain `Unknown`.

#### Verify locally

Use a public IP in `X-Forwarded-For` to simulate a real client (and **do not** send geo headers like `CF-IPCountry`):

```bash
curl -X POST 'http://localhost:8080/api/v1/analytics/collect' \
  -H 'Content-Type: application/json' \
  -H 'X-Forwarded-For: 8.8.8.8' \
  -d '{
    "session_key":"geoip-local-test",
    "path":"/geoip-test",
    "title":"GeoIP Local Test",
    "referrer":"",
    "timezone":"Asia/Shanghai",
    "language":"zh-CN",
    "occurred_at":"2026-03-19T00:00:00Z"
  }'
```

Then check the latest rows in PostgreSQL:

```bash
docker exec -it miku-blog-postgres psql -U miku -d miku_blog -c \
"SELECT country_code, region, city, occurred_at FROM analytics_pageviews WHERE path='/geoip-test' ORDER BY occurred_at DESC LIMIT 5;"
```

Expected:
- If GeoIP fallback works, `country_code` should be a two-letter code (e.g. `US`) instead of `ZZ`.
- If it is still `ZZ`, verify `GEOIP_DB_PATH` and restart backend.

### 4. Seed admin user

```bash
go run cmd/seed/main.go [password]
# Default: username=admin, password=admin123
```

### 5. Start server

```bash
go run main.go
```

Server runs at `http://localhost:8080`.

## API Endpoints

### Public

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/health` | Health check |
| POST | `/api/v1/auth/login` | Admin login |
| POST | `/api/v1/auth/refresh` | Refresh token |
| POST | `/api/v1/auth/logout` | Logout |
| GET | `/api/v1/guestbook/messages` | List guestbook messages |
| POST | `/api/v1/guestbook/messages` | Create guestbook message |
| POST | `/api/v1/guestbook/messages/:id/vote` | Vote on message |
| GET | `/api/v1/moments` | List moments |
| GET | `/api/v1/moments/latest` | Latest moments |
| POST | `/api/v1/moments` | Create moment |
| POST | `/api/v1/moments/:id/like` | Toggle moment like |
| POST | `/api/v1/moments/:id/repost` | Toggle moment repost |
| POST | `/api/v1/moments/:id/comments` | Comment on moment |
| POST | `/api/v1/moments/comments/:id/like` | Like moment comment |
| GET | `/api/v1/friends` | List friend links |
| POST | `/api/v1/analytics/collect` | Collect pageview analytics |
| GET | `/api/v1/posts` | List published posts |
| GET | `/api/v1/posts/hot` | Hot posts |
| GET | `/api/v1/posts/search?q=` | Full-text search |
| GET | `/api/v1/posts/:slug` | Get post by slug |
| POST | `/api/v1/posts/:id/like` | Toggle post like |
| GET | `/api/v1/posts/:id/comments` | List post comments |
| POST | `/api/v1/posts/:id/comments` | Create post comment |

### Admin (requires `Authorization: Bearer <token>`)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/auth/me` | Current admin info |
| GET | `/api/v1/admin/dashboard/stats` | Dashboard stats |
| GET | `/api/v1/admin/dashboard/trend/views` | View trend |
| GET | `/api/v1/admin/dashboard/trend/comments` | Comment trend |
| GET | `/api/v1/admin/dashboard/trend/likes` | Like trend |
| GET | `/api/v1/admin/dashboard/analytics?range=24h|7d|30d&offset=0` | Analytics overview |
| GET | `/api/v1/admin/comments` | List all comments |
| POST | `/api/v1/admin/comments/:id/approve` | Approve comment |
| POST | `/api/v1/admin/comments/:id/reject` | Reject comment |
| DELETE | `/api/v1/admin/comments/:id` | Delete comment |
| GET | `/api/v1/admin/moderation/rate-limit-metrics?minutes=60` | Rate limit metrics |
| GET | `/api/v1/admin/audit-logs` | Audit logs |
| GET | `/api/v1/admin/friends` | List all friends |
| POST | `/api/v1/admin/friends` | Create friend link |
| PUT | `/api/v1/admin/friends/:id` | Update friend link |
| DELETE | `/api/v1/admin/friends/:id` | Delete friend link |
| GET | `/api/v1/admin/posts` | List all posts |
| GET | `/api/v1/admin/posts/:id` | Get post detail |
| POST | `/api/v1/admin/posts` | Create post |
| PUT | `/api/v1/admin/posts/:id` | Update post |
| POST | `/api/v1/admin/posts/:id/publish` | Publish post |
| POST | `/api/v1/admin/posts/:id/unpublish` | Unpublish post |
| POST | `/api/v1/admin/posts/:id/schedule` | Schedule post |
| DELETE | `/api/v1/admin/posts/:id` | Delete post |
| GET | `/api/v1/admin/moments` | List all moments |
| POST | `/api/v1/admin/moments` | Create moment |
| PUT | `/api/v1/admin/moments/:id` | Update moment |
| POST | `/api/v1/admin/moments/:id/publish` | Publish moment |
| POST | `/api/v1/admin/moments/:id/unpublish` | Unpublish moment |
| POST | `/api/v1/admin/moments/:id/schedule` | Schedule moment |
| GET | `/api/v1/admin/backup/export?format=json|sql` | Export backup |

## Development

### Regenerate sqlc code

```bash
sqlc generate
```

### Build

```bash
go build -o miku-blog .
```
