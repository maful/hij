package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) updateConfirm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.deleting {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.deleting = true
			m.deleteIdx = 0
			m.deleteErrs = nil
			return m, tea.Batch(
				m.spinner.Tick,
				m.deleteNextVersion(),
			)
		case "n", "N", "esc":
			m.screen = ScreenVersions
			return m, nil
		}
	}
	return m, nil
}

func (m Model) handleDeleteResult(msg deleteResultMsg) (tea.Model, tea.Cmd) {
	if msg.err != nil {
		m.deleteErrs = append(m.deleteErrs, msg.err)
	}

	m.deleteIdx++

	// Check if we're done
	if m.deleteIdx >= len(m.selectedVersions) {
		m.deleting = false
		if len(m.deleteErrs) == 0 {
			// Success - go back to packages
			m.selectedVersions = make(map[int]struct{})
			m.screen = ScreenPackages
			m.loading = true
			m.loadingMsg = "Refreshing packages..."
			return m, tea.Batch(
				m.spinner.Tick,
				m.fetchPackages(),
			)
		}
		// Stay on confirm screen showing errors
		return m, nil
	}

	// Delete next version
	return m, m.deleteNextVersion()
}

func (m Model) deleteNextVersion() tea.Cmd {
	// Get the Nth selected version
	count := 0
	for _, v := range m.versions {
		if _, ok := m.selectedVersions[v.ID]; ok {
			if count == m.deleteIdx {
				return func() tea.Msg {
					err := m.client.DeletePackageVersion("container", m.selectedPkg.Name, v.ID)
					return deleteResultMsg{idx: v.ID, err: err}
				}
			}
			count++
		}
	}
	return nil
}

func (m Model) viewConfirm() string {
	s := "\n"
	s += "  " + TitleStyle.Render("⚠️  Confirm Deletion") + "\n\n"

	if m.deleting {
		s += "  " + m.spinner.View() + fmt.Sprintf(" Deleting... (%d/%d)\n", m.deleteIdx+1, len(m.selectedVersions))

		if len(m.deleteErrs) > 0 {
			s += "\n  " + ErrorStyle.Render(fmt.Sprintf("%d errors occurred", len(m.deleteErrs))) + "\n"
		}

		return s
	}

	// Show what will be deleted
	s += "  " + WarningStyle.Render("You are about to delete:") + "\n\n"
	s += fmt.Sprintf("  • %d version(s) from %s\n", len(m.selectedVersions), SelectedStyle.Render(m.selectedPkg.Name))

	// List selected versions (max 5)
	count := 0
	for _, v := range m.versions {
		if _, ok := m.selectedVersions[v.ID]; ok && count < 5 {
			name := v.Name
			if len(name) > 20 {
				name = name[:20] + "…"
			}
			s += fmt.Sprintf("    - %s %s\n", name, TagStyle.Render(v.TagsString()))
			count++
		}
	}
	if len(m.selectedVersions) > 5 {
		s += fmt.Sprintf("    ... and %d more\n", len(m.selectedVersions)-5)
	}

	// Show errors if any
	if len(m.deleteErrs) > 0 {
		s += "\n  " + ErrorStyle.Render("Errors:") + "\n"
		for i, err := range m.deleteErrs {
			if i >= 3 {
				s += fmt.Sprintf("    ... and %d more errors\n", len(m.deleteErrs)-3)
				break
			}
			s += "    • " + err.Error() + "\n"
		}
	}

	s += "\n  " + Danger("This action cannot be undone!") + "\n"
	s += "\n  " + Muted("Delete these versions? ") + SelectedStyle.Render("[y/n]") + "\n"

	return s
}
