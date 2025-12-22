<script setup lang="ts">
import { ref } from 'vue'
import { useSourceStore } from './source.store'

const store = useSourceStore()
const url = ref('')
const emit = defineEmits(['submit'])

async function submit() {
  if (!url.value) return
  
  // Basic URL validation
  try {
    new URL(url.value)
  } catch {
    alert('Please enter a valid URL')
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
    <div class="input-group">
      <input 
        v-model="url" 
        type="text" 
        placeholder="https://example.com" 
        :disabled="store.isLoading"
        class="url-input"
      />
      <button 
        type="submit" 
        :disabled="store.isLoading"
        class="submit-btn"
      >
        <span v-if="store.isLoading">Adding...</span>
        <span v-else>Add Source</span>
      </button>
    </div>
    <p v-if="store.error" class="error-msg">{{ store.error }}</p>
  </form>
</template>

<style scoped>
.source-form {
  margin-top: 1rem;
}

.input-group {
  display: flex;
  gap: 0.5rem;
}

.url-input {
  flex: 1;
  padding: 0.75rem;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 1rem;
}

.url-input:focus {
  outline: none;
  border-color: #3498db;
  box-shadow: 0 0 0 2px rgba(52, 152, 219, 0.2);
}

.submit-btn {
  padding: 0.75rem 1.5rem;
  background-color: #3498db;
  color: white;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-weight: bold;
  transition: background-color 0.2s;
}

.submit-btn:hover:not(:disabled) {
  background-color: #2980b9;
}

.submit-btn:disabled {
  background-color: #95a5a6;
  cursor: not-allowed;
}

.error-msg {
  color: #e74c3c;
  margin-top: 0.5rem;
  font-size: 0.9rem;
}
</style>