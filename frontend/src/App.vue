<template>
  <div id="app" :class="theme">
    <nav class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <span class="logo-icon">🔐</span>
        <span class="logo-text" v-show="!sidebarCollapsed">SSH Manager</span>
      </div>
      <ul class="nav-list">
        <li
          v-for="item in navItems"
          :key="item.id"
          :class="{ active: currentView === item.id }"
          @click="currentView = item.id"
          class="nav-item"
        >
          <span class="nav-icon" v-html="item.icon"></span>
          <span class="nav-label" v-show="!sidebarCollapsed">{{ item.label }}</span>
        </li>
      </ul>
      <div class="sidebar-footer">
        <select class="lang-select" v-model="lang" @change="changeLang">
          <option value="zh">中文</option>
          <option value="en">English</option>
        </select>
        <button class="theme-toggle" @click="toggleTheme" :title="t('theme.' + (theme === 'light' ? 'dark' : 'light'))">
          <span v-html="theme === 'light' ? '🌙' : '☀️'"></span>
        </button>
        <button class="sidebar-toggle" @click="sidebarCollapsed = !sidebarCollapsed" :title="t('sidebar.' + (sidebarCollapsed ? 'expand' : 'collapse'))">
          <span v-html="sidebarCollapsed ? '▶' : '◀'"></span>
        </button>
      </div>
    </nav>
    <main class="content">
      <ConfigEditor v-if="currentView === 'config'" />
      <KeyManager v-else-if="currentView === 'keys'" />
      <HostList v-else-if="currentView === 'hosts'" />
      <AuditLogs v-else-if="currentView === 'logs'" />
    </main>
  </div>
</template>

<script>
import { ref, provide, onMounted } from 'vue'
import ConfigEditor from './views/ConfigEditor.vue'
import KeyManager from './views/KeyManager.vue'
import HostList from './views/HostList.vue'
import AuditLogs from './views/AuditLogs.vue'
import { t, getLang, setLang } from './i18n'

export default {
  name: 'App',
  components: { ConfigEditor, KeyManager, HostList, AuditLogs },
  setup() {
    const currentView = ref('config')
    const theme = ref('light')
    const sidebarCollapsed = ref(false)
    const lang = ref(getLang())

    const updateNav = () => {
      navItems.value = [
        { id: 'config', icon: '⚙️', label: t('nav.config') },
        { id: 'keys', icon: '🔑', label: t('nav.keys') },
        { id: 'hosts', icon: '🌐', label: t('nav.hosts') },
        { id: 'logs', icon: '📋', label: t('nav.logs') },
      ]
    }

    const navItems = ref([])
    updateNav()

    const toggleTheme = () => {
      theme.value = theme.value === 'light' ? 'dark' : 'light'
      localStorage.setItem('ssh-manager-theme', theme.value)
    }

    const changeLang = () => {
      setLang(lang.value)
      updateNav()
    }

    onMounted(() => {
      const saved = localStorage.getItem('ssh-manager-theme')
      if (saved) theme.value = saved
    })

    provide('theme', theme)
    provide('t', t)

    return { currentView, theme, sidebarCollapsed, navItems, lang, toggleTheme, changeLang, t }
  }
}
</script>

<style>
html, body {
  height: 100%;
  margin: 0;
  padding: 0;
  overflow: hidden;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

:root {
  --bg-primary: #ffffff;
  --bg-secondary: #f6f8fa;
  --bg-tertiary: #e8ecf0;
  --bg-code: #f0f2f5;
  --border-color: #d0d7de;
  --border-light: #e8ecf0;
  --text-primary: #1f2328;
  --text-secondary: #656d76;
  --text-muted: #8b949e;
  --accent: #0969da;
  --accent-hover: #0550ae;
  --accent-light: rgba(9, 105, 218, 0.08);
  --success: #1a7f37;
  --success-bg: rgba(26, 127, 55, 0.08);
  --warning: #9a6700;
  --warning-bg: rgba(154, 103, 0, 0.08);
  --danger: #cf222e;
  --danger-bg: rgba(207, 34, 46, 0.08);
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.06);
  --shadow-md: 0 2px 8px rgba(0,0,0,0.08);
  --shadow-lg: 0 4px 16px rgba(0,0,0,0.1);
  --radius-sm: 6px;
  --radius-md: 8px;
  --radius-lg: 12px;
  --sidebar-width: 220px;
  --sidebar-collapsed-width: 56px;
  --transition: 0.2s ease;
}

