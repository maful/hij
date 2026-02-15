package ui

import (
	"testing"
	"time"

	"github.com/maful/hij/github"
)

func TestSortVersions(t *testing.T) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	twoDaysAgo := now.Add(-48 * time.Hour)

	versions := []github.PackageVersion{
		{ID: 1, CreatedAt: yesterday},
		{ID: 2, CreatedAt: now},
		{ID: 3, CreatedAt: twoDaysAgo},
	}

	// Use pointer to Model to modify state
	m := &Model{
		sortOrder: "newest",
	}

	// Test Newest First (default)
	m.sortVersions(versions)
	if versions[0].ID != 2 {
		t.Errorf("Newest Sort: Expected first version to be ID 2 (Today), got %d", versions[0].ID)
	}
	if versions[1].ID != 1 {
		t.Errorf("Newest Sort: Expected second version to be ID 1 (Yesterday), got %d", versions[1].ID)
	}
	if versions[2].ID != 3 {
		t.Errorf("Newest Sort: Expected third version to be ID 3 (TwoDaysAgo), got %d", versions[2].ID)
	}

	// Test Oldest First
	m.sortOrder = "oldest"
	m.sortVersions(versions)
	if versions[0].ID != 3 {
		t.Errorf("Oldest Sort: Expected first version to be ID 3 (TwoDaysAgo), got %d", versions[0].ID)
	}
	if versions[1].ID != 1 {
		t.Errorf("Oldest Sort: Expected second version to be ID 1 (Yesterday), got %d", versions[1].ID)
	}
	if versions[2].ID != 2 {
		t.Errorf("Oldest Sort: Expected third version to be ID 2 (Today), got %d", versions[2].ID)
	}
}
