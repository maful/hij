package ui

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/maful/hij/github"
)

func (m Model) updateVersions(msg tea.Msg) (tea.Model, tea.Cmd) {
	// If filter is active, handle text input
	if m.filterActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				m.filterValue = m.filterInput.Value()
				m.filterActive = false
				m.filterInput.Blur()
				m.applyFilter()
				return m, nil
			case "esc":
				m.filterActive = false
				m.filterInput.Blur()
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.filterInput, cmd = m.filterInput.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Clear success message on any key press
		m.successMsg = ""

		switch msg.String() {
		case "up", "k":
			if m.versionCursor > 0 {
				m.versionCursor--
			}
		case "down", "j":
			if m.versionCursor < len(m.filteredVersions)-1 {
				m.versionCursor++
			}
		case " ": // Space to toggle selection
			if len(m.filteredVersions) > 0 {
				idx := m.filteredVersions[m.versionCursor].ID
				if _, ok := m.selectedVersions[idx]; ok {
					delete(m.selectedVersions, idx)
				} else {
					m.selectedVersions[idx] = struct{}{}
				}
			}
		case "a": // Select all (filtered)
			for _, v := range m.filteredVersions {
				m.selectedVersions[v.ID] = struct{}{}
			}
		case "n": // Deselect all
			m.selectedVersions = make(map[int]struct{})
		case "/", ":": // Activate filter
			m.filterActive = true
			m.filterInput.Focus()
			return m, nil
		case "c": // Clear filter
			m.resetFilter()
		case "s": // Toggle sort
			if m.sortOrder == "newest" {
				m.sortOrder = "oldest"
			} else {
				m.sortOrder = "newest"
			}
			m.sortVersions(m.filteredVersions)
		case "d": // Delete selected
			if len(m.selectedVersions) > 0 {
				m.screen = ScreenConfirm
				m.confirmYes = false
			}
		case "esc":
			m.screen = ScreenPackages
			m.selectedVersions = make(map[int]struct{})
			m.filterValue = ""
			m.filteredVersions = nil
			return m, nil
		}
	}
	return m, nil
}

// resetFilter clears the filter and restores all versions
func (m *Model) resetFilter() {
	m.filterValue = ""
	m.filterInput.SetValue("")
	m.filteredVersions = m.versions
	m.selectedVersions = make(map[int]struct{})
	m.versionCursor = 0
}

func (m *Model) applyFilter() {
	filter := strings.TrimSpace(m.filterValue)
	if filter == "" {
		m.filteredVersions = m.versions
		return
	}

	// Reset filtered list and selections
	m.filteredVersions = nil
	m.selectedVersions = make(map[int]struct{})
	m.versionCursor = 0

	// Parse :older N format (colon is optional since : key activates filter mode)
	olderRegex := regexp.MustCompile(`^:?older\s+(\d+)$`)
	if matches := olderRegex.FindStringSubmatch(filter); len(matches) == 2 {
		days, _ := strconv.Atoi(matches[1])
		cutoff := time.Now().AddDate(0, 0, -days)
		for _, v := range m.versions {
			if v.CreatedAt.Before(cutoff) {
				m.filteredVersions = append(m.filteredVersions, v)
				m.selectedVersions[v.ID] = struct{}{}
			}
		}
		return
	}

	// Parse :before DATE or :before DATETIME format (colon is optional)
	beforeRegex := regexp.MustCompile(`^:?before\s+(.+)$`)
	if matches := beforeRegex.FindStringSubmatch(filter); len(matches) == 2 {
		dateStr := matches[1]
		var cutoff time.Time
		var err error

		// Try datetime format first
		cutoff, err = time.Parse("2006-01-02T15:04", dateStr)
		if err != nil {
			// Try date-only format
			cutoff, err = time.Parse("2006-01-02", dateStr)
		}

		if err == nil {
			for _, v := range m.versions {
				if v.CreatedAt.Before(cutoff) {
					m.filteredVersions = append(m.filteredVersions, v)
					m.selectedVersions[v.ID] = struct{}{}
				}
			}
		}
		return
	}

	// If no filter pattern matched, show all versions
	m.filteredVersions = m.versions
	m.sortVersions(m.filteredVersions)
}

