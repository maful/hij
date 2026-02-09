package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/maful/hij/config"
	"github.com/maful/hij/github"
)

// tokenValidatedMsg is sent when token validation succeeds
type tokenValidatedMsg struct {
	token       string
	fromKeychain bool
}

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
			m.pendingToken = token
			m.tokenFromKeychain = false
			return m, tea.Batch(
				m.spinner.Tick,
				m.fetchPackages(),
			)
		case "s": // Save token to keychain when prompted
			if m.showSavePrompt {
				if err := config.SaveToken(m.pendingToken); err != nil {
					m.err = fmt.Errorf("failed to save token: %w", err)
				}
				m.showSavePrompt = false
				m.screen = ScreenPackages
				return m, nil
			}
		case "n": // Skip saving
			if m.showSavePrompt {
				m.showSavePrompt = false
				m.screen = ScreenPackages
				return m, nil
			}
		}
	case packagesMsg:
		// Token validated successfully
		m.loading = false
		m.packages = msg.packages
		// If token came from manual input (not keychain), offer to save
		if !m.tokenFromKeychain && m.pendingToken != "" {
			m.showSavePrompt = true
			return m, nil
		}
		m.screen = ScreenPackages
		return m, nil
	}

	var cmd tea.Cmd
	m.tokenInput, cmd = m.tokenInput.Update(msg)
	return m, cmd
}

func (m Model) viewToken() string {
	s := "\n"
	s += "  " + LogoStyle.Render(" hij ") + "\n"
	s += "  " + TitleStyle.Render("GitHub Packages Cleaner") + "\n\n"

	// Show save prompt if needed
	if m.showSavePrompt {
		s += "  " + Success("✓ Token validated!") + "\n\n"
		s += "  " + SubtitleStyle.Render("Save token to keychain for future use?") + "\n\n"
		s += "\n" + HelpStyle.Render("  s: save • n: skip") + "\n"
		return s
	}

	s += "  " + SubtitleStyle.Render("Enter your GitHub Personal Access Token") + "\n"
	s += "  " + Muted("Required scopes: read:packages, delete:packages") + "\n"
	s += "  " + Muted("Tip: Set HIJ_GITHUB_TOKEN env var to skip this step") + "\n\n"

	if m.loading {
		s += "  " + m.spinner.View() + " " + m.loadingMsg + "\n"
	} else {
		s += "  " + m.tokenInput.View() + "\n"
	}

	if m.err != nil {
		s += "\n  " + ErrorStyle.Render("✗ "+m.err.Error()) + "\n"
	}

	s += "\n" + HelpStyle.Render("  enter: submit • ctrl+c: quit") + "\n"
	s += "\n" + Muted("  built with love by @mafulprayoga") + "\n"

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
