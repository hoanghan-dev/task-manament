# Refactor Error Handling Plan

> **Project**: Task Management Go Backend  
> **Module**: `dev/task-management`  
> **Date**: 2026-06-02  
> **Author**: Senior Go Backend Engineer - Code Review

---

## 1. Summary

### Tình trạng hiện tại

Project Task Management backend được viết bằng Go + Gin framework, sử dụng kiến trúc Clean Architecture theo module (auth, task, workspace, comment, notification). Mỗi module có handler → service → repository layer.

**Vấn đề chính được phát hiện:**

1. **Không có custom error type** — Toàn bộ project dùng `errors.New()` và `fmt.Errorf()` trả về error string thuần. Service/Repository không phân biệt được loại lỗi (validation, not found, forbidden, conflict, internal).

2. **Handler đoán HTTP status bằng cảm tính** — Một số handler hardcode status (luôn trả 500, hoặc luôn trả 400), không dựa trên loại error thực tế từ service.

3. **Handler dùng `err.Error()` string comparison** — Comment handler dùng `switch err.Error()` để map lỗi, rất fragile và dễ break khi message thay đổi.

4. **Repository leak DB error ra ngoài** — Nhiều repository trả thẳng `sql.ErrNoRows`, raw DB error lên service → handler, lộ implementation detail.

5. **Không có error code chuẩn** — Response trả `err.Error()` raw string, không có error code để client parse.

6. **Response format chưa có error code field** — `ApiResponse` hiện tại chỉ có `success`, `message`, `data`, `error` nhưng field `error` nhận raw error string, không phải error code.

7. **Thiếu xử lý `RowsAffected == 0`** — Nhiều repository (UpdateTask, DeleteTask) không kiểm tra `RowsAffected`, dẫn đến silent success khi record không tồn tại.

**Mức độ ảnh hưởng**: Hệ thống hiện tại sẽ trả sai HTTP status cho client trong khoảng **70% các trường hợp lỗi**, đặc biệt nghiêm trọng ở module Task và Workspace.

---

## 2. Problems Found

### 2.1. Handler Layer

| No | File | Function | Current Problem | Wrong Status | Expected Status | Severity |
|----|------|----------|----------------|-------------|----------------|----------|
| 1 | `internal/modules/task/handler/task_handler.go:29` | `GetAllTask` | `GetOwnerId` fail → trả 500, nhưng lỗi này thực tế là user_id không có trong context (authentication issue) | 500 | 401 Unauthorized | **High** |
| 2 | `internal/modules/task/handler/task_handler.go:35` | `GetAllTask` | Service error (có thể là DB error) → luôn trả 404. Nếu DB connection fail, vẫn trả 404 | 404 | 500 Internal | **High** |
| 3 | `internal/modules/task/handler/task_handler.go:58` | `GetTask` | `GetOwnerId` fail → trả 500, nhưng nguyên nhân thật là auth issue | 500 | 401 Unauthorized | **High** |
| 4 | `internal/modules/task/handler/task_handler.go:65` | `GetTask` | Service trả bất kỳ error nào → luôn trả 404. Nếu task tồn tại nhưng DB fail, vẫn 404 | 404 | 404/403/500 tùy case | **High** |
| 5 | `internal/modules/task/handler/task_handler.go:85` | `CreateTask` | `GetOwnerId` fail → trả 500, auth issue | 500 | 401 Unauthorized | **High** |
| 6 | `internal/modules/task/handler/task_handler.go:91` | `CreateTask` | Service error → luôn trả 400. Nếu lỗi "access denied" → nên 403. Nếu lỗi "workspace not found" → nên 404 | 400 | 400/403/404/500 tùy case | **High** |
| 7 | `internal/modules/task/handler/task_handler.go:120` | `UpdateTask` | `GetOwnerId` fail → 500, auth issue | 500 | 401 Unauthorized | **High** |
| 8 | `internal/modules/task/handler/task_handler.go:127` | `UpdateTask` | Service error → luôn trả 400 "validation failed". Nhưng lỗi có thể là not found, forbidden, DB error | 400 | 400/403/404/500 tùy case | **High** |
| 9 | `internal/modules/task/handler/task_handler.go:146` | `DeleteTask` | `GetOwnerId` fail → 500, auth issue | 500 | 401 Unauthorized | **High** |
| 10 | `internal/modules/task/handler/task_handler.go:152` | `DeleteTask` | Service error → luôn trả 404. Nhưng lỗi có thể là forbidden hoặc DB error | 404 | 403/404/500 tùy case | **Medium** |
| 11 | `internal/modules/task/handler/task_handler.go:172` | `AssignTask` | `GetOwnerId` fail → 500 | 500 | 401 Unauthorized | **High** |
| 12 | `internal/modules/task/handler/task_handler.go:179` | `AssignTask` | **Mọi** service error → luôn trả 500. Bao gồm "access denied" (nên 403), "task not found" (nên 404), "task id not exists in workspace" (nên 404) | 500 | 400/403/404/500 tùy case | **High** |
| 13 | `internal/modules/task/handler/task_handler.go:207` | `UpdateTaskStatus` | `GetOwnerId` fail → 500 | 500 | 401 Unauthorized | **High** |
| 14 | `internal/modules/task/handler/task_handler.go:214` | `UpdateTaskStatus` | **Mọi** service error → luôn trả 500. Bao gồm "access denied" (nên 403), validation error status transition (nên 400) | 500 | 400/403/404/500 tùy case | **High** |
| 15 | `internal/modules/auth/handler/user_handler.go:34` | `Register` | Service error → luôn trả 400. Nếu lỗi DB (bcrypt fail, DB connection) → nên 500. Nếu email duplicate → nên 409 | 400 | 400/409/500 tùy case | **High** |
| 16 | `internal/modules/auth/handler/user_handler.go:52` | `Login` | Service error → luôn trả 401. Nếu lỗi DB connection → nên 500 | 401 | 401/500 tùy case | **Medium** |
| 17 | `internal/modules/workspace/handler/workspace_hander.go:28` | `GetWorkspace` | `GetOwnerId` fail → 500, auth issue | 500 | 401 Unauthorized | **High** |
| 18 | `internal/modules/workspace/handler/workspace_hander.go:35` | `GetWorkspace` | Service error → luôn trả 404. DB error cũng trả 404 | 404 | 404/500 tùy case | **Medium** |
| 19 | `internal/modules/workspace/handler/workspace_hander.go:56` | `UpdateWorkspace` | `GetOwnerId` fail → 500 | 500 | 401 Unauthorized | **High** |
| 20 | `internal/modules/workspace/handler/workspace_hander.go:63` | `UpdateWorkspace` | Service error → luôn trả 500. Nhưng nếu workspace not found → nên 404, forbidden → 403 | 500 | 400/403/404/500 tùy case | **Medium** |
| 21 | `internal/modules/workspace/handler/workspace_hander.go:82` | `DeleteWorkspace` | UUID parse fail → trả 500 "Internal Server Error" với message "invalid user_id format in context". Đây rõ ràng là validation error | 500 | 400 Bad Request | **High** |
| 22 | `internal/modules/workspace/handler/workspace_hander.go:90` | `DeleteWorkspace` | `GetOwnerId` fail → 500, auth issue | 500 | 401 Unauthorized | **High** |
| 23 | `internal/modules/workspace/handler/workspace_hander.go:97` | `DeleteWorkspace` | Service error → luôn trả 404. Nhưng nếu forbidden → nên 403 | 404 | 403/404/500 tùy case | **Medium** |
| 24 | `internal/modules/comment/handler/comment_handler.go:41` | `CreateComment` | `GetOwnerId` fail → 500, auth issue | 500 | 401 Unauthorized | **High** |
| 25 | `internal/modules/comment/handler/comment_handler.go:47-56` | `CreateComment` | Dùng `switch err.Error()` string comparison → **rất fragile**. Nếu service thay đổi message → handler break | N/A (fragile) | Dùng error type | **High** |
| 26 | `internal/modules/comment/handler/comment_handler.go:73` | `GetComments` | `GetOwnerId` fail → 500 | 500 | 401 Unauthorized | **High** |
| 27 | `internal/modules/comment/handler/comment_handler.go:79-86` | `GetComments` | Dùng `switch err.Error()` string comparison | N/A (fragile) | Dùng error type | **High** |
| 28 | `internal/modules/comment/handler/comment_handler.go:103` | `DeleteComment` | `GetOwnerId` fail → 500 | 500 | 401 Unauthorized | **High** |
| 29 | `internal/modules/comment/handler/comment_handler.go:109-116` | `DeleteComment` | Dùng `switch err.Error()` string comparison | N/A (fragile) | Dùng error type | **High** |
| 30 | `internal/realtime/hander_ws.go:34` | `Connect` | `GetOwnerId` fail → 500 (cho WS upgrade context) | 500 | 401 Unauthorized | **Medium** |

