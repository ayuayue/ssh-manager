# SSH Manager

[English](README.en.md) | 中文

一款基于 [Wails](https://wails.app/) (Go + Vue 3) 构建的现代桌面 SSH 客户端和管理工具。在一个简洁美观的界面中管理 SSH 配置、密钥、主机，并支持多语言切换。

![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-blue)
![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)
![License](https://img.shields.io/badge/license-MIT-green)

## 功能特性

### SSH 配置管理
- 自动解析并展示 `~/.ssh/config`
- 可视化文本编辑器，支持直接编辑
- 实时语法校验并提示警告信息
- 手动保存，自动持久化到磁盘
- 解析后的主机预览卡片

### SSH 密钥管理
- 列出所有现有密钥（RSA、ED25519、ECDSA）
- 查看公钥/私钥内容，一键复制
- 生成新密钥，支持自定义选项：
  - 类型：RSA / ED25519 / ECDSA
  - 密钥长度 / 曲线选择
  - 邮箱备注
  - 可选密码短语
- 删除密钥需二次确认

### 主机管理
- **收藏主机**：保存和整理常用主机
- **分组与标签**：对主机进行分类，便于快速筛选
- **从配置导入**：自动解析 `~/.ssh/config` 中的主机
- **连接测试**：TCP 连通性测试
- **增删改查**：添加、编辑、删除收藏主机

### 连接历史
- 记录所有 SSH 连接（主机名、IP、端口、用户、时间）
- 按主机名、别名或用户搜索过滤
- 支持删除单条历史记录

### 安全与审计
- 操作审计日志（配置变更、密钥生成/删除、连接记录）
- 基于 SQLite 的本地数据持久化
- 私钥文件权限严格管控（0600）

## 截图

> _(在此添加截图)_

## 环境要求

- **Go** 1.21 或更高版本
- **Node.js** 18 或更高版本
- **Wails CLI**（`go install github.com/wailsapp/wails/v2/cmd/wails@latest`）
- **GCC** 或兼容的 C 编译器（`go-sqlite3` 需要 CGO）

## 安装

### 1. 克隆仓库

```bash
git clone https://github.com/your-username/ssh-manager.git
cd ssh-manager
```

### 2. 安装前端依赖

```bash
cd frontend
npm install
cd ..
```

### 3. 安装 Go 依赖

```bash
go mod tidy
```

### 4. 构建并运行

```bash
# 开发模式（热重载）
wails dev

# 生产构建
wails build
```

编译后的二进制文件位于 `build/bin/` 目录下。

## 项目结构

```
ssh-manager/
├── main.go                      # Wails 应用入口
├── app.go                       # 后端服务绑定（Wails API）
├── wails.json                   # Wails 配置文件
├── go.mod                       # Go 模块定义
├── .gitignore
├── README.md                    # 中文文档（本文件）
├── README.en.md                 # English Documentation
├── AGENTS.md                    # AI 辅助开发指南
├── AGENTS.en.md                 # AI Agent Guide (English)
├── ARCHITECTURE.md              # 架构设计文档
├── CONTRIBUTING.md              # 贡献指南
│
├── backend/
│   ├── config/
│   │   └── config.go            # SSH 配置文件解析、校验、保存
│   ├── hosts/
│   │   └── hosts.go             # 主机管理、导入、连接测试
│   ├── keys/
│   │   └── keys.go              # SSH 密钥生成、列表、删除
│   └── storage/
│       └── database.go          # SQLite 数据库层（历史、收藏、日志）
│
└── frontend/
    ├── index.html               # HTML 入口
    ├── vite.config.js           # Vite 配置
    ├── package.json             # 前端依赖
    └── src/
        ├── main.js              # Vue 应用引导
        ├── i18n.js              # 国际化（中文/英文）
        ├── App.vue              # 根组件（侧边栏布局 + 主题切换）
        └── views/
            ├── ConfigEditor.vue # SSH 配置编辑器
            ├── KeyManager.vue   # SSH 密钥管理
            └── HostList.vue     # 主机收藏与历史
```

## 使用说明

### SSH 配置编辑器

启动时自动加载 `~/.ssh/config`。你可以：

- 直接在编辑器中修改配置文本
- 点击 **校验** 检查语法错误
- 点击 **保存** 将更改写入磁盘
- 在右侧查看解析后的主机预览卡片

### 密钥管理器

- 点击 **生成密钥** 创建新的 SSH 密钥对
- 选择密钥类型（推荐 ED25519）、长度和可选密码短语
- 点击 **Pub** 查看公钥，点击 **Priv** 查看私钥
- 点击 **复制** 将内容复制到剪贴板
- 删除密钥需要确认操作

### 主机管理

- **从配置导入**：解析 `~/.ssh/config` 并将所有主机添加到收藏
- **添加主机**：手动添加主机，设置别名、主机名、端口、用户、标签和分组
- **测试**：执行 TCP 连接测试验证可达性
- **编辑/删除**：管理已收藏的主机

### 连接历史

- 自动记录每次连接信息
- 支持搜索过滤，快速定位目标主机
- 可删除不需要的历史记录

## 配置

### 默认路径

| 项目 | 路径 |
|------|------|
| SSH 配置 | `~/.ssh/config` |
| SSH 密钥 | `~/.ssh/` |
| 应用数据库 | `~/.ssh-manager/data.db` |

### 支持的密钥类型

| 类型 | 默认长度 | 说明 |
|------|---------|------|
| ED25519 | 256 位 | **推荐** — 快速、安全、体积小 |
| RSA | 4096 位 | 兼容性最广，密钥较大 |
| ECDSA | 256/384/521 位 | NIST 曲线，均衡之选 |

## 开发

```bash
# 开发模式（热重载）
wails dev

# 生产构建
wails build

# 带调试信息的构建
wails build -debug
```

## 技术栈

| 层级 | 技术 |
|------|------|
| 框架 | [Wails v2](https://wails.app/) |
| 后端 | Go 1.21+ |
| 前端 | Vue 3 + Vite |
| 数据库 | SQLite (go-sqlite3) |
| SSH | golang.org/x/crypto/ssh |
| 国际化 | 内置 i18n 模块（中文/英文） |

## 安全说明

- 私钥生成时自动设置 `0600` 文件权限
- 主机密钥验证默认使用 `InsecureIgnoreHostKey()` — 生产环境建议实现完整的 known_hosts 检查
- 不存储任何密码 — 认证依赖 SSH 密钥和 Agent
- 审计日志记录所有重要操作

## 开发计划

- [ ] 密码/密码短语认证支持
- [ ] SSH Agent 转发
- [ ] 端口转发 / 隧道管理
- [ ] SFTP 文件浏览器
- [ ] 配置差异对比查看器
- [ ] 主题自定义（更多配色方案）
- [ ] Known hosts 验证
- [ ] 连接模板 / 预设配置
- [ ] 从其他工具导入（PuTTY 等）

## 许可证

MIT

## 贡献

请参阅 [CONTRIBUTING.md](CONTRIBUTING.md) 了解贡献指南。

## 多语言

- 🇨🇳 [中文](README.md)（当前页面）
- 🇺🇸 [English](README.en.md)
