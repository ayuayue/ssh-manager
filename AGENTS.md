# AGENTS.md

> [English Version](AGENTS.en.md)

本文件为 AI 助手（Claude、GitHub Copilot 等）在此代码库中工作时提供指导。

## 项目概述

SSH Manager 是一款基于 [Wails](https://wails.app/) 构建的桌面应用程序 — 一个类似 Electron 但使用原生 WebView 而非 Chromium 的 Go 框架。前端使用 Vue 3，后端使用 Go。

## 架构

```
┌─────────────────────────────────────────┐
│              Wails 应用                  │
│  ┌──────────────┐    ┌───────────────┐  │
│  │   前端       │◄──►│    后端       │  │
│  │  (Vue 3)     │    │     (Go)      │  │
│  │              │    │               │  │
│  │ ConfigEditor │    │ config/       │  │
│  │ KeyManager   │    │ keys/         │  │
│  │ HostList     │    │ hosts/        │  │
│  │              │    │ storage/      │  │
│  └──────────────┘    └───────────────┘  │
│         │                    │          │
│         │   Wails 运行时      │          │
│         │   (JS ↔ Go 桥接)   │          │
│  ┌──────▼────────────────────▼───────┐  │
│  │         WebView (Edge/WebKit)     │  │
│  └───────────────────────────────────┘  │
└─────────────────────────────────────────┘
```

## 前端如何调用后端

Wails 自动将绑定结构体上的所有公开方法暴露给前端，通过 `window.go.main.App.*`。调用模式如下：

```javascript
// 前端 (Vue)
const result = await window.go.main.App.MethodName(arg1, arg2)
if (result.error) {
  // 处理错误
}
// 使用 result.data
```

```go
// 后端 (Go) — 在 app.go 中
func (a *App) MethodName(arg1 string, arg2 int) map[string]interface{} {
    // 执行业务逻辑
    return map[string]interface{}{"data": value}
    // 错误时：
    return map[string]interface{}{"error": "message"}
}
```

**约定**：所有后端方法返回 `map[string]interface{}`。错误返回 `{"error": "message"}`，成功返回 `{"success": true, ...}`。

## 后端包职责

| 包 | 文件 | 职责 |
|---------|------|---------------|
| `main` | `main.go` | Wails 应用引导、窗口配置 |
| `main` | `app.go` | 所有 Wails 暴露方法、服务编排 |
| `config` | `backend/config/config.go` | 解析/写入/校验 `~/.ssh/config` |
| `keys` | `backend/keys/keys.go` | 列出/生成/删除 SSH 密钥 |
| `hosts` | `backend/hosts/hosts.go` | 收藏主机 CRUD、连接历史、TCP 测试 |
| `storage` | `backend/storage/database.go` | SQLite 操作（历史、收藏、审计日志） |

## 关键约定

1. **错误处理**：后端方法返回 `map[string]interface{}`，失败时包含 `"error"` 键。前端检查 `result.error`。
2. **不 panic**：所有错误应优雅返回，后端代码中永远不要 panic。
3. **审计日志**：重要操作（配置保存、密钥生成/删除、SSH 连接）应调用 `a.db.AddAuditLog()`。
4. **文件权限**：SSH 私钥必须 `0600`，公钥 `0644`，配置文件 `0600`。
5. **数据库**：SQLite 在启动时自动初始化表。无需迁移（使用 `CREATE TABLE IF NOT EXISTS`）。

## 前端约定

1. **组件结构**：每个视图是单个 `.vue` 文件，包含 `<template>`、`<script setup>` 和 `<style scoped>`。
2. **状态管理**：使用本地 `ref()`/`computed()` — 当前范围不需要 Pinia/Vuex。
3. **国际化**：通过 `inject('t')` 获取翻译函数，所有文本使用 `t('key.path')`。
4. **样式**：使用 CSS 变量驱动主题（亮色/暗色），每个组件内联 `<style scoped>`。
5. **模态框**：基于遮罩层的模态框，使用 `@click.self` 关闭。

## 添加新功能

### 后端

1. 将逻辑添加到相应的 `backend/` 包
2. 在 `app.go` 中添加一个方法包装逻辑并返回 `map[string]interface{}`
3. 该方法将自动作为 `window.go.main.App.MethodName()` 在前端可用

### 前端

1. 在相应的 `views/*.vue` 组件中添加 UI
2. 调用 `await window.go.main.App.MethodName(...)`
3. 处理 `result.error` 并使用数据

## 环境

- **操作系统**: Windows
- **Wails CLI**: `C:\Users\14012\go\bin\wails.exe`
- **用户目录**: `C:\Users\14012`
- **CGO**: 必须启用（`CGO_ENABLED=1`），需要 GCC 编译器

## 常用命令

```bash
# 开发模式（热重载）
wails dev

# 生产构建
wails build

# 安装前端依赖
cd frontend && npm install

# 安装 Go 依赖
go mod tidy
```

## 依赖

### Go
- `github.com/wailsapp/wails/v2` — 桌面应用框架
- `github.com/mattn/go-sqlite3` — SQLite 驱动（需要 CGO）
- `golang.org/x/crypto/ssh` — SSH 协议实现

### 前端
- `vue@^3.4.0` — UI 框架
- `vite@^5.0.0` — 构建工具

## 注意事项

1. **需要 CGO**：`go-sqlite3` 需要 CGO。确保 `CGO_ENABLED=1` 且有 C 编译器可用。
2. **SSH_AUTH_SOCK**：SSH Agent 支持依赖 `SSH_AUTH_SOCK` 环境变量。
3. **主机密钥验证**：当前使用 `ssh.InsecureIgnoreHostKey()`。作为本地工具可以接受，但生产环境应改进。
4. **Windows 路径**：在 Windows 上，`~` 展开为 `%USERPROFILE%`。`os.UserHomeDir()` 函数正确处理此问题。
5. **Wails 开发服务器**：在开发模式下，前端运行在独立端口（默认 5173）。Wails 自动处理代理。
6. **私钥查看**：私钥内容可通过后端 `GetPrivKeyContent()` 读取，前端显示时带红色警告边框。

## 测试

当前没有自动化测试。添加测试时：

- 后端：使用 Go 的 `testing` 包和表驱动测试
- 前端：考虑使用 Vitest + Vue Test Utils
- 集成测试：Wails 没有内置测试框架 — 独立测试后端逻辑

## 多语言

项目使用内置的 `i18n.js` 模块实现国际化，当前支持：

- 🇨🇳 中文（默认）
- 🇺🇸 English

翻译文件位于 `frontend/src/i18n.js`，通过 `inject('t')` 在组件中使用。
