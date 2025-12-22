<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from './settings.store'
import { Save, Loader2 } from 'lucide-vue-next'

const store = useSettingsStore()

onMounted(() => {
  store.fetchSettings()
})
</script>

<template>
  <div class="settings-form">
    <div v-if="store.error" class="msg error">{{ store.error }}</div>
    <div v-if="store.successMessage" class="msg success">{{ store.successMessage }}</div>

    <div class="form-group">
      <label for="geminiKey">Gemini API Key</label>
      <input 
        id="geminiKey" 
        type="password" 
        v-model="store.geminiApiKey" 
        class="input" 
        placeholder="Enter Gemini API Key"
      />
      <p class="hint">Required for generating embeddings (Google AI Studio).</p>
    </div>

    <div class="form-group">
      <label for="provider">Rerank Provider</label>
      <div class="select-wrapper">
        <select id="provider" v-model="store.rerankProvider" class="input">
          <option value="none">None</option>
          <option value="jina">Jina AI</option>
          <option value="cohere">Cohere</option>
        </select>
      </div>
      <p class="hint">Select an external provider to re-rank search results for better accuracy.</p>
    </div>

    <div class="form-group" v-if="store.rerankProvider !== 'none'">
      <label for="apiKey">API Key</label>
      <input 
        id="apiKey" 
        type="password" 
        v-model="store.rerankApiKey" 
        class="input" 
        placeholder="Enter API Key"
      />
    </div>

    <button 
      class="btn-primary" 
      @click="store.updateSettings" 
      :disabled="store.isLoading"
    >
      <Loader2 v-if="store.isLoading" class="spin" :size="18" />
      <Save v-else :size="18" />
      <span>{{ store.isLoading ? 'Saving...' : 'Save Configuration' }}</span>
    </button>
  </div>
</template>

<style scoped>
.settings-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 500px;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

label {
  font-weight: 500;
  font-size: 0.9rem;
  color: var(--color-text-main);
}

.input {
  padding: 0.75rem;
  background: var(--color-void);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  color: var(--color-text-main);
  font-family: var(--font-mono);
  font-size: 0.95rem;
  width: 100%;
}

.input:focus {
  outline: none;
  border-color: var(--color-primary);
  box-shadow: 0 0 0 1px var(--color-primary);
}

.select-wrapper {
  position: relative;
}

.hint {
  font-size: 0.8rem;
  color: var(--color-text-muted);
  margin: 0;
}

.btn-primary {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  font-weight: 600;
  cursor: pointer;
  align-self: flex-start;
  transition: all 0.2s;
}

.btn-primary:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.btn-primary:disabled {
  background-color: var(--color-border);
  cursor: not-allowed;
  color: var(--color-text-muted);
}

.msg {
  padding: 0.75rem 1rem;
  border-radius: var(--radius-md);
  font-size: 0.9rem;
  border: 1px solid transparent;
}

.msg.error {
  color: var(--color-danger);
  background: rgba(239, 68, 68, 0.1);
  border-color: rgba(239, 68, 68, 0.2);
}

.msg.success {
  color: var(--color-success);
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
