<template>
  <div class="host-list">
    <div class="page-header">
      <h1>{{ t('hosts.title') }}</h1>
      <div class="header-actions">
        <button @click="importFromConfig" class="btn btn-secondary" :disabled="importing">
          {{ importing ? t('hosts.importing') : t('hosts.import') }}
        </button>
        <button @click="openAddModal" class="btn btn-primary">
          + {{ t('hosts.addHost') }}
        </button>
      </div>
    </div>

    <div class="tabs">
      <button
        v-for="tab in tabs"
        :key="tab.id"
        :class="{ active: activeTab === tab.id }"
        @click="activeTab = tab.id"
        class="tab-btn"
      >
        {{ tab.label }}
        <span v-if="tab.count !== null" class="tab-count">{{ tab.count }}</span>
      </button>
    </div>

    <div class="search-bar" v-if="activeTab === 'history'">
      <input
        v-model="searchQuery"
        @input="searchHistory"
        :placeholder="t('hosts.searchPlaceholder')"
        class="search-input"
      />
      <span class="search-icon">🔍</span>
    </div>

    <div class="content-area">
      <div v-if="activeTab === 'favorites'">
        <div v-if="hosts.length === 0" class="empty-state">
          <div class="empty-icon">🌐</div>
          <p>{{ t('hosts.noFavorites') }}</p>
          <div class="empty-actions">
            <button @click="importFromConfig" class="btn btn-secondary">{{ t('hosts.importFromConfig') }}</button>
            <button @click="openAddModal" class="btn btn-primary">{{ t('hosts.addManually') }}</button>
          </div>
        </div>
        <div v-else class="host-grid">
          <div v-for="host in hosts" :key="host.id" class="host-card">
            <div class="host-card-top">
              <div class="host-avatar">{{ host.hostAlias.charAt(0).toUpperCase() }}</div>
              <div class="host-card-info">
                <div class="host-alias">{{ host.hostAlias }}</div>
                <div class="host-addr">{{ host.hostName || 'N/A' }}:{{ host.port }}</div>
                <div class="host-user" v-if="host.user">{{ host.user }}</div>
              </div>
            </div>
            <div class="host-tags" v-if="host.tags || host.group">
              <span v-if="host.group" class="host-tag">{{ host.group }}</span>
              <span v-for="tag2 in (host.tags || '').split(',').filter(Boolean)" :key="tag2" class="host-tag">{{ tag2.trim() }}</span>
            </div>
            <div class="host-card-actions">
              <button @click="testHost(host)" class="btn btn-sm" :title="t('hosts.test')">{{ t('hosts.test') }}</button>
              <button @click="editHost(host)" class="btn btn-sm" :title="t('hosts.edit')">{{ t('hosts.edit') }}</button>
              <button @click="deleteHost(host.id)" class="btn btn-sm btn-danger" :title="t('hosts.del')">{{ t('hosts.del') }}</button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="activeTab === 'history'">
        <div v-if="history.length === 0" class="empty-state">
          <div class="empty-icon">🕐</div>
          <p>{{ t('hosts.noHistory') }}</p>
        </div>
        <div v-else class="history-list">
          <div v-for="h in history" :key="h.id" class="history-item">
            <div class="history-info">
              <div class="history-host">{{ h.hostAlias }}</div>
              <div class="history-detail">{{ h.hostName }}:{{ h.port }} | {{ h.user }}</div>
            </div>
            <div class="history-right">
              <span class="history-time">{{ formatDate(h.lastConnected) }}</span>
              <button @click="deleteHistory(h.id)" class="btn btn-sm btn-danger" title="Delete">&times;</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingId ? t('hosts.editHost') : t('hosts.addHost') }}</h2>
          <button class="modal-close" @click="closeModal">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>{{ t('hosts.alias') }} *</label>
            <input v-model="form.alias" :placeholder="t('hosts.alias')" />
          </div>
          <div class="form-group">
            <label>{{ t('hosts.hostname') }} *</label>
            <input v-model="form.hostname" :placeholder="t('hosts.hostname')" />
          </div>
          <div class="form-row">
            <div class="form-group">
              <label>{{ t('hosts.port') }}</label>
              <input v-model.number="form.port" type="number" placeholder="22" />
            </div>
            <div class="form-group">
              <label>{{ t('hosts.user') }}</label>
              <input v-model="form.user" :placeholder="t('hosts.user')" />
            </div>
          </div>
          <div class="form-group">
            <label>{{ t('hosts.tags') }} <span class="optional">({{ t('hosts.tagsHint') }})</span></label>
            <input v-model="form.tags" :placeholder="t('hosts.tags')" />
          </div>
          <div class="form-group">
            <label>{{ t('hosts.group') }}</label>
            <input v-model="form.group" :placeholder="t('hosts.group')" />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="closeModal" class="btn btn-secondary">{{ t('keys.cancel') }}</button>
          <button @click="saveHost" class="btn btn-primary">{{ editingId ? t('hosts.edit') : t('hosts.addHost') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, inject } from 'vue'

export default {
  name: 'HostList',
  setup() {
    const t = inject('t')
    const activeTab = ref('favorites')
    const hosts = ref([])
    const history = ref([])
    const searchQuery = ref('')
    const showModal = ref(false)
    const editingId = ref(null)
    const importing = ref(false)
    const form = ref({ alias: '', hostname: '', port: 22, user: '', tags: '', group: 'default' })

    const tabs = computed(() => [
      { id: 'favorites', label: t('hosts.favorites'), count: hosts.value.length },
      { id: 'history', label: t('hosts.history'), count: history.value.length },
    ])

    const loadHosts = async () => {
      const r = await window.go.main.App.GetFavoriteHosts()
      if (!r.error) hosts.value = r.hosts || []
    }

    const loadHistory = async () => {
      const r = await window.go.main.App.GetConnectionHistory()
      if (!r.error) history.value = r.history || []
    }

    const searchHistory = async () => {
      if (!searchQuery.value.trim()) { await loadHistory(); return }
      const r = await window.go.main.App.SearchConnectionHistory(searchQuery.value)
      if (!r.error) history.value = r.history || []
    }

    const importFromConfig = async () => {
      importing.value = true
      try {
        const r = await window.go.main.App.ImportHostsFromConfig()
        if (r.error) { alert('Error: ' + r.error); return }
        await loadHosts()
        alert(`Imported ${r.count} hosts.`)
      } catch (e) { alert('Failed: ' + e.message) }
      importing.value = false
    }

    const openAddModal = () => {
      editingId.value = null
      form.value = { alias: '', hostname: '', port: 22, user: '', tags: '', group: 'default' }
      showModal.value = true
    }

    const editHost = (host) => {
      editingId.value = host.id
      form.value = { alias: host.hostAlias, hostname: host.hostName, port: host.port, user: host.user, tags: host.tags, group: host.group }
      showModal.value = true
    }

    const saveHost = async () => {
      if (!form.value.alias || !form.value.hostname) { alert(t('hosts.required')); return }
      let r
      if (editingId.value) {
        r = await window.go.main.App.UpdateFavoriteHost(editingId.value, form.value.alias, form.value.hostname, form.value.port, form.value.user, form.value.tags, form.value.group)
      } else {
        r = await window.go.main.App.AddFavoriteHost(form.value.alias, form.value.hostname, form.value.port, form.value.user, form.value.tags, form.value.group)
      }
      if (r.error) { alert('Error: ' + r.error); return }
      closeModal()
      await loadHosts()
    }

    const deleteHost = async (id) => {
      if (!confirm(t('hosts.deleteConfirm'))) return
      const r = await window.go.main.App.DeleteFavoriteHost(id)
      if (!r.error) await loadHosts()
    }

    const deleteHistory = async (id) => {
      const r = await window.go.main.App.DeleteConnectionHistory(id)
      if (!r.error) await loadHistory()
    }

    const testHost = async (host) => {
      const r = await window.go.main.App.TestConnection(host.hostName, host.port)
      if (r.error) { alert('Failed: ' + r.error); return }
      alert('OK: ' + r.message)
    }

    const closeModal = () => { showModal.value = false; editingId.value = null }

    const formatDate = (d) => {
      if (!d) return 'N/A'
      const date = new Date(d)
      const now = new Date()
      const diff = now - date
      if (diff < 60000) return t('hosts.justNow')
      if (diff < 3600000) return Math.floor(diff / 60000) + t('hosts.minutesAgo')
      if (diff < 86400000) return Math.floor(diff / 3600000) + t('hosts.hoursAgo')
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    }

    onMounted(() => { loadHosts(); loadHistory() })

    return { t, activeTab, hosts, history, searchQuery, showModal, editingId, importing, form, tabs, loadHosts, loadHistory, searchHistory, importFromConfig, openAddModal, editHost, saveHost, deleteHost, deleteHistory, testHost, closeModal, formatDate }
  }
}
</script>

<style scoped>
.host-list {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
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
  flex-wrap: wrap;
}

.tabs {
  display: flex;
  gap: 2px;
  margin-bottom: 16px;
  background: var(--bg-secondary);
  border-radius: var(--radius-sm);
  padding: 3px;
  width: fit-content;
}

.tab-btn {
  padding: 6px 16px;
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  border-radius: 4px;
  font-size: 13px;
  font-weight: 500;
  transition: all var(--transition);
  display: flex;
  align-items: center;
  gap: 6px;
}

.tab-btn.active {
  background: var(--bg-primary);
  color: var(--text-primary);
  box-shadow: var(--shadow-sm);
}

.tab-count {
  background: var(--bg-tertiary);
  color: var(--text-muted);
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 8px;
}

.search-bar {
  position: relative;
  margin-bottom: 16px;
  max-width: 360px;
}

.search-input {
  width: 100%;
  padding: 8px 12px 8px 36px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color var(--transition);
}

.search-input:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-light);
}

