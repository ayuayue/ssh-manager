<template>
  <div class="config-editor">
    <div class="page-header">
      <div class="header-left">
        <h1>{{ t('config.title') }}</h1>
        <span class="file-path" v-if="configInfo.filePath">{{ configInfo.filePath }}</span>
      </div>
      <div class="header-actions">
        <button @click="loadConfig" class="btn btn-secondary" :disabled="loading">{{ t('config.reload') }}</button>
        <button @click="validateConfig" class="btn btn-warning">{{ t('config.validate') }}</button>
        <button @click="saveConfig" class="btn btn-primary" :disabled="!hasChanges || saving">
          {{ saving ? t('config.saving') : t('config.save') }}
        </button>
      </div>
    </div>

    <div v-if="validationResult" class="validation-banner" :class="validationResult.valid ? 'valid' : 'invalid'">
      <span class="banner-icon">{{ validationResult.valid ? '✅' : '⚠️' }}</span>
      <ul class="banner-list">
        <li v-for="(w, i) in validationResult.warnings" :key="i">{{ w }}</li>
      </ul>
      <button class="banner-close" @click="validationResult = null">&times;</button>
    </div>

    <div class="editor-layout">
      <div class="editor-panel">
        <div class="editor-wrapper">
          <textarea
            v-model="configContent"
            @input="hasChanges = true"
            spellcheck="false"
            class="config-textarea"
            :placeholder="t('config.placeholder')"
          ></textarea>
        </div>
      </div>

      <div class="preview-panel" v-if="parsedHosts.length > 0">
        <div class="panel-header">
          <h3>{{ t('config.parsedHosts') }}</h3>
          <span class="host-count">{{ parsedHosts.length }}</span>
        </div>
        <div class="host-tags">
          <div v-for="(host, i) in parsedHosts" :key="i" class="host-tag">
            <span class="tag-name">{{ host.host }}</span>
            <span class="tag-detail" v-if="host.hostName">{{ host.hostName }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, inject } from 'vue'

export default {
  name: 'ConfigEditor',
  setup() {
    const t = inject('t')
    const configContent = ref('')
    const configInfo = ref({})
    const hasChanges = ref(false)
    const loading = ref(false)
    const saving = ref(false)
    const validationResult = ref(null)

    const parsedHosts = computed(() => {
      try {
        const lines = configContent.value.split('\n')
        const hosts = []
        let current = null
        for (const line of lines) {
          const trimmed = line.trim()
          if (trimmed.toLowerCase().startsWith('host ') && !trimmed.toLowerCase().startsWith('hostname')) {
            if (current) hosts.push(current)
            current = { host: trimmed.split(/\s+/)[1], hostName: '', user: '', port: '22' }
          } else if (current) {
            if (trimmed.toLowerCase().startsWith('hostname ')) current.hostName = trimmed.split(/\s+/)[1]
            else if (trimmed.toLowerCase().startsWith('user ')) current.user = trimmed.split(/\s+/)[1]
            else if (trimmed.toLowerCase().startsWith('port ')) current.port = trimmed.split(/\s+/)[1]
          }
        }
        if (current) hosts.push(current)
        return hosts
      } catch {
        return []
      }
    })

    const loadConfig = async () => {
      loading.value = true
      try {
        const result = await window.go.main.App.GetSSHConfig()
        if (result.error) { alert('Error: ' + result.error); return }
        configContent.value = result.rawContent
        configInfo.value = { filePath: result.filePath }
        hasChanges.value = false
        validationResult.value = null
      } catch (e) { alert('Failed: ' + e.message) }
      loading.value = false
    }

    const saveConfig = async () => {
      saving.value = true
      try {
        const result = await window.go.main.App.SaveSSHConfig(configContent.value)
        if (result.error) { alert('Error: ' + result.error); return }
        hasChanges.value = false
      } catch (e) { alert('Failed: ' + e.message) }
      saving.value = false
    }

    const validateConfig = async () => {
      try {
        const result = await window.go.main.App.ValidateSSHConfig(configContent.value)
        validationResult.value = {
          valid: !result.warnings.some(w => w.toLowerCase().includes('unknown') || w.toLowerCase().includes('invalid')),
          warnings: result.warnings
        }
      } catch (e) { alert('Failed: ' + e.message) }
    }

    onMounted(loadConfig)

    return { t, configContent, configInfo, hasChanges, loading, saving, validationResult, parsedHosts, loadConfig, saveConfig, validateConfig }
  }
}
</script>

