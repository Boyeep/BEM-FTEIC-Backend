# BEM FTEIC Backend

REST API production BEM FTEIC untuk authorization admin, validasi request, mutation konten, upload gambar, profil, dan signup whitelist. Supabase menyediakan PostgreSQL dan autentikasi; backend memverifikasi access token Supabase sebelum menjalankan operasi terproteksi.

## Production

- API: <https://api.bem-fteic.com>
- Health: <https://api.bem-fteic.com/health>
- Readiness: <https://api.bem-fteic.com/ready>
- Runtime: Docker Compose di VPS, port grup `8002`
- Database dan Auth: Supabase
- Deployment: GitHub Actions setiap push ke `main`

## Arsitektur

```text
HTTP request -> middleware -> handler -> service -> repository -> PostgreSQL
```

```text
cmd/server/             entrypoint
config/                 environment dan database config
database/migrations/    migrasi SQL berurutan
internal/dto/           validasi request
internal/handlers/      transport HTTP
internal/service/       business logic
internal/repository/    akses database
internal/middleware/    auth, role, CORS, rate limit, request ID
pkg/apperr/             error aplikasi
pkg/response/           format response standar
```

## Keamanan

- Verifikasi token Supabase dan authorization role `admin`
- PostgreSQL RLS dan grant terbatas
- DTO/query validation, pagination server-side, serta response/error terstandardisasi
- Request ID, security headers, rate limiting, dan HTTP timeouts
- Batas ukuran upload dan pemeriksaan MIME
- Rollback deployment jika migration atau readiness gagal

## Menjalankan Lokal

Prasyarat: Go 1.25.12+ dan PostgreSQL/Supabase.

```bash
cp .env.example .env
go mod download
go run ./cmd/server
```

PowerShell:

```powershell
Copy-Item .env.example .env
go mod download
go run ./cmd/server
```

Server default tersedia di <http://localhost:8080>.

## Environment

Gunakan `.env.example` sebagai acuan:

```env
APP_ENV=production
GIN_MODE=release
PORT=8080
DATABASE_URL=postgresql://...
SUPABASE_URL=https://PROJECT_REF.supabase.co
SUPABASE_ANON_KEY=...
ALLOWED_ORIGINS=https://bem-fteic.com,https://www.bem-fteic.com
PUBLIC_API_URL=https://api.bem-fteic.com
```

Jangan commit database password, service-role key, `.env`, atau SSH private key.

## Endpoint

Public:

```text
GET /health
GET /ready
GET /blogs/
GET /blogs/:id
GET /events/
GET /events/:id
GET /gallery/
GET /gallery/:id
GET /profiles?ids=<uuid,uuid>
POST /visitors
GET /visitors/count
GET /uploads/:filename
```

Authenticated user:

```text
GET /me
PUT /me
```

Admin:

```text
GET|POST       /admin/blogs
GET|PUT|DELETE /admin/blogs/:id
GET|POST       /admin/events
GET|PUT|DELETE /admin/events/:id
GET|POST       /admin/gallery
GET|PUT|DELETE /admin/gallery/:id
GET|POST       /admin/whitelist
DELETE         /admin/whitelist/:id
POST|DELETE    /uploads/images
```

Endpoint terproteksi memerlukan:

```http
Authorization: Bearer <supabase_access_token>
```

Format response:

```json
{
  "success": true,
  "data": {},
  "request_id": "..."
}
```

Endpoint list mengembalikan:

```json
{
  "items": [],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total_items": 0,
    "total_pages": 1,
    "has_next_page": false,
    "has_previous_page": false
  }
}
```

Blog menggunakan publication status `DRAFT | PUBLISHED | ARCHIVED`. Event memisahkan lifecycle `UPCOMING | ONGOING | ENDED` dari `publication_status` dengan nilai `DRAFT | PUBLISHED | ARCHIVED`.

## Database dan Migration

`database/migrations` adalah sumber skema resmi. Migration mengelola tabel, role admin, signup whitelist, visitor analytics, Auth hook, profile trigger, content contract, grant, dan RLS.

Di Supabase Dashboard aktifkan Before User Created hook:

```text
Schema: public
Function: hook_validate_signup
```

## Test

```bash
go test ./...
go vet ./...
```

Integration test PostgreSQL:

```bash
TEST_DATABASE_URL=postgresql://... go test ./database -run TestMigrateIntegration -v
```

Gunakan database test terisolasi karena integration test mengubah objek database.

## Deployment

```bash
sudo docker compose -p group1 up -d
sudo docker compose -p group1 ps
```

Push ke `main` menjalankan test, membangun image bertag commit SHA, membuat backup upload, menjalankan migration, memeriksa readiness, dan rollback jika gagal.

Mapping production:

```text
api.bem-fteic.com -> 127.0.0.1:8002
```

Volume `backend_uploads` menyimpan file upload; workflow deployment membuat backup ke direktori `./backups`.
