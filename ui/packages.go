package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updatePackages(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.packageCursor > 0 {
				m.packageCursor--
			}
		case "down", "j":
			if m.packageCursor < len(m.packages)-1 {
				m.packageCursor++
			}
		case "enter":
			if len(m.packages) > 0 {
				m.selectedPkg = &m.packages[m.packageCursor]
				// Reset filter state for new package
				m.filterValue = ""
				m.filterInput.SetValue("")
				m.selectedVersions = make(map[int]struct{})
				m.versionCursor = 0
				m.loading = true
				m.loadingMsg = "Loading versions..."
				return m, tea.Batch(
					m.spinner.Tick,
					m.fetchVersions(),
				)
			}
		}
	}
	return m, nil
}

func (m Model) viewPackages() string {
	s := "\n"
	s += "  " + TitleStyle.Render("📦 Your Packages") + "\n"
	s += "  " + SubtitleStyle.Render("Container images in your account") + "\n\n"

	if m.loading {
		s += "  " + m.spinner.View() + " " + m.loadingMsg + "\n"
		return s
	}

	if len(m.packages) == 0 {
		s += "  " + Muted("No container packages found.") + "\n"
		s += "\n" + HelpStyle.Render("  q: quit") + "\n"
		return s
	}

	for i, pkg := range m.packages {
		cursor := "  "
		if m.packageCursor == i {
			cursor = Cursor() + " "
		}

		name := pkg.Name
		if m.packageCursor == i {
			name = SelectedStyle.Render(name)
		}

		versions := Muted(fmt.Sprintf("(%d versions)", pkg.VersionCount))
		visibility := TagStyle.Render(pkg.Visibility)

		s += fmt.Sprintf("%s%s %s %s\n", cursor, name, versions, visibility)
	}

	if m.err != nil {
		s += "\n  " + ErrorStyle.Render("✗ "+m.err.Error()) + "\n"
	}

	s += "\n" + HelpStyle.Render("  ↑/k: up • ↓/j: down • enter: select • q: quit") + "\n"

	return s
}

func (m Model) fetchVersions() tea.Cmd {
	return func() tea.Msg {
		versions, err := m.client.ListPackageVersions("container", m.selectedPkg.Name)
		if err != nil {
			return errMsg{err}
		}
		return versionsMsg{versions}
	}
}