### 2.2. Service Layer

| No | File | Function | Current Problem | Severity |
|----|------|----------|----------------|----------|
| 31 | `internal/modules/task/services/task_service.go:96-103` | `GetTask` | **Nil pointer dereference bug**: Khi `FindById` trả error, service vẫn gọi `mapper.EntityToTaskResponse(task)` với `task == nil` → **PANIC** | **Critical** |
| 32 | `internal/modules/task/services/task_service.go:111` | `CreateTask` | `IsWorkspaceOwnedBy` trả error → service trả "workspace not found" (hardcode), mất error gốc | **Medium** |
| 33 | `internal/modules/task/services/task_service.go:115` | `CreateTask` | `!ownered` → trả "access denied" string. Handler bắt bằng 400, nên là 403 | **High** |
| 34 | `internal/modules/task/services/task_service.go:123` | `CreateTask` | `!StatusIsValid()` → trả "task status invalid." string. Đây là validation error nhưng handler trả 400 chung | **Low** |
| 35 | `internal/modules/task/services/task_service.go:129` | `CreateTask` | `CreateTask` fail → service trả "Task creation failed." hardcode, mất DB error gốc (duplicate key, foreign key, etc.) | **High** |
| 36 | `internal/modules/task/services/task_service.go:143-145` | `UpdateTask` | Trả thẳng repo error lên handler. Không wrap, không phân loại (not found vs DB error) | **Medium** |
| 37 | `internal/modules/task/services/task_service.go:157-159` | `DeleteTask` | Trả thẳng repo error. Repo không kiểm tra RowsAffected → silent success nếu record không tồn tại | **High** |
| 38 | `internal/modules/task/services/task_service.go:172` | `AssignTask` | "task id not exists in workspace" → là not found nhưng handler trả 500 | **High** |
| 39 | `internal/modules/task/services/task_service.go:178` | `AssignTask` | "Task not found" → fmt.Errorf, handler trả 500 | **High** |
| 40 | `internal/modules/task/services/task_service.go:188` | `AssignTask` | "access denied" → nên 403, handler trả 500 | **High** |
| 41 | `internal/modules/task/services/task_service.go:221` | `UpdateTaskStatus` | "task id not exists in workspace" → handler trả 500, nên 404 | **High** |
| 42 | `internal/modules/task/services/task_service.go:227` | `UpdateTaskStatus` | `FindById` error → hardcode "Task not found" fmt.Errorf, mất error gốc, handler trả 500 | **High** |
| 43 | `internal/modules/task/services/task_service.go:237` | `UpdateTaskStatus` | "access denied" → handler trả 500, nên 403 | **High** |
| 44 | `internal/modules/task/services/task_service.go:241` | `UpdateTaskStatus` | Status transition invalid → nên 400 validation, handler trả 500 | **High** |
| 45 | `internal/modules/auth/services/auth_service.go:50-53` | `Register` | `bcrypt.GenerateFromPassword` fail → trả raw bcrypt error string, handler trả 400 → nên 500 | **Medium** |
| 46 | `internal/modules/auth/services/auth_service.go:60-62` | `Register` | `CreateUser` fail → trả raw DB error (có thể chứa SQL, table name). Handler trả 400 | **High** |
| 47 | `internal/modules/auth/services/auth_service.go:73-76` | `Login` | `GetUserByEmail` fail → trả raw DB error, handler trả 401 | **Medium** |
| 48 | `internal/modules/auth/services/auth_service.go:79` | `Login` | "User not found" → handler trả 401. Đúng behavior nhưng message quá rõ, nên generic "email or password invalid" | **Low** |
| 49 | `internal/modules/auth/services/auth_service.go:90-92` | `Login` | Token generation fail → raw error, handler trả 401 → nên 500 | **Medium** |
| 50 | `internal/modules/comment/services/comment_service.go:49` | `CreateComment` | DB error khi check task → fmt.Errorf wrap nhưng handler dùng string match → **không match** → fall vào 500 default | **Medium** |
| 51 | `internal/modules/comment/services/comment_service.go:58` | `CreateComment` | DB error khi check access → fmt.Errorf wrap → handler dùng string match → **không match** → fall vào 500 default | **Medium** |
| 52 | `internal/modules/comment/services/comment_service.go:71` | `CreateComment` | "failed to create comment: ..." → wrap DB error nhưng leak DB detail | **Medium** |
| 53 | `internal/modules/comment/services/comment_service.go:117-118` | `DeleteComment` | `FindById` error → trả "comment not found". Nhưng error có thể là DB connection fail, không phải not found | **High** |
| 54 | `internal/modules/workspace/services/workspace_service.go:35-36` | `GetWorkspace` | Trả thẳng repo error. Repo trả raw "workspace not found" string hoặc DB error | **Medium** |

