package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/maful/hij/github"
)

func (m Model) updateToken(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			token := m.tokenInput.Value()
			if token == "" {
				m.err = fmt.Errorf("token cannot be empty")
				return m, nil
			}
			m.client = github.NewClient(token)
			m.loading = true
			m.loadingMsg = "Validating token..."
			m.err = nil
			return m, tea.Batch(
				m.spinner.Tick,
				m.fetchPackages(),
			)
		}
	}

	var cmd tea.Cmd
	m.tokenInput, cmd = m.tokenInput.Update(msg)
	return m, cmd
}

func (m Model) viewToken() string {
	s := "\n"
	s += "  " + LogoStyle.Render(" hij ") + "\n"
	s += "  " + TitleStyle.Render("GitHub Packages Cleaner") + "\n\n"

	s += "  " + SubtitleStyle.Render("Enter your GitHub Personal Access Token") + "\n"
	s += "  " + Muted("Required scopes: read:packages, delete:packages") + "\n\n"

	if m.loading {
		s += "  " + m.spinner.View() + " " + m.loadingMsg + "\n"
	} else {
		s += "  " + m.tokenInput.View() + "\n"
	}

	if m.err != nil {
		s += "\n  " + ErrorStyle.Render("✗ "+m.err.Error()) + "\n"
	}

	s += "\n" + HelpStyle.Render("  enter: submit • ctrl+c: quit") + "\n"

	return s
}

func (m Model) fetchPackages() tea.Cmd {
	return func() tea.Msg {
		packages, err := m.client.ListPackages("container")
		if err != nil {
			return errMsg{err}
		}
		return packagesMsg{packages}
	}
}
