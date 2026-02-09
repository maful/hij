package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const baseURL = "https://api.github.com"

// Client is a GitHub API client for package operations
type Client struct {
	token      string
	httpClient *http.Client
}

// NewClient creates a new GitHub client with the given PAT
func NewClient(token string) *Client {
	return &Client{
		token: token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// doRequest performs an authenticated request to GitHub API
func (c *Client) doRequest(method, path string) ([]byte, error) {
	req, err := http.NewRequest(method, baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// ListPackages lists all packages for the authenticated user
func (c *Client) ListPackages(packageType string) ([]Package, error) {
	path := fmt.Sprintf("/user/packages?package_type=%s", packageType)
	body, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var packages []Package
	if err := json.Unmarshal(body, &packages); err != nil {
		return nil, err
	}

	return packages, nil
}

// ListPackageVersions lists all versions for a package
func (c *Client) ListPackageVersions(packageType, packageName string) ([]PackageVersion, error) {
	path := fmt.Sprintf("/user/packages/%s/%s/versions?per_page=100", packageType, packageName)
	body, err := c.doRequest("GET", path)
	if err != nil {
		return nil, err
	}

	var versions []PackageVersion
	if err := json.Unmarshal(body, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

// DeletePackageVersion deletes a specific package version
func (c *Client) DeletePackageVersion(packageType, packageName string, versionID int) error {
	path := fmt.Sprintf("/user/packages/%s/%s/versions/%d", packageType, packageName, versionID)
	_, err := c.doRequest("DELETE", path)
	return err
}

// ValidateToken checks if the token is valid by attempting to list packages
func (c *Client) ValidateToken() error {
	_, err := c.doRequest("GET", "/user")
	return err
}
