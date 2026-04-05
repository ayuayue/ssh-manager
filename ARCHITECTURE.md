# ARCHITECTURE.md

# SSH Manager Architecture

## System Overview

SSH Manager is a cross-platform desktop application using the Wails framework. Wails embeds a Vue 3 SPA inside a native WebView, with Go handling all backend logic. Communication between frontend and backend happens through Wails' automatic JS-Go binding layer.

```
┌──────────────────────────────────────────────────────────┐
│                     Desktop OS                            │
│  ┌────────────────────────────────────────────────────┐  │
│  │                  Wails Binary                       │  │
│  │  ┌─────────────────┐     ┌──────────────────────┐  │  │
│  │  │   Go Backend    │◄───►│   Wails Runtime      │  │  │
│  │  │                 │     │   (JS ↔ Go Bridge)   │  │  │
│  │  │  app.go         │     │                      │  │  │
│  │  │  backend/       │     │   Bind: App struct   │  │  │
│  │  │  storage/       │     │   Embed: frontend/   │  │  │
│  │  └─────────────────┘     └──────────┬───────────┘  │  │
│  │                                     │              │  │
│  │  ┌──────────────────────────────────▼───────────┐  │  │
│  │  │           WebView (Edge/WebKit)              │  │  │
│  │  │  ┌────────────────────────────────────────┐  │  │  │
│  │  │  │         Vue 3 SPA                      │  │  │  │
│  │  │  │  App.vue (sidebar + routing)           │  │  │  │
│  │  │  │  views/ConfigEditor.vue                │  │  │  │
│  │  │  │  views/KeyManager.vue                  │  │  │  │
│  │  │  │  views/HostList.vue                    │  │  │  │
│  │  │  │  views/TerminalView.vue (xterm.js)     │  │  │  │
│  │  │  └────────────────────────────────────────┘  │  │  │
│  │  └──────────────────────────────────────────────┘  │  │
│  └────────────────────────────────────────────────────┘  │
└──────────────────────────────────────────────────────────┘
```

## Data Flow

### SSH Config Flow

```
User opens Config view
       │
       ▼
App.GetSSHConfig() ──► config.ParseSSHConfig(~/.ssh/config)
       │                      │
       │                      ▼
       │              Read file → Parse lines → Build SSHHostEntry[]
       │                      │
       │                      ▼
       │◄───────────── Return {rawContent, entries, filePath}
       │
       ▼
ConfigEditor.vue renders textarea + host cards

User clicks Save
       │
       ▼
App.SaveSSHConfig(content) ──► config.SaveSSHConfig(path, content)
                                      │
                                      ▼
                              Write file (0600 perms)
                                      │
                                      ▼
                              db.AddAuditLog("config_save")
```

### SSH Terminal Flow

```
User clicks Connect on a host
       │
       ▼
App.SSHConnect(host, hostname, port, user, identityFile)
       │
       ├──► ssh.Dial("tcp", addr, config)
       │         │
       │         ├──► loadPrivateKey(identityFile) or trySSHAgent()
       │         │
       │         └──► ssh.ClientConfig{User, Auth, HostKeyCallback}
       │
       ├──► conn.NewSession()
       │
       ├──► session.RequestPty("xterm-256color", 24, 80)
       │
       ├──► session.Shell()
       │
       ├──► Store session in SSHManager.sessions map
       │
       └──► Return sessionId

User types in terminal
       │
       ▼
xterm.js onData ──► App.WriteToTerminal(sessionId, data)
                          │
                          ▼
                   session.Stdin.Write(data)

Remote output
       │
       ▼
session.StdoutPipe ──► (streamed to frontend via Wails events)
```

## Component Details

### Backend Packages

#### `storage` (SQLite Layer)

```
database.go
├── Database struct (sql.DB wrapper)
├── initTables() — Creates 3 tables if not exist
│   ├── connection_history — SSH connection log
│   ├── favorite_hosts — User-saved hosts
│   └── audit_logs — Operation audit trail
├── ConnectionHistory CRUD
├── FavoriteHost CRUD
├── AuditLog (append-only, read)
└── GetGroups() — Distinct group names
```

**Tables**:
```sql
connection_history (id, host_alias, host_name, port, user, last_connected)
favorite_hosts (id, host_alias, host_name, port, user, tags, group_name)
audit_logs (id, action, detail, timestamp)
```

#### `config` (SSH Config Parser)

```
config.go
├── SSHHostEntry struct — Parsed host entry
├── SSHConfig struct — Full config with entries + raw content
├── GetDefaultConfigPath() — ~/.ssh/config
├── ParseSSHConfig(path) — Read and parse config file
├── SaveSSHConfig(path, content) — Write config (0600)
├── ValidateSSHConfig(content) — Check for unknown keywords
├── ExportConfig(path) — Read config as string
└── ImportConfig(path, content) — Write config from string
```

**Parsing logic**: Line-by-line scanner, state machine for Host blocks. Supports standard SSH keywords (Host, HostName, User, Port, IdentityFile, ForwardAgent, ProxyJump, etc.).

#### `keys` (SSH Key Management)

```
keys.go
├── SSHKey struct — Key metadata
├── KeyGenRequest struct — Key generation parameters
├── GetSSHDir() — ~/.ssh/
├── ListKeys() — Scan ~/.ssh/ for key pairs
├── GenerateKey(req) — Generate key pair
│   ├── RSA: crypto/rsa + x509
│   ├── ED25519: ssh.GenerateEd25519Key
│   └── ECDSA: crypto/ecdsa
├── DeleteKey(name) — Remove private + public key files
├── GetPubKeyContent(name) — Read .pub file
└── detectKeyType() / detectPrivateKeyType()
```

