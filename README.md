# hij

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**hij** is a sleek Terminal User Interface (TUI) designed for managing GitHub Packages. It specializes in container image version management, allowing you to browse, filter, and bulk-delete versions with precision and speed.

Built for developers who want to keep their GitHub container registries clean without wrestling with the web UI or complex `gh` CLI commands.

https://github.com/user-attachments/assets/9114c8ff-1b20-423f-ae37-a3aa5d736bca

## ✨ Features

- **🚀 Interactive Browsing**: List all container packages in your account instantly.
- **🔃 Sort Versions**: Toggle between newest and oldest versions (`s`).
- **🔍 Smart Filtering**: Select versions by age (e.g., `:older 30`) or specific dates (e.g., `:before 2024-01-01`).
- **📦 Bulk Operations**: Toggle multiple versions or "Select All" for mass cleanup.
- **🔐 Secure Token Management**: Leverages system keychain for secure storage of your Personal Access Token.
- **⌨️ Keyboard Driven**: Optimized for efficiency with Vim-like keybindings.

## 🚀 Installation

### Automated Install (Recommended)

```bash
curl -fsSL https://raw.githubusercontent.com/maful/hij/main/install.sh | bash
```

### Using Go install

```bash
go install github.com/maful/hij
```

### From Source

```bash
git clone https://github.com/maful/hij.git
cd hij
make build
# Binary will be at ./build/hij
```

## ⚙️ Configuration

**hij** requires a GitHub Personal Access Token (PAT) with the following scopes:
- `read:packages`
- `delete:packages`

You can provide the token in three ways (checked in priority order):
1. `HIJ_GITHUB_TOKEN` environment variable.
2. System Keychain (macOS Keychain, Linux Secret Service, Windows Credential Manager).
3. Interactive prompt upon first run (with an option to save to keychain).

## 🎮 Usage

Launch the TUI:
```bash
hij
```

### Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `j/k` | Navigate lists |
| `Space` | Toggle selection |
| `a` | Select all versions |
| `n` | Deselect all versions |
| `/` or `:` | Open filter input |
| `s` | Toggle sort order (newest/oldest) |
| `d` | Initiate deletion of selected versions |
| `Esc` | Go back |
| `q` | Quit |

### Filtering Commands

Inside the version list, press `:` to filter:
- `:older <days>` — Select versions older than N days (e.g., `:older 10`).
- `:before <date>` — Select versions before a date (e.g., `:before 2024-01-01`).

### CLI Commands

```bash
hij                # Interactive menu (TUI)
hij version        # Show installed version
hij update         # Update to latest version
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`push origin feature/AmazingFeature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.