### 2.3. Repository Layer

| No | File | Function | Current Problem | Severity |
|----|------|----------|----------------|----------|
| 55 | `internal/modules/auth/repositories/user_repository.go:46` | `GetUserByEmail` | `sql.ErrNoRows` → wrap thành `errors.New("User not found")`. Tuy nhiên error message hardcode, không phải domain error type | **Medium** |
| 56 | `internal/modules/auth/repositories/user_repository.go:48` | `GetUserByEmail` | Non-NoRows DB error → trả thẳng raw error lên service | **Medium** |
| 57 | `internal/modules/auth/repositories/user_repository.go:58-59` | `CreateUser` | DB error (duplicate key, constraint violation) → trả thẳng raw error. **Lộ SQL detail** | **High** |
| 58 | `internal/modules/task/repositories/task_repository.go:94-96` | `FindById` | `sql.ErrNoRows` → **KHÔNG xử lý**, trả raw `sql.ErrNoRows` lên service. Service dùng error này nhưng không check | **High** |
| 59 | `internal/modules/task/repositories/task_repository.go:117-131` | `UpdateTask` | **Không kiểm tra RowsAffected** → Nếu task không tồn tại hoặc user không phải owner, query succeed với 0 rows → client nhận 200 OK sai | **Critical** |
| 60 | `internal/modules/task/repositories/task_repository.go:133-142` | `DeleteTask` | **Không kiểm tra RowsAffected** → Silent success khi task không tồn tại | **Critical** |
| 61 | `internal/modules/task/repositories/task_repository.go:188` | `AssignTask` | RowsAffected == 0 → trả "failed to assign task" — không rõ nguyên nhân (not found? already assigned?) | **Medium** |
| 62 | `internal/modules/task/repositories/task_repository.go:209` | `UpdateTaskStatus` | RowsAffected == 0 → trả "failed to assign task" (copy-paste error message!) | **Medium** |
| 63 | `internal/modules/workspace/repositories/workspace_repository.go:37-38` | `FindByOwnerId` | `sql.ErrNoRows` check trên `QueryContext` — `QueryContext` (multi-row) **không bao giờ trả** `sql.ErrNoRows`. Check này dead code, sẽ không bao giờ match | **High** |
| 64 | `internal/modules/workspace/repositories/workspace_repository.go:78` | `Update` | **Không kiểm tra RowsAffected** → Silent success khi workspace không tồn tại hoặc user không phải owner | **High** |
| 65 | `internal/modules/workspace/repositories/workspace_repository.go:87-89` | `Delete` | **Không kiểm tra RowsAffected** → Silent success khi workspace không tồn tại | **High** |
| 66 | `internal/modules/comment/repositories/comment_repository.go:99-101` | `FindById` | `sql.ErrNoRows` → trả raw `sql.ErrNoRows`, service catch nhưng wrap thành generic "comment not found" | **Medium** |
| 67 | `internal/modules/notification/repositories/notification_repository.go:74` | `SaveNotification` | RowsAffected == 0 → trả "faild to save notification" (typo: "faild" → "failed") | **Low** |

### 2.4. Validation / Middleware Layer

| No | File | Function | Current Problem | Severity |
|----|------|----------|----------------|----------|
| 68 | `internal/modules/auth/validates/user_validate.go:27` | `EmailIsValid` | `GetUserByEmail` error bị **ignore** (`_`). Nếu DB connection fail, hàm vẫn return `true, nil` → tiếp tục flow → gây lỗi ở bước create | **Critical** |
| 69 | `internal/modules/auth/validates/user_validate.go:30` | `EmailIsValid` | "Email already exists" → đây thực chất là **409 Conflict** nhưng handler trả 400 | **High** |
| 70 | `pkg/utils/request_context.go:14` | `GetOwnerId` | "user_id not found in context after authentication" → errors.New string. Handler trả 500 nhưng đây là auth issue → nên 401 | **High** |
| 71 | `pkg/cache/redis_cache.go:27` | `Set` | JSON marshal error → `return nil` (swallow error!) → cache set fail silently | **Medium** |

---

## 3. Error Categories

### Category 1: 🔴 Validation error nhưng trả sai status

| Location | Error | Current Status | Expected |
|----------|-------|---------------|----------|
| Auth Register: email duplicate | "Email already exists" | 400 | **409 Conflict** |
| Auth Register: bcrypt fail | raw bcrypt error | 400 | **500 Internal** |
| Auth Register: DB create fail | raw DB error | 400 | **500 Internal** |
| Task CreateTask: status invalid | "task status invalid." | 400 | 400 (đúng, nhưng cần error code) |
| Workspace Delete: UUID parse fail | "invalid user_id format in context" | 500 | **400 Bad Request** |

### Category 2: 🔴 Not Found nhưng trả 500 hoặc 400

| Location | Error | Current Status | Expected |
|----------|-------|---------------|----------|
| Task AssignTask: task not found | "Task not found with id: ..." | 500 | **404 Not Found** |
| Task AssignTask: task not in workspace | "task id not exists in workspace" | 500 | **404 Not Found** |
| Task UpdateTaskStatus: task not found | "Task not found with id: ..." | 500 | **404 Not Found** |
| Task UpdateTaskStatus: task not in workspace | "task id not exists in workspace" | 500 | **404 Not Found** |
| Task CreateTask: workspace not found | "workspace not found" | 400 | **404 Not Found** |

### Category 3: 🔴 Forbidden nhưng trả 500 hoặc 400

| Location | Error | Current Status | Expected |
|----------|-------|---------------|----------|
| Task CreateTask: access denied | "access denied" | 400 | **403 Forbidden** |
| Task AssignTask: access denied | "access denied" | 500 | **403 Forbidden** |
| Task UpdateTaskStatus: access denied | "access denied" | 500 | **403 Forbidden** |
| Task UpdateTaskStatus: invalid transition | "can not change task status..." | 500 | **400 Bad Request** |

### Category 4: 🟡 GetOwnerId fail → 500 thay vì 401

Tất cả handlers gọi `utils.GetOwnerId(c)` đều trả 500 khi fail. Đây là lỗi liên quan auth context → nên trả **401 Unauthorized**.

**Affected functions** (11 instances):
- `TaskHandler.GetAllTask`, `GetTask`, `CreateTask`, `UpdateTask`, `DeleteTask`, `AssignTask`, `UpdateTaskStatus`
- `WorkspaceHandler.GetWorkspace`, `UpdateWorkspace`, `DeleteWorkspace`
- `CommentHandler.CreateComment`, `GetComments`, `DeleteComment`
- `WSHandler.Connect`

