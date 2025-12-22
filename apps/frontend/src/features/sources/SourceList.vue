<script setup lang="ts">
import { useSourceStore } from './source.store'
import { onMounted } from 'vue'

const store = useSourceStore()

onMounted(() => {
  store.fetchSources()
})
</script>

<template>
  <div class="source-list-container">
    <div v-if="store.isLoading && store.sources.length === 0" class="loading">
      Loading sources...
    </div>
    
    <div v-else-if="store.sources.length === 0" class="empty">
      No sources added yet. Add one above!
    </div>

    <ul v-else class="source-list">
      <li v-for="source in store.sources" :key="source.id" class="source-item">
        <div class="source-info">
          <span class="source-url">{{ source.url }}</span>
          <span class="source-status" :class="source.status || 'pending'">
            {{ source.status || 'Pending' }}
          </span>
        </div>
        <!-- Future: Add actions like delete/re-sync here -->
      </li>
    </ul>
  </div>
</template>

<style scoped>
.source-list-container {
  margin-top: 1rem;
}

.loading, .empty {
  text-align: center;
  color: #7f8c8d;
  padding: 2rem;
  background: #f8f9fa;
  border-radius: 4px;
}

.source-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.source-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #eee;
  transition: background-color 0.2s;
}

.source-item:last-child {
  border-bottom: none;
}

.source-item:hover {
  background-color: #f8f9fa;
}

.source-url {
  font-weight: 500;
  color: #2c3e50;
}

.source-status {
  font-size: 0.8rem;
  padding: 0.25rem 0.5rem;
  border-radius: 12px;
  text-transform: uppercase;
  font-weight: bold;
}

.source-status.pending {
  background-color: #f1c40f;
  color: #fff;
}

.source-status.indexed {
  background-color: #2ecc71;
  color: #fff;
}

.source-status.failed {
  background-color: #e74c3c;
  color: #fff;
}
</style>
