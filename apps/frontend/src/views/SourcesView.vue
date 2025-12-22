<script setup lang="ts">
import SourceForm from '../features/sources/SourceForm.vue'
import SourceList from '../features/sources/SourceList.vue'
</script>

<template>
  <div class="view-container">
    <div class="header">
      <h1>Sources</h1>
      <p class="subtitle">Manage ingestion targets for your knowledge base.</p>
    </div>

    <div class="grid">
      <div class="card">
        <h2>Add Source</h2>
        <p class="card-desc">Enter a URL to ingest documentation.</p>
        <SourceForm />
      </div>

      <div class="card full-width">
        <h2>Active Sources</h2>
        <SourceList />
      </div>
    </div>
  </div>
</template>

<style scoped>
.view-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  margin-bottom: 1rem;
}

h1 {
  font-size: 2.5rem;
  color: var(--color-text-main);
}

.subtitle {
  color: var(--color-text-muted);
  font-size: 1.1rem;
}

.grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 1.5rem;
}

.card-desc {
  color: var(--color-text-muted);
  margin-bottom: 1.5rem;
  font-size: 0.9rem;
}

h2 {
  font-size: 1.25rem;
  margin-bottom: 0.5rem;
  color: var(--color-text-main);
}
</style>
