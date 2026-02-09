package ui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Color palette
	primaryColor   = lipgloss.Color("#60A5FA") // Soft blue
	secondaryColor = lipgloss.Color("#06B6D4") // Cyan
	successColor   = lipgloss.Color("#10B981") // Green
	dangerColor    = lipgloss.Color("#EF4444") // Red
	warningColor   = lipgloss.Color("#F59E0B") // Amber
	mutedColor     = lipgloss.Color("#6B7280") // Gray
	textColor      = lipgloss.Color("#F9FAFB") // Light
	bgColor        = lipgloss.Color("#1F2937") // Dark

	// Title style
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Subtitle/description style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginBottom(1)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	// Selected item style
	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor)

	// Cursor style
	CursorStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true)

	// Checked item style
	CheckedStyle = lipgloss.NewStyle().
			Foreground(successColor)

	// Unchecked item style
	UncheckedStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(dangerColor).
			Bold(true)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true)

	// Warning style
	WarningStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Tag style
	TagStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Italic(true)

	// Date style
	DateStyle = lipgloss.NewStyle().
			Foreground(mutedColor)

	// Old version style (for versions older than threshold)
	OldVersionStyle = lipgloss.NewStyle().
			Foreground(warningColor)

	// Box style for sections
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	// Input style
	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(mutedColor).
			Padding(0, 1).
			Width(44).
			MarginLeft(2)

	// Focused input style
	FocusedInputStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(primaryColor).
				Padding(0, 1).
				Width(44).
				MarginLeft(2)

	// Spinner style
	SpinnerStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	// Logo/brand style
	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Background(lipgloss.Color("#1E3A5F")).
			Padding(0, 1)
)

// Helper functions
func Cursor() string {
	return CursorStyle.Render("❯")
}

func Checked() string {
	return CheckedStyle.Render("✓")
}

func Unchecked() string {
	return UncheckedStyle.Render("○")
}

func Danger(s string) string {
	return ErrorStyle.Render(s)
}

func Success(s string) string {
	return SuccessStyle.Render(s)
}

func Muted(s string) string {
	return lipgloss.NewStyle().Foreground(mutedColor).Render(s)
}