#### `hosts` (Host Management)

```
hosts.go
├── HostService struct — Wraps Database
├── ImportFromConfig() — Parse ~/.ssh/config, add to DB
├── GetFavorites() / Add / Update / Delete
├── GetGroups() — Distinct group names
├── GetHistory() / Search / Delete
├── RecordConnection() — Log a connection event
└── TestConnection(hostname, port) — TCP dial timeout test
```

#### `ssh` (SSH Connection Handler)

```
ssh.go
├── ConnectionInfo struct — Connection parameters
├── TerminalSession struct — Active SSH session
├── SSHManager struct — Session registry (map[string]*TerminalSession)
├── Connect(info) — Full SSH handshake + shell
│   ├── createSSHConfig() — Auth methods (key file + agent)
│   ├── ssh.Dial()
│   ├── session.RequestPty()
│   └── session.Shell()
├── WriteToTerminal(sessionId, data) — Send input
├── ResizeTerminal(sessionId, rows, cols) — Window change
├── Close(sessionId) — Cleanup session
├── CloseAll() — Cleanup all sessions
└── IsConnected(sessionId) — Check status
```

### Frontend Views

#### `App.vue`
- Root component with sidebar navigation
- Manages `currentView` state (config/keys/hosts/terminal)
- Conditional rendering of view components
- Global dark theme styles

#### `ConfigEditor.vue`
- Loads `~/.ssh/config` on mount
- Raw text editor (textarea with monospace font)
- Parse-and-preview: shows host cards from parsed config
- Validate button calls backend validator
- Save button writes to disk

#### `KeyManager.vue`
- Lists all keys in `~/.ssh/`
- Generate modal: type selector, size, email, passphrase, name
- View pubkey modal with copy button
- Delete with confirmation
- Color-coded key type badges

#### `HostList.vue`
- Two tabs: Favorites and Connection History
- Favorites: grid of host cards with connect/test/edit/delete
- History: table with search/filter
- Import from config button
- Add/edit host modal

#### `TerminalView.vue`
- xterm.js terminal with FitAddon
- Tab bar for multiple sessions
- Connection modal (alias, hostname, port, user, identity file)
- Status indicators (green dot = connected)
- Auto-resize on container change
- Cleanup on unmount

## Security Model

### Threat Model

| Threat | Mitigation |
|--------|-----------|
| Private key exposure | File permissions 0600, no key content stored in DB |
| Config tampering | Config written with 0600 perms |
| Unauthorized connections | No password storage, relies on SSH keys/agent |
| Audit trail tampering | Append-only audit log (no delete operation) |

### Known Limitations

1. **Host key verification**: Uses `ssh.InsecureIgnoreHostKey()` — vulnerable to MITM. Should implement proper `known_hosts` checking.
2. **No encryption at rest**: SQLite database is unencrypted. Sensitive data (if any) should be encrypted before storage.
3. **No passphrase caching**: Passphrases for encrypted keys are not cached — user must use SSH agent.

## Build Pipeline

```
wails dev
    │
    ├──► npm run dev (Vite dev server on :5173)
    │         │
    │         └──► Hot module replacement for Vue components
    │
    └──► Go backend (Wails dev server)
              │
              └──► Proxy to Vite dev server
              └──► Bind App struct methods to JS

wails build
    │
    ├──► npm run build (Vite → frontend/dist/)
    │
    ├──► Embed frontend/dist/ into Go binary
    │
    └──► Compile Go → build/bin/ssh-manager(.exe)
```

## Database Schema

```
┌────────────────────────────────────────┐
│         connection_history             │
├────────┬───────────┬───────────────────┤
│ id     │ INTEGER   │ PRIMARY KEY       │
│ host_alias │ TEXT  │ NOT NULL          │
│ host_name  │ TEXT  │ NOT NULL          │
│ port       │ INT   │ DEFAULT 22        │
│ user       │ TEXT  │ DEFAULT ''        │
│ last_connected │ DATETIME │ NOT NULL   │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│          favorite_hosts                │
├────────┬───────────┬───────────────────┤
│ id     │ INTEGER   │ PRIMARY KEY       │
│ host_alias │ TEXT  │ NOT NULL          │
│ host_name  │ TEXT  │ NOT NULL          │
│ port       │ INT   │ DEFAULT 22        │
│ user       │ TEXT  │ DEFAULT ''        │
│ tags       │ TEXT  │ DEFAULT ''        │
│ group_name │ TEXT  │ DEFAULT 'default' │
└────────────────────────────────────────┘

┌────────────────────────────────────────┐
│            audit_logs                  │
├────────┬───────────┬───────────────────┤
│ id     │ INTEGER   │ PRIMARY KEY       │
│ action │ TEXT      │ NOT NULL          │
│ detail │ TEXT      │ NOT NULL          │
│ timestamp │ DATETIME │ NOT NULL        │
└────────────────────────────────────────┘
```

## Extension Points

### Adding a New View
1. Create `frontend/src/views/NewView.vue`
2. Import and register in `App.vue`
3. Add nav item in sidebar
4. Add backend methods to `app.go` if needed

### Adding a New Backend Service
1. Create `backend/newservice/service.go`
2. Initialize in `App.startup()`
3. Add methods to `app.go` that delegate to the service
4. Methods will be auto-exposed as `window.go.main.App.*`

### Adding SSH Auth Method
Modify `createSSHConfig()` in `backend/ssh/ssh.go`:
- Add new auth method to `authMethods` slice
- Options: password, keyboard-interactive, certificate-based