func (m *Model) sortVersions(versions []github.PackageVersion) {
	sort.Slice(versions, func(i, j int) bool {
		if m.sortOrder == "oldest" {
			return versions[i].CreatedAt.Before(versions[j].CreatedAt)
		}
		return versions[i].CreatedAt.After(versions[j].CreatedAt)
	})
}

func (m Model) viewVersions() string {
	s := "\n"
	s += "  " + TitleStyle.Render("📋 "+m.selectedPkg.Name) + "\n"
	s += "  " + SubtitleStyle.Render("Select versions to delete") + "\n\n"

	if m.loading {
		s += "  " + m.spinner.View() + " " + m.loadingMsg + "\n"
		return s
	}

	// Show success message if present
	if m.successMsg != "" {
		s += "  " + SuccessStyle.Render("✓ "+m.successMsg) + "\n\n"
	}

	// Filter input
	if m.filterActive {
		s += FocusedInputStyle.Render(m.filterInput.View()) + "\n\n"
	} else if m.filterValue != "" {
		s += "  " + Muted("Filter: ") + TagStyle.Render(m.filterValue) + "  "
	} else {
		s += "  "
	}

	// Sort indicator
	sortIcon := "↓"
	if m.sortOrder == "oldest" {
		sortIcon = "↑"
	}
	s += Muted("Sort: ") + TagStyle.Render(m.sortOrder+" "+sortIcon) + "\n\n"

	if len(m.filteredVersions) == 0 {
		if m.filterValue != "" {
			s += "  " + Muted("No versions match the filter.") + "\n"
		} else {
			s += "  " + Muted("No versions found.") + "\n"
		}
		s += "\n" + HelpStyle.Render("  c: clear filter • esc: back • q: quit") + "\n"
		return s
	}

	// Version list (show max 15)
	start := 0
	end := len(m.filteredVersions)
	maxVisible := 15

	if len(m.filteredVersions) > maxVisible {
		start = m.versionCursor - maxVisible/2
		if start < 0 {
			start = 0
		}
		end = start + maxVisible
		if end > len(m.filteredVersions) {
			end = len(m.filteredVersions)
			start = end - maxVisible
		}
	}

	for i := start; i < end; i++ {
		v := m.filteredVersions[i]
		cursor := "  "
		if m.versionCursor == i {
			cursor = Cursor() + " "
		}

		// Checkbox
		checkbox := Unchecked()
		if _, ok := m.selectedVersions[v.ID]; ok {
			checkbox = Checked()
		}

		// Version name (truncate if too long)
		name := v.Name
		if len(name) > 12 {
			name = name[:12] + "…"
		}

		// Tags
		tags := v.TagsString()

		// Age
		// Age
		ageStr := HumanizeTime(v.CreatedAt)
		if time.Since(v.CreatedAt) > 30*24*time.Hour {
			ageStr = OldVersionStyle.Render(ageStr)
		} else {
			ageStr = DateStyle.Render(ageStr)
		}

		// Format the row
		row := fmt.Sprintf("%s%s %s  %s  %s", cursor, checkbox, name, TagStyle.Render(tags), ageStr)
		s += row + "\n"
	}

	// Selection count
	s += "\n  " + Muted(fmt.Sprintf("Selected: %d of %d", len(m.selectedVersions), len(m.filteredVersions))) + "\n"

	if m.err != nil {
		s += "\n  " + ErrorStyle.Render("✗ "+m.err.Error()) + "\n"
	}

	s += "\n" + HelpStyle.Render("  space: toggle • a: all • n: none • /: filter • s: sort • c: clear • d: delete • esc: back") + "\n"

	return s
}
