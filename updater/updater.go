package updater

import (
	"fmt"
	"os"
	"strings"

	"github.com/blang/semver"
	"github.com/charmbracelet/lipgloss"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("160")).Bold(true)
)

// Update handles the self-update process
func Update(currentVersion string) {
	fmt.Println(infoStyle.Render("Checking for updates..."))

	v, err := parseVersion(currentVersion)
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error parsing current version: %s", err)))
		return
	}

	latest, found, err := selfupdate.DetectLatest("maful/hij")
	if err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error checking for updates: %s", err)))
		return
	}

	if !found || latest.Version.LTE(v) {
		fmt.Println(infoStyle.Render(fmt.Sprintf("Current version %s is the latest", currentVersion)))
		return
	}

	fmt.Println(infoStyle.Render(fmt.Sprintf("Found new version: %s", latest.Version)))
	fmt.Println(infoStyle.Render("Release notes:\n" + latest.ReleaseNotes))
	fmt.Print(infoStyle.Render("Do you want to update? (y/n): "))

	var input string
	if _, err := fmt.Scanln(&input); err != nil {
		// Ignore error and assume no/cancel
	}
	if input != "y" && input != "Y" {
		fmt.Println(infoStyle.Render("Update cancelled."))
		return
	}

	exe, err := os.Executable()
	if err != nil {
		fmt.Println(errorStyle.Render("Could not locate executable path"))
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, exe); err != nil {
		fmt.Println(errorStyle.Render(fmt.Sprintf("Error occurred while updating binary: %s", err)))
		return
	}

	fmt.Println(successStyle.Render(fmt.Sprintf("Successfully updated to version %s", latest.Version)))
}

func parseVersion(v string) (semver.Version, error) {
	return semver.Parse(strings.TrimPrefix(v, "v"))
}
