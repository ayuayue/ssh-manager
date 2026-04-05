# CONTRIBUTING.md

# Contributing to SSH Manager

Thank you for your interest in contributing! This guide will help you get started.

## Getting Started

### Prerequisites

- Go 1.21+
- Node.js 18+
- Wails CLI (`go install github.com/wailsapp/wails/v2/cmd/wails@latest`)
- GCC or compatible C compiler (required for `go-sqlite3`)

### Setup

```bash
git clone https://github.com/your-username/ssh-manager.git
cd ssh-manager
cd frontend && npm install && cd ..
go mod tidy
```

### Development

```bash
# Run in development mode with hot reload
wails dev

# Build for production
wails build
```

## Project Structure

See [ARCHITECTURE.md](ARCHITECTURE.md) for detailed architecture documentation.

Quick reference:
- `app.go` — All Wails-exposed backend methods
- `backend/` — Go packages (config, keys, hosts, ssh, storage)
- `frontend/src/views/` — Vue 3 page components

## Code Conventions

### Go Backend

1. **Return maps, not errors**: All public methods in `app.go` return `map[string]interface{}`. Use `{"error": "message"}` for errors and `{"success": true}` for success.
2. **No panics**: Always return errors gracefully.
3. **Audit log**: Log significant operations with `a.db.AddAuditLog(action, detail)`.
4. **File permissions**: SSH files must follow standard permissions (private keys: 0600, public keys: 0644, config: 0600).
5. **Naming**: Exported methods on `App` struct are PascalCase. They become `window.go.main.App.MethodName()` in JS.

### Vue Frontend

1. **Composition API**: Use `<script setup>` or `setup()` function with `ref()`/`computed()`.
2. **Scoped styles**: Always use `<style scoped>` in view components.
3. **Error handling**: Always check `result.error` after calling backend methods.
4. **No external state manager**: Keep state local with `ref()` for now. No Pinia/Vuex.
5. **Dark theme**: Use the GitHub-dark color palette consistently.

### Color Palette

```
Background:  #0d1117
Surface:     #161b22, #1e2a3a
Border:      #2d3a4a, #30363d
Text:        #c9d1d9, #e0e0e0
Muted:       #8b949e, #6e7681
Primary:     #58a6ff (blue)
Success:     #3fb950, #238636 (green)
Warning:     #d29922, #9e6a03 (yellow)
Danger:      #f85149, #da3633 (red)
```

## Making Changes

### 1. Pick an Issue

Check the [Issues](https://github.com/your-username/ssh-manager/issues) page or create a new one describing your change.

### 2. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/bug-description
```

### 3. Make Your Changes

Follow the conventions above. Keep commits focused and atomic.

### 4. Test

```bash
# Ensure Go code compiles
go build ./...

# Ensure frontend builds
cd frontend && npm run build

# Run in dev mode and manually test
wails dev
```

### 5. Submit a Pull Request

- Describe what you changed and why
- Include screenshots for UI changes
- Reference any related issues

## Adding a New Backend Method

```go
// In app.go
func (a *App) MyNewMethod(param1 string, param2 int) map[string]interface{} {
    // 1. Call the appropriate backend service
    result, err := someService.DoSomething(param1, param2)
    if err != nil {
        return map[string]interface{}{"error": err.Error()}
    }

    // 2. Log the operation (if significant)
    a.db.AddAuditLog("my_action", fmt.Sprintf("Did something with %s", param1))

    // 3. Return success
    return map[string]interface{}{"success": true, "data": result}
}
```

Frontend usage:

```javascript
const result = await window.go.main.App.MyNewMethod('hello', 42)
if (result.error) {
  alert('Error: ' + result.error)
  return
}
// Use result.data
```

## Adding a New Frontend View

1. Create `frontend/src/views/MyView.vue`
2. Add to `App.vue`:

```vue
import MyView from './views/MyView.vue'

// In components:
MyView,

// In template:
<MyView v-else-if="currentView === 'myview'" />

// In sidebar:
<li :class="{ active: currentView === 'myview' }" @click="currentView = 'myview'">
  My View
</li>
```

## Commit Message Format

Use conventional commits:

```
feat: add port forwarding support
fix: resolve terminal resize on Windows
docs: update architecture diagram
refactor: extract SSH config parser into separate function
chore: update dependencies
```

## Code Review Checklist

- [ ] Backend methods return `map[string]interface{}` with proper error handling
- [ ] No panics or unhandled errors in Go code
- [ ] File permissions are correct for SSH files
- [ ] Frontend checks `result.error` after backend calls
- [ ] UI is consistent with the dark theme
- [ ] No hardcoded secrets or credentials
- [ ] Audit logging added for significant operations

## Reporting Bugs

Include:
- OS and version
- Steps to reproduce
- Expected vs actual behavior
- Screenshots if applicable
- Any error messages from the console

## Feature Requests

Describe:
- What problem you're trying to solve
- How you envision the feature working
- Any alternatives you've considered

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
