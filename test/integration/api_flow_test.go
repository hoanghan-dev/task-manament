package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// Main Happy Path: Register → Login → Create Task → Create Comment
// =============================================================================

func TestAPIFlow_RegisterLoginCreateProjectCreateTaskComment(t *testing.T) {
	// Clean database before test to ensure a fresh state
	cleanupDatabase(testDB)

	testEmail := fmt.Sprintf("integrationtest+%s@example.com", uuid.NewString()[:8])
	testPassword := "Secure123"
	testFullName := "Integration Test User"

	// -------------------------------------------------------------------------
	// Step 1: Register a new user
	// -------------------------------------------------------------------------
	var userData registerUserResponse
	var w *httptest.ResponseRecorder
	t.Run("Step1_Register", func(t *testing.T) {
		w, userData = registerUser(t, testRouter, testEmail, testPassword, testFullName)

		require.Equal(t, http.StatusCreated, w.Code, "Register should return 201")

		resp := parseAPIResponse(t, w)
		assert.True(t, resp.Success, "Response success should be true")
		assert.Equal(t, "register account successfully", resp.Message)

		// Verify user fields
		assert.NotEmpty(t, userData.UserID, "user_id should not be empty")
		assert.Equal(t, testEmail, userData.Email)
		assert.Equal(t, testFullName, userData.FullName)

		// Verify workspace was auto-created
		assert.NotEmpty(t, userData.Workspace.WorkspaceID, "workspace_id should not be empty")
		assert.Contains(t, userData.Workspace.Name, testFullName, "workspace name should contain user's full name")
	})

	// -------------------------------------------------------------------------
	// Step 2: Login with the registered user
	// -------------------------------------------------------------------------
	var accessToken string
	var workspaceID string

	t.Run("Step2_Login", func(t *testing.T) {
		w, loginData := loginUser(t, testRouter, testEmail, testPassword)

		require.Equal(t, http.StatusOK, w.Code, "Login should return 200")

		resp := parseAPIResponse(t, w)
		assert.True(t, resp.Success, "Response success should be true")
		assert.Equal(t, "authentication successfully", resp.Message)

		assert.NotEmpty(t, loginData.AccessToken, "access_token should not be empty")
		assert.NotEmpty(t, loginData.RefreshToken, "refresh_token should not be empty")

		accessToken = loginData.AccessToken
	})

	// Get workspace_id from register step (need to re-register or extract from prior step)
	// Since we already have test data, let's get workspace via API
	t.Run("Step2b_GetWorkspace", func(t *testing.T) {
		require.NotEmpty(t, accessToken, "accessToken must be set from Step2")

		w := performRequest(testRouter, http.MethodGet, "/api/workspaces/", nil, accessToken)
		require.Equal(t, http.StatusOK, w.Code, "Get workspace should return 200")

		resp := parseAPIResponse(t, w)
		assert.True(t, resp.Success)

		var wsData struct {
			WorkspaceID string `json:"workspace_id"`
			Name        string `json:"name"`
		}
		parseDataAs(t, resp, &wsData)
		assert.NotEmpty(t, wsData.WorkspaceID, "workspace_id should not be empty")
		workspaceID = wsData.WorkspaceID
	})

	// -------------------------------------------------------------------------
	// Step 3: Create a task in the workspace
	// -------------------------------------------------------------------------
	var taskID string

	t.Run("Step3_CreateTask", func(t *testing.T) {
		require.NotEmpty(t, accessToken, "accessToken must be set from Step2")
		require.NotEmpty(t, workspaceID, "workspaceID must be set from Step2b")
		require.NotEmpty(t, userData.UserID, "userId must be set from Step1")

		body := map[string]interface{}{
			"title":        "Integration Test Task",
			"description":  "This task was created by the integration test suite",
			"status":       "TODO",
			"assignee_id":  userData.UserID,
			"workspace_id": workspaceID,
		}

		w := performRequest(testRouter, http.MethodPost, "/api/tasks/", body, accessToken)
		require.Equal(t, http.StatusCreated, w.Code, "Create task should return 201, body: %s", w.Body.String())

		resp := parseAPIResponse(t, w)
		assert.True(t, resp.Success, "Response success should be true")
		assert.Equal(t, "create task successfully", resp.Message)

		var taskData struct {
			TaskID      string `json:"task_id"`
			Title       string `json:"title"`
			Description string `json:"description"`
			Status      string `json:"status"`
			AssigneeID  string `json:"assignee_id"`
			WorkspaceID string `json:"workspace_id"`
		}
		parseDataAs(t, resp, &taskData)

		assert.NotEmpty(t, taskData.TaskID, "task_id should not be empty")
		assert.Equal(t, "Integration Test Task", taskData.Title)
		assert.Equal(t, "TODO", taskData.Status)
		assert.Equal(t, workspaceID, taskData.WorkspaceID)
		assert.Equal(t, userData.UserID, taskData.AssigneeID)
		taskID = taskData.TaskID
	})

	// -------------------------------------------------------------------------
	// Step 4: Create a comment on the task
	// -------------------------------------------------------------------------
	t.Run("Step4_CreateComment", func(t *testing.T) {
		require.NotEmpty(t, accessToken, "accessToken must be set")
		require.NotEmpty(t, taskID, "taskID must be set from Step3")

		body := map[string]string{
			"content": "This is an integration test comment",
		}

		path := fmt.Sprintf("/api/tasks/%s/comments", taskID)
		w := performRequest(testRouter, http.MethodPost, path, body, accessToken)
		require.Equal(t, http.StatusCreated, w.Code, "Create comment should return 201, body: %s", w.Body.String())

		resp := parseAPIResponse(t, w)
		assert.True(t, resp.Success, "Response success should be true")
		assert.Equal(t, "create comment successfully", resp.Message)

		var commentData struct {
			CommentID string `json:"comment_id"`
			TaskID    string `json:"task_id"`
			UserID    string `json:"user_id"`
			Content   string `json:"content"`
		}
		parseDataAs(t, resp, &commentData)

		assert.NotEmpty(t, commentData.CommentID, "comment_id should not be empty")
		assert.Equal(t, taskID, commentData.TaskID, "comment should belong to the created task")
		assert.NotEmpty(t, commentData.UserID, "user_id should not be empty")
		assert.Equal(t, "This is an integration test comment", commentData.Content)
	})

	// -------------------------------------------------------------------------
	// Step 5: Verify - Get comments for the task
	// -------------------------------------------------------------------------
	t.Run("Step5_GetComments", func(t *testing.T) {
		require.NotEmpty(t, accessToken, "accessToken must be set")
		require.NotEmpty(t, taskID, "taskID must be set")

		path := fmt.Sprintf("/api/tasks/%s/comments", taskID)
		w := performRequest(testRouter, http.MethodGet, path, nil, accessToken)
		require.Equal(t, http.StatusOK, w.Code, "Get comments should return 200")

		resp := parseAPIResponse(t, w)
		assert.True(t, resp.Success)
		assert.Equal(t, "get comments successfully", resp.Message)

		var comments []struct {
			CommentID string `json:"comment_id"`
			Content   string `json:"content"`
		}
		parseDataAs(t, resp, &comments)

		assert.Len(t, comments, 1, "Should have exactly 1 comment")
		assert.Equal(t, "This is an integration test comment", comments[0].Content)
	})
}