.search-icon {
  position: absolute;
  left: 10px;
  top: 50%;
  transform: translateY(-50%);
  font-size: 14px;
  opacity: 0.5;
}

.content-area { flex: 1; overflow-y: auto; }

.host-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: 12px;
}

.host-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 16px;
  transition: all var(--transition);
}

.host-card:hover {
  border-color: var(--accent);
  box-shadow: var(--shadow-sm);
}

.host-card-top {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 10px;
}

.host-avatar {
  width: 40px;
  height: 40px;
  background: var(--accent-light);
  color: var(--accent);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  font-weight: 700;
  flex-shrink: 0;
}

.host-card-info { flex: 1; min-width: 0; }

.host-alias {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.host-addr {
  color: var(--text-secondary);
  font-size: 12px;
  font-family: monospace;
  margin-top: 2px;
}

.host-user {
  color: var(--text-muted);
  font-size: 12px;
  margin-top: 2px;
}

.host-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 10px;
}

.host-tag {
  font-size: 11px;
  padding: 2px 8px;
  background: var(--bg-tertiary);
  color: var(--text-secondary);
  border-radius: 4px;
}

.host-card-actions {
  display: flex;
  gap: 6px;
  justify-content: flex-end;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.history-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 14px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  transition: border-color var(--transition);
  gap: 12px;
  flex-wrap: wrap;
}