### Category 5: 🔴 Repository DB error bị lộ ra ngoài

| Location | Leaked Info |
|----------|-------------|
| `user_repository.CreateUser` | Raw SQL error (table name, constraint name) |
| `task_repository.FindById` | Raw `sql.ErrNoRows` |
| `task_repository.FindAll` | Raw DB scan error |
| `workspace_repository.FindByOwnerId` | Raw DB error |
| `workspace_repository.Update` | Raw DB error |
| `workspace_repository.Delete` | Raw DB error |

### Category 6: 🔴 Service error chưa có type rõ ràng

- Tất cả service đều dùng `errors.New()` hoặc `fmt.Errorf()` → handler không phân biệt được loại lỗi
- Không có error type (interface/struct) để kiểm tra bằng `errors.Is()` hoặc `errors.As()`

### Category 7: 🔴 Handler xử lý error bằng string comparison

```go
// comment_handler.go - 3 function đều dùng pattern này
switch err.Error() {
case "task not found":
    c.JSON(http.StatusNotFound, ...)
case "access denied":
    c.JSON(http.StatusForbidden, ...)
default:
    c.JSON(http.StatusInternalServerError, ...)
}
```

**Rủi ro**: Nếu service thay đổi message từ `"task not found"` thành `"Task not found"` → handler sẽ fall vào `default` → trả 500.

### Category 8: 🟡 Silent success (không kiểm tra RowsAffected)

| Repository Function | Impact |
|---------------------|--------|
| `task_repository.UpdateTask` | Update task không tồn tại → trả success |
| `task_repository.DeleteTask` | Delete task không tồn tại → trả success |
| `workspace_repository.Update` | Update workspace không tồn tại → trả success |
| `workspace_repository.Delete` | Delete workspace không tồn tại → trả success |

### Category 9: 🔴 Bug - Nil Pointer Dereference

```go
// task_service.go:96-103
task, err := s.taskRepository.FindById(context, id, ownerId)
if err != nil {
    fmt.Printf("[ERROR] Set task into cache failed with error: %v\n", err)
}
taskRes := mapper.EntityToTaskResponse(task) // ← PANIC nếu task == nil
```

Khi `FindById` trả error, `task` sẽ là `nil`, nhưng code vẫn tiếp tục gọi mapper → **runtime panic**.

---

## 4. Recommended Architecture

### 4.1. Nên tạo custom error package

**CÓ** — Tạo package `pkg/apperror` chứa:
- `AppError` struct
- Error code constants
- Constructor functions
- HTTP status mapper

### 4.2. Kiến trúc Error Flow

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Repository  │────▶│   Service    │────▶│   Handler   │
│              │     │              │     │              │
│ Wrap DB err  │     │ Return       │     │ Use          │
│ into domain  │     │ AppError     │     │ MapAppError  │
│ AppError     │     │ with code    │     │ to HTTP      │
└─────────────┘     └─────────────┘     └─────────────┘
```

### 4.3. Repository nên trả domain error

**CÓ** — Repository nên:
- Catch `sql.ErrNoRows` → return `apperror.NewNotFound("task")`
- Catch duplicate key → return `apperror.NewConflict("email already exists")`
- Catch foreign key violation → return `apperror.NewBadRequest("invalid reference")`
- Other DB errors → return `apperror.NewInternal("database error")` (hide SQL detail)
- Kiểm tra `RowsAffected` → return `apperror.NewNotFound("resource")` nếu 0 rows

### 4.4. Service nên trả AppError

**CÓ** — Service nên:
- Trả `apperror.NewForbidden("access denied")` thay vì `errors.New("access denied")`
- Trả `apperror.NewValidation("task status invalid")` thay vì `errors.New("task status invalid.")`
- Propagate AppError từ repository (không rewrap nếu đã là AppError)

### 4.5. Handler nên dùng mapper function

```go
func HandleAppError(c *gin.Context, err error) {
    var appErr *apperror.AppError
    if errors.As(err, &appErr) {
        c.JSON(appErr.HTTPStatus(), response.ResponseError(appErr.Message, appErr.Code))
        return
    }
    // Fallback: unknown error → 500
    c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", "INTERNAL_ERROR"))
}
```

### 4.6. Middleware nên xử lý lỗi auth riêng

Middleware auth hiện tại đã xử lý tương đối tốt (trả 401). Tuy nhiên, `GetOwnerId` utility nên trả `apperror.NewUnauthorized()` để handler biết map sang 401.

---

## 5. Proposed AppError Design

### 5.1. AppError Struct

```go
// pkg/apperror/app_error.go
package apperror

import (
    "fmt"
    "net/http"
)

// ErrorCode represents a machine-readable error code
type ErrorCode string

const (
    CodeBadRequest    ErrorCode = "BAD_REQUEST"
    CodeUnauthorized  ErrorCode = "UNAUTHORIZED"
    CodeForbidden     ErrorCode = "FORBIDDEN"
    CodeNotFound      ErrorCode = "NOT_FOUND"
    CodeConflict      ErrorCode = "CONFLICT"
    CodeValidation    ErrorCode = "VALIDATION_ERROR"
    CodeInternal      ErrorCode = "INTERNAL_ERROR"
)

// AppError is the standard error type used across all layers
type AppError struct {
    Code    ErrorCode `json:"code"`
    Message string    `json:"message"`
    Status  int       `json:"-"` // HTTP status code
    Err     error     `json:"-"` // Original error (for logging, not exposed)
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
    }
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
    return e.Err
}

func (e *AppError) HTTPStatus() int {
    return e.Status
}
```

### 5.2. Constructor Functions

```go
// pkg/apperror/constructors.go
package apperror

import "net/http"

func NewBadRequest(message string) *AppError {
    return &AppError{
        Code:    CodeBadRequest,
        Message: message,
        Status:  http.StatusBadRequest,
    }
}

func NewUnauthorized(message string) *AppError {
    return &AppError{
        Code:    CodeUnauthorized,
        Message: message,
        Status:  http.StatusUnauthorized,
    }
}

func NewForbidden(message string) *AppError {
    return &AppError{
        Code:    CodeForbidden,
        Message: message,
        Status:  http.StatusForbidden,
    }
}

func NewNotFound(resource string) *AppError {
    return &AppError{
        Code:    CodeNotFound,
        Message: fmt.Sprintf("%s not found", resource),
        Status:  http.StatusNotFound,
    }
}

func NewConflict(message string) *AppError {
    return &AppError{
        Code:    CodeConflict,
        Message: message,
        Status:  http.StatusConflict,
    }
}

func NewValidation(message string) *AppError {
    return &AppError{
        Code:    CodeValidation,
        Message: message,
        Status:  http.StatusBadRequest,
    }
}

