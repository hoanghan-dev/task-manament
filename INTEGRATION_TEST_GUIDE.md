# Integration Test Guide

## Tổng Quan

Integration test kiểm tra toàn bộ API flow thực tế:

```
Register → Login → Get Workspace → Create Task → Create Comment → Get Comments
```

Test sử dụng `net/http/httptest` với **router thật** (gin.Engine), **database thật** (PostgreSQL), và **Redis thật**.
Không mock handler, service, hay repository — tất cả request đi qua full stack y hệt production.

---

## Flow Được Test

### Happy Path (`TestAPIFlow_RegisterLoginCreateProjectCreateTaskComment`)

| Step | API | Method | Expected Status |
|------|-----|--------|-----------------|
| 1 | `/api/auth/register` | POST | 201 Created |
| 2 | `/api/auth/login` | POST | 200 OK |
| 2b | `/api/workspaces/` | GET | 200 OK |
| 3 | `/api/tasks/` | POST | 201 Created |
| 4 | `/api/tasks/:id/comments` | POST | 201 Created |
| 5 | `/api/tasks/:id/comments` | GET | 200 OK |

### Negative Cases

| Test | Mô tả | Expected Status |
|------|--------|-----------------|
| `TestCreateTask_WithoutToken_Returns401` | Tạo task không có token | 401 |
| `TestCreateComment_NonExistingTask_Returns404` | Comment vào task không tồn tại | 404 |
| `TestCreateComment_WithoutToken_Returns401` | Tạo comment không có token | 401 |
| `TestRegister_DuplicateEmail_Returns409` | Đăng ký email trùng | 409 |

---

## Chuẩn Bị Database Test

### 1. Tạo database test

```sql
-- Chạy trong psql hoặc pgAdmin (kết nối vào PostgreSQL)
CREATE DATABASE task_management_test;
```

### 2. Migration tự động

Test sẽ **tự động chạy migration** (CREATE TABLE IF NOT EXISTS) khi khởi động.
Không cần chạy `migrate` CLI riêng.

### 3. Cleanup

Sau mỗi test run, tất cả bảng được TRUNCATE với `RESTART IDENTITY CASCADE`:
- `comments`
- `notifications`
- `tasks`
- `workspaces`
- `users`

---

## Biến Môi Trường

Có 2 cách cấu hình:

### Cách 1: File `.env.test` (khuyến nghị)

Copy file mẫu và chỉnh sửa nếu cần:

```bash
# File: test/integration/.env.test
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=postgres
TEST_DB_PASSWORD=postgres
TEST_DB_NAME=task_management_test
TEST_DB_SSLMODE=disable
SECRET_KEY=integration-test-secret-key-do-not-use-in-prod
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=1
```

### Cách 2: Set trực tiếp trên terminal

**Windows (PowerShell):**
```powershell
$env:TEST_DB_HOST="localhost"
$env:TEST_DB_PORT="5432"
$env:TEST_DB_USER="postgres"
$env:TEST_DB_PASSWORD="postgres"
$env:TEST_DB_NAME="task_management_test"
$env:TEST_DB_SSLMODE="disable"
$env:SECRET_KEY="integration-test-secret-key-do-not-use-in-prod"
$env:REDIS_ADDR="localhost:6379"
$env:REDIS_PASSWORD=""
$env:REDIS_DB="1"
```

**Linux/macOS:**
```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=task_management_test
export TEST_DB_SSLMODE=disable
export SECRET_KEY=integration-test-secret-key-do-not-use-in-prod
export REDIS_ADDR=localhost:6379
export REDIS_PASSWORD=
export REDIS_DB=1
```

---

## Cách Chạy Test

### Chạy tất cả integration test

```bash
cd <project-root>
go test -v -count=1 ./test/integration/...
```

### Chạy một test cụ thể

```bash
go test -v -count=1 -run TestAPIFlow_RegisterLoginCreateProjectCreateTaskComment ./test/integration/...
```

### Chạy chỉ negative cases

```bash
go test -v -count=1 -run "TestCreate|TestRegister_Duplicate" ./test/integration/...
```

> **Lưu ý**: `-count=1` tắt test caching để đảm bảo test luôn chạy lại.

---

## Prerequisites

1. **PostgreSQL** đang chạy (local hoặc Docker)
   - Database `task_management_test` đã được tạo
   - User/password đúng

2. **Redis** đang chạy (local hoặc Docker)
   - Mặc định: `localhost:6379`
   - Test dùng DB=1 để tách biệt khỏi dev (DB=0)

3. **Docker** (nếu dùng docker-compose):
   ```bash
   docker-compose up -d
   # Sau đó tạo database test:
   docker exec -it task-management-postgres psql -U postgres -c "CREATE DATABASE task_management_test;"
   ```

---

## Cấu Trúc File Test

```
test/
└── integration/
    ├── .env.test          # Template biến môi trường
    ├── test_setup.go      # TestMain, DB/Redis/Router setup, migrations, cleanup
    ├── test_helpers.go    # HTTP helpers, JSON parsing, flow helpers
    └── api_flow_test.go   # Test cases (happy path + negative cases)
```

---

## Lưu Ý Quan Trọng

1. **Không dùng database production** — Test dùng database riêng (`task_management_test`).

2. **Redis DB tách biệt** — Test dùng `REDIS_DB=1`, dev thường dùng `REDIS_DB=0`.

3. **SECRET_KEY cho JWT** — Test set SECRET_KEY riêng. Không dùng production secret.

4. **Idempotent** — Mỗi test tự cleanup database (TRUNCATE), có thể chạy nhiều lần liên tiếp.

5. **Không ảnh hưởng production code** — Chỉ thêm file trong `test/integration/` và file hướng dẫn này.

6. **Workspace auto-create** — Khi register, hệ thống tự tạo workspace mặc định. Không cần API riêng để tạo workspace/project.
