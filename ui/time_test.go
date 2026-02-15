package ui

import (
	"fmt"
	"testing"
	"time"
)

func TestHumanizeTime(t *testing.T) {
	// Fix "now" for testing purposes requires mocking time.Now() which is hard in Go without a library or interface.
	// Instead, we derive input times from time.Now() relative to the test execution time.
	
	now := time.Now()

	tests := []struct {
		name     string
		input    time.Time
		expected string // We might need flexible matching for "seconds ago" if test runs slow
		check    func(string) bool // Custom check function if strict equality is brittle
	}{
		{
			name:     "Just now",
			input:    now,
			expected: "just now",
		},
		{
			name:     "A second ago",
			input:    now.Add(-1500 * time.Millisecond),
			expected: "a second ago",
		},
		{
			name:     "Seconds ago",
			input:    now.Add(-32 * time.Second),
			expected: "32 seconds ago",
		},
		{
			name:     "A minute ago",
			input:    now.Add(-90 * time.Second),
			expected: "a minute ago",
		},
		{
			name:     "Minutes ago",
			input:    now.Add(-15 * time.Minute),
			expected: "15 minutes ago",
		},
		{
			name:     "An hour ago",
			input:    now.Add(-90 * time.Minute),
			expected: "an hour ago",
		},
		{
			name:     "Hours ago",
			input:    now.Add(-14 * time.Hour),
			expected: "14 hours ago",
		},
		{
			name:  "Yesterday",
			input: now.AddDate(0, 0, -1),
			check: func(s string) bool {
				// s should be "yesterday at 5:22pm"
				expectedSuffix := now.AddDate(0, 0, -1).Format("3:04pm")
				return s == fmt.Sprintf("yesterday at %s", expectedSuffix)
			},
		},
		{
			name:  "This week", // e.g., 3 days ago
			input: now.AddDate(0, 0, -3),
			check: func(s string) bool {
				// e.g. "Tuesday at 12:48am"
				expectedTime := now.AddDate(0, 0, -3)
				// If 3 days ago was yesterday (e.g. now is Tuesday, 3 days ago is Saturday? No. 3 days ago is Saturday.)
				// Wait, if 3 days ago was yesterday, it would be caught by "Yesterday".
				// So we pick 3 days ago.
				// However, if "Yesterday" covers 1 day ago.
				return s == fmt.Sprintf("%s at %s", expectedTime.Format("Monday"), expectedTime.Format("3:04pm"))
			},
		},
		{
			name:  "This year",
			input: time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, time.Local),
			check: func(s string) bool {
				// If today is Jan 1st, this might be "just now" or "hours ago".
				// Let's force it to be earlier this year but > 7 days ago.
				// If we are in Jan, use Jan 1st?
				target := time.Date(now.Year(), time.January, 1, 10, 0, 0, 0, time.Local)
				if now.Sub(target) < 7*24*time.Hour {
					// Too close to now, might fall into "This week" or "Yesterday"
					// Skip this test case dynamically if we are early in the year?
					// Let's skip strict check if we assume the logic assumes > 7 days.
					return true 
				}
				return s == "on Jan 01"
			},
		},
		// Case for "This year" that is definitely > 7 days ago.
		// If current date is < Jan 8th, we can't test "This year > 7 days ago" easily without mocking "now".
		// But let's assume standard case.
		
		{
			name:     "Last year",
			input:    now.AddDate(-1, 0, 0), // Exactly 1 year ago
			check: func(s string) bool {
				target := now.AddDate(-1, 0, 0)
				// "on Jan 31, 2012"
				return s == fmt.Sprintf("on %s", target.Format("Jan 02, 2006"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := HumanizeTime(tc.input)
			if tc.expected != "" {
				if got != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, got)
				}
			} else if tc.check != nil {
				if !tc.check(got) {
					t.Errorf("check failed for %q", got)
				}
			}
		})
	}
}
