# SSH Manager

A modern desktop SSH client and manager built with [Wails](https://wails.app/) (Go + Vue 3). Manage your SSH configurations, keys, hosts, and connect to servers with an integrated terminal — all from a single, beautiful interface.

![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)
![License](https://img.shields.io/badge/license-MIT-green)

## Features

### SSH Config Management
- Parse and display `~/.ssh/config` automatically
- Visual editor with syntax highlighting
- Real-time config validation with warnings
- Auto-save and manual save
- Parsed host preview cards

### SSH Key Management
- List all existing SSH keys (RSA, ED25519, ECDSA)
- View public key content with one-click copy
- Generate new keys with customizable options:
  - Type: RSA / ED25519 / ECDSA
  - Key size / curve selection
  - Email comment
  - Optional passphrase
- Delete keys with confirmation dialog

### Host Management
- **Favorites**: Save and organize frequently used hosts
- **Groups & Tags**: Categorize hosts for quick filtering
- **Import from Config**: Auto-parse hosts from `~/.ssh/config`
- **Connection Testing**: TCP ping / handshake test
- **CRUD Operations**: Add, edit, delete favorite hosts

### Connection History
- Track all SSH connections (host, IP, port, user, timestamp)
- Search and filter by hostname, alias, or user
- Delete individual history entries

### SSH Terminal
- Built-in terminal emulator powered by xterm.js
- Multi-tab support for simultaneous connections
- Connection status indicators (green dot = connected)
- Automatic PTY allocation with xterm-256color
- Terminal resize handling
- GitHub-dark theme

### Security & Auditing
- Operation audit log (config changes, key generation, connections)
- SQLite-based local data persistence
- Private key file permissions (0600)

## Screenshots

> _(Add screenshots here)_

## Prerequisites

- **Go** 1.21 or later
- **Node.js** 18 or later
- **Wails CLI** (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- **GCC** or compatible C compiler (for `go-sqlite3` CGO)

## Installation

### 1. Clone the repository

```bash
git clone https://github.com/your-username/ssh-manager.git
cd ssh-manager
```

### 2. Install frontend dependencies

```bash
cd frontend
npm install
cd ..
```

### 3. Install Go dependencies

```bash
go mod tidy
```

### 4. Build and run

```bash
# Development mode (hot reload)
wails dev

# Production build
wails build
```

The built binary will be in `build/bin/`.

## Project Structure

```
ssh-manager/
├── main.go                      # Wails app entry point
├── app.go                       # Backend service bindings (Wails API)
├── wails.json                   # Wails configuration
├── go.mod                       # Go module definition
├── .gitignore
├── README.md
├── AGENTS.md
├── ARCHITECTURE.md
├── CONTRIBUTING.md
│
├── backend/
│   ├── config/
│   │   └── config.go            # SSH config file parsing, validation, save
│   ├── hosts/
│   │   └── hosts.go             # Host management, import, test connection
│   ├── keys/
│   │   └── keys.go              # SSH key generation, listing, deletion
│   ├── ssh/
│   │   └── ssh.go               # SSH connection, terminal session management
│   └── storage/
│       └── database.go          # SQLite database layer (history, favorites, logs)
│
└── frontend/
    ├── index.html               # HTML entry point
    ├── vite.config.js           # Vite configuration
    ├── package.json             # Frontend dependencies
    └── src/
        ├── main.js              # Vue app bootstrap
        ├── App.vue              # Root component with sidebar layout
        └── views/
            ├── ConfigEditor.vue # SSH config editor view
            ├── KeyManager.vue   # SSH key management view
            ├── HostList.vue     # Host favorites & history view
            └── TerminalView.vue # SSH terminal with xterm.js
```

## Usage

### SSH Config Editor

The editor automatically loads your `~/.ssh/config` on startup. You can:

- Edit the raw config text directly
- Click **Validate** to check for syntax errors
- Click **Save** to write changes back to disk
- View parsed hosts as cards below the editor

### Key Manager

- Click **Generate New Key** to create a new SSH key pair
- Select key type (ED25519 recommended), size, and optional passphrase
- Click **View PubKey** on any key to see and copy the public key
- Delete keys with the **Delete** button (requires confirmation)

### Host Management

- **Import from Config**: Parses `~/.ssh/config` and adds all hosts to favorites
- **Add Host**: Manually add a host with alias, hostname, port, user, tags, and group
- **Test**: Performs a TCP connection test to verify reachability
- **Connect**: Opens the terminal tab and connects to the host

### Terminal

- Click **+** to open a new connection dialog
- Enter host details and click **Connect**
- Multiple tabs supported — each tab is an independent SSH session
- Status indicator shows connection state (green = connected)

## Configuration

### Default Paths

| Item | Path |
|------|------|
| SSH Config | `~/.ssh/config` |
| SSH Keys | `~/.ssh/` |
| Application DB | `~/.ssh-manager/data.db` |

### Supported SSH Key Types

| Type | Default Size | Notes |
|------|-------------|-------|
| ED25519 | 256 bits | **Recommended** — fast, secure, compact |
| RSA | 4096 bits | Widely compatible, larger keys |
| ECDSA | 256/384/521 bits | NIST curves, good balance |

## Development

```bash
# Run in development mode with hot reload
wails dev

# Build for production
wails build

# Build with debug info
wails build -debug
```

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Framework | [Wails v2](https://wails.app/) |
| Backend | Go 1.21+ |
| Frontend | Vue 3 + Vite |
| Terminal | xterm.js + xterm-addon-fit |
| Database | SQLite (go-sqlite3) |
| SSH | golang.org/x/crypto/ssh |

## Security Notes

- Private keys are generated with `0600` file permissions
- Host key verification uses `InsecureIgnoreHostKey()` by default — for production use, implement proper known_hosts checking
- No passwords are stored — authentication relies on SSH keys and agent
- Audit logs record all significant operations

## Roadmap

- [ ] Password/passphrase authentication support
- [ ] SSH agent forwarding in terminal
- [ ] Port forwarding / tunnel management
- [ ] SFTP file browser
- [ ] Config diff viewer
- [ ] Theme customization
- [ ] Known hosts verification
- [ ] Connection profiles / templates
- [ ] Import from other tools (PuTTY, etc.)

## License

MIT

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.
