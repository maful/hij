package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/maful/hij/github"
)

// Screen represents the current screen in the app
type Screen int

const (
	ScreenToken Screen = iota
	ScreenPackages
	ScreenVersions
	ScreenConfirm
)

// Model is the main application model
type Model struct {
	screen       Screen
	client       *github.Client
	err          error
	loading      bool
	loadingMsg   string
	spinner      spinner.Model
	quitting     bool

	// Token screen
	tokenInput textinput.Model

	// Packages screen
	packages       []github.Package
	packageCursor  int
	selectedPkg    *github.Package

	// Versions screen
	versions        []github.PackageVersion
	versionCursor   int
	selectedVersions map[int]struct{}
	filterInput     textinput.Model
	filterActive    bool
	filterValue     string

	// Confirm screen
	confirmYes bool
	deleting   bool
	deleteIdx  int
	deleteErrs []error
}

// New creates a new application model
func New() Model {
	ti := textinput.New()
	ti.Placeholder = "ghp_xxxxxxxxxxxxxxxxxxxx"
	ti.Focus()
	ti.CharLimit = 100
	ti.Width = 50
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	fi := textinput.New()
	fi.Placeholder = ":older 10 or :before 2024-01-01"
	fi.CharLimit = 50
	fi.Width = 40

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = SpinnerStyle

	return Model{
		screen:           ScreenToken,
		tokenInput:       ti,
		filterInput:      fi,
		spinner:          s,
		selectedVersions: make(map[int]struct{}),
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "q":
			if m.screen != ScreenToken && !m.filterActive {
				m.quitting = true
				return m, tea.Quit
			}
		case "esc":
			if m.filterActive {
				m.filterActive = false
				m.filterInput.Blur()
				return m, nil
			}
			// Go back
			switch m.screen {
			case ScreenVersions:
				m.screen = ScreenPackages
				m.selectedVersions = make(map[int]struct{})
				return m, nil
			case ScreenConfirm:
				m.screen = ScreenVersions
				return m, nil
			}
		}
	case errMsg:
		m.loading = false
		m.err = msg.err
		return m, nil
	case packagesMsg:
		m.loading = false
		m.packages = msg.packages
		m.screen = ScreenPackages
		return m, nil
	case versionsMsg:
		m.loading = false
		m.versions = msg.versions
		m.screen = ScreenVersions
		return m, nil
	case deleteResultMsg:
		return m.handleDeleteResult(msg)
	case spinner.TickMsg:
		if m.loading || m.deleting {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	// Delegate to screen-specific update
	switch m.screen {
	case ScreenToken:
		return m.updateToken(msg)
	case ScreenPackages:
		return m.updatePackages(msg)
	case ScreenVersions:
		return m.updateVersions(msg)
	case ScreenConfirm:
		return m.updateConfirm(msg)
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	if m.quitting {
		return "\n  " + Success("Goodbye!") + "\n\n"
	}

	switch m.screen {
	case ScreenToken:
		return m.viewToken()
	case ScreenPackages:
		return m.viewPackages()
	case ScreenVersions:
		return m.viewVersions()
	case ScreenConfirm:
		return m.viewConfirm()
	}

	return ""
}

// Custom messages
type errMsg struct{ err error }
type packagesMsg struct{ packages []github.Package }
type versionsMsg struct{ versions []github.PackageVersion }
type deleteResultMsg struct {
	idx int
	err error
}
