<script setup lang="ts">
import Settings from '../features/settings/Settings.vue'
</script>

<template>
  <div class="view-container">
    <div class="header">
      <h1>Settings</h1>
      <p class="subtitle">Configure AI providers and system preferences.</p>
    </div>

    <div class="card">
      <Settings />
    </div>
  </div>
</template>

<style scoped>
.view-container {
  display: flex;
  flex-direction: column;
  gap: 2rem;
  max-width: 800px; /* Narrower for settings */
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

.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: 2rem;
}
</style>
