package integration

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type permissionTestUser struct {
	UserID      uuid.UUID
	WorkspaceID uuid.UUID
	Token       string
}

type permissionTestSetup struct {
	UserA  permissionTestUser
	UserB  permissionTestUser
	TaskID uuid.UUID
}

func TestPermission_UserBCannotUpdateUserATask(t *testing.T) {
	cleanupDatabase(testDB)

	setup := setupTwoUsersAndUserATask(t)

	body := map[string]any{
		"title":        "User B tries to update User A task",
		"description":  "This update must be forbidden",
		"status":       "TODO",
		"assignee_id":  setup.UserA.UserID.String(),
		"workspace_id": setup.UserA.WorkspaceID.String(),
	}

	res := performRequest(
		testRouter,
		http.MethodPut,
		"/api/tasks/"+setup.TaskID.String(),
		body,
		setup.UserB.Token,
	)

	assert.Equal(t, http.StatusNotFound, res.Code, res.Body.String())
}

func TestPermission_UserBCannotDeleteUserATask(t *testing.T) {
	cleanupDatabase(testDB)

	setup := setupTwoUsersAndUserATask(t)

	res := performRequest(
		testRouter,
		http.MethodDelete,
		"/api/tasks/"+setup.TaskID.String(),
		nil,
		setup.UserB.Token,
	)

	assert.Equal(t, http.StatusNotFound, res.Code, res.Body.String())
}

func TestPermission_UserBCannotCommentOnUserATask(t *testing.T) {
	cleanupDatabase(testDB)

	setup := setupTwoUsersAndUserATask(t)

	body := map[string]any{
		"content": "User B tries to comment on User A task",
	}

	res := performRequest(
		testRouter,
		http.MethodPost,
		"/api/tasks/"+setup.TaskID.String()+"/comments",
		body,
		setup.UserB.Token,
	)

	assert.Equal(t, http.StatusForbidden, res.Code, res.Body.String())
}

func TestPermission_MissingTokenReturnsUnauthorized(t *testing.T) {
	cleanupDatabase(testDB)

	setup := setupTwoUsersAndUserATask(t)

	body := map[string]any{
		"title":        "Missing token update",
		"description":  "This request must be unauthorized",
		"status":       "TODO",
		"assignee_id":  setup.UserA.UserID.String(),
		"workspace_id": setup.UserA.WorkspaceID.String(),
	}

	res := performRequest(
		testRouter,
		http.MethodPut,
		"/api/tasks/"+setup.TaskID.String(),
		body,
		"",
	)

	assert.Equal(t, http.StatusUnauthorized, res.Code, res.Body.String())
}

func TestPermission_InvalidTokenReturnsUnauthorized(t *testing.T) {
	cleanupDatabase(testDB)

	setup := setupTwoUsersAndUserATask(t)

	body := map[string]any{
		"title":        "Invalid token update",
		"description":  "This request must be unauthorized",
		"status":       "TODO",
		"assignee_id":  setup.UserA.UserID.String(),
		"workspace_id": setup.UserA.WorkspaceID.String(),
	}

	res := performRequest(
		testRouter,
		http.MethodPut,
		"/api/tasks/"+setup.TaskID.String(),
		body,
		"invalid-token",
	)

	assert.Equal(t, http.StatusUnauthorized, res.Code, res.Body.String())
}

func setupTwoUsersAndUserATask(t *testing.T) permissionTestSetup {
	t.Helper()

	userA := registerPermissionTestUser(t, "permission-user-a@test.com", "User A Permission")
	taskID := createTaskForUser(t, userA)

	userB := registerPermissionTestUser(t, "permission-user-b@test.com", "User B Permission")

	return permissionTestSetup{
		UserA:  userA,
		UserB:  userB,
		TaskID: taskID,
	}
}

func registerPermissionTestUser(t *testing.T, email string, name string) permissionTestUser {
	t.Helper()

	/*
		This function assumes the existing helper registerAndLogin has this shape:

			authData := registerAndLogin(t, email, password, name)

		and authData contains:

			authData.UserID
			authData.WorkspaceID
			authData.AccessToken

		If your existing helper returns a raw response instead, replace only this function
		with parseAPIResponse() + parseDataAs() according to your current test helpers.
	*/

	AccessToken, WorkspaceID, UserID := registerAndLogin(t, testRouter, email, "Password123!", name)

	return permissionTestUser{
		UserID:      uuid.MustParse(UserID),
		WorkspaceID: uuid.MustParse(WorkspaceID),
		Token:       AccessToken,
	}
}

func createTaskForUser(t *testing.T, user permissionTestUser) uuid.UUID {
	t.Helper()

	body := map[string]any{
		"title":        "User A Task",
		"description":  "Task owned by User A",
		"status":       "TODO",
		"assignee_id":  user.UserID.String(),
		"workspace_id": user.WorkspaceID.String(),
	}

	res := performRequest(
		testRouter,
		http.MethodPost,
		"/api/tasks/",
		body,
		user.Token,
	)

	require.Equal(t, http.StatusCreated, res.Code, res.Body.String())

	apiRes := parseAPIResponse(t, res)

	var data struct {
		TaskID string `json:"task_id"`
	}
	parseDataAs(t, apiRes, &data)

	require.NotEmpty(t, data.TaskID)

	return uuid.MustParse(data.TaskID)
}
