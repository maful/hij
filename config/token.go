package config

import (
	"os"

	"github.com/zalando/go-keyring"
)

const (
	envVarName     = "HIJ_GITHUB_TOKEN"
	keyringService = "hij"
	keyringUser    = "github-token"
)

// GetToken retrieves the token from environment variable or keychain.
// Returns the token and its source ("env", "keychain", or "" if not found).
func GetToken() (string, string) {
	// Check environment variable first
	if token := os.Getenv(envVarName); token != "" {
		return token, "env"
	}

	// Check keychain
	if token, err := keyring.Get(keyringService, keyringUser); err == nil && token != "" {
		return token, "keychain"
	}

	return "", ""
}

// SaveToken saves the token to the system keychain.
func SaveToken(token string) error {
	return keyring.Set(keyringService, keyringUser, token)
}

// DeleteToken removes the token from the system keychain.
func DeleteToken() error {
	return keyring.Delete(keyringService, keyringUser)
}
