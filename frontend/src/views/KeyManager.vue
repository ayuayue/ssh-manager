<template>
  <div class="key-manager">
    <div class="page-header">
      <h1>{{ t('keys.title') }}</h1>
      <button @click="showGenModal = true" class="btn btn-primary">
        + {{ t('keys.generate') }}
      </button>
    </div>

    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <span>{{ t('keys.generating').replace('生成', '加载') }}...</span>
    </div>

    <div class="key-list" v-if="!loading">
      <div v-if="keys.length === 0" class="empty-state">
        <div class="empty-icon">🔑</div>
        <p>{{ t('keys.noKeys') }}</p>
        <button @click="showGenModal = true" class="btn btn-primary">{{ t('keys.generateFirst') }}</button>
      </div>

      <div v-for="key in keys" :key="key.name" class="key-card">
        <div class="key-card-header">
          <div class="key-icon-wrapper">
            <span class="key-icon">🔑</span>
          </div>
          <div class="key-card-info">
            <div class="key-name">{{ key.name }}</div>
            <div class="key-meta">
              <span class="key-type-badge" :class="'type-' + key.type.toLowerCase()">{{ key.type }}</span>
              <span v-if="key.fingerprint" class="key-fp" :title="key.fingerprint">{{ key.fingerprint }}</span>
            </div>
          </div>
          <div class="key-card-actions">
            <button @click="viewKey(key, 'pub')" class="btn btn-sm" :title="t('keys.viewPubKey')">👁 Pub</button>
            <button @click="viewKey(key, 'priv')" class="btn btn-sm" :title="t('keys.viewPrivKey')">🔒 Priv</button>
            <button @click="deleteKey(key.name)" class="btn btn-sm btn-danger" :title="t('keys.deleteKey')">🗑</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Generate Modal -->
    <div v-if="showGenModal" class="modal-overlay" @click.self="showGenModal = false">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ t('keys.generateKey') }}</h2>
          <button class="modal-close" @click="showGenModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="form-group">
            <label>{{ t('keys.keyType') }}</label>
            <div class="type-selector">
              <button
                v-for="t2 in keyTypes"
                :key="t2.value"
                :class="{ active: genForm.type === t2.value }"
                @click="genForm.type = t2.value"
                class="type-btn"
              >
                {{ t2.label }}
              </button>
            </div>
          </div>
          <div class="form-group" v-if="genForm.type === 'rsa'">
            <label>{{ t('keys.keySize') }}</label>
            <select v-model.number="genForm.bits">
              <option :value="2048">2048 bits</option>
              <option :value="4096">4096 bits</option>
            </select>
          </div>
          <div class="form-group" v-if="genForm.type === 'ecdsa'">
            <label>{{ t('keys.curveSize') }}</label>
            <select v-model.number="genForm.bits">
              <option :value="256">P-256</option>
              <option :value="384">P-384</option>
              <option :value="521">P-521</option>
            </select>
          </div>
          <div class="form-group">
            <label>{{ t('keys.keyName') }}</label>
            <input v-model="genForm.name" :placeholder="'id_' + genForm.type" />
          </div>
          <div class="form-group">
            <label>{{ t('keys.email') }}</label>
            <input v-model="genForm.email" placeholder="you@example.com" />
          </div>
          <div class="form-group">
            <label>{{ t('keys.passphrase') }} <span class="optional">({{ t('keys.optional') }})</span></label>
            <input v-model="genForm.passphrase" type="password" :placeholder="t('keys.passphraseHint')" />
          </div>
        </div>
        <div class="modal-footer">
          <button @click="showGenModal = false" class="btn btn-secondary">{{ t('keys.cancel') }}</button>
          <button @click="generateKey" class="btn btn-primary" :disabled="generating">
            {{ generating ? t('keys.generating') : t('keys.generateBtn') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Key Content Modal -->
    <div v-if="showKeyModal" class="modal-overlay" @click.self="showKeyModal = false">
      <div class="modal modal-lg">
        <div class="modal-header">
          <h2>{{ keyModalTitle }}: {{ selectedKey?.name }}</h2>
          <button class="modal-close" @click="showKeyModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="key-content-box" :class="{ 'is-priv': keyModalType === 'priv' }">
            <pre>{{ keyContent }}</pre>
          </div>
        </div>
        <div class="modal-footer">
          <button @click="copyKey" class="btn btn-primary">📋 {{ t('keys.copy') }}</button>
          <button @click="showKeyModal = false" class="btn btn-secondary">{{ t('keys.close') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, inject } from 'vue'

export default {
  name: 'KeyManager',
  setup() {
    const t = inject('t')
    const keys = ref([])
    const loading = ref(true)
    const showGenModal = ref(false)
    const showKeyModal = ref(false)
    const selectedKey = ref(null)
    const generating = ref(false)
    const keyModalType = ref('pub')
    const keyContent = ref('')

    const keyTypes = [
      { value: 'ed25519', label: 'ED25519' },
      { value: 'rsa', label: 'RSA' },
      { value: 'ecdsa', label: 'ECDSA' },
    ]

    const genForm = ref({ type: 'ed25519', bits: 4096, email: '', passphrase: '', name: '' })

    const keyModalTitle = computed(() => {
      return keyModalType.value === 'priv' ? t('keys.privateKey') : t('keys.publicKey')
    })

    const loadKeys = async () => {
      loading.value = true
      try {
        const result = await window.go.main.App.ListKeys()
        if (!result.error) keys.value = result.keys || []
      } catch (e) { console.error(e) }
      loading.value = false
    }

    const generateKey = async () => {
      generating.value = true
      try {
        const result = await window.go.main.App.GenerateKey(
          genForm.value.type, genForm.value.bits, genForm.value.email,
          genForm.value.passphrase, genForm.value.name
        )
        if (result.error) { alert('Error: ' + result.error); return }
        showGenModal.value = false
        genForm.value = { type: 'ed25519', bits: 4096, email: '', passphrase: '', name: '' }
        await loadKeys()
      } catch (e) { alert('Failed: ' + e.message) }
      generating.value = false
    }

    const deleteKey = async (name) => {
      if (!confirm(`${t('keys.deleteConfirm')} "${name}"?\n${t('keys.deleteWarning')}`)) return
      try {
        const result = await window.go.main.App.DeleteKey(name)
        if (result.error) { alert('Error: ' + result.error); return }
        await loadKeys()
      } catch (e) { alert('Failed: ' + e.message) }
    }

    const viewKey = async (key, type) => {
      selectedKey.value = key
      keyModalType.value = type
      try {
        if (type === 'priv') {
          const result = await window.go.main.App.GetPrivKeyContent(key.name)
          if (result.error) { alert('Error: ' + result.error); return }
          keyContent.value = result.content
        } else {
          if (key.pubKeyContent) {
            keyContent.value = key.pubKeyContent
          } else {
            const result = await window.go.main.App.GetPubKeyContent(key.name)
            if (result.error) { alert('Error: ' + result.error); return }
            keyContent.value = result.content
          }
        }
        showKeyModal.value = true
      } catch (e) { alert('Failed: ' + e.message) }
    }

    const copyKey = async () => {
      try {
        await navigator.clipboard.writeText(keyContent.value)
        alert(t('keys.copied'))
      } catch { alert(t('keys.copyFailed')) }
    }

    onMounted(loadKeys)

    return { t, keys, loading, showGenModal, showKeyModal, selectedKey, generating, keyTypes, genForm, keyModalType, keyContent, keyModalTitle, loadKeys, generateKey, deleteKey, viewKey, copyKey }
  }
}
</script>

<style scoped>
.key-manager {
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

.key-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.key-card {
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-md);
  padding: 14px 16px;
  transition: all var(--transition);
}

.key-card:hover {
  border-color: var(--accent);
  box-shadow: var(--shadow-sm);
}

.key-card-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.key-icon-wrapper {
  width: 40px;
  height: 40px;
  background: var(--accent-light);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.key-icon { font-size: 18px; }

.key-card-info { flex: 1; min-width: 0; }

.key-name {
  font-weight: 600;
  color: var(--text-primary);
  font-size: 14px;
}

.key-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  flex-wrap: wrap;
}

.key-type-badge {
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 11px;
  font-weight: 600;
}

.type-rsa { background: #dbeafe; color: #1d4ed8; }
.type-ed25519 { background: #dcfce7; color: #15803d; }
.type-ecdsa { background: #fef3c7; color: #92400e; }

.dark .type-rsa { background: rgba(96,165,250,0.15); color: #60a5fa; }
.dark .type-ed25519 { background: rgba(74,222,128,0.15); color: #4ade80; }
.dark .type-ecdsa { background: rgba(251,191,36,0.15); color: #fbbf24; }

.key-fp {
  font-family: monospace;
  font-size: 12px;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}

.key-card-actions {
  display: flex;
  gap: 6px;
  flex-shrink: 0;
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

.modal-lg { width: 600px; }

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border-color);
}

.modal-header h2 {
  font-size: 16px;
  font-weight: 600;
  color: var(--text-primary);
}

.modal-close {
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 22px;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
}

.modal-close:hover { color: var(--text-primary); }

.modal-body {
  padding: 20px;
  overflow-y: auto;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 16px 20px;
  border-top: 1px solid var(--border-color);
}

.form-group { margin-bottom: 16px; }

.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-primary);
  margin-bottom: 6px;
}

.optional { font-weight: 400; color: var(--text-muted); }

.form-group input,
.form-group select {
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

.form-group input:focus,
.form-group select:focus {
  border-color: var(--accent);
  box-shadow: 0 0 0 3px var(--accent-light);
}

.type-selector {
  display: flex;
  gap: 6px;
}

.type-btn {
  flex: 1;
  padding: 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 13px;
  font-weight: 500;
  transition: all var(--transition);
}

.type-btn.active {
  background: var(--accent-light);
  color: var(--accent);
  border-color: var(--accent);
}

.key-content-box {
  background: var(--bg-code);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-sm);
  padding: 12px;
  max-height: 300px;
  overflow: auto;
}

.key-content-box.is-priv {
  background: var(--danger-bg);
  border-color: var(--danger);
}

.key-content-box pre {
  font-family: monospace;
  font-size: 12px;
  color: var(--text-primary);
  white-space: pre-wrap;
  word-break: break-all;
  margin: 0;
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
.btn-sm { padding: 6px 8px; font-size: 12px; }
.btn-primary { background: var(--accent); color: #fff; border-color: var(--accent); }
.btn-primary:hover:not(:disabled) { background: var(--accent-hover); }
.btn-secondary:hover:not(:disabled) { background: var(--bg-tertiary); }
.btn-danger { color: var(--danger); border-color: var(--danger); background: var(--danger-bg); }
.btn-danger:hover:not(:disabled) { opacity: 0.8; }

@media (max-width: 600px) {
  .key-meta { flex-direction: column; align-items: flex-start; gap: 4px; }
  .key-fp { max-width: 100%; }
  .modal { margin: 10px; }
}
</style>
