package ui

import (
	"fmt"
	"time"
)

// HumanizeTime returns a human-readable string representing the time elapsed since the given time.
// It mimics the behavior of Rails' local_time helper.
//
// Examples:
// Recent: "a second ago", "32 seconds ago", "an hour ago", "14 hours ago"
// Yesterday: "yesterday at 5:22pm"
// This week: "Tuesday at 12:48am"
// This year: "on Nov 17"
// Last year: "on Jan 31, 2012"
func HumanizeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	// Future time (shouldn't happen for version creation, but handle gracefully)
	if diff < 0 {
		return "in the future"
	}

	// Recent (less than 1 day)
	if diff < 24*time.Hour {
		if diff < time.Second {
			return "just now"
		}
		if diff < 2*time.Second {
			return "a second ago"
		}
		if diff < time.Minute {
			return fmt.Sprintf("%d seconds ago", int(diff.Seconds()))
		}
		if diff < 2*time.Minute {
			return "a minute ago"
		}
		if diff < time.Hour {
			return fmt.Sprintf("%d minutes ago", int(diff.Minutes()))
		}
		if diff < 2*time.Hour {
			return "an hour ago"
		}
		return fmt.Sprintf("%d hours ago", int(diff.Hours()))
	}

	// Yesterday
	// Check if it was yesterday by truncating to day
	yesterday := now.AddDate(0, 0, -1)
	if isSameDay(t, yesterday) {
		return fmt.Sprintf("yesterday at %s", t.Format("3:04pm"))
	}

	// This week (within last 6 days, excluding today and yesterday which are handled above)
	// Actually, "This week" usually means within the last 7 days or since Monday.
	// The example says "Tuesday at 12:48am". Let's assume it means within the last 6-7 days.
	if diff < 7*24*time.Hour {
		return fmt.Sprintf("%s at %s", t.Format("Monday"), t.Format("3:04pm"))
	}

	// This year
	if t.Year() == now.Year() {
		return fmt.Sprintf("on %s", t.Format("Jan 02"))
	}

	// older
	return fmt.Sprintf("on %s", t.Format("Jan 02, 2006"))
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
