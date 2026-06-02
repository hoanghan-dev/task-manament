# Implementation Plan - API Error Test Suite

This plan details the design and implementation of a comprehensive API error test suite for the Go backend. It verifies that all main HTTP endpoints, middleware, and handlers correctly process, map, and return the appropriate HTTP status codes and custom `ErrorCode` values after the error handling refactor.

## Goal Description

Create a robust test suite that executes HTTP handler requests through the Gin router (near-production setup) using mock services. Verify that all standard error cases return the correct HTTP status codes and JSON response bodies:

```json
{
  "success": false,
  "message": "...",
  "error": "ERROR_CODE"
}
```

This plan covers:
1. Designing all API test cases in `API_ERROR_TEST_CASES.md`.
2. Setting up a modular mocking architecture using the existing service interfaces.
3. Implementing automated handler tests for Auth, Middleware, Tasks, Workspaces, and Comments.
4. Running the test suite and checking for race conditions.
5. Summarizing the results in `API_ERROR_TEST_REPORT.md`.

---

## User Review Required

> [!IMPORTANT]
> - **Mock Services vs Real DB**: We will utilize mock service implementations to test the handler and router layers. This avoids flaky DB connections, avoids polluting the database, and keeps the test suite fast.
> - **Testing Protected Endpoints**: We will use the production JWT utilities to sign valid tokens using a mock/test secret key, verifying the authentication middleware end-to-end.

---

## Open Questions

None. The requirements and architecture are fully understood.

---

## Proposed Changes

We will group our changes/additions by modules and packages.

### 1. Documentation
#### [NEW] [API_ERROR_TEST_CASES.md](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/API_ERROR_TEST_CASES.md)
Contains a complete table of all error test cases, their request body, expected HTTP status, expected `ErrorCode`, and notes.

---

### 2. Testing Utilities & Setup
#### [NEW] [testutil.go](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/internal/testutil/testutil.go)
Provides:
- Valid token generation for test users (setting `SECRET_KEY` in environment).
- Helper function to parse Gin response bodies into a standard structure.
- Router instantiation and mock registration utility.

---

### 3. Middleware Tests
#### [NEW] [auth_test.go](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/internal/middleware/auth_test.go)
Tests for:
- Missing `Authorization` header -> `401 Unauthorized` / `"missing authorization header"`
- Incorrect header format (no `Bearer ` prefix) -> `401 Unauthorized` / `"authorization header format invalid"`
- Expired or invalid token -> `401 Unauthorized`

---

### 4. Auth Module Tests
#### [NEW] [user_handler_test.go](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/internal/modules/auth/handler/user_handler_test.go)
Tests for:
- `/api/auth/register`:
  - Missing field -> `400 Bad Request` with `VALIDATION_ERROR`
  - Duplicate email conflict -> `409 Conflict` with `CONFLICT`
- `/api/auth/login`:
  - Missing field in body -> `400 Bad Request` with `VALIDATION_ERROR`
  - Wrong password/email (mock returns unauthorized) -> `401 Unauthorized` with `UNAUTHORIZED`

---

### 5. Workspace Module Tests
#### [NEW] [workspace_handler_test.go](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/internal/modules/workspace/handler/workspace_handler_test.go)
Tests for:
- `/api/workspaces/` (GET):
  - Request with no token -> `401 Unauthorized` (via AuthMiddleware)
  - Not Found -> `404 Not Found` with `NOT_FOUND`
- `/api/workspaces/` (PUT):
  - Missing validation field -> `400 Bad Request` with `VALIDATION_ERROR`
  - Not Found -> `404 Not Found` with `NOT_FOUND`
- `/api/workspaces/:wpId` (DELETE):
  - Invalid UUID parameter -> `400 Bad Request` with `VALIDATION_ERROR`
  - Not Found -> `404 Not Found` with `NOT_FOUND`

---

### 6. Task Module Tests
#### [NEW] [task_handler_test.go](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/internal/modules/task/handler/task_handler_test.go)
Tests for:
- `/api/tasks/` (GET / POST):
  - Invalid UUID/Missing required field -> `400 Bad Request` with `VALIDATION_ERROR`
  - Not workspace owner -> `403 Forbidden` with `FORBIDDEN`
- `/api/tasks/:id` (GET / PUT / DELETE):
  - Invalid UUID format -> `400 Bad Request` with `VALIDATION_ERROR`
  - Task not found -> `404 Not Found` with `NOT_FOUND`
- `/api/tasks/assign` (PATCH):
  - Assignee not found / task not found -> `404 Not Found` with `NOT_FOUND`
  - Assign task in workspace not owned by requester -> `403 Forbidden` with `FORBIDDEN`
- `/api/tasks/:id/status` (PATCH):
  - Invalid status transition -> `400 Bad Request` with `BAD_REQUEST`

---

### 7. Comment Module Tests
#### [NEW] [comment_handler_test.go](file:///d:/Documents/FPT_Education/SU26/OJT/Let's%20GO/task-management-antigavity/task-manament/internal/modules/comment/handler/comment_handler_test.go)
Tests for:
- `/api/tasks/:id/comments` (POST):
  - Invalid UUID parameter -> `400 Bad Request` with `VALIDATION_ERROR`
  - Task not found -> `404 Not Found` with `NOT_FOUND`
  - Empty comment content -> `400 Bad Request` with `VALIDATION_ERROR`
- `/api/comments/:commentId` (DELETE):
  - Comment not found -> `404 Not Found` with `NOT_FOUND`
  - Comment not owned by deleting user -> `403 Forbidden` with `FORBIDDEN`

---

## Verification Plan

### Automated Tests
Run the test suite:
```bash
go test ./...
```
Verify no race conditions:
```bash
go test -race ./...
```

### Manual Verification
Ensure all written test logs compile and execute correctly, outputting green statuses.
Write the test execution outcomes into `API_ERROR_TEST_REPORT.md`.