.dark {
  --bg-primary: #0d1117;
  --bg-secondary: #161b22;
  --bg-tertiary: #21262d;
  --bg-code: #161b22;
  --border-color: #30363d;
  --border-light: #21262d;
  --text-primary: #e6edf3;
  --text-secondary: #8b949e;
  --text-muted: #6e7681;
  --accent: #58a6ff;
  --accent-hover: #79c0ff;
  --accent-light: rgba(88, 166, 255, 0.1);
  --success: #3fb950;
  --success-bg: rgba(63, 185, 80, 0.1);
  --warning: #d29922;
  --warning-bg: rgba(210, 153, 34, 0.1);
  --danger: #f85149;
  --danger-bg: rgba(248, 81, 73, 0.1);
  --shadow-sm: 0 1px 2px rgba(0,0,0,0.2);
  --shadow-md: 0 2px 8px rgba(0,0,0,0.3);
  --shadow-lg: 0 4px 16px rgba(0,0,0,0.4);
}

#app {
  display: flex;
  height: 100%;
  width: 100%;
  background: var(--bg-primary);
  transition: background var(--transition);
}

.sidebar {
  width: var(--sidebar-width);
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  flex-shrink: 0;
  transition: width var(--transition);
  overflow: hidden;
}

.sidebar.collapsed {
  width: var(--sidebar-collapsed-width);
}

.sidebar-header {
  padding: 16px;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 10px;
  min-height: 56px;
}

.logo-icon { font-size: 20px; flex-shrink: 0; }
.logo-text {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-primary);
  white-space: nowrap;
  transition: opacity var(--transition);
}

.nav-list {
  list-style: none;
  padding: 8px;
  flex: 1;
}

.nav-item {
  padding: 10px 12px;
  cursor: pointer;
  transition: all var(--transition);
  font-size: 14px;
  color: var(--text-secondary);
  display: flex;
  align-items: center;
  gap: 10px;
  border-radius: var(--radius-sm);
  margin-bottom: 2px;
  white-space: nowrap;
}

.nav-item:hover {
  background: var(--accent-light);
  color: var(--text-primary);
}

.nav-item.active {
  background: var(--accent-light);
  color: var(--accent);
  font-weight: 500;
}

.nav-icon {
  font-size: 16px;
  width: 20px;
  text-align: center;
  flex-shrink: 0;
}

.nav-label { transition: opacity var(--transition); }

.sidebar-footer {
  padding: 8px;
  border-top: 1px solid var(--border-color);
  display: flex;
  gap: 4px;
}

.lang-select {
  flex: 1;
  padding: 6px 4px;
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 12px;
  cursor: pointer;
  outline: none;
  transition: all var(--transition);
}

.lang-select:hover {
  border-color: var(--accent);
}

.theme-toggle,
.sidebar-toggle {
  flex: 1;
  padding: 8px;
  background: transparent;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 14px;
  transition: all var(--transition);
}

.theme-toggle:hover,
.sidebar-toggle:hover {
  background: var(--accent-light);
  color: var(--accent);
  border-color: var(--accent);
}

.content {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 24px;
  transition: background var(--transition);
}

.content::-webkit-scrollbar { width: 8px; }
.content::-webkit-scrollbar-track { background: transparent; }
.content::-webkit-scrollbar-thumb { background: var(--border-color); border-radius: 4px; }
.content::-webkit-scrollbar-thumb:hover { background: var(--text-muted); }

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0; top: 0; bottom: 0;
    z-index: 100;
    box-shadow: var(--shadow-lg);
  }
  .sidebar.collapsed { width: 0; border: none; }
  .content { padding: 16px; }
}

@media (max-width: 480px) {
  .sidebar { width: 100%; }
  .sidebar.collapsed { width: 0; }
}
</style>
