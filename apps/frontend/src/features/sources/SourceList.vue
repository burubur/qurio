<script setup lang="ts">
import { useSourceStore } from './source.store'
import { onMounted } from 'vue'
import { RefreshCw, Trash2, ExternalLink } from 'lucide-vue-next'
import StatusBadge from '../../components/ui/StatusBadge.vue'

const store = useSourceStore()

const handleDelete = async (id: string) => {
  if (confirm('Are you sure you want to delete this source?')) {
    await store.deleteSource(id)
  }
}

const handleResync = async (id: string) => {
  await store.resyncSource(id)
}

onMounted(() => {
  store.fetchSources()
})
</script>

<template>
  <div class="source-list-container">
    <div v-if="store.isLoading && store.sources.length === 0" class="loading">
      <div class="loading-spinner"></div>
      <span>Retrieving knowledge sources...</span>
    </div>
    
    <div v-else-if="store.sources.length === 0" class="empty">
      No sources configured. Ingest documentation to begin.
    </div>

    <ul v-else class="source-list">
      <li v-for="source in store.sources" :key="source.id" class="source-item">
        <div class="source-main">
          <div class="source-header">
            <span class="source-url">{{ source.url }}</span>
            <a :href="source.url" target="_blank" class="external-link">
              <ExternalLink :size="14" />
            </a>
          </div>
          <div class="source-meta">
             <span class="id">ID: {{ source.id.substring(0, 8) }}</span>
          </div>
        </div>
        
        <div class="source-controls">
          <StatusBadge :status="source.status || 'pending'" />
          
          <div class="actions">
            <button @click="handleResync(source.id)" class="btn-icon" title="Re-sync">
              <RefreshCw :size="18" />
            </button>
            <button @click="handleDelete(source.id)" class="btn-icon delete" title="Delete">
              <Trash2 :size="18" />
            </button>
          </div>
        </div>
      </li>
    </ul>
  </div>
</template>

<style scoped>
.source-list-container {
  width: 100%;
}

.loading, .empty {
  text-align: center;
  color: var(--color-text-muted);
  padding: 3rem;
  background: rgba(255, 255, 255, 0.02);
  border-radius: var(--radius-md);
  border: 1px dashed var(--color-border);
  font-family: var(--font-mono);
  font-size: 0.9rem;
}

.loading-spinner {
  width: 24px;
  height: 24px;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-primary);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.source-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.source-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  transition: all 0.2s ease;
}

.source-item:hover {
  background: rgba(255, 255, 255, 0.04);
  border-color: var(--color-primary);
}

.source-main {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.source-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.source-url {
  font-family: var(--font-mono);
  font-weight: 500;
  color: var(--color-text-main);
  font-size: 0.95rem;
}

.external-link {
  color: var(--color-text-muted);
  opacity: 0;
  transition: opacity 0.2s;
  display: flex;
  align-items: center;
}

.source-item:hover .external-link {
  opacity: 1;
}

.source-meta {
  display: flex;
  gap: 1rem;
  font-size: 0.75rem;
  color: var(--color-text-muted);
  font-family: var(--font-mono);
}

.source-controls {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.actions {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: transparent;
  border: 1px solid transparent;
  color: var(--color-text-muted);
  cursor: pointer;
  padding: 0.5rem;
  border-radius: var(--radius-sm);
  transition: all 0.2s;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  background: rgba(255, 255, 255, 0.1);
  color: var(--color-text-main);
}

.btn-icon.delete:hover {
  background: rgba(239, 68, 68, 0.1);
  color: var(--color-danger);
  border-color: rgba(239, 68, 68, 0.2);
}
</style>
