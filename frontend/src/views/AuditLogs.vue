<template>
  <div class="audit-logs">
    <div class="page-header">
      <h1>{{ t('logs.title') }}</h1>
      <div class="header-actions">
        <select v-model="filterAction" @change="loadLogs" class="filter-select">
          <option value="">{{ t('logs.allActions') }}</option>
          <option v-for="a in actionTypes" :key="a" :value="a">{{ a }}</option>
        </select>
        <button @click="loadLogs" class="btn btn-secondary" :disabled="loading">
          {{ t('logs.reload') }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <span>{{ t('logs.loading') }}</span>
    </div>

    <div v-else-if="logs.length === 0" class="empty-state">
      <div class="empty-icon">📋</div>
      <p>{{ t('logs.empty') }}</p>
    </div>

    <div v-else class="log-list">
      <div v-for="log in filteredLogs" :key="log.id" class="log-item" :class="'action-' + log.action">
        <div class="log-icon-wrapper">
          <span class="log-icon">{{ getActionIcon(log.action) }}</span>
        </div>
        <div class="log-content">
          <div class="log-header">
            <span class="log-action">{{ getActionLabel(log.action) }}</span>
            <span class="log-time">{{ formatTime(log.timestamp) }}</span>
          </div>
          <div class="log-detail">{{ log.detail }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, inject } from 'vue'

export default {
  name: 'AuditLogs',
  setup() {
    const t = inject('t')
    const logs = ref([])
    const loading = ref(true)
    const filterAction = ref('')

    const actionTypes = computed(() => {
      const actions = new Set(logs.value.map(l => l.action))
      return [...actions]
    })

    const filteredLogs = computed(() => {
      if (!filterAction.value) return logs.value
      return logs.value.filter(l => l.action === filterAction.value)
    })

    const getActionIcon = (action) => {
      const icons = {
        config_save: '💾',
        config_import: '📥',
        key_generate: '🔑',
        key_delete: '🗑',
        ssh_connect: '🔌',
      }
      return icons[action] || '📝'
    }

    const getActionLabel = (action) => {
      const labels = {
        config_save: t('logs.configSave'),
        config_import: t('logs.configImport'),
        key_generate: t('logs.keyGenerate'),
        key_delete: t('logs.keyDelete'),
        ssh_connect: t('logs.sshConnect'),
      }
      return labels[action] || action
    }

    const loadLogs = async () => {
      loading.value = true
      try {
        const result = await window.go.main.App.GetAuditLogs()
        if (!result.error) logs.value = result.logs || []
      } catch (e) { console.error(e) }
      loading.value = false
    }

    const formatTime = (d) => {
      if (!d) return 'N/A'
      const date = new Date(d)
      const now = new Date()
      const diff = now - date
      if (diff < 60000) return t('logs.justNow')
      if (diff < 3600000) return Math.floor(diff / 60000) + t('logs.minutesAgo')
      if (diff < 86400000) return Math.floor(diff / 3600000) + t('logs.hoursAgo')
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }

    onMounted(loadLogs)

    return { t, logs, loading, filterAction, actionTypes, filteredLogs, getActionIcon, getActionLabel, loadLogs, formatTime }
  }
}
</script>

<style scoped>
.audit-logs {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  flex-wrap: wrap;
  gap: 12px;
}

.page-header h1 {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
}

.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.filter-select {
  padding: 8px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  cursor: pointer;
}

.log-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.log-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-md);
  transition: all var(--transition);
}

.log-item:hover {
  border-color: var(--border-color);
  box-shadow: var(--shadow-sm);
}

.log-icon-wrapper {
  width: 36px;
  height: 36px;
  background: var(--bg-tertiary);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  font-size: 16px;
}

.log-content {
  flex: 1;
  min-width: 0;
}

.log-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  flex-wrap: wrap;
}

.log-action {
  font-weight: 600;
  font-size: 13px;
  color: var(--text-primary);
}

.log-time {
  font-size: 12px;
  color: var(--text-muted);
  white-space: nowrap;
}

.log-detail {
  font-size: 13px;
  color: var(--text-secondary);
  margin-top: 4px;
  font-family: monospace;
  word-break: break-all;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 60px;
  color: var(--text-secondary);
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid var(--border-color);
  border-top-color: var(--accent);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-secondary);
}

.empty-icon { font-size: 48px; margin-bottom: 16px; opacity: 0.4; }
.empty-state p { margin-bottom: 20px; font-size: 15px; }

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all var(--transition);
  background: var(--bg-primary);
  color: var(--text-primary);
}

.btn:disabled { opacity: 0.5; cursor: not-allowed; }
.btn-secondary:hover:not(:disabled) { background: var(--bg-tertiary); }

@media (max-width: 600px) {
  .log-header { flex-direction: column; align-items: flex-start; gap: 4px; }
}
</style>
