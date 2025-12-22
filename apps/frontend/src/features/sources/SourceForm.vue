<script setup lang="ts">
import { ref } from 'vue'
import { useSourceStore } from './source.store'
import { Plus, Loader2 } from 'lucide-vue-next'

const store = useSourceStore()
const url = ref('')
const emit = defineEmits(['submit'])

async function submit() {
  if (!url.value) return
  
  // Basic URL validation
  try {
    new URL(url.value)
  } catch {
    alert('Please enter a valid URL (e.g., https://docs.example.com)')
    return
  }

  await store.addSource({ name: url.value, url: url.value })
  if (!store.error) {
    url.value = ''
    emit('submit')
  }
}
</script>

<template>
  <form @submit.prevent="submit" class="source-form">
    <div class="input-wrapper">
      <input 
        v-model="url" 
        type="text" 
        placeholder="https://docs.example.com" 
        :disabled="store.isLoading"
        class="url-input"
      />
      <button 
        type="submit" 
        :disabled="store.isLoading"
        class="submit-btn"
      >
        <Loader2 v-if="store.isLoading" class="spin" :size="18" />
        <Plus v-else :size="18" />
        <span>Ingest</span>
      </button>
    </div>
    <div v-if="store.error" class="error-msg">
      <span class="error-icon">!</span>
      {{ store.error }}
    </div>
  </form>
</template>

<style scoped>
.source-form {
  width: 100%;
}

.input-wrapper {
  display: flex;
  gap: 0.75rem;
  background: var(--color-void);
  padding: 0.5rem;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: border-color 0.2s;
}

.input-wrapper:focus-within {
  border-color: var(--color-primary);
  box-shadow: 0 0 0 1px var(--color-primary);
}

.url-input {
  flex: 1;
  background: transparent;
  border: none;
  padding: 0.5rem 1rem;
  color: var(--color-text-main);
  font-family: var(--font-mono);
  font-size: 0.95rem;
}

.url-input::placeholder {
  color: var(--color-border);
}

.url-input:focus {
  outline: none;
}

.submit-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1.25rem;
  background-color: var(--color-primary);
  color: white;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background-color: var(--color-primary-hover);
}

.submit-btn:disabled {
  background-color: var(--color-border);
  color: var(--color-text-muted);
  cursor: not-allowed;
}

.spin {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-msg {
  margin-top: 0.75rem;
  color: var(--color-danger);
  font-size: 0.85rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-family: var(--font-mono);
}
</style>