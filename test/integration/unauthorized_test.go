package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestProtectedAPIs_WithoutToken_Returns401(t *testing.T) {
	cleanupDatabase(testDB)

	tests := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{
			name:   "Get workspaces without token",
			method: http.MethodGet,
			path:   "/api/workspaces/",
			body:   nil,
		},
		{
			name:   "Create task without token",
			method: http.MethodPost,
			path:   "/api/tasks/",
			body: map[string]any{
				"title":        "Unauthorized Task",
				"description":  "Should not be created",
				"status":       "TODO",
				"workspace_id": uuid.NewString(),
			},
		},
		{
			name:   "Create comment without token",
			method: http.MethodPost,
			path:   fmt.Sprintf("/api/tasks/%s/comments", uuid.NewString()),
			body: map[string]any{
				"content": "Unauthorized comment",
			},
		},
		{
			name:   "Get comments without token",
			method: http.MethodGet,
			path:   fmt.Sprintf("/api/tasks/%s/comments", uuid.NewString()),
			body:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(testRouter, tt.method, tt.path, tt.body, "")

			assert.Equal(t, http.StatusUnauthorized, w.Code,
				"Expected 401 for %s %s, body: %s",
				tt.method,
				tt.path,
				w.Body.String(),
			)

			resp := parseAPIResponse(t, w)
			assert.False(t, resp.Success)
		})
	}
}

func TestProtectedAPIs_InvalidToken_Returns401(t *testing.T) {
	cleanupDatabase(testDB)

	invalidToken := "invalid.jwt.token"

	tests := []struct {
		name   string
		method string
		path   string
		body   any
	}{
		{
			name:   "Get workspaces with invalid token",
			method: http.MethodGet,
			path:   "/api/workspaces/",
			body:   nil,
		},
		{
			name:   "Create task with invalid token",
			method: http.MethodPost,
			path:   "/api/tasks/",
			body: map[string]any{
				"title":        "Invalid Token Task",
				"description":  "Should not be created",
				"status":       "TODO",
				"workspace_id": uuid.NewString(),
			},
		},
		{
			name:   "Create comment with invalid token",
			method: http.MethodPost,
			path:   fmt.Sprintf("/api/tasks/%s/comments", uuid.NewString()),
			body: map[string]any{
				"content": "Invalid token comment",
			},
		},
		{
			name:   "Get comments with invalid token",
			method: http.MethodGet,
			path:   fmt.Sprintf("/api/tasks/%s/comments", uuid.NewString()),
			body:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequest(testRouter, tt.method, tt.path, tt.body, invalidToken)

			assert.Equal(t, http.StatusUnauthorized, w.Code,
				"Expected 401 for %s %s, body: %s",
				tt.method,
				tt.path,
				w.Body.String(),
			)

			resp := parseAPIResponse(t, w)
			assert.False(t, resp.Success)
		})
	}
}

func TestProtectedAPIs_MalformedAuthorizationHeader_Returns401(t *testing.T) {
	cleanupDatabase(testDB)

	tests := []struct {
		name       string
		authHeader string
	}{
		{
			name:       "No Bearer prefix",
			authHeader: "invalid-token",
		},
		{
			name:       "Basic auth instead of Bearer",
			authHeader: "Basic abcxyz",
		},
		{
			name:       "Bearer without token",
			authHeader: "Bearer ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := performRequestWithAuthHeader(
				testRouter,
				http.MethodGet,
				"/api/workspaces/",
				nil,
				tt.authHeader,
			)

			assert.Equal(t, http.StatusUnauthorized, w.Code,
				"Expected 401, body: %s", w.Body.String())

			resp := parseAPIResponse(t, w)
			assert.False(t, resp.Success)
		})
	}
}