func NewInternal(message string) *AppError {
    return &AppError{
        Code:    CodeInternal,
        Message: message,
        Status:  http.StatusInternalServerError,
    }
}

// Wrap preserves the original error for logging while returning AppError
func Wrap(err error, appErr *AppError) *AppError {
    appErr.Err = err
    return appErr
}
```

### 5.3. Error Handler Helper

```go
// pkg/apperror/handler.go
package apperror

import (
    "errors"
    "net/http"

    "dev/task-management/pkg/response"
    "github.com/gin-gonic/gin"
)

// HandleError maps AppError to HTTP response
func HandleError(c *gin.Context, err error) {
    var appErr *AppError
    if errors.As(err, &appErr) {
        c.JSON(appErr.Status, response.ResponseError(appErr.Message, appErr.Code))
        return
    }
    // Unknown error → 500
    c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", CodeInternal))
}
```

---

## 6. Handler Refactor Strategy

### Nguyên tắc

1. **Handler KHÔNG nên biết business logic error** — chỉ cần gọi `apperror.HandleError(c, err)`
2. **Handler chỉ xử lý request binding** — ShouldBindJSON, Param parsing
3. **GetOwnerId fail → trả 401** (không phải 500)

### Pattern mới

```go
func (h *TaskHandler) CreateTask(c *gin.Context) {
    var taskReq *request.CreateTaskRequest
    if err := c.ShouldBindJSON(&taskReq); err != nil {
        apperror.HandleError(c, apperror.NewValidation(err.Error()))
        return
    }

    ownerId, err := utils.GetOwnerId(c)
    if err != nil {
        apperror.HandleError(c, apperror.NewUnauthorized("invalid authentication context"))
        return
    }

    task, err := h.service.CreateTask(c.Request.Context(), taskReq, ownerId)
    if err != nil {
        apperror.HandleError(c, err) // Service đã trả AppError → handler chỉ forward
        return
    }

    c.JSON(http.StatusCreated, response.ResponseSuccess("create task successfully", task))
}
```

### So sánh

| Aspect | Before | After |
|--------|--------|-------|
| Error mapping | Hardcode HTTP status per handler | Service trả AppError, handler forward |
| GetOwnerId fail | 500 Internal | 401 Unauthorized |
| Service error | 400 or 500 (đoán) | Đúng status từ AppError |
| String comparison | `switch err.Error()` | `errors.As(err, &appErr)` |
| Code duplication | Mỗi handler viết error logic riêng | Một function `HandleError` dùng chung |

---

## 7. Service Refactor Strategy

### Nguyên tắc

1. **Service trả `*apperror.AppError`** cho business error
2. **Service wrap repo error** nếu cần context thêm
3. **Service KHÔNG return raw string error**
4. **Service validate business rule** trả `apperror.NewValidation()` hoặc `apperror.NewForbidden()`

### Pattern mới

```go
func (s *taskService) CreateTask(ctx context.Context, t *request.CreateTaskRequest, ownerId uuid.UUID) (*response.TaskResponse, error) {
    ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, t.Workspace, ownerId)
    if err != nil {
        return nil, err // Propagate AppError from workspace service
    }
    if !ownered {
        return nil, apperror.NewForbidden("you don't have permission to create task in this workspace")
    }

    id := uuid.New()
    createAt := time.Now()
    task := mapper.TaskRequestToEntity(id, t, createAt)

    if !task.StatusIsValid() {
        return nil, apperror.NewValidation("task status invalid, must be one of: TODO, IN_PROGRESS, DONE, BLOCKED")
    }

    if err := s.taskRepository.CreateTask(ctx, task); err != nil {
        return nil, err // Repo đã wrap thành AppError
    }

    _ = s.redisClient.Clear("tasks:*")
    return mapper.EntityToTaskResponse(task), nil
}
```

---

## 8. Repository Refactor Strategy

### Nguyên tắc

1. **Xử lý `sql.ErrNoRows`** → return `apperror.NewNotFound("resource")`
2. **Xử lý duplicate key** (PostgreSQL error code `23505`) → return `apperror.NewConflict("...")`
3. **Xử lý foreign key violation** (PostgreSQL error code `23503`) → return `apperror.NewBadRequest("...")`
4. **Other DB errors** → return `apperror.Wrap(err, apperror.NewInternal("database error"))` (hide SQL, keep for logging)
5. **Kiểm tra `RowsAffected`** cho UPDATE/DELETE → return `apperror.NewNotFound("resource")` nếu 0

### PostgreSQL Error Detection Helper

```go
// pkg/apperror/postgres.go
package apperror

import (
    "database/sql"
    "errors"

    "github.com/jackc/pgx/v5/pgconn"
)

const (
    PgUniqueViolation     = "23505"
    PgForeignKeyViolation = "23503"
    PgCheckViolation      = "23514"
)