<style scoped>
.config-editor {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 12px;
}

.header-left {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}

.page-header h1 {
  font-size: 22px;
  font-weight: 600;
  color: var(--text-primary);
}

.file-path {
  color: var(--text-muted);
  font-size: 13px;
  font-family: monospace;
}

.header-actions {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.validation-banner {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  border-radius: var(--radius-md);
  margin-bottom: 16px;
  position: relative;
}

.validation-banner.valid {
  background: var(--success-bg);
  border: 1px solid var(--success);
  color: var(--success);
}

.validation-banner.invalid {
  background: var(--warning-bg);
  border: 1px solid var(--warning);
  color: var(--warning);
}

.banner-icon { font-size: 16px; flex-shrink: 0; margin-top: 1px; }

.banner-list {
  margin: 0;
  padding-left: 16px;
  font-size: 13px;
  line-height: 1.6;
  flex: 1;
}

.banner-close {
  background: none;
  border: none;
  color: inherit;
  font-size: 18px;
  cursor: pointer;
  opacity: 0.6;
  padding: 0 4px;
}

.banner-close:hover { opacity: 1; }

.editor-layout {
  display: flex;
  gap: 16px;
  flex: 1;
  min-height: 0;
}

.editor-panel {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.editor-wrapper {
  flex: 1;
  background: var(--bg-code);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 12px;
  transition: border-color var(--transition);
}

.editor-wrapper:focus-within {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-light);
}

.config-textarea {
  width: 100%;
  height: 100%;
  background: transparent;
  border: none;
  color: var(--text-primary);
  font-family: 'JetBrains Mono', 'Fira Code', 'Consolas', monospace;
  font-size: 14px;
  line-height: 1.6;
  resize: none;
  outline: none;
}

.preview-panel {
  width: 280px;
  flex-shrink: 0;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 16px;
  overflow-y: auto;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
}

.panel-header h3 {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-primary);
}

.host-count {
  background: var(--accent-light);
  color: var(--accent);
  font-size: 12px;
  font-weight: 600;
  padding: 2px 8px;
  border-radius: 10px;
}

.host-tags {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.host-tag {
  padding: 8px 10px;
  background: var(--bg-primary);
  border: 1px solid var(--border-light);
  border-radius: var(--radius-sm);
  transition: border-color var(--transition);
}

.host-tag:hover {
  border-color: var(--accent);
}

.tag-name {
  font-weight: 600;
  color: var(--accent);
  font-size: 13px;
}

.tag-detail {
  display: block;
  color: var(--text-muted);
  font-size: 12px;
  margin-top: 2px;
  font-family: monospace;
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

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--accent);
  color: #fff;
  border-color: var(--accent);
}

.btn-primary:hover:not(:disabled) {
  background: var(--accent-hover);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--bg-tertiary);
}

.btn-warning {
  background: var(--warning-bg);
  color: var(--warning);
  border-color: var(--warning);
}

.btn-warning:hover:not(:disabled) {
  opacity: 0.85;
}

@media (max-width: 900px) {
  .editor-layout {
    flex-direction: column;
  }

  .preview-panel {
    width: 100%;
    max-height: 200px;
  }

  .host-tags {
    flex-direction: row;
    flex-wrap: wrap;
  }
}

@media (max-width: 600px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
  }

  .header-actions .btn {
    flex: 1;
    justify-content: center;
  }
}
</style>
