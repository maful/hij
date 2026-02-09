package ui

import (
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/maful/hij/github"
)

// createTestVersions creates test versions with specific ages for filter testing
func createTestVersions() []github.PackageVersion {
	now := time.Now()
	return []github.PackageVersion{
		{ID: 1, Name: "v1", CreatedAt: now.Add(-5 * 24 * time.Hour)},   // 5 days old
		{ID: 2, Name: "v2", CreatedAt: now.Add(-15 * 24 * time.Hour)},  // 15 days old
		{ID: 3, Name: "v3", CreatedAt: now.Add(-30 * 24 * time.Hour)},  // 30 days old
		{ID: 4, Name: "v4", CreatedAt: now.Add(-60 * 24 * time.Hour)},  // 60 days old
	}
}

func TestModel_ApplyFilter_OlderN(t *testing.T) {
	tests := []struct {
		name           string
		filter         string
		expectedCount  int
		expectedIDs    []int
	}{
		{
			name:          "older 10 - filters versions older than 10 days",
			filter:        "older 10",
			expectedCount: 3,
			expectedIDs:   []int{2, 3, 4},
		},
		{
			name:          ":older 10 - with colon prefix",
			filter:        ":older 10",
			expectedCount: 3,
			expectedIDs:   []int{2, 3, 4},
		},
		{
			name:          "older 25 - filters versions older than 25 days",
			filter:        "older 25",
			expectedCount: 2,
			expectedIDs:   []int{3, 4},
		},
		{
			name:          "older 100 - no versions older than 100 days",
			filter:        "older 100",
			expectedCount: 0,
			expectedIDs:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{
				versions:         createTestVersions(),
				selectedVersions: make(map[int]struct{}),
				filterInput:      textinput.New(),
				filterValue:      tt.filter,
			}

			m.applyFilter()

			if len(m.filteredVersions) != tt.expectedCount {
				t.Errorf("filtered count = %d, want %d", len(m.filteredVersions), tt.expectedCount)
			}

			// Check that correct versions were selected
			if len(m.selectedVersions) != tt.expectedCount {
				t.Errorf("selected count = %d, want %d", len(m.selectedVersions), tt.expectedCount)
			}

			for _, id := range tt.expectedIDs {
				if _, ok := m.selectedVersions[id]; !ok {
					t.Errorf("expected version ID %d to be selected", id)
				}
			}
		})
	}
}

func TestModel_ApplyFilter_BeforeDate(t *testing.T) {
	// Create versions with specific dates
	versions := []github.PackageVersion{
		{ID: 1, Name: "v1", CreatedAt: time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "v2", CreatedAt: time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)},
		{ID: 3, Name: "v3", CreatedAt: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)},
		{ID: 4, Name: "v4", CreatedAt: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC)},
	}

	tests := []struct {
		name          string
		filter        string
		expectedCount int
		expectedIDs   []int
	}{
		{
			name:          "before 2024-04-01 - filters March and earlier",
			filter:        "before 2024-04-01",
			expectedCount: 3,
			expectedIDs:   []int{2, 3, 4},
		},
		{
			name:          ":before 2024-02-01 - with colon prefix",
			filter:        ":before 2024-02-01",
			expectedCount: 2,
			expectedIDs:   []int{3, 4},
		},
		{
			name:          "before 2023-01-01 - no versions before this date",
			filter:        "before 2023-01-01",
			expectedCount: 0,
			expectedIDs:   []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Model{
				versions:         versions,
				selectedVersions: make(map[int]struct{}),
				filterInput:      textinput.New(),
				filterValue:      tt.filter,
			}

			m.applyFilter()

			if len(m.filteredVersions) != tt.expectedCount {
				t.Errorf("filtered count = %d, want %d", len(m.filteredVersions), tt.expectedCount)
			}

			for _, id := range tt.expectedIDs {
				if _, ok := m.selectedVersions[id]; !ok {
					t.Errorf("expected version ID %d to be selected", id)
				}
			}
		})
	}
}

func TestModel_ApplyFilter_BeforeDateTime(t *testing.T) {
	versions := []github.PackageVersion{
		{ID: 1, Name: "v1", CreatedAt: time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC)},
		{ID: 2, Name: "v2", CreatedAt: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)},
	}

	m := &Model{
		versions:         versions,
		selectedVersions: make(map[int]struct{}),
		filterInput:      textinput.New(),
		filterValue:      "before 2024-01-15T12:00",
	}

	m.applyFilter()

	if len(m.filteredVersions) != 1 {
		t.Errorf("filtered count = %d, want 1", len(m.filteredVersions))
	}

	if _, ok := m.selectedVersions[2]; !ok {
		t.Error("expected version ID 2 to be selected")
	}
}

func TestModel_ApplyFilter_EmptyFilter(t *testing.T) {
	versions := createTestVersions()

	m := &Model{
		versions:         versions,
		selectedVersions: make(map[int]struct{}),
		filterInput:      textinput.New(),
		filterValue:      "",
	}

	m.applyFilter()

	if len(m.filteredVersions) != len(versions) {
		t.Errorf("filtered count = %d, want %d (all versions)", len(m.filteredVersions), len(versions))
	}
}

func TestModel_ApplyFilter_InvalidFilter(t *testing.T) {
	versions := createTestVersions()

	m := &Model{
		versions:         versions,
		selectedVersions: make(map[int]struct{}),
		filterInput:      textinput.New(),
		filterValue:      "invalid filter",
	}

	m.applyFilter()

	// Invalid filter should return all versions
	if len(m.filteredVersions) != len(versions) {
		t.Errorf("filtered count = %d, want %d (all versions for invalid filter)", len(m.filteredVersions), len(versions))
	}
}

func TestModel_ResetFilter(t *testing.T) {
	versions := createTestVersions()

	m := &Model{
		versions:         versions,
		filteredVersions: versions[:2], // Partial list
		selectedVersions: map[int]struct{}{1: {}, 2: {}},
		filterInput:      textinput.New(),
		filterValue:      "older 10",
		versionCursor:    5,
	}

	m.resetFilter()

	if m.filterValue != "" {
		t.Errorf("filterValue = %q, want empty", m.filterValue)
	}
	if len(m.filteredVersions) != len(versions) {
		t.Errorf("filteredVersions count = %d, want %d", len(m.filteredVersions), len(versions))
	}
	if len(m.selectedVersions) != 0 {
		t.Errorf("selectedVersions count = %d, want 0", len(m.selectedVersions))
	}
	if m.versionCursor != 0 {
		t.Errorf("versionCursor = %d, want 0", m.versionCursor)
	}
}