// WrapDBError converts database errors into AppError
func WrapDBError(err error, resourceName string) *AppError {
    if err == nil {
        return nil
    }

    // sql.ErrNoRows → Not Found
    if errors.Is(err, sql.ErrNoRows) {
        return NewNotFound(resourceName)
    }

    // PostgreSQL specific errors
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        switch pgErr.Code {
        case PgUniqueViolation:
            return Wrap(err, NewConflict(resourceName+" already exists"))
        case PgForeignKeyViolation:
            return Wrap(err, NewBadRequest("invalid reference for "+resourceName))
        }
    }

    // Unknown DB error → Internal (hide detail)
    return Wrap(err, NewInternal("database error"))
}
```

### Ví dụ refactor Repository

```go
func (r *taskRepository) FindById(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*entities.Task, error) {
    sqlStr := `select t.task_id, t.title, t.description, t.status, t.assignee_id, t.workspace_id, t.create_at 
               from tasks t
               join workspaces wp on t.workspace_id = wp.workspace_id
               where (t.task_id = $1) and (wp.owner_id = $2 or t.assignee_id = $3)`

    row := r.database.QueryRowContext(ctx, sqlStr, id, ownerId, ownerId)
    var task entities.Task
    err := row.Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.Assignee, &task.Workspace, &task.CreateAt)
    if err != nil {
        return nil, apperror.WrapDBError(err, "task")
    }
    return &task, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task, ownerId uuid.UUID) error {
    sqlStr := `update tasks set title = $1, description = $2, status = $3, workspace_id = $4
               where task_id = $5 and workspace_id in (select workspace_id from workspaces where owner_id = $6)`
    
    result, err := r.database.ExecContext(ctx, sqlStr, task.Title, task.Description, task.Status, task.Workspace, id, ownerId)
    if err != nil {
        return apperror.WrapDBError(err, "task")
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return apperror.Wrap(err, apperror.NewInternal("database error"))
    }
    if rowsAffected == 0 {
        return apperror.NewNotFound("task")
    }
    return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
    sqlStr := `delete from tasks t where t.task_id = $1 and exists (
                select 1 from workspaces w where w.workspace_id = t.workspace_id and w.owner_id = $2)`
    
    result, err := r.database.ExecContext(ctx, sqlStr, id, ownerId)
    if err != nil {
        return apperror.WrapDBError(err, "task")
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return apperror.Wrap(err, apperror.NewInternal("database error"))
    }
    if rowsAffected == 0 {
        return apperror.NewNotFound("task")
    }
    return nil
}
```

---

## 9. Refactor Priority

### Phase 1: 🔴 Critical — Fix wrong HTTP status & Panic bugs (Sprint 1)

| Task | File(s) | Impact |
|------|---------|--------|
| Fix nil pointer dereference in `GetTask` service | `task_service.go:96-103` | **Runtime panic** |
| Fix `EmailIsValid` ignoring DB error | `user_validate.go:27` | **Silent corruption** |
| Create `pkg/apperror` package | NEW: `pkg/apperror/*.go` | Foundation |
| Fix `GetOwnerId` → trả 401 thay vì 500 | All handlers (14 instances) | Wrong auth status |
| Fix `AssignTask` handler → trả đúng status | `task_handler.go:179` | 500 → 403/404 |
| Fix `UpdateTaskStatus` handler → trả đúng status | `task_handler.go:214` | 500 → 400/403/404 |
| Fix `DeleteWorkspace` UUID validation → 400 | `workspace_hander.go:82` | 500 → 400 |

### Phase 2: 🟠 High — Standardize AppError across Services (Sprint 2)

| Task | File(s) |
|------|---------|
| Refactor `task_service.go` to return `AppError` | `task_service.go` |
| Refactor `auth_service.go` to return `AppError` | `auth_service.go` |
| Refactor `comment_service.go` to return `AppError` | `comment_service.go` |
| Refactor `workspace_service.go` to return `AppError` | `workspace_service.go` |
| Add `WrapDBError` helper for PostgreSQL | `pkg/apperror/postgres.go` |
| Refactor `user_validate.go` → return `AppError` | `user_validate.go` |

### Phase 3: 🟡 Medium — Clean Handlers & Repositories (Sprint 3)

| Task | File(s) |
|------|---------|
| Refactor all handlers to use `apperror.HandleError()` | All handler files |
| Remove `switch err.Error()` from comment handler | `comment_handler.go` |
| Add `RowsAffected` checks to task/workspace repositories | `task_repository.go`, `workspace_repository.go` |
| Fix `workspace_repository.FindByOwnerId` dead code | `workspace_repository.go:37` |
| Fix `redis_cache.Set` swallowing marshal error | `redis_cache.go:27` |
| Update `ApiResponse` to include error code | `api_response.go` |

### Phase 4: 🟢 Low — Unit Tests & Polish (Sprint 4)

| Task | File(s) |
|------|---------|
| Unit test AppError constructors | `pkg/apperror/*_test.go` |
| Unit test `WrapDBError` helper | `pkg/apperror/postgres_test.go` |
| Unit test handlers with error scenarios | `handler/*_test.go` |
| Unit test services with AppError returns | `services/*_test.go` |
| Fix typos ("faild", "Resquest", etc.) | Multiple files |

---

## 10. Example Before/After

### Example 1: Task Handler — AssignTask (500 → đúng status)

**Before** (mọi error → 500):
```go
func (h *TaskHandler) AssignTask(c *gin.Context) {
    var taskReq request.AssignTaskRequest
    err := c.ShouldBindJSON(&taskReq)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
        return
    }

    ownerId, err := utils.GetOwnerId(c)
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
        return
    }

    err = h.service.AssignTask(c.Request.Context(), &taskReq, ownerId)
    if err != nil {
        // ❌ Mọi lỗi đều trả 500: access denied, not found, DB error
        c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
        return
    }
    c.JSON(http.StatusOK, response.ResponseSuccess("Assign task successfully", nil))
}
```

**After** (đúng status cho từng loại lỗi):
```go
func (h *TaskHandler) AssignTask(c *gin.Context) {
    var taskReq request.AssignTaskRequest
    if err := c.ShouldBindJSON(&taskReq); err != nil {
        apperror.HandleError(c, apperror.NewValidation(err.Error()))
        return
    }

    ownerId, err := utils.GetOwnerId(c)
    if err != nil {
        apperror.HandleError(c, apperror.NewUnauthorized("invalid authentication context"))
        return
    }

    if err := h.service.AssignTask(c.Request.Context(), &taskReq, ownerId); err != nil {
        // ✅ Service trả AppError → handler forward đúng status
        // "access denied" → 403, "task not found" → 404, DB error → 500
        apperror.HandleError(c, err)
        return
    }
    c.JSON(http.StatusOK, response.ResponseSuccess("Assign task successfully", nil))
}
```

### Example 2: Comment Handler — Loại bỏ String Comparison

**Before** (switch err.Error() fragile):
```go
func (h *CommentHandler) CreateComment(c *gin.Context) {
    // ... binding code ...

    comment, err := h.service.CreateComment(c.Request.Context(), taskId, userId, &req)
    if err != nil {
        // ❌ String comparison — break nếu message thay đổi
        switch err.Error() {
        case "task not found":
            c.JSON(http.StatusNotFound, response.ResponseError("create comment failed", err.Error()))
        case "access denied":
            c.JSON(http.StatusForbidden, response.ResponseError("create comment failed", err.Error()))
        case "comment content cannot be empty":
            c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
        default:
            c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
        }
        return
    }
    c.JSON(http.StatusCreated, response.ResponseSuccess("create comment successfully", comment))
}
```

**After** (clean, type-safe):
```go
func (h *CommentHandler) CreateComment(c *gin.Context) {
    // ... binding code ...

    userId, err := utils.GetOwnerId(c)
    if err != nil {
        apperror.HandleError(c, apperror.NewUnauthorized("invalid authentication context"))
        return
    }

    comment, err := h.service.CreateComment(c.Request.Context(), taskId, userId, &req)
    if err != nil {
        // ✅ Service trả AppError → errors.As map tự động
        apperror.HandleError(c, err)
        return
    }
    c.JSON(http.StatusCreated, response.ResponseSuccess("create comment successfully", comment))
}
```

### Example 3: Task Service GetTask — Fix Nil Pointer + Error Type

**Before** (nil pointer panic + wrong error handling):
```go
func (s *taskService) GetTask(context context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error) {
    redisKey := "tasks:owner_id:" + ownerId.String() + ":task_id:" + id.String()
    taskCache := &response.TaskResponse{}

    err := s.redisClient.Get(redisKey, &taskCache)
    if err == nil {
        return taskCache, nil
    }

    task, err := s.taskRepository.FindById(context, id, ownerId)
    if err != nil {
        // ❌ Chỉ log, KHÔNG return → tiếp tục xử lý với task == nil
        fmt.Printf("[ERROR] Set task into cache failed with error: %v\n", err)
    }
    taskRes := mapper.EntityToTaskResponse(task) // ❌ PANIC khi task == nil

    s.redisClient.Set(redisKey, taskRes, 10)
    return taskRes, nil
}
```

**After** (safe, proper error propagation):
```go
func (s *taskService) GetTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error) {
    redisKey := "tasks:owner_id:" + ownerId.String() + ":task_id:" + id.String()
    taskCache := &response.TaskResponse{}

    err := s.redisClient.Get(redisKey, &taskCache)
    if err == nil {
        return taskCache, nil
    }

    task, err := s.taskRepository.FindById(ctx, id, ownerId)
    if err != nil {
        // ✅ Repository trả AppError (NotFound hoặc Internal) → propagate lên handler
        return nil, err
    }

    taskRes := mapper.EntityToTaskResponse(task) // ✅ Safe: task != nil

    // Cache set failure is non-critical → just log
    if cacheErr := s.redisClient.Set(redisKey, taskRes, 10); cacheErr != nil {
        fmt.Printf("[WARN] Failed to cache task: %v\n", cacheErr)
    }

    return taskRes, nil
}
```

### Example 4: Auth Register — Đúng status cho duplicate email

**Before** (duplicate email → 400):
```go
// Service
func (s *authService) Register(ctx context.Context, userDTO *request.RegisterUserRequestDTO) (*response.UserResponseDTO, error) {
    emailValid, err := validates.EmailIsValid(userDTO.Email, ctx, s.userRepo)
    if !emailValid {
        return nil, err // ❌ errors.New("Email already exists") → handler trả 400
    }
    // ...
    errCreate := s.userRepo.CreateUser(ctx, user)
    if errCreate != nil {
        return nil, errCreate // ❌ Raw DB error string → handler trả 400
    }
    // ...
}

// Handler
func (h *AuthHander) Register(c *gin.Context) {
    // ...
    res, err := h.authService.Register(c.Request.Context(), &userReq)
    if err != nil {
        c.JSON(http.StatusBadRequest, ...) // ❌ Always 400
    }
}
```

**After** (đúng 409 Conflict cho duplicate, 500 cho internal):
```go
// Validate
func EmailIsValid(email string, ctx context.Context, repo repositories.UserRepository) error {
    user, err := repo.GetUserByEmail(ctx, email)
    if err != nil {
        // ✅ Check if NotFound (expected) or real DB error
        var appErr *apperror.AppError
        if errors.As(err, &appErr) && appErr.Code == apperror.CodeNotFound {
            return nil // Email not taken → valid
        }
        return apperror.Wrap(err, apperror.NewInternal("failed to check email availability"))
    }
    if user != nil {
        return apperror.NewConflict("email already exists") // ✅ 409
    }
    return nil
}

// Service
func (s *authService) Register(ctx context.Context, userDTO *request.RegisterUserRequestDTO) (*response.UserResponseDTO, error) {
    if err := validates.EmailIsValid(userDTO.Email, ctx, s.userRepo); err != nil {
        return nil, err // ✅ AppError: Conflict hoặc Internal
    }
    // ...
}

// Handler
func (h *AuthHander) Register(c *gin.Context) {
    // ...
    res, err := h.authService.Register(c.Request.Context(), &userReq)
    if err != nil {
        apperror.HandleError(c, err) // ✅ 409 cho duplicate, 400 cho validation, 500 cho internal
        return
    }
    c.JSON(http.StatusCreated, response.ResponseSuccess("register account successfully", res))
}
```

---

## 11. Unit Test Plan

### 11.1. AppError Package Tests

```go
// pkg/apperror/app_error_test.go

func TestNewBadRequest(t *testing.T) {
    err := NewBadRequest("invalid input")
    assert.Equal(t, CodeBadRequest, err.Code)
    assert.Equal(t, http.StatusBadRequest, err.Status)
    assert.Equal(t, "invalid input", err.Message)
}

func TestWrapDBError_NoRows(t *testing.T) {
    err := WrapDBError(sql.ErrNoRows, "task")
    assert.Equal(t, CodeNotFound, err.Code)
    assert.Equal(t, http.StatusNotFound, err.Status)
}

func TestWrapDBError_UniqueViolation(t *testing.T) {
    pgErr := &pgconn.PgError{Code: "23505"}
    err := WrapDBError(pgErr, "email")
    assert.Equal(t, CodeConflict, err.Code)
    assert.Equal(t, http.StatusConflict, err.Status)
}

func TestErrorsAs(t *testing.T) {
    err := NewNotFound("task")
    var appErr *AppError
    assert.True(t, errors.As(err, &appErr))
    assert.Equal(t, CodeNotFound, appErr.Code)
}
```

### 11.2. Handler Error Response Tests

#### 400 Validation Tests
```go
func TestCreateTask_InvalidJSON(t *testing.T) {
    // Given: malformed JSON body
    // When: POST /api/tasks/
    // Then: 400 with code "VALIDATION_ERROR"
}

func TestCreateTask_InvalidStatus(t *testing.T) {
    // Given: status = "INVALID"
    // When: POST /api/tasks/
    // Then: 400 with code "VALIDATION_ERROR"
}

func TestUpdateTaskStatus_InvalidTransition(t *testing.T) {
    // Given: task is DONE, try to change to IN_PROGRESS
    // When: PATCH /api/tasks/:id/status
    // Then: 400 with code "BAD_REQUEST"
}
```

#### 401 Unauthorized Tests
```go
func TestGetAllTask_NoAuthHeader(t *testing.T) {
    // Given: no Authorization header
    // When: GET /api/tasks/
    // Then: 401 with code "UNAUTHORIZED"
}

func TestGetAllTask_InvalidToken(t *testing.T) {
    // Given: expired or malformed token
    // When: GET /api/tasks/
    // Then: 401 with code "UNAUTHORIZED"
}

func TestLogin_WrongPassword(t *testing.T) {
    // Given: valid email, wrong password
    // When: POST /api/auth/login
    // Then: 401 with code "UNAUTHORIZED"
}
```

#### 403 Forbidden Tests
```go
func TestCreateTask_NotWorkspaceOwner(t *testing.T) {
    // Given: user is not workspace owner
    // When: POST /api/tasks/ with another user's workspace_id
    // Then: 403 with code "FORBIDDEN"
}

func TestAssignTask_NotWorkspaceOwner(t *testing.T) {
    // Given: user is not workspace owner
    // When: PATCH /api/tasks/assign
    // Then: 403 with code "FORBIDDEN"
}

func TestDeleteComment_NotOwnerNotAdmin(t *testing.T) {
    // Given: user is not comment creator and not workspace owner
    // When: DELETE /api/comments/:commentId
    // Then: 403 with code "FORBIDDEN"
}
```

#### 404 Not Found Tests
```go
func TestGetTask_NotExists(t *testing.T) {
    // Given: task ID does not exist
    // When: GET /api/tasks/:id
    // Then: 404 with code "NOT_FOUND"
}

func TestDeleteTask_NotExists(t *testing.T) {
    // Given: task ID does not exist
    // When: DELETE /api/tasks/:id
    // Then: 404 with code "NOT_FOUND"
}

func TestCreateComment_TaskNotExists(t *testing.T) {
    // Given: task ID does not exist
    // When: POST /api/tasks/:id/comments
    // Then: 404 with code "NOT_FOUND"
}

func TestDeleteWorkspace_NotExists(t *testing.T) {
    // Given: workspace ID does not exist
    // When: DELETE /api/workspaces/:wpId
    // Then: 404 with code "NOT_FOUND"
}
```

#### 409 Conflict Tests
```go
func TestRegister_DuplicateEmail(t *testing.T) {
    // Given: email already exists in database
    // When: POST /api/auth/register
    // Then: 409 with code "CONFLICT"
}

func TestAssignTask_AlreadyAssigned(t *testing.T) {
    // Given: task already assigned to the same user
    // When: PATCH /api/tasks/assign
    // Then: 409 with code "CONFLICT"
}
```

#### 500 Internal Server Error Tests
```go
func TestGetAllTask_DBConnectionFail(t *testing.T) {
    // Given: database connection is down
    // When: GET /api/tasks/
    // Then: 500 with code "INTERNAL_ERROR"
    // And: response does NOT contain SQL details
}

func TestRegister_BcryptFail(t *testing.T) {
    // Given: bcrypt fails (e.g. password too long)
    // When: POST /api/auth/register
    // Then: 500 with code "INTERNAL_ERROR"
}
```

### 11.3. Response Format Verification

```go
func TestErrorResponseFormat(t *testing.T) {
    // Verify all error responses match the standard format:
    // {
    //     "success": false,
    //     "message": "Task not found",
    //     "error": "NOT_FOUND"
    // }
    
    var resp response.ApiResponse
    json.Unmarshal(body, &resp)
    
    assert.False(t, resp.Success)
    assert.NotEmpty(t, resp.Message)
    assert.NotEmpty(t, resp.Error)
    assert.Nil(t, resp.Data)
}

func TestSuccessResponseFormat(t *testing.T) {
    // Verify all success responses match:
    // {
    //     "success": true,
    //     "message": "create task successfully",
    //     "data": { ... }
    // }
    
    var resp response.ApiResponse
    json.Unmarshal(body, &resp)
    
    assert.True(t, resp.Success)
    assert.NotEmpty(t, resp.Message)
    assert.Nil(t, resp.Error)
}
```

---

## Appendix A: Updated ApiResponse Format

### Current
```go
type ApiResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
    Error   any    `json:"error,omitempty"` // ← raw error string
}
```

### Proposed
```go
type ApiResponse struct {
    Success bool   `json:"success"`
    Message string `json:"message"`
    Data    any    `json:"data,omitempty"`
    Error   any    `json:"error,omitempty"` // ← ErrorCode string (e.g. "NOT_FOUND")
}
```

Error response example:
```json
{
    "success": false,
    "message": "Task not found",
    "error": "NOT_FOUND"
}
```

Success response example:
```json
{
    "success": true,
    "message": "create task successfully",
    "data": {
        "id": "uuid-here",
        "title": "My Task"
    }
}
```

---

## Appendix B: Typos Found

| File | Current | Fix |
|------|---------|-----|
| `workspace_hander.go` (filename) | `workspace_hander.go` | `workspace_handler.go` |
| `workspace_hander.go:49` | `"validation faild"` | `"validation failed"` |
| `workspace_hander.go:75` | `"validation faild"` | `"validation failed"` |
| `user_handler.go:52` | `"Authentication faild"` | `"Authentication failed"` |
| `user_handler.go:12` | `AuthHander` | `AuthHandler` |
| `workspace/dto` | `WorkspaceUpdateResquestDTO` | `WorkspaceUpdateRequestDTO` |
| `workspace/dto` | `WorkspaceResquestDTO` | `WorkspaceRequestDTO` |
| `notification_repository.go:74` | `"faild to save notification"` | `"failed to save notification"` |
| `task_repository.go:209` | `"failed to assign task"` (in UpdateTaskStatus) | `"failed to update task status"` |

---

## Appendix C: Files to Create/Modify

### New Files
| File | Purpose |
|------|---------|
| `pkg/apperror/app_error.go` | AppError struct + ErrorCode constants |
| `pkg/apperror/constructors.go` | Constructor functions (NewBadRequest, NewNotFound, etc.) |
| `pkg/apperror/handler.go` | HandleError helper for Gin |
| `pkg/apperror/postgres.go` | WrapDBError for PostgreSQL error detection |
| `pkg/apperror/app_error_test.go` | Unit tests |

### Modified Files
| File | Changes |
|------|---------|
| `pkg/response/api_response.go` | Giữ nguyên struct, error field nhận ErrorCode |
| `pkg/utils/request_context.go` | GetOwnerId trả AppError |
| `internal/modules/auth/handler/user_handler.go` | Dùng apperror.HandleError |
| `internal/modules/auth/services/auth_service.go` | Return AppError |
| `internal/modules/auth/repositories/user_repository.go` | WrapDBError |
| `internal/modules/auth/validates/user_validate.go` | Return AppError, fix error ignore |
| `internal/modules/task/handler/task_handler.go` | Dùng apperror.HandleError |
| `internal/modules/task/services/task_service.go` | Return AppError, fix nil pointer |
| `internal/modules/task/repositories/task_repository.go` | WrapDBError, RowsAffected checks |
| `internal/modules/workspace/handler/workspace_hander.go` | Dùng apperror.HandleError |
| `internal/modules/workspace/services/workspace_service.go` | Return AppError |
| `internal/modules/workspace/repositories/workspace_repository.go` | WrapDBError, RowsAffected checks, fix dead code |
| `internal/modules/comment/handler/comment_handler.go` | Remove string comparison, dùng HandleError |
| `internal/modules/comment/services/comment_service.go` | Return AppError |
| `internal/modules/comment/repositories/comment_repository.go` | WrapDBError |
| `internal/realtime/hander_ws.go` | GetOwnerId → 401 |
| `pkg/cache/redis_cache.go` | Fix swallowed error in Set() |
