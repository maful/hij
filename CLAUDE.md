# hij

A terminal user interface (TUI) for managing GitHub Packages. Delete container image versions individually or in bulk with an intuitive, keyboard-driven interface.

## Overview

**hij** connects to the GitHub Packages API using a Personal Access Token (PAT) and provides:

- **Browse packages** — View all container packages in your account
- **Version management** — List and select package versions for deletion
- **Smart filtering** — Select versions by age (`:older 10`) or date (`:before 2024-01-01`)
- **Bulk deletion** — Delete multiple versions in a single operation

## Project Structure

```
hij/
├── main.go              # Entry point
├── github/
│   ├── client.go        # GitHub API client
│   └── types.go         # Package & version types
└── ui/
    ├── app.go           # Main app model & navigation
    ├── token.go         # PAT input screen
    ├── packages.go      # Package list screen
    ├── versions.go      # Version list & filtering
    ├── confirm.go       # Deletion confirmation
    └── styles.go        # TUI styling
```

## Tech Stack

- **Go 1.23**
- **[Bubble Tea](https://github.com/charmbracelet/bubbletea)** — TUI framework (Elm architecture)
- **[Bubbles](https://github.com/charmbracelet/bubbles)** — Input components (text input, spinner)
- **[Lip Gloss](https://github.com/charmbracelet/lipgloss)** — Terminal styling

## Key Commands

```bash
# Run
go run .

# Build
go build -o hij
```

## Navigation

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate |
| `Space` | Toggle selection |
| `a` | Select all |
| `n` | Deselect all |
| `/` or `:` | Open filter |
| `d` | Delete selected |
| `Esc` | Back |
| `q` | Quit |

## Filters

- `:older N` — Select versions older than N days
- `:before DATE` — Select versions before date (e.g., `2024-01-01`)
- `:before DATETIME` — Select versions before datetime (e.g., `2024-01-01T15:00`)

## Required Permissions

PAT needs `read:packages` and `delete:packages` scopes.
