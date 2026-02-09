package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-token")

	if client.token != "test-token" {
		t.Errorf("token = %q, want %q", client.token, "test-token")
	}
	if client.httpClient == nil {
		t.Error("httpClient should not be nil")
	}
	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("timeout = %v, want %v", client.httpClient.Timeout, 30*time.Second)
	}
	if client.baseURL != defaultBaseURL {
		t.Errorf("baseURL = %q, want %q", client.baseURL, defaultBaseURL)
	}
}

func TestClient_ListPackages(t *testing.T) {
	tests := []struct {
		name           string
		responseCode   int
		responseBody   string
		wantErr        bool
		wantCount      int
		wantErrContain string
	}{
		{
			name:         "success",
			responseCode: http.StatusOK,
			responseBody: `[{"id":1,"name":"pkg1"},{"id":2,"name":"pkg2"}]`,
			wantErr:      false,
			wantCount:    2,
		},
		{
			name:           "unauthorized",
			responseCode:   http.StatusUnauthorized,
			responseBody:   `{"message":"Bad credentials"}`,
			wantErr:        true,
			wantErrContain: "invalid or expired token",
		},
		{
			name:           "forbidden",
			responseCode:   http.StatusForbidden,
			responseBody:   `{"message":"Must have admin rights"}`,
			wantErr:        true,
			wantErrContain: "access denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request headers
				if auth := r.Header.Get("Authorization"); !strings.HasPrefix(auth, "Bearer ") {
					t.Errorf("Authorization header = %q, want Bearer prefix", auth)
				}
				if accept := r.Header.Get("Accept"); accept != "application/vnd.github+json" {
					t.Errorf("Accept header = %q, want %q", accept, "application/vnd.github+json")
				}

				w.WriteHeader(tt.responseCode)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			client := NewClient("test-token")
			client.baseURL = server.URL

			packages, err := client.ListPackages("container")

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if tt.wantErrContain != "" && !strings.Contains(err.Error(), tt.wantErrContain) {
					t.Errorf("error = %q, want containing %q", err.Error(), tt.wantErrContain)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if len(packages) != tt.wantCount {
					t.Errorf("package count = %d, want %d", len(packages), tt.wantCount)
				}
			}
		})
	}
}

func TestClient_ListPackageVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		versions := []PackageVersion{
			{ID: 1, Name: "sha256:abc123"},
			{ID: 2, Name: "sha256:def456"},
		}
		json.NewEncoder(w).Encode(versions)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.baseURL = server.URL

	versions, err := client.ListPackageVersions("container", "test-pkg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(versions) != 2 {
		t.Errorf("version count = %d, want 2", len(versions))
	}
}

func TestClient_DeletePackageVersion(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		wantErr      bool
	}{
		{
			name:         "success",
			responseCode: http.StatusNoContent,
			wantErr:      false,
		},
		{
			name:         "not found",
			responseCode: http.StatusNotFound,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "DELETE" {
					t.Errorf("method = %q, want DELETE", r.Method)
				}
				w.WriteHeader(tt.responseCode)
			}))
			defer server.Close()

			client := NewClient("test-token")
			client.baseURL = server.URL

			err := client.DeletePackageVersion("container", "test-pkg", 123)

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestClient_ValidateToken(t *testing.T) {
	tests := []struct {
		name         string
		responseCode int
		wantErr      bool
	}{
		{
			name:         "valid token",
			responseCode: http.StatusOK,
			wantErr:      false,
		},
		{
			name:         "invalid token",
			responseCode: http.StatusUnauthorized,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.responseCode)
				w.Write([]byte(`{}`))
			}))
			defer server.Close()

			client := NewClient("test-token")
			client.baseURL = server.URL

			err := client.ValidateToken()

			if tt.wantErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestParseAPIError(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		body       string
		wantContain string
	}{
		{
			name:       "401 unauthorized",
			statusCode: 401,
			body:       `{}`,
			wantContain: "invalid or expired token",
		},
		{
			name:       "403 with message",
			statusCode: 403,
			body:       `{"message":"Resource protected"}`,
			wantContain: "access denied: Resource protected",
		},
		{
			name:       "403 without message",
			statusCode: 403,
			body:       `{}`,
			wantContain: "read:packages",
		},
		{
			name:       "404 not found",
			statusCode: 404,
			body:       `{}`,
			wantContain: "resource not found",
		},
		{
			name:       "422 with message",
			statusCode: 422,
			body:       `{"message":"Validation failed"}`,
			wantContain: "request failed: Validation failed",
		},
		{
			name:       "429 rate limited",
			statusCode: 429,
			body:       `{}`,
			wantContain: "rate limited",
		},
		{
			name:       "500 with message",
			statusCode: 500,
			body:       `{"message":"Internal error"}`,
			wantContain: "GitHub API error: Internal error",
		},
		{
			name:       "500 without message",
			statusCode: 500,
			body:       `{}`,
			wantContain: "status 500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parseAPIError(tt.statusCode, []byte(tt.body))
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if !strings.Contains(err.Error(), tt.wantContain) {
				t.Errorf("error = %q, want containing %q", err.Error(), tt.wantContain)
			}
		})
	}
}
