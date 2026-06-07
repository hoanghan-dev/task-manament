package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// =============================================================================
// HTTP Request Helpers
// =============================================================================

// performRequest sends an HTTP request through the gin router and returns the recorder.
// - method: HTTP method (GET, POST, PUT, DELETE, PATCH)
// - path: API path (e.g. "/api/auth/register")
// - body: request body (will be marshaled to JSON), pass nil for no body
// - token: Bearer token string, pass "" for no auth header
func performRequest(router *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal request body: %v", err))
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		panic(fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// =============================================================================
// JSON Response Helpers
// =============================================================================

// apiResponse represents the standard API response structure.
type apiResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data,omitempty"`
	Error   interface{}     `json:"error,omitempty"`
}

// parseAPIResponse parses the recorder body into the standard apiResponse.
func parseAPIResponse(t *testing.T, w *httptest.ResponseRecorder) apiResponse {
	t.Helper()

	var resp apiResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Failed to parse API response JSON: %v\nBody: %s", err, w.Body.String())
	}
	return resp
}

// parseDataAs unmarshals the apiResponse.Data field into the target struct.
func parseDataAs(t *testing.T, resp apiResponse, target interface{}) {
	t.Helper()

	if resp.Data == nil {
		t.Fatal("Expected response data but got nil")
	}
	if err := json.Unmarshal(resp.Data, target); err != nil {
		t.Fatalf("Failed to parse response data: %v\nRaw: %s", err, string(resp.Data))
	}
}

// =============================================================================
// Flow Helpers (compose multiple API calls)
// =============================================================================

// registerUserResponse holds parsed fields from a successful register response.
type registerUserResponse struct {
	UserID    string `json:"user_id"`
	Email     string `json:"email"`
	FullName  string `json:"full_name"`
	Workspace struct {
		WorkspaceID string `json:"workspace_id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"workspace"`
}

// loginResponse holds parsed fields from a successful login response.
type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// registerUser sends a POST /api/auth/register request and returns the parsed response.
func registerUser(t *testing.T, r *gin.Engine, email, password, fullName string) (*httptest.ResponseRecorder, registerUserResponse) {
	t.Helper()

	body := map[string]string{
		"email":     email,
		"password":  password,
		"full_name": fullName,
	}

	w := performRequest(r, http.MethodPost, "/api/auth/register", body, "")

	var userData registerUserResponse
	if w.Code == http.StatusCreated {
		resp := parseAPIResponse(t, w)
		parseDataAs(t, resp, &userData)
	}

	return w, userData
}

// loginUser sends a POST /api/auth/login request and returns the access token.
func loginUser(t *testing.T, r *gin.Engine, email, password string) (*httptest.ResponseRecorder, loginResponse) {
	t.Helper()

	body := map[string]string{
		"email":    email,
		"password": password,
	}

	w := performRequest(r, http.MethodPost, "/api/auth/login", body, "")

	var loginData loginResponse
	if w.Code == http.StatusOK {
		resp := parseAPIResponse(t, w)
		parseDataAs(t, resp, &loginData)
	}

	return w, loginData
}

// registerAndLogin is a convenience helper that registers a user and logs them in,
// returning the access token, workspace ID, and user ID.
func registerAndLogin(t *testing.T, r *gin.Engine, email, password, fullName string) (token string, workspaceID string, userID string) {
	t.Helper()

	// Register
	regW, regData := registerUser(t, r, email, password, fullName)
	if regW.Code != http.StatusCreated {
		t.Fatalf("Register failed: status=%d, body=%s", regW.Code, regW.Body.String())
	}

	// Login
	loginW, loginData := loginUser(t, r, email, password)
	if loginW.Code != http.StatusOK {
		t.Fatalf("Login failed: status=%d, body=%s", loginW.Code, loginW.Body.String())
	}

	return loginData.AccessToken, regData.Workspace.WorkspaceID, regData.UserID
}

func performRequestWithAuthHeader(
	router *gin.Engine,
	method string,
	path string,
	body interface{},
	authHeader string,
) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer

	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal request body: %v", err))
		}
		reqBody = bytes.NewBuffer(jsonBytes)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		panic(fmt.Sprintf("failed to create request: %v", err))
	}

	req.Header.Set("Content-Type", "application/json")

	if authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}
