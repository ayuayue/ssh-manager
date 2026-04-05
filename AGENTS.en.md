# AGENTS.md

This file provides guidance for AI agents (Claude, GitHub Copilot, etc.) working on this codebase.

## Project Overview

SSH Manager is a desktop application built with [Wails](https://wails.app/) — a Go-based framework similar to Electron but using native WebView instead of Chromium. The frontend is Vue 3, the backend is Go.

## Architecture

```
┌─────────────────────────────────────────┐
│              Wails App                   │
│  ┌──────────────┐    ┌───────────────┐  │
│  │   Frontend   │◄──►│    Backend    │  │
│  │  (Vue 3)     │    │     (Go)      │  │
│  │              │    │               │  │
│  │ ConfigEditor │    │ config/       │  │
│  │ KeyManager   │    │ keys/         │  │
│  │ HostList     │    │ hosts/        │  │
│  │ TerminalView │    │ ssh/          │  │
│  │              │    │ storage/      │  │
│  └──────────────┘    └───────────────┘  │
│         │                    │          │
│         │   Wails Runtime    │          │
│         │   (JS ↔ Go bridge) │          │
│  ┌──────▼────────────────────▼───────┐  │
│  │         WebView (Edge/WebKit)     │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## How the Frontend Calls the Backend

Wails automatically exposes all public methods on bound structs to the frontend via `window.go.main.App.*`. The pattern is:

```javascript
// Frontend (Vue)
const result = await window.go.main.App.MethodName(arg1, arg2)
if (result.error) {
  // handle error
}
// use result.data
```

```go
// Backend (Go) — in app.go
func (a *App) MethodName(arg1 string, arg2 int) map[string]interface{} {
    // do work
    return map[string]interface{}{"data": value}
    // or on error:
    return map[string]interface{}{"error": "message"}
}
```

**Convention**: All backend methods return `map[string]interface{}`. Errors are returned as `{"error": "message"}`, success as `{"success": true, ...}`.

## Backend Package Responsibilities

| Package | File | Responsibility |
|---------|------|---------------|
| `main` | `main.go` | Wails app bootstrap, window config |
| `main` | `app.go` | All Wails-exposed methods, service orchestration |
| `config` | `backend/config/config.go` | Parse/write/validate `~/.ssh/config` |
| `keys` | `backend/keys/keys.go` | List/generate/delete SSH keys |
| `hosts` | `backend/hosts/hosts.go` | Favorite hosts CRUD, connection history, TCP test |
| `ssh` | `backend/ssh/ssh.go` | SSH connections, terminal sessions, PTY |
| `storage` | `backend/storage/database.go` | SQLite operations (history, favorites, audit logs) |

## Key Conventions

1. **Error handling**: Backend methods return `map[string]interface{}` with an `"error"` key on failure. Frontend checks `result.error`.
2. **No panics**: All errors should be returned gracefully, never panic in backend code.
3. **Audit logging**: Significant operations (config save, key gen/delete, SSH connect) should call `a.db.AddAuditLog()`.
4. **File permissions**: SSH private keys must be `0600`, public keys `0644`, config `0600`.
5. **Database**: SQLite auto-initializes tables on startup. No migrations needed (uses `CREATE TABLE IF NOT EXISTS`).

## Frontend Conventions

1. **Component structure**: Each view is a single `.vue` file with `<template>`, `<script setup>`, and `<style scoped>`.
2. **State management**: Local `ref()`/`computed()` — no Pinia/Vuex needed for current scope.
3. **Styling**: Dark theme (GitHub-dark palette). Inline `<style scoped>` per component.
4. **Modals**: Overlay-based modals with `@click.self` to dismiss.

## Adding a New Feature

### Backend

1. Add logic to the appropriate `backend/` package
2. Add a method to `app.go` that wraps the logic and returns `map[string]interface{}`
3. The method will be automatically available as `window.go.main.App.MethodName()`

### Frontend

1. Add UI to the relevant `views/*.vue` component
2. Call `await window.go.main.App.MethodName(...)`
3. Handle `result.error` and use the data

## Environment

- **OS**: Windows
- **Wails CLI**: `C:\Users\14012\go\bin\wails.exe`
- **User home**: `C:\Users\14012`

## Common Commands

```bash
# Development (hot reload)
wails dev

# Production build
wails build

# Install frontend deps
cd frontend && npm install

# Install Go deps
go mod tidy
```

## Dependencies

### Go
- `github.com/wailsapp/wails/v2` — Desktop app framework
- `github.com/mattn/go-sqlite3` — SQLite driver (requires CGO)
- `golang.org/x/crypto/ssh` — SSH protocol implementation

### Frontend
- `vue@^3.4.0` — UI framework
- `xterm@^5.3.0` — Terminal emulator
- `xterm-addon-fit@^0.8.0` — Terminal resize addon
- `vite@^5.0.0` — Build tool

## Gotchas

1. **CGO required**: `go-sqlite3` needs CGO. Ensure `CGO_ENABLED=1` and a C compiler is available.
2. **SSH_AUTH_SOCK**: SSH agent support depends on the `SSH_AUTH_SOCK` environment variable.
3. **Terminal resize**: The terminal resize uses `term._core._renderService._dimensions` which is internal xterm API. This may break on xterm major version updates.
4. **Host key verification**: Currently uses `ssh.InsecureIgnoreHostKey()`. This is fine for a local tool but should be improved for production.
5. **Windows paths**: On Windows, `~` expands to `%USERPROFILE%`. The `os.UserHomeDir()` function handles this correctly.
6. **Wails dev server**: In dev mode, the frontend runs on a separate port (default 5173). Wails handles the proxy automatically.

## Testing

Currently no automated tests. When adding tests:

- Backend: Use Go's `testing` package with table-driven tests
- Frontend: Consider Vitest + Vue Test Utils
- Integration: Wails has no built-in test framework — test backend logic independently
