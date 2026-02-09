package config

import (
	"os"
	"testing"
)

func TestGetToken_FromEnvVar(t *testing.T) {
	// Set up
	const testToken = "ghp_test_token_123"
	os.Setenv("HIJ_GITHUB_TOKEN", testToken)
	defer os.Unsetenv("HIJ_GITHUB_TOKEN")

	// Test
	token, source := GetToken()

	// Verify
	if token != testToken {
		t.Errorf("token = %q, want %q", token, testToken)
	}
	if source != "env" {
		t.Errorf("source = %q, want %q", source, "env")
	}
}

func TestGetToken_NotSet(t *testing.T) {
	// Ensure env var is not set
	os.Unsetenv("HIJ_GITHUB_TOKEN")

	// Test
	token, source := GetToken()

	// When running in test environment without keychain access,
	// we expect empty values if no token is configured
	// Note: We cannot test keychain in unit tests as it requires system access
	if source == "env" {
		t.Error("source should not be 'env' when HIJ_GITHUB_TOKEN is not set")
	}
	// Token should be empty or from keychain (if configured on the system)
	_ = token // avoid unused variable warning
}