// =============================================================================
// Negative Case 1: Create task without token → 401
// =============================================================================

func TestCreateTask_WithoutToken_Returns401(t *testing.T) {
	cleanupDatabase(testDB)

	body := map[string]interface{}{
		"title":        "Unauthorized Task",
		"description":  "Should not be created",
		"status":       "TODO",
		"workspace_id": uuid.NewString(),
	}

	w := performRequest(testRouter, http.MethodPost, "/api/tasks/", body, "")
	assert.Equal(t, http.StatusUnauthorized, w.Code, "Create task without token should return 401")

	resp := parseAPIResponse(t, w)
	assert.False(t, resp.Success, "Response success should be false")
}

// =============================================================================
// Negative Case 2: Create comment on non-existing task → 404
// =============================================================================

func TestCreateComment_NonExistingTask_Returns404(t *testing.T) {
	cleanupDatabase(testDB)

	// Register and login to get a valid token
	token, _, _ := registerAndLogin(t, testRouter,
		fmt.Sprintf("negativetest+%s@example.com", uuid.NewString()[:8]),
		"Secure123",
		"Negative Test User",
	)

	nonExistingTaskID := uuid.NewString()
	body := map[string]string{
		"content": "Comment on a ghost task",
	}

	path := fmt.Sprintf("/api/tasks/%s/comments", nonExistingTaskID)
	w := performRequest(testRouter, http.MethodPost, path, body, token)

	assert.Equal(t, http.StatusNotFound, w.Code,
		"Create comment on non-existing task should return 404, body: %s", w.Body.String())

	resp := parseAPIResponse(t, w)
	assert.False(t, resp.Success, "Response success should be false")
}

// =============================================================================
// Negative Case 3: Create task without token (empty Authorization) → 401
// =============================================================================

func TestCreateComment_WithoutToken_Returns401(t *testing.T) {
	cleanupDatabase(testDB)

	body := map[string]string{
		"content": "Unauthorized comment",
	}

	path := fmt.Sprintf("/api/tasks/%s/comments", uuid.NewString())
	w := performRequest(testRouter, http.MethodPost, path, body, "")

	assert.Equal(t, http.StatusUnauthorized, w.Code,
		"Create comment without token should return 401")

	resp := parseAPIResponse(t, w)
	assert.False(t, resp.Success, "Response success should be false")
}

// =============================================================================
// Negative Case 4: Register duplicate email → 409
// =============================================================================

func TestRegister_DuplicateEmail_Returns409(t *testing.T) {
	cleanupDatabase(testDB)

	email := fmt.Sprintf("duplicate+%s@example.com", uuid.NewString()[:8])
	password := "Secure123"

	// Register first time - should succeed
	w1, _ := registerUser(t, testRouter, email, password, "First User")
	require.Equal(t, http.StatusCreated, w1.Code)

	// Register again with same email - should fail with 409
	w2, _ := registerUser(t, testRouter, email, password, "Second User")
	assert.Equal(t, http.StatusConflict, w2.Code,
		"Register with duplicate email should return 409, body: %s", w2.Body.String())

	resp := parseAPIResponse(t, w2)
	assert.False(t, resp.Success)
}