.history-item:hover { border-color: var(--border-color); }

.history-host {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 13px;
}

.history-detail {
  color: var(--text-secondary);
  font-size: 12px;
  font-family: monospace;
  margin-top: 2px;
}

.history-right {
  display: flex;
  align-items: center;
  gap: 10px;
}

.history-time {
  color: var(--text-muted);
  font-size: 12px;
  white-space: nowrap;
}

.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: var(--text-secondary);
}

.empty-icon { font-size: 48px; margin-bottom: 16px; opacity: 0.4; }
.empty-state p { margin-bottom: 20px; font-size: 15px; }
.empty-actions { display: flex; gap: 8px; justify-content: center; flex-wrap: wrap; }

.modal-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0,0,0,0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
  backdrop-filter: blur(2px);
}

.modal {
  background: var(--bg-primary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-lg);
  width: 440px;
  max-width: 100%;
  box-shadow: var(--shadow-lg);
  max-height: 90vh;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 { font-size: 16px; font-weight: 600; color: var(--text-primary); }
.modal-close { background: none; border: none; color: var(--text-secondary); font-size: 22px; cursor: pointer; padding: 4px; line-height: 1; }
.modal-close:hover { color: var(--text-primary); }
.modal-body { padding: 20px; overflow-y: auto; }
.modal-footer { display: flex; justify-content: flex-end; gap: 8px; padding: 16px 20px; border-top: 1px solid var(--border-color); }

.form-group { margin-bottom: 14px; }
.form-row { display: flex; gap: 12px; }
.form-row .form-group { flex: 1; }

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.optional { font-weight: 400; color: var(--text-muted); }

.form-group input {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-primary);
  font-size: 13px;
  outline: none;
  transition: border-color var(--transition);
}

.form-group input:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-light);
}

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
.btn-sm { padding: 5px 10px; font-size: 12px; }
.btn-primary { background: var(--accent); color: #fff; border-color: var(--accent); }
.btn-primary:hover:not(:disabled) { background: var(--accent-hover); }
.btn-secondary:hover:not(:disabled) { background: var(--bg-tertiary); }
.btn-danger { color: var(--danger); border-color: var(--danger); background: var(--danger-bg); }
.btn-danger:hover:not(:disabled) { opacity: 0.8; }

@media (max-width: 600px) {
  .host-grid { grid-template-columns: 1fr; }
  .history-item { flex-direction: column; align-items: flex-start; }
  .history-right { width: 100%; justify-content: space-between; }
  .form-row { flex-direction: column; gap: 0; }
}
</style>
